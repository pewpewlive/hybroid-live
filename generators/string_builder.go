// Copyright 2017 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//

package generators

import (
	"unicode/utf8"
	"unsafe"
)

// A BetterBuilder is used to efficiently build a string using [BetterBuilder.Write] methods.
// It minimizes memory copying. The zero value is ready to use.
// Do not copy a non-zero BetterBuilder.
type BetterBuilder struct {
	addr *BetterBuilder // of receiver, to detect copies by value
	buf  []byte
}

// noescape hides a pointer from escape analysis. It is the identity function
// but escape analysis doesn't think the output depends on the input.
// noescape is inlined and currently compiles down to zero instructions.
// USE CAREFULLY!
// This was copied from the runtime; see issues 23382 and 7921.
//
//go:nosplit
//go:nocheckptr
func noescape(p unsafe.Pointer) unsafe.Pointer {
	x := uintptr(p)
	return unsafe.Pointer(x ^ 0)
}

func (b *BetterBuilder) copyCheck() {
	if b.addr == nil {
		// This hack works around a failing of Go's escape analysis
		// that was causing b to escape and be heap allocated.
		// See issue 23382.
		// TODO: once issue 7921 is fixed, this should be reverted to
		// just "b.addr = b".
		b.addr = (*BetterBuilder)(noescape(unsafe.Pointer(b)))
	} else if b.addr != b {
		panic("strings: illegal use of non-zero BetterBuilder copied by value")
	}
}

// String returns the accumulated string.
func (b *BetterBuilder) String() string {
	return unsafe.String(unsafe.SliceData(b.buf), len(b.buf))
}

// Len returns the number of accumulated bytes; b.Len() == len(b.String()).
func (b *BetterBuilder) Len() int { return len(b.buf) }

// Cap returns the capacity of the BetterBuilder's underlying byte slice. It is the
// total space allocated for the string being built and includes any bytes
// already written.
func (b *BetterBuilder) Cap() int { return cap(b.buf) }

// Reset resets the [BetterBuilder] to be empty.
func (b *BetterBuilder) Reset() {
	b.addr = nil
	b.buf = nil
}

/*
// grow copies the buffer to a new, larger buffer so that there are at least n
// bytes of capacity beyond len(b.buf).
func (b *BetterBuilder) grow(n int) {
	buf := bytealg.MakeNoZero(2*cap(b.buf) + n)[:len(b.buf)]
	copy(buf, b.buf)
	b.buf = buf
}

// Grow grows b's capacity, if necessary, to guarantee space for
// another n bytes. After Grow(n), at least n bytes can be written to b
// without another allocation. If n is negative, Grow panics.
func (b *BetterBuilder) Grow(n int) {
	b.copyCheck()
	if n < 0 {
		panic("strings.BetterBuilder.Grow: negative count")
	}
	if cap(b.buf)-len(b.buf) < n {
		b.grow(n)
	}

}

*/

func (b *BetterBuilder) Append(strs ...string) (int, error) {
	b.copyCheck()
	for i := range strs {
		b.buf = append(b.buf, strs[i]...)
	}

	return len(strs), nil
}

// Write appends the contents of p to b's buffer.
// Write always returns len(p), nil.
func (b *BetterBuilder) Write(p []byte) (int, error) {
	b.copyCheck()
	b.buf = append(b.buf, p...)
	return len(p), nil
}

// WriteByte appends the byte c to b's buffer.
// The returned error is always nil.
func (b *BetterBuilder) WriteByte(c byte) error {
	b.copyCheck()
	b.buf = append(b.buf, c)
	return nil
}

// WriteRune appends the UTF-8 encoding of Unicode code point r to b's buffer.
// It returns the length of r and a nil error.
func (b *BetterBuilder) WriteRune(r rune) (int, error) {
	b.copyCheck()
	n := len(b.buf)
	b.buf = utf8.AppendRune(b.buf, r)
	return len(b.buf) - n, nil
}

// WriteString appends the contents of s to b's buffer.
// It returns the length of s and a nil error.
func (b *BetterBuilder) WriteString(s string) (int, error) {
	b.copyCheck()
	b.buf = append(b.buf, s...)
	return len(s), nil
}

func (b *BetterBuilder) GetBuf() *[]byte {
	b.copyCheck()
	return &b.buf
}
