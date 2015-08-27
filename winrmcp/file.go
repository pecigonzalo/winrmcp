package winrmcp

import (
	"encoding/base64"
	"io"
)

type File struct {
	position int
	path     string
	reader   io.Reader
}

func (f *File) Read() (string, int, error) {
	var err error
	var content string
	return content, f.position, err
}

func getChunk(reader io.Reader, filePath string) (string, bool, error) {
	// Upload the file in chunks to get around the Windows command line size limit.
	// Base64 encodes each set of three bytes into four bytes. In addition the output
	// is padded to always be a multiple of four.
	//
	//   ceil(n / 3) * 4 = m1 - m2
	//
	//   where:
	//     n  = bytes
	//     m1 = max (8192 character command limit.)
	//     m2 = len(filePath)

	chunkSize := ((8000 - len(filePath)) / 4) * 3
	chunk := make([]byte, chunkSize)

	n, err := reader.Read(chunk)
	if err != nil && err != io.EOF {
		return "", false, err
	}
	if n == 0 {
		return "", true, nil
	}

	content := base64.StdEncoding.EncodeToString(chunk[:n])

	return content, false, nil

}
