package main

import (
	"flag"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"text/template"
)

var (
	excludeDir  string
	outputFile  string
	packageName string
	scanDir     string
	mod         string
	cwd         string
)

func init() {
	flag.StringVar(&excludeDir, "exclude", "tool", "Directory to exclude")
	flag.StringVar(&outputFile, "output", "gen/functions_gen.go", "Output file name")
	flag.StringVar(&packageName, "package", "gen", "Package name for generated file")
	flag.StringVar(&scanDir, "dir", ".", "Directory to scan")
	flag.StringVar(&mod, "mod", "myProject/hsq/", "Directory to scan")
	flag.StringVar(&cwd, "cwd", "/home/gohook/", "Directory to scan")
	flag.Parse()
}

type Function struct {
	Package string
	Name    string
}

type Method struct {
	Name string
}

type Type struct {
	Package string
	Name    string
}

type TemplateData struct {
	PackageName string
	Imports     []string
	Functions   []Function
	Methods     map[Type][]Method
}

const templateText = `
package {{.PackageName}}

import (
    {{range .Imports}}
    "{{.}}"
    {{end}}
)

var FuncMap = make(map[string]interface{})

func init() {
    {{range .Functions}}{{if eq .Package "main"}}{{else if eq .Package $.PackageName}}
	    FuncMap["{{.Name}}"] = {{.Name}}{{else}}
    FuncMap["{{.Package}}.{{.Name}}"] = {{.Package}}.{{.Name}}{{end}}{{end}}
    {{range $type, $methods := .Methods}}{{range $methods}}{{if eq $type.Package $.PackageName}}
    FuncMap["{{$type.Name}}.{{.Name}}"] = (*{{$type.Name}}).{{.Name}}{{else}}
    FuncMap["{{$type.Package}}.{{$type.Name}}.{{.Name}}"] = (*{{$type.Package}}.{{$type.Name}}).{{.Name}}{{end}}{{end}}{{end}}
}

`

func main() {
	if packageName == "" {
		fmt.Fprintf(os.Stderr, "Package name is required\n")
		os.Exit(1)
	}

	data := TemplateData{
		PackageName: packageName,
		Imports:     []string{},
		Functions:   []Function{},
		Methods:     make(map[Type][]Method),
	}

	imports := make(map[string]bool)

	err := filepath.Walk(scanDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// 排除指定目录
		if excludeDir != "" && strings.HasPrefix(path, excludeDir) {
			if info.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}

		// 排除 _gen.go 结尾的文件
		if strings.HasSuffix(path, "_gen.go") {
			return nil
		}

		if !info.IsDir() && strings.HasSuffix(path, ".go") {
			// 解析文件
			fset := token.NewFileSet()
			f, err := parser.ParseFile(fset, path, nil, 0)
			if err != nil {
				return nil
			}

			// 获取包名
			pkgName := f.Name.Name
			if pkgName == "main" || pkgName == "internal" {
				return nil
			}

			// 遍历 AST
			ast.Inspect(f, func(n ast.Node) bool {
				switch x := n.(type) {
				case *ast.FuncDecl:
					// 获取文件路径
					filePath := fset.Position(x.Pos()).Filename

					// 获取绝对路径
					absPath, _ := filepath.Abs(filePath)
					absPath = strings.TrimPrefix(absPath, cwd)
					if strings.LastIndexByte(absPath, '/') >= 0 {
						absPath = absPath[:strings.LastIndexByte(absPath, '/')]
					}
					pkgPath := mod + absPath
					fmt.Println(pkgPath)
					if x.Recv == nil {
						// 这是一个函数
						data.Functions = append(data.Functions, Function{
							Package: pkgName,
							Name:    x.Name.Name,
						})
						if strings.HasSuffix(pkgPath, pkgName) && pkgName != packageName {
							imports[pkgPath] = true
						}
					} else {
						// 这是一个方法
						if len(x.Recv.List) > 0 {
							var recvType ast.Expr
							switch t := x.Recv.List[0].Type.(type) {
							case *ast.StarExpr:
								// 指针接收器
								recvType = t.X
							case *ast.Ident:
								// 值接收器
								recvType = t
							default:
								// 其他类型，跳过
								return true
							}

							if ident, ok := recvType.(*ast.Ident); ok {
								t := Type{
									Package: pkgName,
									Name:    ident.Name,
								}
								data.Methods[t] = append(data.Methods[t], Method{Name: x.Name.Name})
								if strings.HasSuffix(pkgPath, pkgName) && pkgName != packageName {
									imports[pkgPath] = true
								}
							}
						}
					}
				}
				return true
			})
		}

		return nil
	})

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error walking the path: %v\n", err)
		os.Exit(1)
	}

	// 将导入的包转换为排序后的切片
	for imp := range imports {
		data.Imports = append(data.Imports, imp)
	}
	sort.Strings(data.Imports)
	fmt.Println(imports)

	// 解析模板
	tmpl, err := template.New("funcmap").Parse(templateText)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing template: %v\n", err)
		os.Exit(1)
	}

	// 创建输出文件
	out, err := os.Create(outputFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating output file: %v\n", err)
		os.Exit(1)
	}
	defer out.Close()

	// 执行模板
	err = tmpl.Execute(out, data)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error executing template: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Generated function map in %s\n", outputFile)
}
