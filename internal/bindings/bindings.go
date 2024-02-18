// None

// WARNING: This file has automatically been generated
// Code generated by https://git.io/c-for-go. DO NOT EDIT.

package bindings

/*
#cgo LDFLAGS: C:/home/dev/src/github.com/Project-Faster/qpep-faster/backend/quicly-go/internal/deps/lib/libquicly.a C:/home/dev/src/github.com/Project-Faster/qpep-faster/backend/quicly-go/internal/deps/lib/libcrypto.a C:/home/dev/src/github.com/Project-Faster/qpep-faster/backend/quicly-go/internal/deps/lib/libssl.a -lm -lmswsock -lws2_32
#cgo CPPFLAGS: -DWIN32 -IC:/home/dev/src/github.com/Project-Faster/qpep-faster/backend/quicly-go/internal/deps/include/
#include "quicly.h"
#include "quicly_wrapper.h"
#include "quicly/streambuf.h"
#include <stdlib.h>
#include "cgo_helpers.h"
*/
import "C"
import (
	"runtime"
	"unsafe"
)

// QuiclyInitializeEngine function as declared in include/quicly_wrapper.h:35
func QuiclyInitializeEngine(Is_client uint64, Alpn string, Certificate_file string, Key_file string, Idle_timeout_ms uint64, Cc_algo uint64, Trace_quicly uint64) int32 {
	cIs_client, cIs_clientAllocMap := (C.uint64_t)(Is_client), cgoAllocsUnknown
	Alpn = safeString(Alpn)
	cAlpn, cAlpnAllocMap := unpackPCharString(Alpn)
	Certificate_file = safeString(Certificate_file)
	cCertificate_file, cCertificate_fileAllocMap := unpackPCharString(Certificate_file)
	Key_file = safeString(Key_file)
	cKey_file, cKey_fileAllocMap := unpackPCharString(Key_file)
	cIdle_timeout_ms, cIdle_timeout_msAllocMap := (C.uint64_t)(Idle_timeout_ms), cgoAllocsUnknown
	cCc_algo, cCc_algoAllocMap := (C.uint64_t)(Cc_algo), cgoAllocsUnknown
	cTrace_quicly, cTrace_quiclyAllocMap := (C.uint64_t)(Trace_quicly), cgoAllocsUnknown
	__ret := C.QuiclyInitializeEngine(cIs_client, cAlpn, cCertificate_file, cKey_file, cIdle_timeout_ms, cCc_algo, cTrace_quicly)
	runtime.KeepAlive(cTrace_quiclyAllocMap)
	runtime.KeepAlive(cCc_algoAllocMap)
	runtime.KeepAlive(cIdle_timeout_msAllocMap)
	runtime.KeepAlive(Key_file)
	runtime.KeepAlive(cKey_fileAllocMap)
	runtime.KeepAlive(Certificate_file)
	runtime.KeepAlive(cCertificate_fileAllocMap)
	runtime.KeepAlive(Alpn)
	runtime.KeepAlive(cAlpnAllocMap)
	runtime.KeepAlive(cIs_clientAllocMap)
	__v := (int32)(__ret)
	return __v
}

// QuiclyCloseEngine function as declared in include/quicly_wrapper.h:40
func QuiclyCloseEngine() int32 {
	__ret := C.QuiclyCloseEngine()
	__v := (int32)(__ret)
	return __v
}

// QuiclyProcessMsg function as declared in include/quicly_wrapper.h:42
func QuiclyProcessMsg(Is_client int32, Address string, Port int32, Msg []byte, Dgram_len Size_t, Id *Size_t) int32 {
	cIs_client, cIs_clientAllocMap := (C.int)(Is_client), cgoAllocsUnknown
	Address = safeString(Address)
	cAddress, cAddressAllocMap := unpackPCharString(Address)
	cPort, cPortAllocMap := (C.int)(Port), cgoAllocsUnknown
	cMsg, cMsgAllocMap := copyPCharBytes((*sliceHeader)(unsafe.Pointer(&Msg)))
	cDgram_len, cDgram_lenAllocMap := (C.size_t)(Dgram_len), cgoAllocsUnknown
	cId, cIdAllocMap := (*C.size_t)(unsafe.Pointer(Id)), cgoAllocsUnknown
	__ret := C.QuiclyProcessMsg(cIs_client, cAddress, cPort, cMsg, cDgram_len, cId)
	runtime.KeepAlive(cIdAllocMap)
	runtime.KeepAlive(cDgram_lenAllocMap)
	runtime.KeepAlive(cMsgAllocMap)
	runtime.KeepAlive(cPortAllocMap)
	runtime.KeepAlive(Address)
	runtime.KeepAlive(cAddressAllocMap)
	runtime.KeepAlive(cIs_clientAllocMap)
	__v := (int32)(__ret)
	return __v
}

// QuiclyConnect function as declared in include/quicly_wrapper.h:44
func QuiclyConnect(Address string, Port int32, Id *Size_t) int32 {
	Address = safeString(Address)
	cAddress, cAddressAllocMap := unpackPCharString(Address)
	cPort, cPortAllocMap := (C.int)(Port), cgoAllocsUnknown
	cId, cIdAllocMap := (*C.size_t)(unsafe.Pointer(Id)), cgoAllocsUnknown
	__ret := C.QuiclyConnect(cAddress, cPort, cId)
	runtime.KeepAlive(cIdAllocMap)
	runtime.KeepAlive(cPortAllocMap)
	runtime.KeepAlive(Address)
	runtime.KeepAlive(cAddressAllocMap)
	__v := (int32)(__ret)
	return __v
}

// QuiclyOpenStream function as declared in include/quicly_wrapper.h:46
func QuiclyOpenStream(Conn_id Size_t, Stream_id *Size_t) int32 {
	cConn_id, cConn_idAllocMap := (C.size_t)(Conn_id), cgoAllocsUnknown
	cStream_id, cStream_idAllocMap := (*C.size_t)(unsafe.Pointer(Stream_id)), cgoAllocsUnknown
	__ret := C.QuiclyOpenStream(cConn_id, cStream_id)
	runtime.KeepAlive(cStream_idAllocMap)
	runtime.KeepAlive(cConn_idAllocMap)
	__v := (int32)(__ret)
	return __v
}

// QuiclyCloseStream function as declared in include/quicly_wrapper.h:48
func QuiclyCloseStream(Conn_id Size_t, Stream_id Size_t, Error int32) int32 {
	cConn_id, cConn_idAllocMap := (C.size_t)(Conn_id), cgoAllocsUnknown
	cStream_id, cStream_idAllocMap := (C.size_t)(Stream_id), cgoAllocsUnknown
	cError, cErrorAllocMap := (C.int)(Error), cgoAllocsUnknown
	__ret := C.QuiclyCloseStream(cConn_id, cStream_id, cError)
	runtime.KeepAlive(cErrorAllocMap)
	runtime.KeepAlive(cStream_idAllocMap)
	runtime.KeepAlive(cConn_idAllocMap)
	__v := (int32)(__ret)
	return __v
}

// QuiclyClose function as declared in include/quicly_wrapper.h:50
func QuiclyClose(Conn_id Size_t, Error int32) int32 {
	cConn_id, cConn_idAllocMap := (C.size_t)(Conn_id), cgoAllocsUnknown
	cError, cErrorAllocMap := (C.int)(Error), cgoAllocsUnknown
	__ret := C.QuiclyClose(cConn_id, cError)
	runtime.KeepAlive(cErrorAllocMap)
	runtime.KeepAlive(cConn_idAllocMap)
	__v := (int32)(__ret)
	return __v
}

// QuiclyOutgoingMsgQueue function as declared in include/quicly_wrapper.h:52
func QuiclyOutgoingMsgQueue(Id Size_t, Dgram []Iovec, Num_dgrams *Size_t) int32 {
	cId, cIdAllocMap := (C.size_t)(Id), cgoAllocsUnknown
	cDgram, cDgramAllocMap := unpackArgSIovec(Dgram)
	cNum_dgrams, cNum_dgramsAllocMap := (*C.size_t)(unsafe.Pointer(Num_dgrams)), cgoAllocsUnknown
	__ret := C.QuiclyOutgoingMsgQueue(cId, cDgram, cNum_dgrams)
	runtime.KeepAlive(cNum_dgramsAllocMap)
	packSIovec(Dgram, cDgram)
	runtime.KeepAlive(cDgramAllocMap)
	runtime.KeepAlive(cIdAllocMap)
	__v := (int32)(__ret)
	return __v
}

// QuiclyWriteStream function as declared in include/quicly_wrapper.h:54
func QuiclyWriteStream(Conn_id Size_t, Stream_id Size_t, Msg []byte, Dgram_len Size_t) int32 {
	cConn_id, cConn_idAllocMap := (C.size_t)(Conn_id), cgoAllocsUnknown
	cStream_id, cStream_idAllocMap := (C.size_t)(Stream_id), cgoAllocsUnknown
	cMsg, cMsgAllocMap := copyPCharBytes((*sliceHeader)(unsafe.Pointer(&Msg)))
	cDgram_len, cDgram_lenAllocMap := (C.size_t)(Dgram_len), cgoAllocsUnknown
	__ret := C.QuiclyWriteStream(cConn_id, cStream_id, cMsg, cDgram_len)
	runtime.KeepAlive(cDgram_lenAllocMap)
	runtime.KeepAlive(cMsgAllocMap)
	runtime.KeepAlive(cStream_idAllocMap)
	runtime.KeepAlive(cConn_idAllocMap)
	__v := (int32)(__ret)
	return __v
}

// QuiclyCanSendStream function as declared in include/quicly_wrapper.h:56
func QuiclyCanSendStream(Conn_id Size_t, Stream_id Size_t) int32 {
	cConn_id, cConn_idAllocMap := (C.size_t)(Conn_id), cgoAllocsUnknown
	cStream_id, cStream_idAllocMap := (C.size_t)(Stream_id), cgoAllocsUnknown
	__ret := C.QuiclyCanSendStream(cConn_id, cStream_id)
	runtime.KeepAlive(cStream_idAllocMap)
	runtime.KeepAlive(cConn_idAllocMap)
	__v := (int32)(__ret)
	return __v
}
