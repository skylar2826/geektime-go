package main

import (
	_ "embed"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

//go:embed tpl.gohtml
var genOrm string

func gen(w io.Writer, srcFile string) error {
	// ast语法树解析
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, srcFile, nil, parser.ParseComments)
	if err != nil {
		return err
	}

	s := &SingleFileEntryVisitor{}
	ast.Walk(s, f)
	file := s.get()

	// 操作模板
	tpl := template.New("gen-orm")
	tpl, err = tpl.Parse(genOrm)
	if err != nil {
		return err
	}
	return tpl.Execute(w, data{
		File: file,
		Ops:  []string{"Lt", "Eq"},
	})
}

type data struct {
	*File
	Ops []string
}

//命令行跑这个main函数
//
//1. package 改为main
//2. 当前目录下执行go install .
//3. cd testdata
//4. gen(包名) user.go(参数 os.Args[1])

func main() {
	src := os.Args[1]
	dstDir := filepath.Dir(src)
	fileName := filepath.Base(src)
	idx := strings.LastIndexByte(fileName, '.')
	dst := filepath.Join(dstDir, fileName[:idx]+"_gen.go")
	f, err := os.Create(dst)
	if err != nil {
		fmt.Println(err)
		return
	}
	err = gen(f, src)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("生成成功")
}
