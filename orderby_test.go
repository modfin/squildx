package squildx

import (
	"errors"
	"reflect"
	"testing"
)

func TestOrderByMultiple(t *testing.T) {
	q, _, err := New().
		Select("*").
		From("users").
		OrderBy("name ASC").
		OrderBy("age DESC").
		OrderBy("id").
		Build()

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := "SELECT * FROM users ORDER BY name ASC, age DESC, id"
	if q != expected {
		t.Errorf("SQL mismatch\n got: %s\nwant: %s", q, expected)
	}
}

func TestOrderByWithParams(t *testing.T) {
	vec := []float64{0.1, 0.2, 0.3}
	q, params, err := New().
		Select("id", "title").
		From("documents").
		OrderBy("similarity(embedding, :query_vec) DESC", vec).
		Build()

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := "SELECT id, title FROM documents ORDER BY similarity(embedding, :query_vec) DESC"
	if q != expected {
		t.Errorf("SQL mismatch\n got: %s\nwant: %s", q, expected)
	}

	if !reflect.DeepEqual(params["query_vec"], vec) {
		t.Errorf("param mismatch: got %v, want %v", params["query_vec"], vec)
	}
}

func TestOrderByParamMismatch(t *testing.T) {
	_, _, err := New().
		Select("*").
		From("documents").
		OrderBy("similarity(embedding, :query_vec) DESC").
		Build()

	if !errors.Is(err, ErrParamMismatch) {
		t.Fatalf("expected ErrParamMismatch, got: %v", err)
	}
}

func TestOrderByWithWhere(t *testing.T) {
	vec := []float64{0.1, 0.2, 0.3}
	q, params, err := New().
		Select("id", "title").
		From("documents").
		Where("category = :cat", "science").
		OrderBy("similarity(embedding, :query_vec) DESC", vec).
		Build()

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := "SELECT id, title FROM documents WHERE category = :cat ORDER BY similarity(embedding, :query_vec) DESC"
	if q != expected {
		t.Errorf("SQL mismatch\n got: %s\nwant: %s", q, expected)
	}

	if params["cat"] != "science" {
		t.Errorf("param 'cat' mismatch: got %v, want %s", params["cat"], "science")
	}
	if !reflect.DeepEqual(params["query_vec"], vec) {
		t.Errorf("param 'query_vec' mismatch: got %v, want %v", params["query_vec"], vec)
	}
}
