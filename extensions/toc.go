package extensions

import (
	"bytes"
	"strconv"

	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
)

var TOCKind = ast.NewNodeKind("TableOfContents")
var _ ast.Node = (*TOCPlaceholderNode)(nil)

type TOCPlaceholderNode struct {
	ast.BaseInline

	Level   int
	Content []byte
}

func (n *TOCPlaceholderNode) Dump(src []byte, level int) {
	ast.DumpHelper(n, src, level, map[string]string{
		"Level":   strconv.Itoa(n.Level),
		"Content": string(n.Content),
	}, nil)
}

// Kind reports the kind of this node.
func (n *TOCPlaceholderNode) Kind() ast.NodeKind {
	return TOCKind
}

type TOCEntry struct {
	Level     int
	Name      string
	Reference string
}

type TOCParser struct {
	Entries []TOCEntry
}

var _ parser.InlineParser = (*TOCParser)(nil)

var (
	_open   = []byte("((")
	_indent = []byte{'+'}
	_close  = []byte("))")
	trimset = " \t"
)

// Trigger returns characters that trigger this parser.
func (p *TOCParser) Trigger() []byte {
	return []byte{' ', '(', '\t'}
}

func (p *TOCParser) Parse(_ ast.Node, block text.Reader, _ parser.Context) ast.Node {
	line, seg := block.PeekLine()
	line = bytes.TrimLeft(line, trimset)
	if !bytes.HasPrefix(line, _open) {
		return nil
	}

	stop := bytes.Index(line, _close)
	if stop < 0 {
		return nil // must close on the same line
	}

	seg = text.NewSegment(seg.Start+len(_open), seg.Start+stop)

	n := &TOCPlaceholderNode{Content: block.Value(seg)}

	p.Entries = append(p.Entries, TOCEntry{Level: 1, Name: string(line), Reference: ""})

	n.AppendChild(n, ast.NewTextSegment(seg))
	block.Advance(seg.Stop)
	return n
}
