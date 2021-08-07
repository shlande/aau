package parser

// Parser is a tools used to parse useful info from common
type Parser interface {
	Parse(name string) (*Result, error)
}
