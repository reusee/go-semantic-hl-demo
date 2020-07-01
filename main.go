package main

import (
	"go/ast"
	"go/token"
	"os"
	"time"

	"golang.org/x/tools/go/packages"
)

func main() {
	t0 := time.Now()
	fset := new(token.FileSet)
	pkgs, err := packages.Load(
		&packages.Config{
			Mode: packages.NeedSyntax |
				packages.NeedTypes |
				packages.NeedTypesInfo,
			Fset: fset,
		},
		os.Args[1:]...,
	)
	ce(err)
	if packages.PrintErrors(pkgs) > 0 {
		return
	}
	pt("loaded in %v\n", time.Since(t0))

	pkg := pkgs[0]
	for _, file := range pkg.Syntax {
		ast.Inspect(file, func(node ast.Node) bool {
			if node == nil {
				return true
			}
			begin := fset.Position(node.Pos())
			end := fset.Position(node.End())
			pt("%v\n\tto %v\n", begin, end)
			pt("\tsyntax %T\n", node)
			if expr, ok := node.(ast.Expr); ok {
				typeAndValue := pkg.TypesInfo.Types[expr]
				if typeAndValue.Type != nil {
					pt("\ttype %v\n", typeAndValue.Type)
				}
				if typeAndValue.Value != nil {
					pt("\tvalue %v\n", typeAndValue.Value)
				}
			}
			return true
		})
	}

	// print all identifiers and their references
	for ident, obj := range pkg.TypesInfo.Defs {
		pt("%s at %v\n", ident.Name, fset.Position(ident.NamePos))
		for ref, refObj := range pkg.TypesInfo.Uses {
			if refObj != obj {
				continue
			}
			pt("\t%s at %v\n", ref.Name, fset.Position(ref.NamePos))
		}
	}

}
