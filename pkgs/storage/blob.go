package storage

import "io"

type Blob struct {
	Content []string
	R       io.Reader
}

func (b *Blob) Read(p []byte) (n int, err error) {
	return b.R.Read(p)
}
