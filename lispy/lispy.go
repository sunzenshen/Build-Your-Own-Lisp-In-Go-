package lispy

import (
	"math/big"
	"strconv"
	"strings"

	mpc "github.com/sunzenshen/golispy/mpcinterface"
)

// Lispy is a collection of the Lispy parser definitions
type Lispy struct {
	numberParser, operatorParser, exprParser, lispyParser mpc.MpcParser
}

// CleanLispy is used after parsers initiated by InitLispy are not longer to be used
func CleanLispy(l Lispy) {
	mpc.MpcCleanup(l.numberParser, l.operatorParser, l.exprParser, l.lispyParser)
}

// Eval translates an AST into the final result of the represented instructions
func Eval(tree mpc.MpcAst) lval {
	if strings.Contains(mpc.GetTag(tree), "number") {
		num, err := strconv.ParseInt(mpc.GetContents(tree), 10, 0)
		if err != nil {
			return lvalErr(lerrBadNum)
		}
		return lvalNum(num)
	}
	op := mpc.GetOperator(tree)
	x := Eval(mpc.GetChild(tree, 2))
	i := 3
	for strings.Contains(mpc.GetTag(mpc.GetChild(tree, i)), "expr") {
		x = evalOp(x, op, Eval(mpc.GetChild(tree, i)))
		i++
	}
	return x
}

func evalOp(x lval, op string, y lval) lval {
	if x.ltype == lvalErrType {
		return x
	}
	if x.ltype == lvalErrType {
		return y
	}
	if strings.Contains(op, "+") || strings.Contains(op, "add") {
		return lvalNum(x.num + y.num)
	}
	if strings.Contains(op, "-") || strings.Contains(op, "sub") {
		return lvalNum(x.num - y.num)
	}
	if strings.Contains(op, "*") || strings.Contains(op, "mul") {
		return lvalNum(x.num * y.num)
	}
	if strings.Contains(op, "/") || strings.Contains(op, "div") {
		if y.num == 0 {
			return lvalErr(lerrDivZero)
		}
		return lvalNum(x.num / y.num)
	}
	if strings.Contains(op, "%") || strings.Contains(op, "mod") {
		return lvalNum(x.num % y.num)
	}
	if strings.Contains(op, "^") || strings.Contains(op, "pow") {
		z := big.NewInt(0)
		z.Exp(big.NewInt(x.num), big.NewInt(y.num), nil)
		return lvalNum(z.Int64())
	}
	return lvalErr(lerrBadOp)
}

// InitLispy returns the parsers for the Lispy language definition
func InitLispy() Lispy {
	number := mpc.MpcNew("number")
	operator := mpc.MpcNew("operator")
	expr := mpc.MpcNew("expr")
	lispy := mpc.MpcNew("lispy")
	language := "" +
		"number : /-?[0-9]+/                                                  ; " +
		"operator :   '+'   |   '-'   |   '*'   |   '/'   |   '%'   |   '^'     " +
		"         | \"add\" | \"sub\" | \"mul\" | \"div\" | \"mod\" | \"pow\" ; " +
		"expr     : <number> | '(' <operator> <expr>+ ')'                     ; " +
		"lispy    : /^/ <operator> <expr>+ /$/                                ; "
	mpc.MpcaLang(language, number, operator, expr, lispy)
	parserSet := Lispy{}
	parserSet.numberParser = number
	parserSet.operatorParser = operator
	parserSet.exprParser = expr
	parserSet.lispyParser = lispy
	return parserSet
}

// PrintAst prints the AST of a Lispy expression.
func (l *Lispy) PrintAst(input string) {
	mpc.PrintAst(input, l.lispyParser)
}

// ReadEval takes a string, tries to interpret it in Lispy
func (l *Lispy) ReadEval(input string, printErrors bool) lval {
	r, err := mpc.MpcParse(input, l.lispyParser)
	if err != nil {
		if printErrors {
			mpc.MpcErrPrint(&r)
		}
		mpc.MpcErrDelete(&r)
		return lvalErr(lerrParseFail)
	}
	defer mpc.MpcAstDelete(&r)
	return Eval(mpc.GetOutput(&r))
}

// ReadEvalPrint takes a string, tries to interpret it in Lispy, or returns an parsing error
func (l *Lispy) ReadEvalPrint(input string) {
	evalResult := l.ReadEval(input, true)
	lvalPrintLn(evalResult)
}
