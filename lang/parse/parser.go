package parse

import "io"

type Parser interface {
	Parse(io.Reader) (int64, error)
	ParseBuffer([]byte) (int64, error)
	ParseString(string) (int64, error)
}
