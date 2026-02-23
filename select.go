package squildx

import (
	"reflect"
	"strings"
)

func (b *builder) Select(columns ...string) Builder {
	cp := b.clone()
	cp.columns = append(cp.columns, columns...)
	return cp
}

func (b *builder) SelectObject(obj any, table ...string) Builder {
	cp := b.clone()
	prefix := ""
	if len(table) > 0 {
		prefix = table[0]
	}
	cols, err := structColumns(obj, prefix)
	if err != nil {
		cp.err = err
		return cp
	}
	cp.columns = append(cp.columns, cols...)
	return cp
}

func structColumns(obj any, table string) ([]string, error) {
	t := reflect.TypeOf(obj)
	if t == nil {
		return nil, ErrNotAStruct
	}
	for t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	if t.Kind() != reflect.Struct {
		return nil, ErrNotAStruct
	}

	var cols []string
	collectColumns(t, table, &cols)
	return cols, nil
}

func collectColumns(t reflect.Type, table string, cols *[]string) {
	for i := range t.NumField() {
		f := t.Field(i)
		if !f.IsExported() {
			continue
		}

		ft := f.Type
		for ft.Kind() == reflect.Ptr {
			ft = ft.Elem()
		}

		// Anonymous (embedded) struct: recurse to flatten
		if f.Anonymous && ft.Kind() == reflect.Struct {
			collectColumns(ft, table, cols)
			continue
		}

		// All other fields (including named structs like time.Time)
		name := ""
		for _, tag := range []string{"squildx", "db", "json"} {
			if v, ok := f.Tag.Lookup(tag); ok {
				v, _, _ = strings.Cut(v, ",")
				if v == "-" {
					name = "-"
					break
				}
				if v != "" {
					name = v
					break
				}
			}
		}
		if name == "-" {
			continue
		}
		if name == "" {
			name = toSnakeCase(f.Name)
		}
		if table != "" {
			name = table + "." + name
		}
		*cols = append(*cols, name)
	}
}

func (b *builder) RemoveSelect(columns ...string) Builder {
	cp := b.clone()
	remove := make(map[string]struct{}, len(columns))
	for _, c := range columns {
		remove[c] = struct{}{}
	}
	filtered := make([]string, 0, len(cp.columns))
	for _, c := range cp.columns {
		if _, ok := remove[c]; !ok {
			filtered = append(filtered, c)
		}
	}
	cp.columns = filtered
	return cp
}
