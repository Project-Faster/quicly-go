// None

// WARNING: This file has automatically been generated
// Code generated by https://git.io/c-for-go. DO NOT EDIT.

package bindings

/*
#cgo LDFLAGS: C:/home/dev/src/github.com/parvit/quicly-go/internal/deps/lib/libquicly.a C:/home/dev/src/github.com/parvit/quicly-go/internal/deps/lib/libcrypto.a C:/home/dev/src/github.com/parvit/quicly-go/internal/deps/lib/libssl.a -lm -lmswsock -lws2_32
#cgo CPPFLAGS: -IC:/home/dev/src/github.com/parvit/quicly-go/internal/deps/include/
#include "quicly.h"
#include "quicly_wrapper.h"
#include "quicly/streambuf.h"
#include <stdlib.h>
#include "cgo_helpers.h"
*/
import "C"

// Size_t type as declared in include/crtdefs.h:35
type Size_t uint64

// Iovec as declared in include/quicly.h:42
type Iovec struct {
	Iov_base      []byte
	Iov_len       Size_t
	ref4b778f8    *C.struct_iovec
	allocs4b778f8 interface{}
}