audio
=====

This project is meant to facilitate audio decoding by tracking implementations
of a common decoding interface for the various audio file formats. The
[audio.Decoder][] interface from the [audio][] package of [Azul3D][] has very
much in common with the [image.Image][] interface from the standard library.

Stumbling upon the audio package was a pleasant surprise. It feels well
engineered, and is able to provide a common audio API without compromising
performance or being too restrictive. Don't take my word for it however! Have a
[look for yourself][audio] and start a discussion in the [issue tracker][] if
you have any constructive criticism. Lets design an audio API that is as mature
and usable as the image.Image interface of the standard library.

[audio.Decoder]: https://godoc.org/gopkg.in/azul3d/audio.v1#Decoder
[audio]: https://godoc.org/gopkg.in/azul3d/audio.v1
[Azul3D]: https://azul3d.github.io/
[image.Image]: https://golang.org/pkg/image/#Image
[issue tracker]: https://github.com/azul3d/audio/issues

Documentation
-------------

Documentation provided by GoDoc.

* [flac][]: provides a FLAC audio decoder.

		go get github.com/mewkiz/audio/flac

* [wav][]: decodes wav audio files.

		go get azul3d.org/audio/wav.v1

[flac]: https://godoc.org/github.com/mewkiz/audio/flac
[wav]: https://godoc.org/gopkg.in/azul3d/audio-wav.v1

public domain
-------------

This code is hereby released into the *[public domain][]*.

[public domain]: https://creativecommons.org/publicdomain/zero/1.0/
