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

type writer struct {
	simple     *C.pa_simple
	sampleSpec C.pa_sample_spec

	config Config
}

func (s writer) UnsafeWrite(ptr unsafe.Pointer, size int) error {
	var errCode C.int

	if C.pa_simple_write(
		s.simple,
		ptr,
		C.size_t(size),
		&errCode,
	) >= 0 {
		return nil
	}

	return fmt.Errorf("pulseaudio: bad write, err code %d", int(errCode))
}

func (s writer) Close() {
	C.pa_simple_free(s.simple)
}

func newWriter(cfg Config) (*writer, error) {
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

	if cfg.MaxLength != 0 {
		attr = &C.pa_buffer_attr{}
		attr.maxlength = C.uint32_t(cfg.MaxLength)
	}

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
		C.PA_STREAM_PLAYBACK,
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

	return &writer{
		simple:     simple,
		sampleSpec: sampleSpec,
		config:     cfg,
	}, nil
}

// NewWriter creates a new pulseaudio.Writer object, allowing audio to be
// written into Pulse.
func NewWriter(cfg Config) (*Writer, error) {
	w, err := newWriter(cfg)
	if err != nil {
		return nil, err
	}
	return &Writer{writer: *w}, nil
}

// Writer is an encapsulation of the underlying pulseaudio stream, allowing
// for an idiomatic Go library interface to that stream.
type Writer struct {
	writer writer
}

// Close will close the stream, and preform all cleanup required.
func (w Writer) Close() {
	w.writer.Close()
}

// write2F32NE is called from `Write` when the provided input type is of
// type [][2]float32
func (w Writer) write2F32NE(samples [][2]float32) error {
	if w.writer.config.Format != SampleFormatFloat32NE {
		return fmt.Errorf("pulseaudio: stream format isn't float32 NE")
	}
	if w.writer.config.Channels != 2 {
		return fmt.Errorf("pulseaudio: wrong number of channels provided")
	}
	return w.writer.UnsafeWrite(
		unsafe.Pointer(&samples[0]),
		len(samples)*int(unsafe.Sizeof(samples[0])),
	)
}

// writeF32NE is called from `Write` when the provided input type is of
// type []float32
func (w Writer) writeF32NE(samples []float32) error {
	if w.writer.config.Format != SampleFormatFloat32NE {
		return fmt.Errorf("pulseaudio: stream format isn't float32 NE")
	}
	if w.writer.config.Channels != 1 {
		return fmt.Errorf("pulseaudio: wrong number of channels provided")
	}
	return w.writer.UnsafeWrite(
		unsafe.Pointer(&samples[0]),
		len(samples)*int(unsafe.Sizeof(samples[0])),
	)
}

// Write will write the provided data to the pulseaudio stream. This function
// accepts any type, but will return an error if th type is not something that
// this library understands, the number of channels and the stream format.
//
// Currently accepted types:
//
//   - []float32    (Format: SampleFormatFloat32NE, Channels: 1)
//   - [][2]float32 (Format: SampleFormatFloat32NE, Channels: 2)
//
func (w Writer) Write(samples interface{}) error {
	switch s := samples.(type) {
	case [][2]float32:
		return w.write2F32NE(s)
	case []float32:
		return w.writeF32NE(s)
	default:
		return fmt.Errorf("pulseaudio: unknown sample type")
	}
}

// vim: foldmethod=marker
