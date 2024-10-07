package c

import (
	"fmt"

	Text "github.com/linkdotnet/golang-stringbuilder"
	"github.com/magnetenstad/dragon-compiler/pkg/ast"
)

type Context struct {
	tabs            int
	block           int
	uniqueIndex     int
	instancesToFree map[int][]string
	sb              *Text.StringBuilder
}

func Generate(root *ast.RootNode) string {
	sb := Text.StringBuilder{}
	ctx := Context{
		sb:              &sb,
		instancesToFree: make(map[int][]string),
	}
	ctx.sb.Append("#include <stdio.h>\n\n")
	for _, declaration := range root.Declarations {
		generate(declaration, &ctx)
	}
	ctx.sb.Append("\nint main(int argc, char *argv[]) {;\n")
	ctx.tabs += 1
	for _, child := range root.Children {
		generate(child, &ctx)
	}
	writeTabs(ctx.sb, ctx.tabs)
	ctx.sb.Append("return 0;\n")
	ctx.tabs -= 1
	ctx.sb.Append("}\n")
	return ctx.sb.ToString()
}

func generate(node *ast.Node, ctx *Context) {

	switch node.Type {

	case ast.TypeBlock:
		ctx.block += 1
		block := ctx.block
		writeTabs(ctx.sb, ctx.tabs)
		ctx.sb.Append(fmt.Sprintf("__StartBlock_%d__: {;\n", block))
		ctx.tabs += 1
		for _, child := range node.Children {
			generate(child, ctx)
		}
		instancesToFree, exist := ctx.instancesToFree[block]
		if exist {
			for _, id := range instancesToFree {
				writeTabs(ctx.sb, ctx.tabs)
				ctx.sb.Append(fmt.Sprintf("free(%s);\n", id))
			}
		}
		ctx.tabs -= 1
		writeTabs(ctx.sb, ctx.tabs)
		ctx.sb.Append("}\n")
		writeTabs(ctx.sb, ctx.tabs)
		ctx.sb.Append(fmt.Sprintf("__EndBlock_%d__: {}\n", block))

	case ast.TypeBlocks:
		fallthrough
	case ast.TypeStatements:
		fallthrough
	case ast.TypeStatement:
		for _, child := range node.Children {
			generate(child, ctx)
		}

	case ast.TypePrintStatement:
		writeTabs(ctx.sb, ctx.tabs)
		ctx.sb.Append("printf(")
		generate(node.Children[0], ctx)
		ctx.sb.Append(");\n")

	case ast.TypeAssignmentStatement:
		writeTabs(ctx.sb, ctx.tabs) // TODO: handle identifiers and types
		typeHint := "int"
		if node.Children[1].Children[0].Type == ast.TypeConstructor {
			typeHint = node.Children[1].Children[0].Lexeme
		}
		ctx.sb.Append(fmt.Sprintf(
			"%s %s = ", typeHint, node.Children[0].Lexeme))

		generate(node.Children[1], ctx)
		ctx.sb.Append(";\n")

	case ast.TypeSkipStatement:
		writeTabs(ctx.sb, ctx.tabs)
		ctx.sb.Append(fmt.Sprintf("goto __EndBlock_%d__;\n", ctx.block))

	case ast.TypeSkipIfStatement:
		writeTabs(ctx.sb, ctx.tabs)
		ctx.sb.Append("if (")
		generate(node.Children[0], ctx)
		ctx.sb.Append(fmt.Sprintf(") goto __EndBlock_%d__;\n", ctx.block))

	case ast.TypeExpression:
		generate(node.Children[0], ctx)

	case ast.TypeOperator:
		ctx.sb.Append("(")
		generate(node.Children[0], ctx)
		ctx.sb.Append(node.Lexeme)
		generate(node.Children[1], ctx)
		ctx.sb.Append(")")

	case ast.TypeLiteral:
		ctx.sb.Append(fmt.Sprintf("\"%s\"", node.Lexeme))

	case ast.TypeNumber:
		ctx.sb.Append(fmt.Sprintf("%d", node.Number))

	case ast.TypeBoolean:
		ctx.sb.Append(fmt.Sprintf("%d", node.Number))

	case ast.TypeIdentifier:
		ctx.sb.Append(node.Lexeme) // TODO: handle identifiers

	case ast.TypeNot:
		ctx.sb.Append("!(")
		generate(node.Children[0], ctx)
		ctx.sb.Append(")")

	case ast.TypeStructDeclaration:
		writeTabs(ctx.sb, ctx.tabs)
		ctx.sb.Append("typedef struct {\n")
		ctx.tabs += 1
		for _, child := range node.Children {
			generate(child, ctx)
		}
		ctx.tabs -= 1
		writeTabs(ctx.sb, ctx.tabs)
		ctx.sb.Append(fmt.Sprintf("} %s;\n", node.Lexeme))
		writeTabs(ctx.sb, ctx.tabs)
		ctx.sb.Append(fmt.Sprintf(
			"void __Construct_%s__(%s *o) {\n", node.Lexeme, node.Lexeme))
		ctx.tabs += 1
		for _, child := range node.Children {
			writeTabs(ctx.sb, ctx.tabs)
			if len(child.Children) > 0 {
				ctx.sb.Append(fmt.Sprintf("o->%s = ", child.Lexeme))
				generate(child.Children[0], ctx)
			} else {
				defaultValue := getDefaultValue(child)
				if len(defaultValue) == 0 {
					ctx.sb.Append(fmt.Sprintf(
						"__Construct_%s__(&o->%s)", child.TypeHint, child.Lexeme))
				} else {
					ctx.sb.Append(fmt.Sprintf("o->%s = ", child.Lexeme))
					ctx.sb.Append(getDefaultValue(child))
				}
			}
			ctx.sb.Append(";\n")
		}
		ctx.tabs -= 1
		writeTabs(ctx.sb, ctx.tabs)
		ctx.sb.Append("}\n")

	case ast.TypeStructField:
		writeTabs(ctx.sb, ctx.tabs)
		ctx.sb.Append(
			fmt.Sprintf("%s %s;\n",
				typeHintToString(node.TypeHint), node.Lexeme))

	case ast.TypeConstructor:
		sb := Text.StringBuilder{}
		ctx.uniqueIndex += 1
		uniqueIndex := ctx.uniqueIndex
		sb.AppendRune('\n')
		instanceId := fmt.Sprintf("__Instance_%d__", uniqueIndex)
		writeTabs(&sb, ctx.tabs)
		sb.Append(fmt.Sprintf(
			"%s %s;\n", node.Lexeme, instanceId))
		writeTabs(&sb, ctx.tabs)
		sb.Append(fmt.Sprintf(
			"__Construct_%s__(&%s);\n", node.Lexeme, instanceId))
		sbPrev := ctx.sb
		ctx.sb = &sb
		for _, child := range node.Children {
			generate(child, ctx)
		}
		ctx.sb = sbPrev
		index := ctx.sb.FindLast(";")
		sb.TrimEnd('\n')
		ctx.sb.Insert(index+1, sb.ToString())
		ctx.sb.Append(instanceId)
		// list, exist := ctx.instancesToFree[ctx.block]
		// if !exist {
		// 	list = make([]string, 0)
		// 	ctx.instancesToFree[ctx.block] = list
		// }
		// list = append(list, instanceIdPointer)
		// ctx.instancesToFree[ctx.block] = list

	case ast.TypeStructArgument:
		writeTabs(ctx.sb, ctx.tabs)
		instanceId := fmt.Sprintf("__Instance_%d__",
			ctx.uniqueIndex)
		ctx.sb.Append(fmt.Sprintf("%s.%s = ", instanceId, node.Lexeme))
		generate(node.Children[0], ctx)
		ctx.sb.Append(";\n")
	}
}

func getDefaultValue(node *ast.Node) string {
	switch node.TypeHint {
	case "Int":
		return "0"
	case "Float":
		return "0.0"
	case "Bool":
		return "false"
	case "String":
		return "\"\""
	default:
		return ""
	}
}

func typeHintToString(lexeme string) string {
	switch lexeme {
	case "Int":
		return "int"
	case "Float":
		return "float"
	case "Bool":
		return "bool"
	case "String":
		return "char*"
	default:
		return lexeme
	}
}

func writeTabs(sb *Text.StringBuilder, tabs int) {
	for i := 0; i < tabs; i++ {
		sb.AppendRune('\t')
	}
}
