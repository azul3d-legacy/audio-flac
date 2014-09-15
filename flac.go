// Copyright 2014 The Azul3D Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package flac provides a FLAC audio decoder. It implements the audio.Decoder
// interface of azul3d.org/audio.v1.
//
// NOTE: This package is a work in progress. The implementation is incomplete
// and subject to change. The documentation can be inaccurate.
package flac

import (
	"errors"
	"fmt"
	"io"

	"azul3d.org/audio.v1"
	"github.com/mewkiz/flac"
	"github.com/mewkiz/flac/frame"
)

func init() {
	// Register the FLAC audio decoder.
	audio.RegisterFormat("flac", "fLaC", newDecoder)
}

// decoder is capable of decoding the audio samples of a FLAC stream.
type decoder struct {
	// The FLAC audio stream.
	stream *flac.Stream
	// The previous audio frame is used as a buffer when the last call to Read
	// read fewer audio samples than contained within the audio frame. It is nil
	// otherwise.
	prev *frame.Frame
	// Points to the first unread audio sample in prev.
	i int
}

// newDecoder returns a FLAC audio decoder, which may be used to decode the
// encoded audio samples of the io.Reader or io.ReadSeeker r.
//
// It returns either [audio.Decoder, nil] or [nil, audio.ErrInvalidData] upon
// being called where the returned decoder is used to decode the encoded audio
// data of r.
func newDecoder(r interface{}) (audio.Decoder, error) {
	rr, ok := r.(io.Reader)
	if !ok {
		return nil, fmt.Errorf("flac.newDecoder: unable to decode r; expected io.Reader, got %T", r)
	}

	stream, err := flac.New(rr)
	if err != nil {
		return nil, audio.ErrInvalidData
	}

	return &decoder{
		stream: stream,
	}, nil
}

// Config returns the audio stream configuration of the decoder.
func (dec *decoder) Config() audio.Config {
	return audio.Config{
		SampleRate: int(dec.stream.Info.SampleRate),
		Channels:   int(dec.stream.Info.NChannels),
	}
}

// Read tries to read into the audio slice, b, filling it with at most b.Len()
// audio samples.
//
// Returned is the number of samples that were read into the slice, and an
// error if any occurred.
//
// It is possible for the number of samples read to be non-zero; and for an
// error to be returned at the same time (E.g. read 300 audio samples, but also
// encountered audio.EOS).
func (dec *decoder) Read(b audio.Slice) (n int, err error) {
	// The set closure sets the i:th sample of b to sample.
	var set func(i int, sample int32)
	switch v := b.(type) {
	case audio.PCM8Samples:
		set = func(i int, sample int32) {
			// Unsigned 8-bit PCM audio sample.
			v[i] = audio.PCM8(0x80 + sample)
		}
	case audio.PCM16Samples:
		set = func(i int, sample int32) {
			// Signed 16-bit PCM audio sample.
			v[i] = audio.PCM16(sample)
		}
	case audio.PCM32Samples:
		set = func(i int, sample int32) {
			// Signed 32-bit PCM audio sample.
			v[i] = audio.PCM32(sample)
		}
	default:
		set = func(i int, sample int32) {
			// Generic implementation.
			f := pcmToF64(sample, dec.stream.Info.BitsPerSample)
			b.Set(i, f)
		}
	}

	// Fill b with audio samples from the previous decoded audio frame.
	if dec.prev != nil {
		for ; dec.i < int(dec.prev.BlockSize); dec.i++ {
			for _, subframe := range dec.prev.Subframes {
				sample := subframe.Samples[dec.i]
				if n >= b.Len() {
					return n, nil
				}
				set(n, sample)
				n++
			}
		}
	}
	dec.prev = nil

	// Fill b with audio samples from decoded audio frames.
	for {
		frame, err := dec.stream.ParseNext()
		if err != nil {
			if err == io.EOF {
				return n, audio.EOS
			}
			return n, err
		}
		for i := 0; i < int(frame.BlockSize); i++ {
			for _, subframe := range frame.Subframes {
				sample := subframe.Samples[i]
				if n >= b.Len() {
					if i != int(frame.BlockSize)-1 {
						// Fewer audio samples were read than contained within the audio
						// frame. Store the decoded audio frame and the current sample
						// position for future read operations.
						dec.prev = frame
						dec.i = i
					}
					return n, nil
				} else {
					set(n, sample)
					n++
				}
			}
		}
	}
}

// pcmToF64 converts a signed bps-bit linear PCM audio sample to a 64-bit
// floating-point linear audio sample in the range of -1 to +1.
func pcmToF64(sample int32, bps uint8) audio.F64 {
	switch bps {
	case 8:
		// Unsigned 8-bit PCM audio sample.
		return audio.PCM8ToF64(audio.PCM8(0x80 + sample))
	case 16:
		// Signed 16-bit PCM audio sample.
		return audio.PCM16ToF64(audio.PCM16(sample))
	case 24:
		// Signed 32-bit PCM audio sample.
		return audio.PCM32ToF64(audio.PCM32(sample))
	default:
		panic(fmt.Sprintf("not yet implemented; conversion from %d-bit PCM to F64", bps))
	}
}

// Seek seeks to the specified sample number, relative to the start of the
// stream. As such, subsequent Read() calls on the Reader, begin reading at the
// specified sample.
//
// If any error is returned, it means it was impossible to seek to the specified
// audio sample for some reason, and that the current playhead is unchanged.
func (dec *decoder) Seek(sample uint64) error {
	return errors.New("flac.decoder.Seek: not yet implemented")
}
