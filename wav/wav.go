// Package wav provides a WAV audio encoder.
//
// A brief introduction of the WAV audio format [1][2] follows. A WAV file
// consists of a sequence of chunks as specified by the RIFF format. Each chunk
// has a header and a body. The header specifies the type of the chunk and the
// size of its body.
//
// The first chunk of a WAV file is the standard RIFF type chunk, with a "WAVE"
// type ID. It is followed by a mandatory format chunk, which describes the
// basic properties of the audio stream; such as its sample rate and the number
// of channels used. Subsequent chunks may appear in any order and several
// chunks are optional. The only other chunk that is mandatory is the data
// chunk, which contains the encoded audio samples.
//
// Below follows an overview of a basic WAV file.
//
//    Header: {id: "RIFF", size: 0004}
//    Body:   "WAVE"
//    Header: {id: "fmt ", size: NNNN}
//    Body:   format of the audio samples
//    Header: {id: "data", size: NNNN}
//    Body:   audio samples
//
// Please refer to the WAV specification for more in-depth details about its
// file format.
//
//    [1]: http://www.sonicspot.com/guide/wavefiles.html
//    [2]: https://ccrma.stanford.edu/courses/422/projects/WaveFormat/
package wav

import (
	"io"

	"azul3d.org/audio.v1"
)

// encoder is capable of encoding audio samples to a WAV file.
type encoder struct {
	// Underlying io.WriteSeeker to which the WAV file is written to.
	w io.WriteSeeker
	// Audio configuration; including sample rate and number of channels.
	conf audio.Config
}

// bps represents the number of bits-per-sample used to encode audio samples.
const bps = 16

// NewEncoder creates a new WAV encoder, which stores the audio configuration in
// a WAV header and encodes any audio samples written to it. The contents of the
// WAV header and the encoded audio samples are written to w.
//
// Note: The Close method of the encoder must be called when finished using it.
func NewEncoder(w io.WriteSeeker, conf audio.Config) (enc audio.Writer, err error) {
	// Write WAV file header to w, based on the audio configuration.
	err = writeHeader(w, conf)
	if err != nil {
		return nil, err
	}

	// Return encoder which encodes the audio samples written to it and stores
	// writes those to w.
	return &encoder{w: w, conf: conf}, nil
}

// Write attempts to write all, b.Len(), samples in the slice to the
// writer.
//
// Returned is the number of samples from the slice that where wrote to
// the writer, and an error if any occurred.
//
// If the number of samples wrote is less than buf.Len() then the returned
// error must be non-nil. If any error occurs it should be considered fatal
// with regards to the writer: no more data can be subsequently wrote after
// an error.
func (enc *encoder) Write(b audio.Slice) (n int, err error) {
	// TODO(u): Implement fast-paths for PCM8, PCM16 and PCM24 (PCM32).

	var buf [2]byte
	for ; n < b.Len(); n++ {
		f := b.At(n)
		sample := audio.F64ToPCM16(f)
		buf[0] = uint8(sample)
		buf[1] = uint8(sample >> 8)
		m, err := enc.w.Write(buf[:])
		if err != nil {
			return n, err
		}
		if m < len(buf) {
			return n, io.ErrShortWrite
		}
	}

	return n, nil
}

// Close signals to the encoder that encoding has been completed, thereby
// allowing it to update the placeholder values in the WAV file header.
func Close() error {
	panic("not yet implemented.")
}
