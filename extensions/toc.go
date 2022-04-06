package extensions

import (
	"bytes"
	"log"
	"strconv"

	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
)

var TOCKind = ast.NewNodeKind("TableOfContents")
var _ ast.Node = (*TOCNode)(nil)

type TOCNode struct {
	ast.BaseInline

	Level   int
	Content []byte
}

func (n *TOCNode) Dump(src []byte, level int) {
	ast.DumpHelper(n, src, level, map[string]string{
		"Level":   strconv.Itoa(n.Level),
		"Content": string(n.Content),
	}, nil)
}

// Kind reports the kind of this node.
func (n *TOCNode) Kind() ast.NodeKind {
	return TOCKind
}

type TOCParser struct {
}

var _ parser.InlineParser = (*TOCParser)(nil)

var (
	_open  = []byte("((")
	_hash  = []byte{'+'}
	_close = []byte("))")
)

// Trigger returns characters that trigger this parser.
func (p *TOCParser) Trigger() []byte {
	return []byte{'('}
}

func (p *TOCParser) Parse(_ ast.Node, block text.Reader, _ parser.Context) ast.Node {
	line, seg := block.PeekLine()
	if !bytes.HasPrefix(line, _open) {
		return nil
	}

	stop := bytes.Index(line, _close)
	if stop < 0 {
		return nil // must close on the same line
	}
	log.Printf("Line: %q, Seg: %#v, Start: %c, Stop: %v", string(line), seg, seg.Start, stop)

	seg = text.NewSegment(seg.Start+2, seg.Start+stop)

	n := &TOCNode{Content: block.Value(seg)}

	n.AppendChild(n, ast.NewTextSegment(seg))
	block.Advance(seg.Stop)
	return n
}
