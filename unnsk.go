package main

import (
	"bytes"
	"encoding/binary"
	"flag"
	"fmt"
	"io"
	"os"
	"reflect"

	"github.com/JoshVarga/blast"
)

const (
	name = "NaShrinK"
	ext  = "NSK"
)

var (
	fileID  = []byte("NSK")
	fset    = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	arcFile = fset.String("e", "", fmt.Sprintf("The %s file to extract.", ext))
	dstPath = fset.String("d", "", "Optional output directory to extract to.")
)

func errpanic(e error) {
	if e != nil {
		panic(e)
	}
}

type header struct {
	id     []byte // Stores the file sig `NSK`
	cSize  int    // Compressed file size
	uSize  int    // Uncompressed file size
	fnSize int    // Filename length
	fSize  int    // Total file size
	fn     string // Filename for this header.
}

func main() {
	fset.Parse(os.Args[1:])

	if *arcFile == "" {
		fset.Usage()
		os.Exit(0)
	}

	if *dstPath != "" {
		if _, err := os.Stat(*dstPath); os.IsNotExist(err) {
			fmt.Fprintf(os.Stderr, "destination folder (%s) doesn't exist: %v", *dstPath, err)
			os.Exit(1)
		}
	}

	f, err := os.Open(*arcFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error opening %s: %v", *arcFile, err)
		os.Exit(1)
	}
	defer f.Close()

	s, err := f.Stat()
	errpanic(err)

	h := &header{
		id:    make([]byte, 3),
		fSize: int(s.Size()),
	}

	// Read archive members in a loop until done.
	for {
		_, err = f.Read(h.id)
		errpanic(err)

		if !reflect.DeepEqual(h.id, fileID) {
			fmt.Fprintf(os.Stderr, "file is not a %s file", ext)
			os.Exit(1)
		}

		// Read metadate from file.
		metadata := make([]byte, 14)
		_, err = f.Read(metadata)
		errpanic(err)
		h.cSize = int(binary.LittleEndian.Uint32(metadata[:4]))
		h.uSize = int(binary.LittleEndian.Uint32(metadata[9:13]))
		h.fnSize = int(metadata[13])

		if h.fnSize > 12 {
			fmt.Fprintf(os.Stderr, "filename length is > 12: %d", h.fnSize)
			os.Exit(1)
		}

		fn := make([]byte, h.fnSize)
		_, err = f.Read(fn)
		errpanic(err)
		h.fn = string(fn)

		fmt.Printf("Extracting: %s compressed / uncompressed: %d / %d bytes\n", h.fn, h.cSize, h.uSize)

		// Read all the DCL compressed data.
		dcl := make([]byte, h.cSize)
		_, err = f.Read(dcl)
		errpanic(err)

		b := bytes.NewReader(dcl)
		r, err := blast.NewReader(b)
		errpanic(err)

		if *dstPath != "" {
			h.fn = *dstPath + "/" + h.fn
		}

		o, err := os.Create(h.fn)
		errpanic(err)

		_, err = io.Copy(o, r)
		errpanic(err)
		o.Close()

		cur, err := f.Seek(0, io.SeekCurrent)
		errpanic(err)
		if cur == int64(h.fSize) {
			break
		}
	}
}
