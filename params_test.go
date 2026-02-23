package squildx

import (
	"errors"
	"testing"
)

func TestDuplicateParamSameValue(t *testing.T) {
	_, _, err := New().Select("*").
		From("users").
		Where("age > :val", 18).
		Where("score > :val", 18).
		Build()

	if err != nil {
		t.Errorf("expected no error for duplicate param with same value, got: %v", err)
	}
}

func TestDuplicateParamDifferentValue(t *testing.T) {
	_, _, err := New().Select("*").
		From("users").
		Where("age > :val", 18).
		Where("score > :val", 99).
		Build()

	if !errors.Is(err, ErrDuplicateParam) {
		t.Errorf("expected ErrDuplicateParam, got: %v", err)
	}
}

func TestParamCountMismatch(t *testing.T) {
	_, _, err := New().Select("*").
		From("users").
		Where("age > :min AND age < :max", 18).
		Build()

	if !errors.Is(err, ErrParamMismatch) {
		t.Errorf("expected ErrParamMismatch, got: %v", err)
	}
}

func TestParamCountMismatchTooManyValues(t *testing.T) {
	_, _, err := New().Select("*").
		From("users").
		Where("age > :min", 18, 65).
		Build()

	if !errors.Is(err, ErrParamMismatch) {
		t.Errorf("expected ErrParamMismatch, got: %v", err)
	}
}
