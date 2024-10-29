package main

import (
	"go/ast"
)

type SingleFileEntryVisitor struct {
	file *FileVisitor
}

func (s *SingleFileEntryVisitor) get() *File {
	types := make([]Type, 0, len(s.file.Types))
	for _, typ := range s.file.Types {
		types = append(types, Type{
			Name:   typ.Name,
			Fields: typ.Fields,
		})
	}
	return &File{
		Package: s.file.Package,
		Imports: s.file.Imports,
		Types:   types,
	}
}

func (s *SingleFileEntryVisitor) Visit(node ast.Node) (w ast.Visitor) {
	fn, ok := node.(*ast.File)
	if !ok {
		// 不是我们需要的file
		return s
	}
	s.file = &FileVisitor{
		Package: fn.Name.String(),
	}
	return s.file
}

type Type struct {
	Name   string
	Fields []Field
}

type File struct {
	Package string
	Imports []string
	Types   []Type
}

type TypeVisitor struct {
	Name   string
	Fields []Field
}

func (t *TypeVisitor) Visit(node ast.Node) (w ast.Visitor) {
	n, ok := node.(*ast.Field)
	if !ok {
		return t
	}

	var typ string
	switch nt := n.Type.(type) {
	case *ast.Ident:
		typ = nt.String()
	case *ast.StarExpr:
		switch nx := nt.X.(type) {
		case *ast.SelectorExpr:
			typ = "*" + nx.X.(*ast.Ident).String() + "." + nx.Sel.String()
		}

	case *ast.ArrayType:
		typ = "[]" + nt.Elt.(*ast.Ident).Name
	default:
		typ = "string"
	}

	//type user struct {
	//	a,b,c string // names
	//}
	for _, name := range n.Names {
		t.Fields = append(t.Fields, Field{
			Name: name.String(),
			Type: typ,
		})
	}
	return t
}

type Field struct {
	Name string
	Type string
}

type FileVisitor struct {
	Package string
	Imports []string
	Types   []*TypeVisitor
}

func (f *FileVisitor) Visit(node ast.Node) (w ast.Visitor) {
	switch n := node.(type) {
	case *ast.ImportSpec:
		path := n.Path.Value
		if n.Name != nil && n.Name.String() != "" {
			path = n.Name.String() + " " + path
		}
		f.Imports = append(f.Imports, path)
	case *ast.TypeSpec:
		v := &TypeVisitor{
			Name: n.Name.String(),
		}
		f.Types = append(f.Types, v)
		return v
	}

	return f
}
