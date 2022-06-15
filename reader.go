// {{{ Copyright (c) Paul R. Tagliamonte <paul@k3xec.com>, 2021
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE. }}}

package pulseaudio

import (
	"fmt"
	"unsafe"
)

// #cgo pkg-config: libpulse libpulse-simple
//
// #include <pulse/simple.h>
import "C"

// Reader is a handle to the underlying pulseaudio stream to allow for reading
// from the microphone on the system.
type Reader struct {
	simple     *C.pa_simple
	sampleSpec C.pa_sample_spec

	config Config
}

// Flush will empty the buffer held by pulseaudio.
func (r Reader) Flush() error {
	var errCode C.int
	if C.pa_simple_flush(r.simple, &errCode) >= 0 {
		return nil
	}
	return fmt.Errorf("pulseaudio: bad flush, err code %d", int(errCode))
}

// UnsafeRead will perform a very unsafe read into the target location in
// memory. Please avoid using this.
func (r Reader) UnsafeRead(ptr unsafe.Pointer, size int) error {
	var errCode C.int

	if C.pa_simple_read(
		r.simple,
		ptr,
		C.size_t(size),
		&errCode,
	) >= 0 {
		return nil
	}

	return fmt.Errorf("pulseaudio: bad write, err code %d", int(errCode))
}

// Close will free the underlying handle.
func (r Reader) Close() {
	C.pa_simple_free(r.simple)
}

// Read will read the audio samples into the provided object.
func (r Reader) Read(in interface{}) error {
	switch data := in.(type) {
	case []float32:
		return r.UnsafeRead(
			unsafe.Pointer(&(data[0])),
			int(unsafe.Sizeof(data[0]))*len(data),
		)
	default:
		return fmt.Errorf("pulseaudio: unknown read type")
	}
}

// NewReader will create a new Reader with the provided configuration.
func NewReader(cfg Config) (*Reader, error) {
	var (
		paErr C.int
		attr  *C.pa_buffer_attr
	)

	switch cfg.Format {
	case SampleFormatFloat32NE:
		break
	default:
		return nil, fmt.Errorf("pulseaudio: unsupported stream format")
	}

	sampleSpec := C.pa_sample_spec{}
	sampleSpec.format = C.pa_sample_format_t(cfg.Format)
	sampleSpec.channels = C.uint8_t(cfg.Channels)
	sampleSpec.rate = C.uint32_t(cfg.Rate)

	// TODO(paultag): Add in support for frag, in particular.
	// attr = &C.pa_buffer_attr{}

	if C.pa_channels_valid(sampleSpec.channels) == 0 {
		return nil, fmt.Errorf(
			"pulseaudio: channels '%d' is invalid",
			sampleSpec.channels,
		)
	}

	if C.pa_sample_rate_valid(sampleSpec.rate) == 0 {
		return nil, fmt.Errorf(
			"pulseaudio: sample rate '%d' is invalid",
			sampleSpec.rate,
		)
	}

	simple := C.pa_simple_new(
		nil,
		C.CString(cfg.AppName),
		C.PA_STREAM_RECORD,
		nil,
		C.CString(cfg.StreamName),
		&sampleSpec,
		nil,
		attr,
		&paErr,
	)

	if simple == nil {
		return nil, rvToErr(paErr)
	}

	return &Reader{
		simple:     simple,
		sampleSpec: sampleSpec,
		config:     cfg,
	}, nil
}

// vim: foldmethod=marker
