package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"time"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func randString(n int) string {
	rand.Seed(time.Now().UnixNano())
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

func addOutputToInit(filePath string, node *ast.File) {
	// 遍历文件中的每个声明
	for _, decl := range node.Decls {
		// 如果是函数声明
		if fnDecl, ok := decl.(*ast.FuncDecl); ok {
			// 如果是init函数
			if fnDecl.Name.Name == "init" {
				// 在init函数的开头添加输出语句
				fnDecl.Body.List = append([]ast.Stmt{
					&ast.ExprStmt{
						X: &ast.CallExpr{
							Fun: &ast.SelectorExpr{
								X:   &ast.Ident{Name: "fmt"},
								Sel: &ast.Ident{Name: "Println"},
							},
							Args: []ast.Expr{
								&ast.BasicLit{Kind: token.STRING, Value: fmt.Sprintf(`"%s %s"`, node.Name.Name, filepath.Base(filePath))},
							},
						},
					},
				}, fnDecl.Body.List...)
			}
		}
	}
}

func duplicateInit(node *ast.File) {
	// 遍历文件中的每个声明
	for _, decl := range node.Decls {
		// 如果是函数声明
		if fnDecl, ok := decl.(*ast.FuncDecl); ok {
			// 如果是init函数
			if fnDecl.Name.Name == "init" {
				// 创建一个新的init函数，名字加上随机字符串
				newInitName := "init" + randString(10)
				newInit := *fnDecl
				newInit.Name = &ast.Ident{Name: newInitName}

				// 将新init函数添加到文件中
				node.Decls = append(node.Decls, &newInit)
			}
		}
	}
}

func scanAndModifyFiles(dirPath string) error {
	// 遍历文件夹
	err := filepath.Walk(dirPath, func(filePath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		// 如果是.go文件
		if !info.IsDir() && filepath.Ext(filePath) == ".go" {
			// 打开并解析文件
			fset := token.NewFileSet()
			node, err := parser.ParseFile(fset, filePath, nil, parser.ParseComments)
			if err != nil {
				log.Printf("Error parsing file %s: %v", filePath, err)
				return nil
			}

			// 为每个文件中的init方法添加输出语句
			addOutputToInit(filePath, node)

			// 复制每个文件中的init方法
			duplicateInit(node)

			// 将修改后的AST写回文件
			outputFile, err := os.Create(filePath)
			if err != nil {
				log.Printf("Error creating file %s: %v", filePath, err)
				return nil
			}
			defer outputFile.Close()

			if err := printer.Fprint(outputFile, fset, node); err != nil {
				log.Printf("Error writing to file %s: %v", filePath, err)
				return nil
			}
		}
		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

func main() {
	// 要扫描的文件夹路径
	dirPath := "your_directory_path"

	// 扫描文件夹下所有的.go文件，并为每个文件中的init方法添加输出语句，并生成一个新方法
	err := scanAndModifyFiles(dirPath)
	if err != nil {
		log.Fatalf("Error scanning directory: %v", err)
	}

	log.Println("Done!")
}
