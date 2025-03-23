package quiclylib

import "C"
import (
	"context"
	"github.com/Project-Faster/quicly-go/internal"
	"github.com/Project-Faster/quicly-go/quiclylib/errors"
	"github.com/Project-Faster/quicly-go/quiclylib/types"
	log "github.com/rs/zerolog"
	"math/rand"
	"net"
	"runtime"
	"sync"
	"time"
)

type QClientSession struct {
	// exported fields
	NetConn *net.UDPConn
	Ctx     context.Context
	Logger  log.Logger

	// callback
	types.Callbacks

	// unexported fields
	id             uint64
	connected      bool
	closing        bool
	ctxCancel      context.CancelFunc
	handlersWaiter sync.WaitGroup

	streams     map[uint64]types.Stream
	streamsLock sync.RWMutex

	exclusiveLock sync.RWMutex

	incomingQueue chan *types.Packet
}

var _ net.Listener = &QClientSession{}
var _ types.Session = &QClientSession{}

func (s *QClientSession) enterCritical(readonly bool) {
	// s.Logger.Warn().Msgf("Will Enter Critical section (%v)", readonly)
	if readonly {
		s.streamsLock.RLock()
	} else {
		s.streamsLock.Lock()
	}
	// s.Logger.Warn().Msgf("Enter Critical section (%v)", readonly)
}
func (s *QClientSession) exitCritical(readonly bool) {
	// s.Logger.Warn().Msgf("Will Exit Critical section (%v)", readonly)
	if readonly {
		s.streamsLock.RUnlock()
	} else {
		s.streamsLock.Unlock()
	}
	// s.Logger.Warn().Msgf("Exit Critical section (%v)", readonly)
}

func (s *QClientSession) init() {
	if s.incomingQueue == nil {
		s.Ctx, s.ctxCancel = context.WithCancel(s.Ctx)

		s.id = rand.Uint64()
		s.Logger.Info().Msgf("Client handler init: %v", s.id)

		s.incomingQueue = make(chan *types.Packet, 1024)
		s.streams = make(map[uint64]types.Stream)
	} else {
		s.Logger.Warn().Msgf("Client handler was already init: %v", s.id)
	}
}

func (s *QClientSession) connect() int {
	if s.connected {
		return errors.QUICLY_OK
	}

	s.init()

	udpAddr := s.Addr().(*net.UDPAddr)

	if ret := QuiclyConnect(udpAddr.IP.String(), int32(udpAddr.Port), s.id); ret != errors.QUICLY_OK {
		return int(ret)
	}

	internal.RegisterConnection(s, s.id)

	s.connected = true
	s.closing = false

	s.handlersWaiter.Add(2)
	go s.connectionInHandler()
	go s.connectionProcessHandler()

	go func() {
		s.handlersWaiter.Wait()
		s.Close()
	}()

	if s.OnConnectionOpen != nil {
		s.OnConnectionOpen(s)
	}

	return errors.QUICLY_OK
}

// --- Handlers routines --- //

func (s *QClientSession) connectionInHandler() {
	defer func() {
		_ = recover()
		runtime.UnlockOSThread()
		s.ctxCancel()
		s.handlersWaiter.Done()
	}()

	runtime.LockOSThread()

	var buffList = make([][]byte, 0, 32)
	for i := 0; i < 4096; i++ {
		buffList = append(buffList, make([]byte, SMALL_BUFFER_SIZE))
	}

	s.Logger.Debug().Msgf("CONN IN START %v", s.id)
	defer s.Logger.Debug().Msgf("CONN IN END %v", s.id)

	var counter = 0

	for {
		select {
		case <-s.Ctx.Done():
			s.Logger.Debug().Msgf("CONN IN STOP %v", s.id)
			return
		default:
			break
		}

		s.NetConn.SetReadDeadline(time.Now().Add(1 * time.Second))

		n, addr, err := s.NetConn.ReadFromUDP(buffList[0])
		s.Logger.Debug().Msgf("[%v] UDP packet %d %v (%d)", s.id, n, addr, counter)
		if n == 0 || (n == 0 && err != nil) {
			s.Logger.Debug().Msgf("QUICLY No packet")
			continue
		}
		counter++

		buf := buffList[0]
		pkt := &types.Packet{
			Data:       buf[:n],
			DataLen:    n,
			RetAddress: addr,
		}
		s.incomingQueue <- pkt

		buffList = buffList[1:]
		buffList = append(buffList, make([]byte, SMALL_BUFFER_SIZE))
	}
}

func (s *QClientSession) connectionProcessHandler() {
	returnAddr := s.NetConn.RemoteAddr().(*net.UDPAddr)
	defer func() {
		_ = recover()
		s.ctxCancel()
		s.handlersWaiter.Done()
	}()

	s.Logger.Debug().Msgf("CONN PROC START %v", s.id)
	defer s.Logger.Debug().Msgf("CONN PROC END %v", s.id)

	buffer := make([]*types.Packet, 0, 32)

	for {
		select {
		case <-s.Ctx.Done():
			s.Logger.Debug().Msgf("CONN PROC STOP %v", s.id)
			return

		case pkt := <-s.incomingQueue:
			s.Logger.Debug().Msgf("[%v] RECV packet (%d : %d)", s.id, pkt.Streamid, pkt.DataLen)
			buffer = append(buffer, pkt)
			break

		case <-time.After(5 * time.Millisecond):
			for _, pkt := range buffer {
				if len(s.streams) == 0 {
					s.Logger.Debug().Msgf("[%v] No active streams", s.id)
					break
				}

				s.Logger.Debug().Msgf("[%v] PROC packet %v %d(%v)", s.id, s.id, pkt.DataLen, pkt.Streamid)
				if pkt == nil {
					break
				}
				addr, port := pkt.Address()
				if len(addr) == 0 || port == -1 {
					addr, port = returnAddr.IP.String(), returnAddr.Port
				}

				err := QuiclyProcessMsg(int32(1), addr, int32(port), pkt.Data, uint64(pkt.DataLen), s.id)
				if err != errors.QUICLY_OK {
					if err == errors.QUICLY_ERROR_PACKET_IGNORED {
						s.Logger.Error().Msgf("[%v] Process error %d bytes (ignored processing %v)", s.id, pkt.DataLen, err)
					} else {
						s.Logger.Error().Msgf("[%v] Received %d bytes (failed processing %v)", s.id, pkt.DataLen, err)
					}
				}
			}
			buffer = buffer[:0]

			if ret := s.flushOutgoingQueue(); ret != errors.QUICLY_OK {
				return
			}
			break
		}
	}
}

func (s *QClientSession) flushOutgoingQueue() int32 {
	num_packets := uint64(4096)
	packets_buf := make([]Iovec, 4096)

	var ret = QuiclyOutgoingMsgQueue(s.id, packets_buf, &num_packets)
	if int(num_packets) == 0 {
		s.Logger.Debug().Msgf("QUICLY flushOutgoingQueue %d: 0", s.id)
		return ret
	}

	switch ret {
	case errors.QUICLY_ERROR_NOT_OPEN:
		s.Logger.Error().Msgf("QUICLY Send failed: errors.QUICLY_ERROR_NOT_OPEN")
		return ret
	default:
		s.Logger.Debug().Msgf("QUICLY Send failed: %d - %v", num_packets, ret)
		return ret
	case errors.QUICLY_OK:
		break
	}

	s.Logger.Debug().Msgf("CONN flush (%d) %v", num_packets, s.id)

	s.enterCritical(true)
	for i := 0; i < int(num_packets); i++ {
		packets_buf[i].Deref() // realize the struct copy from C -> go

		data := IovecToBytes(packets_buf[i])

		_ = s.NetConn.SetWriteDeadline(time.Now().Add(WRITE_TIMEOUT))

		n, err := s.NetConn.Write(data)
		s.Logger.Debug().Msgf("[%v] SEND packet %d bytes [%v]", s.id, n, err)
	}
	s.exitCritical(true)

	runtime.KeepAlive(num_packets)
	runtime.KeepAlive(packets_buf)

	return ret
}

// --- Session interface --- //

func (s *QClientSession) ID() uint64 {
	return s.id
}

func (s *QClientSession) OpenStream() types.Stream {
	s.enterCritical(false)
	if err := s.connect(); err != errors.QUICLY_OK {
		s.Logger.Error().Msgf("connect error: %d", err)
		s.exitCritical(false)
		return nil
	}
	s.exitCritical(false)

	var streamId uint64 = 0

	if ret := QuiclyOpenStream(s.id, &streamId); ret != errors.QUICLY_OK {
		s.Logger.Debug().Msgf("open stream err")
		return nil
	}

	st := &QStream{
		session: s,
		conn:    s.NetConn,
		id:      streamId,
		Logger:  s.Logger,
	}
	st.init()

	s.enterCritical(false)
	s.streams[st.id] = st
	s.exitCritical(false)

	return st
}

func (s *QClientSession) GetStream(id uint64) types.Stream {
	s.enterCritical(true)
	defer s.exitCritical(true)
	return s.streams[id]
}

func (s *QClientSession) StreamPacket(packet *types.Packet) {
	//defer func() {
	//	_ = recover()
	//}()
	//if packet == nil || s.outgoingQueue == nil {
	//	return
	//}
	//select {
	//case s.outgoingQueue <- packet:
	//	break
	//case <-time.After(3 * time.Millisecond):
	//	return
	//}
}

func (s *QClientSession) OnStreamOpen(streamId uint64) {
	s.enterCritical(true)
	st, ok := s.streams[streamId]
	s.exitCritical(true)
	if ok {
		if s.OnStreamOpenCallback != nil {
			s.OnStreamOpenCallback(st)
		}
		st.OnOpened()
	}
}

func (s *QClientSession) OnStreamClose(streamId uint64, error int) {
	s.Logger.Debug().Msgf("STREAM CLOSE START: %d\n", streamId)
	defer s.Logger.Debug().Msgf("STREAM CLOSE END: %d\n", streamId)

	s.enterCritical(false)
	st, ok := s.streams[streamId]
	delete(s.streams, streamId)
	shouldTerm := len(s.streams) == 0
	s.exitCritical(false)
	if ok {
		_ = st.OnClosed()
	}

	if ok && s.OnStreamCloseCallback != nil {
		s.OnStreamCloseCallback(st, error)
	}

	if shouldTerm {
		s.Logger.Debug().Msgf(">> Closing parent: %d\n", s.id)
		s.ctxCancel()
		return
	}
}

// --- Listener interface --- //

func (s *QClientSession) Accept() (net.Conn, error) {
	return nil, net.ErrClosed
}

func (s *QClientSession) Close() error {
	defer func() {
		_ = recover()
	}()
	if !s.connected || s.closing || s == nil || s.NetConn == nil {
		return nil
	}
	s.Logger.Debug().Msgf("== Connections %v WaitEnd ==\"", s.id)
	defer s.Logger.Info().Msgf("== Connections %v End ==\"", s.id)

	s.enterCritical(false)

	s.closing = true

	s.ctxCancel()
	// copy of stream list is to workaround lock issues
	tmpStreams := make([]types.Stream, len(s.streams))
	pos := 0
	for _, stream := range s.streams {
		tmpStreams[pos] = stream
		pos++
	}
	s.exitCritical(false)

	wg := &sync.WaitGroup{}
	wg.Add(len(tmpStreams))

	for _, stream := range tmpStreams {
		go func(st types.Stream) {
			defer wg.Done()

			s.Logger.Debug().Msgf(">> Trying to close stream %d / %p", st.ID(), st)
			st.Close()
			s.Logger.Debug().Msgf(">> Closed stream %d / %p", st.ID(), st)
		}(stream)
	}

	go func() {
		defer func() {
			s.closing = false
			s.connected = false
		}()
		wg.Wait()

		s.enterCritical(false)
		s.Logger.Debug().Msgf(">> Close queues %d(%v)", s.id, s.id)
		safeClose(s.incomingQueue)
		_ = s.NetConn.Close()
		s.exitCritical(false)

		if s.OnConnectionClose != nil {
			s.Logger.Debug().Msgf("Close connection: %d\n", s.id)

			s.OnConnectionClose(s)
		}
		s.Logger.Debug().Msgf(">> Wait routines %d(%v)", s.id, s.id)
		s.handlersWaiter.Wait()

		var err = QuiclyClose(s.id, 0)

		internal.RemoveConnection(s.id)
		s.Logger.Debug().Msgf(">> Quicly Close %d: %v", s.id, err)
	}()

	return nil
}

func (s *QClientSession) IsClosed() bool {
	return !s.connected
}

func (s *QClientSession) Addr() net.Addr {
	return s.NetConn.RemoteAddr()
}
