package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"

	"azul3d.org/audio.v1"
	_ "github.com/mewkiz/audio/flac"
	"github.com/mewkiz/audio/wav"
	"github.com/mewkiz/pkg/osutil"
	"github.com/mewkiz/pkg/pathutil"
)

func main() {
	flag.Parse()
	for _, path := range flag.Args() {
		err := flac2wav(path)
		if err != nil {
			log.Fatalln(err)
		}
	}
}

// flac2wav converts the provided FLAC file to a WAV file.
func flac2wav(path string) error {
	// Open FLAC file.
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()
	br := bufio.NewReader(f)

	// Create FLAC decoder.
	dec, magic, err := audio.NewDecoder(br)
	if err != nil {
		return err
	}
	fmt.Println("magic:", magic)
	conf := dec.Config()
	fmt.Println(conf)

	// Create WAV file.
	wavPath := pathutil.TrimExt(path) + ".wav"
	exists, err := osutil.Exists(wavPath)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("the file %q exists already.", wavPath)
	}
	w, err := os.Create(wavPath)

	// Create WAV encoder.
	enc, err := wav.NewEncoder(w, conf)
	if err != nil {
		return err
	}
	defer enc.Close()

	_, err = audio.Copy(enc, dec)
	if err != nil {
		return err
	}

	return nil
}
