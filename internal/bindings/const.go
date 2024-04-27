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

const (
	// QUICLY_ERROR_PACKET_IGNORED as defined in quicly/constants.h:104
	QUICLY_ERROR_PACKET_IGNORED = 0xff01
	// QUICLY_ERROR_SENDBUF_FULL as defined in quicly/constants.h:105
	QUICLY_ERROR_SENDBUF_FULL = 0xff02
	// QUICLY_ERROR_FREE_CONNECTION as defined in quicly/constants.h:106
	QUICLY_ERROR_FREE_CONNECTION = 0xff03
	// QUICLY_ERROR_RECEIVED_STATELESS_RESET as defined in quicly/constants.h:107
	QUICLY_ERROR_RECEIVED_STATELESS_RESET = 0xff04
	// QUICLY_ERROR_NO_COMPATIBLE_VERSION as defined in quicly/constants.h:108
	QUICLY_ERROR_NO_COMPATIBLE_VERSION = 0xff05
	// QUICLY_ERROR_IS_CLOSING as defined in quicly/constants.h:109
	QUICLY_ERROR_IS_CLOSING = 0xff06
	// QUICLY_ERROR_STATE_EXHAUSTION as defined in quicly/constants.h:110
	QUICLY_ERROR_STATE_EXHAUSTION = 0xff07
	// QUICLY_ERROR_INVALID_INITIAL_VERSION as defined in quicly/constants.h:111
	QUICLY_ERROR_INVALID_INITIAL_VERSION = 0xff08
	// QUICLY_ERROR_DECRYPTION_FAILED as defined in quicly/constants.h:112
	QUICLY_ERROR_DECRYPTION_FAILED = 0xff09
)

const ()

const (
	// QUICLY_OK as declared in include/quicly_wrapper.h:11
	QUICLY_OK = iota
	// QUICLY_ERROR_NOTINITILIZED as declared in include/quicly_wrapper.h:12
	QUICLY_ERROR_NOTINITILIZED = 1
	// QUICLY_ERROR_ALREADY_INIT as declared in include/quicly_wrapper.h:13
	QUICLY_ERROR_ALREADY_INIT = 2
	// QUICLY_ERROR_FAILED as declared in include/quicly_wrapper.h:14
	QUICLY_ERROR_FAILED = 3
	// QUICLY_ERROR_CERT_LOAD_FAILED as declared in include/quicly_wrapper.h:15
	QUICLY_ERROR_CERT_LOAD_FAILED = 4
	// QUICLY_ERROR_DECODE_FAILED as declared in include/quicly_wrapper.h:16
	QUICLY_ERROR_DECODE_FAILED = 5
	// QUICLY_ERROR_DESTINATION_NOT_FOUND as declared in include/quicly_wrapper.h:17
	QUICLY_ERROR_DESTINATION_NOT_FOUND = 6
	// QUICLY_ERROR_NOT_OPEN as declared in include/quicly_wrapper.h:18
	QUICLY_ERROR_NOT_OPEN = 7
	// QUICLY_ERROR_STREAM_NOT_FOUND as declared in include/quicly_wrapper.h:19
	QUICLY_ERROR_STREAM_NOT_FOUND = 8
	// QUICLY_ERROR_UNKNOWN_CC_ALGO as declared in include/quicly_wrapper.h:20
	QUICLY_ERROR_UNKNOWN_CC_ALGO = 9
	// QUICLY_ERROR_CANNOT_SEND as declared in include/quicly_wrapper.h:21
	QUICLY_ERROR_CANNOT_SEND = 10
	// QUICLY_ERROR_STREAM_BUSY as declared in include/quicly_wrapper.h:22
	QUICLY_ERROR_STREAM_BUSY = 11
)

const ()

const ()

const ()

const ()

const ()

const ()
