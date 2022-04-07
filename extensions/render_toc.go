package extensions

import (
	"fmt"
	"sync"

	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/util"
)

//var _ renderer.Renderer = (*TOCRenderer)(nil)

type TOCRenderer struct {
	once sync.Once // guards init
}

func (r *TOCRenderer) init() {
	r.once.Do(func() {})
}

// RegisterFuncs registers the rendering functions with the provided
// goldmark registerer. This teaches goldmark to call us when it encounters a
// TOC in the AST.
func (r *TOCRenderer) RegisterFuncs(reg renderer.NodeRendererFuncRegisterer) {
	reg.Register(TOCKind, r.Render)
}

// Render renders the provided Node.
func (r *TOCRenderer) Render(w util.BufWriter, _ []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	r.init()

	_, ok := node.(*TOCPlaceholderNode)
	if !ok {
		return ast.WalkStop, fmt.Errorf("unexpected node %T, expected *extensions.TOCPlaceholderNode", node)
	}

	w.WriteString(``)

	return ast.WalkSkipChildren, nil
}
