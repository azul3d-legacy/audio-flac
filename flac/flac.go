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
)

func init() {
	// Register the FLAC audio decoder.
	audio.RegisterFormat("flac", "fLaC", newDecoder)
}

// newDecoder returns a FLAC audio decoder, which may be used to decode the
// encoded audio sample of the io.Reader or io.ReadSeeker r.
//
// It returns either [audio.Decoder, nil] or [nil, audio.ErrInvalidData] upon
// being called where the returned decoder is used to decode the encoded audio
// data of r.
func newDecoder(r interface{}) (audio.Decoder, error) {
	rr, ok := r.(io.Reader)
	if !ok {
		return nil, fmt.Errorf("flac.newDecoder: unable to decode r; expected io.Reader, got %T", r)
	}
	dec := &decoder{buf: make([]int32, 0)}
	var err error
	dec.stream, err = flac.New(rr)
	if err != nil {
		return nil, audio.ErrInvalidData
	}
	return dec, nil
}

// decoder is capable of decoding the audio samples of a FLAC stream.
type decoder struct {
	// The FLAC audio stream.
	stream *flac.Stream
	// Buffer superfluous decoded samples from previous Read operations, for
	// future calls to Read.
	buf []int32
}

// Config returns the audio stream configuration of this decoder. It may block
// until at least the configuration part of the stream has been read.
func (dec *decoder) Config() audio.Config {
	return audio.Config{
		SampleRate: int(dec.stream.Info.SampleRate),
		Channels:   int(dec.stream.Info.NChannels),
	}
}

// Read tries to read into the audio slice, b, filling it with at max b.Len()
// audio samples.
//
// Returned is the number of samples that where read into the slice, and an
// error if any occurred.
//
// It is possible for the number of samples read to be non-zero; and for an
// error to be returned at the same time (E.g. read 300 audio samples, but also
// encountered audio.EOS).
func (dec *decoder) Read(b audio.Slice) (n int, err error) {
	// TODO(u): Implement fast paths for common audio sample formats:
	//    * PCM8
	//    * PCM16
	//    * PCM24 (PCM32)

	// Generic implementation.

	// Store buffered audio samples from the previous frame.
	bps := dec.stream.Info.BitsPerSample
	for i := 0; i < len(dec.buf); i++ {
		if n >= b.Len() {
			break
		}
		sample := pcmToF64(dec.buf[i], bps)
		b.Set(n, sample)
		n++
	}
	dec.buf = dec.buf[n:]
	if n >= b.Len() {
		return n, nil
	}

	// Decode audio samples of the next frame.
	frame, err := dec.stream.ParseNext()
	if err == io.EOF {
		err = audio.EOS
	}
	if len(frame.Subframes) < 1 {
		return n, err
	}

	// Store decoded audio samples in b.
	nsamples := len(frame.Subframes[0].Samples)
	for i := 0; i < nsamples; i++ {
		for _, subframe := range frame.Subframes {
			if n >= b.Len() {
				break
			}
			sample := pcmToF64(subframe.Samples[i], bps)
			b.Set(n, sample)
			n++
		}
	}

	// Buffer superfluous decoded audio samples in dec.buf.
	extra := nsamples - n
	if extra > 0 {
		if cap(dec.buf) >= extra {
			dec.buf = dec.buf[:extra]
		} else {
			dec.buf = make([]int32, extra)
		}
		for i := range dec.buf {
			for _, subframe := range frame.Subframes {
				dec.buf[i] = subframe.Samples[n+i]
			}
		}
	}

	return n, err
}

// pcmToF64 converts a signed bps-bit linear PCM audio sample to a 64-bit
// floating-point linear audio sample in the range of -1 to +1.
func pcmToF64(sample int32, bps uint8) audio.F64 {
	switch bps {
	case 16:
		return audio.PCM16ToF64(audio.PCM16(sample))
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
