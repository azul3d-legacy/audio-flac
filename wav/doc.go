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
//
// NOTE: This package is a work in progress. The implementation is incomplete
// and subject to change. The documentation can be inaccurate.
package wav
