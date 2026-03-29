package squildx

import (
	"reflect"
	"strings"
)

func (b *updateBuilder) Set(sql string, params ...Params) UpdateBuilder {
	cp := b.clone()
	extracted, err := extractParams(params)
	if err != nil {
		cp.err = err
		return cp
	}
	parsed, prefix, err := parseParams(sql, extracted)
	if err != nil {
		cp.err = err
		return cp
	}
	if err := cp.setPrefix(prefix); err != nil {
		cp.err = err
		return cp
	}
	cp.sets = append(cp.sets, paramClause{sql: sql, params: parsed})
	return cp
}

func (b *updateBuilder) SetObject(obj any) UpdateBuilder {
	cp := b.clone()
	sql, params, err := structSetSQL(obj)
	if err != nil {
		cp.err = err
		return cp
	}
	if sql == "" {
		return cp
	}
	if err := cp.setPrefix(detectPrefix(sql)); err != nil {
		cp.err = err
		return cp
	}
	cp.sets = append(cp.sets, paramClause{sql: sql, params: params})
	return cp
}

func structSetSQL(obj any) (string, Params, error) {
	v := reflect.ValueOf(obj)
	t := reflect.TypeOf(obj)
	if t == nil {
		return "", nil, ErrNotAStruct
	}
	for t.Kind() == reflect.Ptr {
		t = t.Elem()
		v = v.Elem()
	}
	if t.Kind() != reflect.Struct {
		return "", nil, ErrNotAStruct
	}

	var assignments []string
	params := Params{}
	collectSetClauses(t, v, &assignments, params)
	return strings.Join(assignments, ", "), params, nil
}

func collectSetClauses(t reflect.Type, v reflect.Value, assignments *[]string, params Params) {
	for i := range t.NumField() {
		f := t.Field(i)
		if !f.IsExported() {
			continue
		}

		ft := f.Type
		fv := v.Field(i)

		tagName := fieldTagName(f)
		if tagName == "-" {
			continue
		}

		if f.Anonymous && ft.Kind() == reflect.Ptr {
			if fv.IsNil() {
				continue
			}
			ft = ft.Elem()
			fv = fv.Elem()
		}

		if f.Anonymous && ft.Kind() == reflect.Struct && tagName == "" {
			collectSetClauses(ft, fv, assignments, params)
			continue
		}

		if ft.Kind() == reflect.Ptr {
			if fv.IsNil() {
				continue
			}
			fv = fv.Elem()
		}

		name := tagName
		if name == "" {
			name = toSnakeCase(f.Name)
		}
		*assignments = append(*assignments, name+" = :"+name)
		params[name] = fv.Interface()
	}
}
