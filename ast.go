package main

import (
	"fmt"
	"strings"
)

type AST interface {
	eval(*Evaluator)
	show()
}

type ASTBase struct {
	filename string
	lineno   int
}

type AssignAST struct {
	ASTBase
	lhs string
	rhs string
	op  string
}

func (ast *AssignAST) eval(ev *Evaluator) {
	ev.evalAssign(ast)
}

func (ast *AssignAST) evalRHS(ev *Evaluator, lhs string) Var {
	switch ast.op {
	case ":=":
		return SimpleVar{value: ev.evalExpr(ast.rhs), origin: "file"}
	case "=":
		return RecursiveVar{expr: ast.rhs, origin: "file"}
	case "+=":
		prev := ev.LookupVar(lhs)
		if !prev.IsDefined() {
			return RecursiveVar{expr: ast.rhs, origin: "file"}
		}
		return prev.Append(ev, ast.rhs)
	case "?=":
		prev := ev.LookupVar(lhs)
		if prev.IsDefined() {
			return prev
		}
		return RecursiveVar{expr: ast.rhs, origin: "file"}
	default:
		panic(fmt.Sprintf("unknown assign op: %q", ast.op))
	}
}

func (ast *AssignAST) show() {
	Log("%s=%q", ast.lhs, ast.rhs)
}

// Note we cannot be sure what this is, until all variables in |expr|
// are expanded.
type MaybeRuleAST struct {
	ASTBase
	expr      string
	cmds      []string
	cmdLineno int
}

func (ast *MaybeRuleAST) eval(ev *Evaluator) {
	ev.evalMaybeRule(ast)
}

func (ast *MaybeRuleAST) show() {
	Log("%s", ast.expr)
	for _, cmd := range ast.cmds {
		Log("\t%s", strings.Replace(cmd, "\n", `\n`, -1))
	}
}

type IncludeAST struct {
	ASTBase
	expr string
	op   string
}

func (ast *IncludeAST) eval(ev *Evaluator) {
	ev.evalInclude(ast)
}

func (ast *IncludeAST) show() {
	Log("include %s", ast.expr)
}

type IfAST struct {
	ASTBase
	op         string
	lhs        string
	rhs        string // Empty if |op| is ifdef or ifndef.
	trueStmts  []AST
	falseStmts []AST
}

func (ast *IfAST) eval(ev *Evaluator) {
	ev.evalIf(ast)
}

func (ast *IfAST) show() {
	// TODO
	Log("if")
}
