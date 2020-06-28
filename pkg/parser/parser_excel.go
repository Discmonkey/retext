package parser

type XlsxParser struct {
}

func (x XlsxParser) Convert([]byte) (Document, error) {
	panic("implement me")
}

var _ Parser = XlsxParser{}
