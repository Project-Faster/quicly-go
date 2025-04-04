// None

// WARNING: This file has automatically been generated
// Code generated by https://git.io/c-for-go. DO NOT EDIT.

package bindings

/*
#cgo LDFLAGS: -Llib/ -l:libquicly.a -l:libcrypto.a -l:libssl.a -lm -lmswsock -lws2_32
#cgo CPPFLAGS: -DWIN32 -I../../../ -Iinclude/ -Iinclude/quicly/ -Wno-format
#include "quicly.h"
#include "quicly_wrapper.h"
#include "quicly/streambuf.h"
#include <stdlib.h>
#include "cgo_helpers.h"
*/
import "C"
import (
	"fmt"
	"runtime"
	"sync"
	"unsafe"
)

// cgoAllocMap stores pointers to C allocated memory for future reference.
type cgoAllocMap struct {
	mux sync.RWMutex
	m   map[unsafe.Pointer]struct{}
}

var cgoAllocsUnknown = new(cgoAllocMap)

func (a *cgoAllocMap) Add(ptr unsafe.Pointer) {
	a.mux.Lock()
	if a.m == nil {
		a.m = make(map[unsafe.Pointer]struct{})
	}
	a.m[ptr] = struct{}{}
	a.mux.Unlock()
}

func (a *cgoAllocMap) IsEmpty() bool {
	a.mux.RLock()
	isEmpty := len(a.m) == 0
	a.mux.RUnlock()
	return isEmpty
}

func (a *cgoAllocMap) Borrow(b *cgoAllocMap) {
	if b == nil || b.IsEmpty() {
		return
	}
	b.mux.Lock()
	a.mux.Lock()
	for ptr := range b.m {
		if a.m == nil {
			a.m = make(map[unsafe.Pointer]struct{})
		}
		a.m[ptr] = struct{}{}
		delete(b.m, ptr)
	}
	a.mux.Unlock()
	b.mux.Unlock()
}

func (a *cgoAllocMap) Free() {
	a.mux.Lock()
	for ptr := range a.m {
		C.free(ptr)
		delete(a.m, ptr)
	}
	a.mux.Unlock()
}

// allocStruct_iovecMemory allocates memory for type C.struct_iovec in C.
// The caller is responsible for freeing the this memory via C.free.
func allocStruct_iovecMemory(n int) unsafe.Pointer {
	mem, err := C.calloc(C.size_t(n), (C.size_t)(sizeOfStruct_iovecValue))
	if mem == nil {
		panic(fmt.Sprintln("memory alloc error: ", err))
	}
	return mem
}

const sizeOfStruct_iovecValue = unsafe.Sizeof([1]C.struct_iovec{})

// copyPCharBytes copies the data from Go slice as *C.char.
func copyPCharBytes(slice *sliceHeader) (*C.char, *cgoAllocMap) {
	allocs := new(cgoAllocMap)
	defer runtime.SetFinalizer(allocs, func(a *cgoAllocMap) {
		go a.Free()
	})

	mem0 := unsafe.Pointer(C.CBytes(*(*[]byte)(unsafe.Pointer(&sliceHeader{
		Data: slice.Data,
		Len:  int(sizeOfCharValue) * slice.Len,
		Cap:  int(sizeOfCharValue) * slice.Len,
	}))))
	allocs.Add(mem0)

	return (*C.char)(mem0), allocs
}

type sliceHeader struct {
	Data unsafe.Pointer
	Len  int
	Cap  int
}

// allocCharMemory allocates memory for type C.char in C.
// The caller is responsible for freeing the this memory via C.free.
func allocCharMemory(n int) unsafe.Pointer {
	mem, err := C.calloc(C.size_t(n), (C.size_t)(sizeOfCharValue))
	if mem == nil {
		panic(fmt.Sprintln("memory alloc error: ", err))
	}
	return mem
}

const sizeOfCharValue = unsafe.Sizeof([1]C.char{})

// Ref returns the underlying reference to C object or nil if struct is nil.
func (x *Iovec) Ref() *C.struct_iovec {
	if x == nil {
		return nil
	}
	return x.ref4b778f8
}

// Free invokes alloc map's free mechanism that cleanups any allocated memory using C free.
// Does nothing if struct is nil or has no allocation map.
func (x *Iovec) Free() {
	if x != nil && x.allocs4b778f8 != nil {
		x.allocs4b778f8.(*cgoAllocMap).Free()
		x.ref4b778f8 = nil
	}
}

// NewIovecRef creates a new wrapper struct with underlying reference set to the original C object.
// Returns nil if the provided pointer to C object is nil too.
func NewIovecRef(ref unsafe.Pointer) *Iovec {
	if ref == nil {
		return nil
	}
	obj := new(Iovec)
	obj.ref4b778f8 = (*C.struct_iovec)(unsafe.Pointer(ref))
	return obj
}

// PassRef returns the underlying C object, otherwise it will allocate one and set its values
// from this wrapping struct, counting allocations into an allocation map.
func (x *Iovec) PassRef() (*C.struct_iovec, *cgoAllocMap) {
	if x == nil {
		return nil, nil
	} else if x.ref4b778f8 != nil {
		return x.ref4b778f8, nil
	}
	mem4b778f8 := allocStruct_iovecMemory(1)
	ref4b778f8 := (*C.struct_iovec)(mem4b778f8)
	allocs4b778f8 := new(cgoAllocMap)
	allocs4b778f8.Add(mem4b778f8)

	var ciov_base_allocs *cgoAllocMap
	ref4b778f8.iov_base, ciov_base_allocs = copyPCharBytes((*sliceHeader)(unsafe.Pointer(&x.Iov_base)))
	allocs4b778f8.Borrow(ciov_base_allocs)

	var ciov_len_allocs *cgoAllocMap
	ref4b778f8.iov_len, ciov_len_allocs = (C.size_t)(x.Iov_len), cgoAllocsUnknown
	allocs4b778f8.Borrow(ciov_len_allocs)

	x.ref4b778f8 = ref4b778f8
	x.allocs4b778f8 = allocs4b778f8
	return ref4b778f8, allocs4b778f8

}

// PassValue does the same as PassRef except that it will try to dereference the returned pointer.
func (x Iovec) PassValue() (C.struct_iovec, *cgoAllocMap) {
	if x.ref4b778f8 != nil {
		return *x.ref4b778f8, nil
	}
	ref, allocs := x.PassRef()
	return *ref, allocs
}

// Deref uses the underlying reference to C object and fills the wrapping struct with values.
// Do not forget to call this method whether you get a struct for C object and want to read its values.
func (x *Iovec) Deref() {
	if x.ref4b778f8 == nil {
		return
	}
	hxfc4425b := (*sliceHeader)(unsafe.Pointer(&x.Iov_base))
	hxfc4425b.Data = unsafe.Pointer(x.ref4b778f8.iov_base)
	hxfc4425b.Cap = 0x7fffffff
	hxfc4425b.Len = int(x.ref4b778f8.iov_len)

	x.Iov_len = (Size_t)(x.ref4b778f8.iov_len)
}

// safeString ensures that the string is NULL-terminated, a NULL-terminated copy is created otherwise.
func safeString(str string) string {
	if len(str) > 0 && str[len(str)-1] != '\x00' {
		str = str + "\x00"
	} else if len(str) == 0 {
		str = "\x00"
	}
	return str
}

// unpackPCharString copies the data from Go string as *C.char.
func unpackPCharString(str string) (*C.char, *cgoAllocMap) {
	allocs := new(cgoAllocMap)
	defer runtime.SetFinalizer(allocs, func(a *cgoAllocMap) {
		go a.Free()
	})

	str = safeString(str)
	mem0 := unsafe.Pointer(C.CString(str))
	allocs.Add(mem0)
	return (*C.char)(mem0), allocs
}

type stringHeader struct {
	Data unsafe.Pointer
	Len  int
}

const sizeOfPtr = unsafe.Sizeof(&struct{}{})

// unpackArgSIovec transforms a sliced Go data structure into plain C format.
func unpackArgSIovec(x []Iovec) (unpacked *C.struct_iovec, allocs *cgoAllocMap) {
	if x == nil {
		return nil, nil
	}
	allocs = new(cgoAllocMap)
	defer runtime.SetFinalizer(allocs, func(a *cgoAllocMap) {
		go a.Free()
	})

	len0 := len(x)
	mem0 := allocStruct_iovecMemory(len0)
	allocs.Add(mem0)
	h0 := &sliceHeader{
		Data: mem0,
		Cap:  len0,
		Len:  len0,
	}
	v0 := *(*[]C.struct_iovec)(unsafe.Pointer(h0))
	for i0 := range x {
		allocs0 := new(cgoAllocMap)
		v0[i0], allocs0 = x[i0].PassValue()
		allocs.Borrow(allocs0)
	}
	h := (*sliceHeader)(unsafe.Pointer(&v0))
	unpacked = (*C.struct_iovec)(h.Data)
	return
}

// packSIovec reads sliced Go data structure out from plain C format.
func packSIovec(v []Iovec, ptr0 *C.struct_iovec) {
	const m = 0x7fffffff
	for i0 := range v {
		ptr1 := (*(*[m / sizeOfStruct_iovecValue]C.struct_iovec)(unsafe.Pointer(ptr0)))[i0]
		v[i0] = *NewIovecRef(unsafe.Pointer(&ptr1))
	}
}
