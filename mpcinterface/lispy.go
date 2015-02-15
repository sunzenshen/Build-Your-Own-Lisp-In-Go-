package mpcinterface

// #cgo LDFLAGS: -ledit -lm
// #include "mpc_interface.h"
import "C"
import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"unsafe"
)

// CleanLispy is used after parsers initiated by InitLispy are not longer to be used
func CleanLispy(number *C.struct_mpc_parser_t,
	operator *C.struct_mpc_parser_t,
	expr *C.struct_mpc_parser_t,
	lispy *C.struct_mpc_parser_t) {
	defer C.mpc_cleanup_if(C.int(4), number, operator, expr, lispy)
}

// Eval translates an AST into the final result of the represented instructions
func Eval(tree MpcAst) int64 {
	if strings.Contains(getTag(tree), "number") {
		num, _ := strconv.ParseInt(getContents(tree), 10, 0)
		return num
	}
	op := getOperator(tree)
	x := Eval(getChild(tree, 2))
	i := 3
	for strings.Contains(getTag(getChild(tree, i)), "expr") {
		x = evalOp(x, op, Eval(getChild(tree, i)))
		i++
	}
	return x
}

func evalOp(x int64, op string, y int64) int64 {
	if strings.Contains(op, "+") {
		return x + y
	}
	if strings.Contains(op, "-") {
		return x - y
	}
	if strings.Contains(op, "*") {
		return x * y
	}
	if strings.Contains(op, "/") {
		return x / y
	}
	return 0
}

// InitLispy returns the parsers for the Lispy language definition
func InitLispy() (
	*C.struct_mpc_parser_t,
	*C.struct_mpc_parser_t,
	*C.struct_mpc_parser_t,
	*C.struct_mpc_parser_t) {
	Number := mpcNew("number")
	Operator := mpcNew("operator")
	Expr := mpcNew("expr")
	Lispy := mpcNew("lispy")
	language := "" +
		"number : /-?[0-9]+/                               ; " +
		"operator : '+' | '-' | '*' | '/'                  ; " +
		"expr     : <number> | '(' <operator> <expr>+ ')'  ; " +
		"lispy    : /^/ <operator> <expr>+ /$/             ; "
	cLanguage := C.CString(language)
	defer C.free(unsafe.Pointer(cLanguage))
	C.mpca_lang_if(C.MPCA_LANG_DEFAULT,
		cLanguage,
		Number, Operator, Expr, Lispy)
	return Number, Operator, Expr, Lispy
}

// ReadEval takes a string, tries to interpret it in Lispy
func ReadEval(input string, mpcParser *C.struct_mpc_parser_t) (int64, error) {
	var r C.mpc_result_t
	cInput := C.CString(input)
	defer C.free(unsafe.Pointer(cInput))
	stdin := C.CString("<stdin>")
	defer C.free(unsafe.Pointer(stdin))
	if C.mpc_parse(stdin, cInput, mpcParser, &r) != C.int(0) {
		defer C.mpc_ast_delete(C.get_output(&r))
		return Eval(C.get_output(&r)), nil
	}
	return 0, errors.New("mpc: failed to parse input string")
}

// ReadEvalPrint takes a string, tries to interpret it in Lispy, or returns an parsing error
func ReadEvalPrint(input string, mpcParser *C.struct_mpc_parser_t) {
	var r C.mpc_result_t
	cInput := C.CString(input)
	defer C.free(unsafe.Pointer(cInput))
	stdin := C.CString("<stdin>")
	defer C.free(unsafe.Pointer(stdin))
	if C.mpc_parse(stdin, cInput, mpcParser, &r) != C.int(0) {
		fmt.Println(Eval(C.get_output(&r)))
		C.mpc_ast_delete(C.get_output(&r))
	} else {
		C.mpc_err_print(C.get_error(&r))
		C.mpc_err_delete(C.get_error(&r))
	}
}
