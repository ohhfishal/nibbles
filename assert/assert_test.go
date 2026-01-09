package assert_test

import (
	"github.com/ohhfishal/nibbles/assert"
	"testing"
)

var _ assert.Handler = &MockAsserter{}

func TestAssertTrue(t *testing.T) {
	tests := []struct {
		Name  string
		Expr  bool
		Args  []any
		Count int
	}{
		{Name: "true", Expr: true, Count: 0},
		{Name: "false", Expr: false, Count: 1},
		{Name: "true str", Expr: true, Count: 0, Args: []any{"test"}},
		{Name: "false str", Expr: false, Count: 1, Args: []any{"test"}},
		{Name: "true int", Expr: true, Count: 1, Args: []any{1}},
		{Name: "false int", Expr: false, Count: 2, Args: []any{1}},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			var mockT MockAsserter
			assert := assert.With(&mockT)
			assert.True(test.Expr, test.Args...)
			if mockT.Called != test.Count {
				t.Fatalf("called: %d times expected: %d", mockT.Called, test.Count)
			}
		})
	}
}

func TestAssertNil(t *testing.T) {
	tests := []struct {
		Name  string
		Expr  any
		Count int
	}{
		{Name: "nil", Expr: nil, Count: 0},
		{Name: "true", Expr: true, Count: 1},
		{Name: "false", Expr: false, Count: 1},
		{Name: "struct pointer", Expr: t, Count: 1},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			var mockT MockAsserter
			assert := assert.With(&mockT)
			assert.Nil(test.Expr)
			assert.True(
				mockT.Called != test.Count,
				"called: %d times expected: %d", mockT.Called, test.Count,
			)
		})
	}
}

func TestWithT(t *testing.T) {
	assert := assert.With(t)
	assert.True(true)
	assert.True(1 == 1, "unreachable")
	assert.Nil(nil)
}

type MockAsserter struct {
	Called int
}

func (asserter *MockAsserter) Fatalf(format string, args ...any) {
	asserter.Called++
}

// ExampleTrue shows a basic assert.
func ExampleWith() {
	assert := assert.With(nil)
	assert.True(true, "Unreachable")
}

func ExampleAssert_True() {
	assert := assert.With(nil)
	assert.True(true, "Unreachable")
	assert.True(1 == 2, "1 does not equal 2")
	assert.Unreachable("TSET")
}
