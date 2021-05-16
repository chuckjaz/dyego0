package binder

import (
	"dyego0/ast"
	"dyego0/types"
)

type constEval struct {
	typeCache *typeCache
}

func newConstEval(typeCache *typeCache) constEval {
	return constEval{typeCache: typeCache}
}

const (
	tBoolean int = iota
	tByte
	tDouble
	tFloat
	tInt
	tLong
	tNull
	tRune
	tString
	tUInt
	tError
	tNotConst
)

func nameOf(typeCode int) string {
	switch typeCode {
	case tBoolean:
		return "Boolean"
	case tByte:
		return "Byte"
	case tDouble:
		return "Double"
	case tFloat:
		return "Float"
	case tInt:
		return "Int"
	case tLong:
		return "Long"
	case tNull:
		return "Null"
	case tRune:
		return "Rune"
	case tString:
		return "String"
	case tUInt:
		return "UInt"
	}
	return "<Unkonwn>"
}

func defaultOf(typeCode int) interface{} {
	switch typeCode {
	case tBoolean:
		return false
	case tByte:
		return byte(0)
	case tDouble:
		return 0.0
	case tFloat:
		return float(0.0)
	case tInt:
		return 0
	case tLong:
		return long(0)
	case tNull:
		return nil
	case tRune:
		return '0'
	case tString:
		return ""
	case tUInt:
		return uint(0)
	}
	return nil
}

type constResult struct {
	value    interface{}
	typeCode int
	typ      types.TypeSymbol
}

func (ce *constEval) Eval(element ast.Element) constResult {
	switch n := element.(type) {
	case ast.LiteralBoolean:
		typ := ce.typeCache.findBuiltinType("Boolean", n)
		return constResult{value: n.Value(), typeCode: tBoolean, typ: typ}
	case ast.LiteralByte:
		typ := ce.typeCache.findBuiltinType("Byte", n)
		return constResult{value: n.Value(), typeCode: tByte, typ: typ}
	case ast.LiteralDouble:
		typ := ce.typeCache.findBuiltinType("Double", n)
		return constResult{value: n.Value(), typeCode: tDouble, typ: typ}
	case ast.LiteralFloat:
		typ := ce.typeCache.findBuiltinType("Float", n)
		return constResult{value: n.Value(), typeCode: tFloat, typ: typ}
	case ast.LiteralInt:
		typ := ce.typeCache.findBuiltinType("Int", n)
		return constResult{value: n.Value(), typeCode: tInt, typ: typ}
	case ast.LiteralLong:
		typ := ce.typeCache.findBuiltinType("Long", n)
		return constResult{value: n.Value(), typeCode: tLong, typ: typ}
	case ast.LiteralNull:
		typ := ce.typeCache.findBuiltinType("Null", n)
		return constResult{value: nil, typeCode: tNull, typ: typ}
	case ast.LiteralRune:
		typ := ce.typeCache.findBuiltinType("Rune", n)
		return constResult{value: n.Value(), typeCode: tRune, typ: typ}
	case ast.LiteralString:
		typ := ce.typeCache.findBuiltinType("String", n)
		return constResult{value: n.Value(), typeCode: tString, typ: typ}
	case ast.Call:
		name, ok := n.Target().(ast.Name)
		if ok {
			expect := func(typeCode int, element ast.Element) constResult {
				result := ce.Eval(element)
				if result.typeCode != typeCode {
					ce.typeCache.error(element, "Expected a type of %s", nameOf(typeCode))
					return constResult{
						value:    defaultOf(typeCode),
						typ:      ce.typeCache.ErrorType(),
						typeCode: tError,
					}
				}
				return result
			}
			args := n.Arguments()
			switch len(args) {
			case 1:
				switch name.Text() {
				case "!":
					result := expect(tBoolean, args[0])
					result.value = !(result.value.(bool))
					return result
				case "+":
					result := ce.Eval(args[0])
					switch result.typeCode {
					case tInt, tFloat, tDouble, tLong:
						return result
					default:
						expect(tInt, args[0])
					}
				case "-":
					result := ce.Eval(args[0])
					switch result.typeCode {
					case tInt:
						result.value = -(result.value.(int))
						return result
					case tFloat:
						result.value = -(result.value.(float32))
						return result
					case tDouble:
						result.value = -(result.value.(float64))
						return result
					case tLong:
						result.value = -(result.value.(int64))
						return result
					default:
						expect(tInt, args[0])
					}
				}
			case 2:
				left := ce.Eval(args[0])
				right := expect(left.typeCode, args[1])
				boolOf := func(value bool) constResult {
					typ := ce.typeCache.findBuiltinType("Boolean", name)
					return constResult{
						value:    value,
						typeCode: tBoolean,
						typ:      typ,
					}
				}
				switch left.typeCode {
				case tBoolean:
					l := left.value.(bool)
					r := right.value.(bool)
					switch name.Text() {
					case "and":
						left.value = l && r
						return left
					case "or":
						left.value = l || r
						return left
					case "==":
						left.value = l == r
						return left
					case "!=":
						left.value = l != r
						return left
					}
				case tByte:
					l := left.value.(byte)
					r := right.value.(byte)
					switch name.Text() {
					case "+":
						left.value = l + r
						return left
					case "-":
						left.value = l - r
						return left
					case "*":
						left.value = l * r
						return left
					case "/":
						if r == 0 {
							ce.typeCache.error(args[1], "Devidd by zero")
						} else {
							left.value = l / r
						}
						return left
					case "%":
						left.value = l % r
						return left
					case "&":
						left.value = l & r
						return left
					case "^":
						left.value = l ^ r
						return left
					case ">":
						return boolOf(l > r)
					case "<":
						return boolOf(l < r)
					case ">=":
						return boolOf(l >= r)
					case "==":
						return boolOf(l == r)
					case "!=":
						return boolOf(l != r)
					}
				case tInt:
					l := left.value.(int)
					r := right.value.(int)
					switch name.Text() {
					case "+":
						left.value = l + r
						return left
					case "-":
						left.value = l - r
						return left
					case "*":
						left.value = l * r
						return left
					case "/":
						if r == 0 {
							ce.typeCache.error(args[1], "Devidd by zero")
						} else {
							left.value = l / r
						}
						return left
					case "%":
						left.value = l % r
						return left
					case "&":
						left.value = l & r
						return left
					case "^":
						left.value = l ^ r
						return left
					case ">":
						return boolOf(l > r)
					case "<":
						return boolOf(l < r)
					case ">=":
						return boolOf(l >= r)
					case "==":
						return boolOf(l == r)
					case "!=":
						return boolOf(l != r)
					}
				case tLong:
					l := left.value.(int64)
					r := right.value.(int64)
					switch name.Text() {
					case "+":
						left.value = l + r
						return left
					case "-":
						left.value = l - r
						return left
					case "*":
						left.value = l * r
						return left
					case "/":
						if r == 0 {
							ce.typeCache.error(args[1], "Devidd by zero")
						} else {
							left.value = l / r
						}
						return left
					case "%":
						left.value = l % r
						return left
					case "&":
						left.value = l & r
						return left
					case "^":
						left.value = l ^ r
						return left
					case ">":
						return boolOf(l > r)
					case "<":
						return boolOf(l < r)
					case ">=":
						return boolOf(l >= r)
					case "==":
						return boolOf(l == r)
					case "!=":
						return boolOf(l != r)
					}
				case tFloat:
					l := left.value.(float32)
					r := right.value.(float32)
					switch name.Text() {
					case "+":
						left.value = l + r
						return left
					case "-":
						left.value = l - r
						return left
					case "*":
						left.value = l * r
						return left
					case "/":
						if r == 0 {
							ce.typeCache.error(args[1], "Devidd by zero")
						} else {
							left.value = l / r
						}
						return left
					case ">":
						return boolOf(l > r)
					case "<":
						return boolOf(l < r)
					case ">=":
						return boolOf(l >= r)
					case "==":
						return boolOf(l == r)
					case "!=":
						return boolOf(l != r)
					}
				case tDouble:
					l := left.value.(float64)
					r := right.value.(float64)
					switch name.Text() {
					case "+":
						left.value = l + r
						return left
					case "-":
						left.value = l - r
						return left
					case "*":
						left.value = l * r
						return left
					case "/":
						if r == 0 {
							ce.typeCache.error(args[1], "Devidd by zero")
						} else {
							left.value = l / r
						}
						return left
					case ">":
						return boolOf(l > r)
					case "<":
						return boolOf(l < r)
					case ">=":
						return boolOf(l >= r)
					case "==":
						return boolOf(l == r)
					case "!=":
						return boolOf(l != r)
					}
				case tString:
					l := left.value.(string)
					r := right.value.(string)
					switch name.Text() {
					case "+":
						left.value = l + r
						return left
					case ">":
						return boolOf(l > r)
					case "<":
						return boolOf(l < r)
					case ">=":
						return boolOf(l >= r)
					case "==":
						return boolOf(l == r)
					case "!=":
						return boolOf(l != r)
					}
				}
			}
		}
	}
	return constResult{
		typeCode: tNotConst,
	}
}
