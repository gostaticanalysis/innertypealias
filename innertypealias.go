package innertypealias

import (
	"fmt"
	"go/ast"
	"go/token"
	"go/types"

	"github.com/gostaticanalysis/analysisutil"
	"golang.org/x/tools/go/analysis"
)

const doc = "innertypealias find a type which is an alias for exported same package's type"

var Analyzer = &analysis.Analyzer{
	Name: "innertypealias",
	Doc:  doc,
	Run:  run,
}

func run(pass *analysis.Pass) (interface{}, error) {
	embeddeds := make(map[types.Type]bool)
	for _, st := range analysisutil.Structs(pass.Pkg) {
		for i := 0; i < st.NumFields(); i++ {
			if field := st.Field(i); field.Embedded() {
				embeddeds[field.Type()] = true
			}
		}
	}
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
				if typ == nil || typ.Obj().Pkg() != pass.Pkg || !typ.Obj().Exported() || embeddeds[typ] {
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
					Message:        fmt.Sprintf("%s is a alias for %s but it is exported type", x, y),
					SuggestedFixes: []analysis.SuggestedFix{fix},
				})
			}
		}
	}

	return nil, nil
}
