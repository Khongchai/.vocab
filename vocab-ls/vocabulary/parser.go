package vocabulary

import "context"

type Parser struct {
	Ctx context.Context
	Uri string
}

func (p *Parser) Parse(text string) *VocabAst {
	return &VocabAst{
		Documents: []*Document{
			{},
		},
	}
}
