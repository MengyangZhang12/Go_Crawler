package engine

type Request struct {
	Url        string
	Parser     Parser
	ParserFunc func([]byte) ParseResult
}

type ParseResult struct {
	Requests []Request
	Items    []Item
}

type Parser interface {
	Parse(contents []byte) ParseResult
	Serialize() (name string, args interface{})
}

type Item struct {
	Url     string
	Id      string
	Payload interface{}
}

type NilParser struct{}

func (n NilParser) Parse(_ []byte) ParseResult {
	return ParseResult{}
}

func (n NilParser) Serialize() (name string, args interface{}) {
	return "NilParser", nil
}

type FuncParser struct {
	parser func([]byte) ParseResult
	name   string
}

func (f *FuncParser) Parse(contents []byte) ParseResult {
	return f.parser(contents)
}

func (f *FuncParser) Serialize() (name string, args interface{}) {
	return f.name, nil
}

func NewFuncParser(p func([]byte) ParseResult, name string) *FuncParser {
	return &FuncParser{
		parser: p,
		name:   name,
	}
}
