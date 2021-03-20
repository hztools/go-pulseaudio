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

// #cgo pkg-config: libpulse libpulse-simple
//
// #include <pulse/simple.h>
import "C"

// Config will allow you to set the parameters that pulseaudio will be initialized
// with at creation time.
type Config struct {
	// Format denoes the Audio format type. Currently only "NE" or
	// "Native Endian" is supported, since we're providing a way to write
	// things like `[]float` into the audio stream, and it's unlikely this
	// will be used to play back flat files on the system without processing
	// them.
	Format SampleFormat

	// Rate is the numbers of samples per second.
	Rate uint

	// Channels is the number of channels in the audio stream.
	Channels uint

	// AppName is the name of the Application doing the streaming.
	AppName string

	// StreamName is the name of the audio stream from the Application.
	StreamName string

	// MaxLength will set the pa_buffer_attr's maxlength to the set number
	// of bytes.
	MaxLength uint
}

// SampleFormat denotes the audio format to be used by the pulse streamer.
type SampleFormat int

var (
	// SampleFormatFloat32NE is used when the provided data is comprised of
	// 32 bit floating point numbers in the native endian format. This means
	// the data can be passed to `Write` as either []float32 or [][2]float32,
	// depending on the channels provided.
	SampleFormatFloat32NE SampleFormat = C.PA_SAMPLE_FLOAT32NE
)

// vim: foldmethod=marker
