package utils

import (
	"strconv"
	"strings"
)

type EvalOperator string

const (
	OperatorEqual            EvalOperator = "="
	OperatorGreaterThan      EvalOperator = ">"
	OperatorLessThan         EvalOperator = "<"
	OperatorGreaterThanEqual EvalOperator = ">="
	OperatorLessThanEqual    EvalOperator = "<="
	OperatorUndefined        EvalOperator = ""
)

type EvalDataType string

const (
	DataTypeFloat64   EvalDataType = "float64"
	DataTypeFloat32   EvalDataType = "float32"
	DataTypeInt64     EvalDataType = "int64"
	DataTypeInt32     EvalDataType = "int32"
	DataTypeInt       EvalDataType = "int"
)

func (eo EvalOperator) Value() string {
	return string(eo)
}

type Eval struct {
	expr     string
	operator EvalOperator
	value    float64
}

func NewEval() *Eval {
	return &Eval{}
}

func (e *Eval) SetExpression(expr string) *Eval {
	e.expr = strings.Trim(expr, " ")
	e.operator = OperatorUndefined
	for _, op := range e.supportedOperators() {
		if !strings.HasPrefix(e.expr, op.Value()) {
			continue
		}
		vv, err := strconv.ParseFloat(strings.TrimLeft(e.expr, op.Value()), 64)
		if err != nil {
			continue
		}
		e.operator = op
		e.value = vv
		break
	}
	return e
}

func (e *Eval) supportedOperators() []EvalOperator {
	return []EvalOperator{
		OperatorEqual,
		OperatorGreaterThan, OperatorLessThan,
		OperatorGreaterThanEqual, OperatorLessThanEqual,
	}
}

func (e *Eval) supportedDataTypes() []EvalDataType {
	return []EvalDataType{
		DataTypeFloat64, DataTypeFloat32,
		DataTypeInt64, DataTypeInt32,
		DataTypeInt,
	}
}

func (e *Eval) Evaluate(value interface{}) bool {
	if e.operator == OperatorUndefined {
		return false
	}
	if value == nil {
		return false
	}
	for _, dt := range e.supportedDataTypes() {
		switch dt {
		case DataTypeFloat64:
			if v, ok := value.(float64); ok {
				return e.eval(e.operator, v)
			}
		case DataTypeFloat32:
			if v, ok := value.(float32); ok {
				return e.eval(e.operator, float64(v))
			}
		case DataTypeInt64:
			if v, ok := value.(int64); ok {
				return e.eval(e.operator, float64(v))
			}
		case DataTypeInt32:
			if v, ok := value.(int32); ok {
				return e.eval(e.operator, float64(v))
			}
		case DataTypeInt:
			if v, ok := value.(int); ok {
				return e.eval(e.operator, float64(v))
			}
		}
	}
	return false
}

func (e *Eval) eval(op EvalOperator, value float64) bool {
	switch op {
	case OperatorEqual:
		return value == e.value
	case OperatorLessThan:
		return value < e.value
	case OperatorGreaterThan:
		return value > e.value
	case OperatorLessThanEqual:
		return value <= e.value
	case OperatorGreaterThanEqual:
		return value >= e.value
	default:
		return false
	}
}

