package innertypealias

import (
	"fmt"
	"go/ast"
	"go/token"
	"go/types"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

const doc = "innertypealias find a type which is an alias for exported same package's type"

var Analyzer = &analysis.Analyzer{
	Name: "innertypealias",
	Doc:  doc,
	Run:  run,
	Requires: []*analysis.Analyzer{
		inspect.Analyzer,
	},
}

func run(pass *analysis.Pass) (interface{}, error) {
	inspect := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)
	nodes := []ast.Node{(*ast.StructType)(nil)}
	embeddeds := make(map[string]bool)
	inspect.Preorder(nodes, func(n ast.Node) {
		st, _ := n.(*ast.StructType)
		if st == nil {
			return
		}

		for _, f := range st.Fields.List {
			id, _ := f.Type.(*ast.Ident)
			if id == nil || f.Names != nil {
				continue
			}

			obj := pass.TypesInfo.ObjectOf(id)
			if obj.Pkg() == pass.Pkg {
				embeddeds[id.Name] = true
			}
		}
	})

	for _, f := range pass.Files {
		for _, decl := range f.Decls {
			decl, _ := decl.(*ast.GenDecl)
			if decl == nil || decl.Tok != token.TYPE {
				continue
			}

			for _, spec := range decl.Specs {
				spec, _ := spec.(*ast.TypeSpec)
				if spec == nil || spec.Assign == token.NoPos || !spec.Name.IsExported() {
					continue
				}

				typ, _ := pass.TypesInfo.TypeOf(spec.Type).(*types.Named)
				if typ == nil || typ.Obj().Pkg() != pass.Pkg || !typ.Obj().Exported() || embeddeds[spec.Name.Name] {
					continue
				}

				// type X = Y => type X Y
				x, y := spec.Name.Name, typ.Obj().Name()
				fix := analysis.SuggestedFix{
					Message: "fix type alias to defined type",
					TextEdits: []analysis.TextEdit{{
						Pos:     spec.Pos(),
						End:     spec.End(),
						NewText: []byte(x + " " + y),
					}},
				}
				pass.Report(analysis.Diagnostic{
					Pos:            spec.Pos(),
					End:            spec.End(),
					Message:        fmt.Sprintf("%s is an alias for %s but it is exported type", x, y),
					SuggestedFixes: []analysis.SuggestedFix{fix},
				})
			}
		}
	}

	return nil, nil
}
