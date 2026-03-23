package squildx

import (
	"reflect"
	"strings"
)

func (b *insertBuilder) Values(sql string, params ...Params) InsertBuilder {
	cp := b.clone()
	extracted, err := extractParams(params)
	if err != nil {
		cp.err = err
		return cp
	}
	merged, prefix, err := parseParams(sql, extracted)
	if err != nil {
		cp.err = err
		return cp
	}
	if err := cp.setPrefix(prefix); err != nil {
		cp.err = err
		return cp
	}
	cp.valueRows = append(cp.valueRows, paramClause{sql: sql, params: merged})
	return cp
}

func (b *insertBuilder) ValuesObject(obj any) InsertBuilder {
	cp := b.clone()
	cols, sql, params, err := structFieldValues(obj)
	if err != nil {
		cp.err = err
		return cp
	}
	if err := cp.setPrefix(':'); err != nil {
		cp.err = err
		return cp
	}
	switch {
	case len(cp.columns) == 0:
		cp.columns = cols
	case !reflect.DeepEqual(cp.columns, cols):
		cp.err = ErrColumnMismatch
		return cp
	}
	cp.valueRows = append(cp.valueRows, paramClause{sql: sql, params: params})
	return cp
}

func structFieldValues(obj any) (columns []string, sql string, params Params, err error) {
	v := reflect.ValueOf(obj)
	t := reflect.TypeOf(obj)
	if t == nil {
		return nil, "", nil, ErrNotAStruct
	}
	for t.Kind() == reflect.Ptr {
		t = t.Elem()
		v = v.Elem()
	}
	if t.Kind() != reflect.Struct {
		return nil, "", nil, ErrNotAStruct
	}

	params = Params{}
	collectFieldValues(t, v, "", &columns, params)

	placeholders := make([]string, len(columns))
	for i, col := range columns {
		placeholders[i] = ":" + col
	}
	sql = strings.Join(placeholders, ", ")
	return columns, sql, params, nil
}

func collectFieldValues(t reflect.Type, v reflect.Value, table string, cols *[]string, params Params) {
	for i := range t.NumField() {
		f := t.Field(i)
		if !f.IsExported() {
			continue
		}

		ft := f.Type
		fv := v.Field(i)
		for ft.Kind() == reflect.Ptr {
			if fv.IsNil() {
				break
			}
			ft = ft.Elem()
			fv = fv.Elem()
		}

		tagName := fieldTagName(f)
		if tagName == "-" {
			continue
		}

		if f.Anonymous && ft.Kind() == reflect.Struct && tagName == "" {
			collectFieldValues(ft, fv, table, cols, params)
			continue
		}

		name := tagName
		if name == "" {
			name = toSnakeCase(f.Name)
		}
		if table != "" {
			name = table + "." + name
		}
		*cols = append(*cols, name)
		params[name] = fv.Interface()
	}
}
