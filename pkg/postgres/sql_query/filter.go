package sql_query

import (
	"fmt"
	"reflect"
	"strings"
)

type WhereCondition struct {
	parts   []string
	args    []any
	counter int
}

func NewWhere() *WhereCondition {
	return &WhereCondition{}
}

func (w *WhereCondition) And(cond string, args ...any) *WhereCondition {
	return w.add("AND", cond, args...)
}

func (w *WhereCondition) Or(cond string, args ...any) *WhereCondition {
	return w.add("OR", cond, args...)
}

func (w *WhereCondition) AndGroup(fn func(*WhereCondition)) *WhereCondition {
	return w.group("AND", fn)
}

func (w *WhereCondition) OrGroup(fn func(*WhereCondition)) *WhereCondition {
	return w.group("OR", fn)
}

func (w *WhereCondition) SQL() (string, []any) {
	if len(w.parts) == 0 {
		return "", nil
	}

	return " WHERE " + strings.Join(w.parts, " "), w.args
}

func (w *WhereCondition) AndNot(cond string, args ...any) *WhereCondition {
	return w.add("AND", "NOT ("+cond+")", args...)
}

func (w *WhereCondition) OrNot(cond string, args ...any) *WhereCondition {
	return w.add("OR", "NOT ("+cond+")", args...)
}

func (w *WhereCondition) AndNotGroup(fn func(*WhereCondition)) *WhereCondition {
	return w.group("AND_NOT", fn)
}

func (w *WhereCondition) OrNotGroup(fn func(*WhereCondition)) *WhereCondition {
	return w.group("OR_NOT", fn)
}

func (w *WhereCondition) group(op string, fn func(*WhereCondition)) *WhereCondition {
	sub := &WhereCondition{counter: w.counter}
	fn(sub)
	w.counter = sub.counter

	if len(sub.parts) > 0 {
		groupSQL := "(" + strings.Join(sub.parts, " ") + ")"

		switch op {
		case "AND_NOT":
			w.groupAndNot(groupSQL)
		case "OR_NOT":
			w.groupOrNot(groupSQL)
		default:
			w.groupDefault(groupSQL, op)
		}

		w.args = append(w.args, sub.args...)
	}

	return w
}

func (w *WhereCondition) groupAndNot(groupSQL string) {
	groupSQL = "NOT " + groupSQL
	if len(w.parts) == 0 {
		w.parts = append(w.parts, groupSQL)
	} else {
		w.parts = append(w.parts, "AND "+groupSQL)
	}
}

func (w *WhereCondition) groupOrNot(groupSQL string) {
	groupSQL = "NOT " + groupSQL
	if len(w.parts) == 0 {
		w.parts = append(w.parts, groupSQL)
	} else {
		w.parts = append(w.parts, "OR "+groupSQL)
	}
}

func (w *WhereCondition) groupDefault(groupSQL string, op string) {
	if len(w.parts) == 0 {
		w.parts = append(w.parts, groupSQL)
	} else {
		w.parts = append(w.parts, op+" "+groupSQL)
	}
}

func (w *WhereCondition) add(op, cond string, args ...any) *WhereCondition {
	var placeholders []string

	var flatArgs []any

	for _, arg := range args {
		val := reflect.ValueOf(arg)
		if val.Kind() == reflect.Slice || val.Kind() == reflect.Array {
			for i := range val.Len() {
				placeholders = append(placeholders, w.nextPlaceholder())
				flatArgs = append(flatArgs, val.Index(i).Interface())
			}

			cond = strings.Replace(cond, "?", strings.Join(placeholders, ", "), 1)
			placeholders = placeholders[:0]
		} else {
			ph := w.nextPlaceholder()
			cond = strings.Replace(cond, "?", ph, 1)

			flatArgs = append(flatArgs, arg)
		}
	}

	if len(w.parts) == 0 {
		w.parts = append(w.parts, cond)
	} else {
		w.parts = append(w.parts, op+" "+cond)
	}

	w.args = append(w.args, flatArgs...)

	return w
}

func (w *WhereCondition) nextPlaceholder() string {
	w.counter++

	return fmt.Sprintf("$%d", w.counter)
}
