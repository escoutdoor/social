package types

import "io"

type File struct {
	Name    string
	Payload io.Reader
	Size    int64
}
