package main

type Position struct {
	Index, Line, Col int
	FileName, FileText string
}

func NewPosition(idx, ln, col int, fn, ftxt string) *Position {
	p := &Position{
		Index: idx,
		Line: ln,
		Col: col,
		FileName: fn,
		FileText: ftxt,
	}

	return p
}

func (p *Position) Advance(c string) {
	p.Index += 1
	p.Col += 1

	if c == "\n" {
		p.Line += 1
		p.Col = 0
	}
}
