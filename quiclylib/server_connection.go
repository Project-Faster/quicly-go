package quiclylib

import (
	"context"
	"github.com/Project-Faster/quicly-go/internal"
	"github.com/Project-Faster/quicly-go/quiclylib/errors"
	"github.com/Project-Faster/quicly-go/quiclylib/types"
	"io"
	"math/rand"
	"net"
	"runtime"
	"runtime/debug"
	"sync"
	"time"

	log "github.com/rs/zerolog"
)

type QServerConnection struct {
	// exported fields
	Ctx     context.Context
	NetConn *net.UDPConn
	Logger  log.Logger

	// unexported fields
	id         uint64
	started    bool
	accepted   bool
	closing    bool
	session    *QServerSession
	returnAddr *net.UDPAddr
	returnHash uint64

	exclusiveLock sync.RWMutex

	streams           map[uint64]types.Stream
	streamsLock       sync.RWMutex
	streamAcceptQueue chan types.Stream

	cancelFunc     context.CancelFunc
	routinesWaiter sync.WaitGroup

	incomingQueue chan *types.Packet
}

const (
	MAX_CONNECTIONS = 8192
)

var _ types.Session = &QServerConnection{}

func (r *QServerConnection) enterCritical(readonly bool) {
	//r.Logger.Warn().Msgf("Enter Critical section >>")
	if readonly {
		r.streamsLock.RLock()
	} else {
		r.streamsLock.Lock()
	}
	//r.Logger.Warn().Msgf("Enter Critical section <<")
}
func (r *QServerConnection) exitCritical(readonly bool) {
	//r.Logger.Warn().Msgf("Exit Critical section >>")
	if readonly {
		r.streamsLock.RUnlock()
	} else {
		r.streamsLock.Unlock()
	}
	//r.Logger.Warn().Msgf("Exit Critical section <<")
}

func (r *QServerConnection) init(session *QServerSession, addr *net.UDPAddr, addrHash uint64) {
	if r.started {
		session.Logger.Warn().Msgf("Server handler was already initialized: %v (%v)", addr, r.id)
		return
	}

	r.id = rand.Uint64()
	session.Logger.Info().Msgf("QServerConnection init: %v (%v)", addr, r.id)

	if session == nil || addr == nil {
		panic(errors.QUICLY_ERROR_FAILED)
	}

	r.started = true
	r.accepted = false
	r.closing = false

	r.session = session
	r.Ctx, r.cancelFunc = context.WithCancel(context.Background())
	r.NetConn = session.NetConn
	r.Logger = session.Logger
	r.returnAddr = addr
	r.returnHash = addrHash

	r.streams = make(map[uint64]types.Stream)
	r.streamAcceptQueue = make(chan types.Stream, 256)

	r.incomingQueue = make(chan *types.Packet, 4092)

	r.routinesWaiter.Add(2)

	internal.RegisterConnection(r, r.id)

	go r.connectionProcess()
	go r.connectionOutgoing()

	go func() {
		r.routinesWaiter.Wait()
		r.Close()
	}()
}

func (r *QServerConnection) receiveIncomingPacket(pkt *types.Packet) {
	defer func() {
		_ = recover()
	}()
	r.enterCritical(false)
	defer r.exitCritical(false)

	if !r.started || pkt == nil {
		return
	}
	if pkt.RetAddress == nil {
		pkt.RetAddress = r.returnAddr
	}
	r.incomingQueue <- pkt
}

func (r *QServerConnection) handleProcessPacket(pkt *types.Packet) int32 {
	addr, port := pkt.Address()

	return QuiclyProcessMsg(int32(0), addr, int32(port), pkt.Data, uint64(pkt.DataLen), r.id)
}

func (r *QServerConnection) flushOutgoingQueue() int32 {
	if !r.started {
		return errors.QUICLY_ERROR_NOTINITILIZED
	}

	num_packets := uint64(4096)
	packets_buf := make([]Iovec, 4096)

	connId := r.id

	var ret = QuiclyOutgoingMsgQueue(connId, packets_buf, &num_packets)
	if int(num_packets) == 0 {
		r.Logger.Debug().Msgf("QUICLY flushOutgoingQueue %d: 0", connId)
		return errors.QUICLY_OK
	}

	switch ret {
	default:
		r.Logger.Debug().Msgf("QUICLY flushOutgoingQueue failed: %v", ret)
		return ret
	case errors.QUICLY_OK:
		break
	}

	r.Logger.Debug().Msgf("CONN flush (%d) %v", num_packets, r.id)

	r.enterCritical(false)
	for i := 0; i < int(num_packets); i++ {
		packets_buf[i].Deref() // realize the struct copy from C -> go

		data := IovecToBytes(packets_buf[i])

		_ = r.NetConn.SetWriteDeadline(time.Now().Add(WRITE_TIMEOUT))

		n, err := r.NetConn.WriteToUDP(data, r.returnAddr)
		r.Logger.Debug().Msgf("[%v] WRITE packet %d bytes [%v]", r.id, n, err)
	}
	r.exitCritical(false)

	runtime.KeepAlive(num_packets)
	runtime.KeepAlive(packets_buf)

	return errors.QUICLY_OK
}

// --- Handlers routines --- //

func (r *QServerConnection) connectionProcess() {
	defer func() {
		r.routinesWaiter.Done()
		r.cancelFunc()
		if err := recover(); err != nil {
			r.Logger.Error().Msgf("PANIC: %v %v\n", err, string(debug.Stack()))
			//_ = r.Close()
		}
	}()

	r.Logger.Debug().Msgf("CONN PROC START %v", r.id)
	defer r.Logger.Debug().Msgf("CONN PROC END %v", r.id)

	buffer := make([]*types.Packet, 0, 32)

	for {
		select {
		case <-r.Ctx.Done():
			return

		case pkt := <-r.incomingQueue:
			if pkt != nil {
				r.Logger.Debug().Msgf("CONN IN BUFF %v", pkt.DataLen)
				buffer = append(buffer, pkt)
			}
			break

		case <-time.After(5 * time.Millisecond):
			for i, pkt := range buffer {
				if pkt == nil {
					r.Logger.Error().Msgf("CONN PROC ERR %v", pkt)
					break
				}
				r.Logger.Debug().Msgf("CONN PROC %v (%d)", pkt.DataLen, i)

				ret := r.handleProcessPacket(pkt)
				switch ret {
				case errors.QUICLY_ERROR_NOT_OPEN:
					r.Logger.Error().Msgf("QUICLY Send failed: errors.QUICLY_ERROR_NOT_OPEN")
					return
				case errors.QUICLY_ERROR_PACKET_IGNORED:
					r.Logger.Error().Msgf("[%v] Process error %d bytes (ignored processing %v)", r.id, pkt.DataLen, ret)
					continue
				default:
					r.Logger.Error().Msgf("[%v] Received %d bytes (failed processing %v)", r.id, pkt.DataLen, ret)
					break
				case errors.QUICLY_OK:
					r.accepted = true
					break
				}
			}
			buffer = buffer[:0]

			if ret := r.flushOutgoingQueue(); ret != errors.QUICLY_OK {
				r.Logger.Debug().Msgf("CONN PROC ERR %v", ret)
				return
			}
			break
		}
	}
}

func (r *QServerConnection) connectionOutgoing() {
	defer func() {
		r.routinesWaiter.Done()
		r.cancelFunc()
		if err := recover(); err != nil {
			r.Logger.Error().Msgf("PANIC: %v %v\n", err, string(debug.Stack()))
			debug.PrintStack()
		}
	}()

	for !r.accepted {
		<-time.After(1 * time.Millisecond)
	}

	r.Logger.Debug().Msgf("CONN OUT START %v", r.id)
	defer r.Logger.Debug().Msgf("CONN OUT END %v", r.id)

	for {
		<-time.After(1 * time.Millisecond)

		select {
		case <-r.Ctx.Done():
			return
		default:
			break
		}

		ret := r.flushOutgoingQueue()
		switch ret {
		default:
			continue
		case errors.QUICLY_ERROR_NOT_OPEN:
			fallthrough
		case errors.QUICLY_ERROR_DESTINATION_NOT_FOUND:
			r.Logger.Error().Msgf("QUICLY Send failed: %v", ret)
			return
		case errors.QUICLY_OK:
			break
		}
	}
}

// --- Session interface --- //

func (r *QServerConnection) StreamPacket(packet *types.Packet) {
	if !r.started || packet == nil {
		return
	}
	//r.Logger.Debug().Msgf("ON SEND PACKET")
	//defer r.Logger.Debug().Msgf("ON PACKET SENT")
	//r.outgoingQueue <- packet
}

func (r *QServerConnection) ID() uint64 {
	return r.id
}

func (r *QServerConnection) OpenStream() types.Stream {
	return nil
}

func (r *QServerConnection) OnStreamOpen(streamId uint64) {
	if len(r.streams) == 0 {
		r.session.enqueueConnAccept(r)
		r.session.OnConnectionOpen(r)
	}

	st := &QStream{
		session: r,
		conn:    r.NetConn,
		id:      streamId,
		Logger:  r.Logger,
	}
	r.enterCritical(false)
	r.streams[streamId] = st
	st.init()
	st.OnOpened()
	r.exitCritical(false)

	r.streamAcceptQueue <- st

	r.session.OnStreamOpen(streamId)
}

func (r *QServerConnection) OnStreamClose(streamId uint64, code int) {
	r.Logger.Debug().Msgf(">> On close stream: %d\n", streamId)
	defer r.Logger.Debug().Msgf("<< On close stream: %d\n", streamId)

	st := r.GetStream(streamId)
	if st == nil {
		return
	}

	_ = st.OnClosed()

	r.enterCritical(false)
	delete(r.streams, streamId)
	if len(r.streams) == 0 {
		r.cancelFunc()
	}
	r.exitCritical(false)

	if r.session.OnStreamCloseCallback != nil {
		r.session.OnStreamCloseCallback(st, code)
	}
}

func (r *QServerConnection) GetStream(id uint64) types.Stream {
	r.enterCritical(true)
	defer r.exitCritical(true)
	return r.streams[id]
}

// --- Listener interface --- //

func (r *QServerConnection) Accept() (types.Stream, error) {
	for {
		select {
		case st := <-r.streamAcceptQueue:
			return st, nil
		case <-r.Ctx.Done():
			return nil, io.ErrClosedPipe
		case <-time.After(1 * time.Millisecond):
			break
		}
	}
}

func (r *QServerConnection) Close() error {
	r.enterCritical(false)
	if !r.started || r.closing {
		r.exitCritical(false)
		return nil
	}
	r.closing = true

	r.Logger.Debug().Msgf("== Connections %v WaitEnd ==\"", r.id)

	r.cancelFunc()

	wg := &sync.WaitGroup{}
	wg.Add(len(r.streams))

	for _, stream := range r.streams {
		go func(st types.Stream) {
			defer wg.Done()

			r.Logger.Debug().Msgf(">> Trying to close stream %d:%d", r.id, st.ID())
			err := st.Close()
			r.Logger.Debug().Msgf(">> Closed stream %d:%d (%v)", r.id, st.ID(), err)
		}(stream)
	}

	go func() {
		defer func() {
			r.closing = false
			r.started = false
			r.Logger.Info().Msgf("== Connections %v End ==\"", r.id)
		}()
		wg.Wait()

		safeClose(r.incomingQueue)
		safeClose(r.streamAcceptQueue)

		r.exitCritical(false)

		r.routinesWaiter.Wait()

		<-time.After(2 * time.Second) // linger connection

		var ptrId = r.id
		var err = QuiclyClose(ptrId, 0)

		r.session.connectionDelete(r.id)

		r.Logger.Warn().Msgf(">> Quicly Close %d(%v): %v", r.id, r.id, err)
	}()

	return nil
}

func (r *QServerConnection) IsClosed() bool {
	return !r.started
}

func (r *QServerConnection) Addr() net.Addr {
	return r.returnAddr
}
