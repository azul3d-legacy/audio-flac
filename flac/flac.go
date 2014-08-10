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
	dec := &decoder{buf: make(audio.F64Samples, 0)}
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
	buf audio.F64Samples
	// Points to the first buffered sample in buf.
	first int
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

	// Drain buffered samples from previous read operations.
	for _, f := range dec.buf[dec.first:] {
		if n >= b.Len() {
			dec.first = n
			return n, nil
		}
		b.Set(n, f)
		n++
	}
	dec.first = 0
	dec.buf = dec.buf[:0]

	// Decode the audio samples of a frame.
	frame, err := dec.stream.ParseNext()
	if err != nil {
		if err == io.EOF {
			return n, audio.EOS
		}
		return n, err
	}
	bps := dec.stream.Info.BitsPerSample
	for i := 0; i < int(frame.BlockSize); i++ {
		for _, subframe := range frame.Subframes {
			sample := subframe.Samples[i]
			f := pcmToF64(sample, bps)
			if n < b.Len() {
				b.Set(n, f)
				n++
			} else {
				// Buffer superfluous audio samples.
				dec.buf = append(dec.buf, f)
			}
		}
	}

	return n, nil
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
