package bindings

/*
#include "quicly_wrapper.h"
*/
import "C"

import (
	"fmt"
	"github.com/Project-Faster/quicly-go/quiclylib/types"
	"sync"
)

var connectionsRegistry map[uint64]types.Session
var callbackLock sync.Mutex

func ResetRegistry() {
	callbackLock.Lock()
	defer callbackLock.Unlock()
	connectionsRegistry = make(map[uint64]types.Session)
}

func RegisterConnection(s types.Session, id uint64) {
	callbackLock.Lock()
	defer callbackLock.Unlock()
	connectionsRegistry[id] = s
}

func RemoveConnection(id uint64) {
	callbackLock.Lock()
	defer callbackLock.Unlock()
	delete(connectionsRegistry, id)
}

//export goQuiclyOnStreamOpen
func goQuiclyOnStreamOpen(conn_id C.uint64_t, stream_id C.uint64_t) {
	fmt.Printf("open stream: %d %d\n", uint64(conn_id), uint64(stream_id))

	callbackLock.Lock()
	conn, ok := connectionsRegistry[uint64(conn_id)]
	callbackLock.Unlock()

	if !ok {
		return
	}

	conn.OnStreamOpen(uint64(stream_id))
}

//export goQuiclyOnStreamClose
func goQuiclyOnStreamClose(conn_id C.uint64_t, stream_id C.uint64_t, error C.int) {
	fmt.Printf("close stream: %d %d\n", uint64(conn_id), uint64(stream_id))

	callbackLock.Lock()
	conn, ok := connectionsRegistry[uint64(conn_id)]
	callbackLock.Unlock()

	if !ok {
		return
	}

	conn.OnStreamClose(uint64(stream_id), int(error))
}

//export goQuiclyOnStreamReceived
func goQuiclyOnStreamReceived(conn_id C.uint64_t, stream_id C.uint64_t, data *C.struct_iovec) {
	callbackLock.Lock()
	conn, ok := connectionsRegistry[uint64(conn_id)]
	callbackLock.Unlock()

	if !ok {
		fmt.Printf("could not find connection: %d %d\n", uint64(conn_id), uint64(stream_id))
		return
	}

	st := conn.GetStream(uint64(stream_id))
	if st == nil {
		fmt.Printf("could not find stream: %d %d\n", uint64(conn_id), uint64(stream_id))
		return
	}

	vIn := Iovec{}
	vIn.ref4b778f8 = data
	vIn.Deref()

	buf := IovecToBytes(vIn)
	st.OnReceived(buf, int(vIn.Iov_len))
}

func IovecToBytes(data Iovec) []byte {
	return data.Iov_base
}
