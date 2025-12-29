// Package assert has syntax sugar for asserting conditions.
// Designed to be used for [testing] and general runtime in production
// and devel environments
package assert

import (
	"fmt"
	"testing"
)

// Assert our interface implementations are correct
var _ Handler = &testing.T{}
var _ Handler = &funcHandler{}

var global = With(PanicHandler())

// Global gets the global [Assert] instance. If not set, the default global uses [PanicHandler].
func Global() *Assert {
	return global
}

// SetGlobal sets the global [Assert] instance.
func SetGlobal(a *Assert) {
	global = a
}

// Shorthandle for [Global]().True
func True(condition bool, formatArgs ...any) {
	Global().True(condition, formatArgs...)
}

// Shorthandle for [Global]().Nil
func Nil(item any, formatArgs ...any) {
	Global().Nil(item, formatArgs...)
}

// Shorthandle for [Global]().Unreachable
func Unreachable() {
	Global().Unreachable()
}

/*
Handler contains the method to call when an assert fails.
In the standard library, [testing.T] implements this interface.
*/
type Handler interface {
	Fatalf(string, ...any)
}

type Assert struct {
	handler Handler
}

func With(handler Handler) *Assert {
	return &Assert{
		handler: handler,
	}
}

/*
True asserts that the provided the condition is true.
You may provide the args to Fatalf should the assertion fail.
Note, those arguments are asserted to be valid regardless of the value of condition.
*/
func (assert *Assert) True(condition bool, formatArgs ...any) {
	var fail func()
	switch len(formatArgs) {
	case 0:
		fail = func() { assert.handler.Fatalf("assert failed") }
	default:
		str, ok := formatArgs[0].(string)
		assert.True(ok)
		fail = func() { assert.handler.Fatalf(str, formatArgs[1:]...) }
	}

	if !condition {
		fail()
	}
}

// Nil asserts that item is nil. See [Assert.True] for formatArgs.
func (assert *Assert) Nil(item any, formatArgs ...any) {
	assert.helper(
		item == nil,
		wrap("Expected nil. Got: %v", item),
		formatArgs...,
	)
}

func (assert *Assert) Unreachable() {
	assert.helper(
		false,
		wrap("Unreachable code reached."),
	)
}

type message struct {
	args []any
}

func wrap(format string, args ...any) message {
	wrapped := []any{format}
	for _, arg := range args {
		wrapped = append(wrapped, arg)
	}
	return message{wrapped}
}

func (assert *Assert) helper(cond bool, defaultMsg message, formatArgs ...any) {
	switch len(formatArgs) {
	case 0:
		assert.True(cond, defaultMsg.args...)
	default:
		assert.True(cond, formatArgs...)
	}
}

// Returns a [Handler] that panics when called.
func PanicHandler() Handler {
	return funcHandler{
		F: func(format string, args ...any) {
			panic(fmt.Errorf(format, args...))
		},
	}
}

// Returns a [Handler] that calls f when called.
func FuncHandler(f func(string, ...any)) Handler {
	return funcHandler{
		F: f,
	}
}

type funcHandler struct {
	F func(string, ...any)
}

func (handler funcHandler) Fatalf(format string, args ...any) {
	handler.F(format, args...)
}
