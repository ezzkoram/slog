package slog

import (
	"context"
)

type Logger interface {
	// Debug means a potentially noisy log.
	Debug(ctx context.Context, msg string, fields ...interface{})
	// Info means an informational log.
	Info(ctx context.Context, msg string, fields ...interface{})
	// Warn means something may be going wrong.
	Warn(ctx context.Context, msg string, fields ...interface{})
	// Error means the an error occured but does not require immediate attention.
	Error(ctx context.Context, msg string, fields ...interface{})
	// Critical means an error occured and requires immediate attention.
	Critical(ctx context.Context, msg string, fields ...interface{})
	// Fatal is the same as critical but calls os.Exit(1) afterwards.
	Fatal(ctx context.Context, msg string, fields ...interface{})

	// With returns a logger that will merge the given fields with all fields logged.
	// Fields logged with one of the above methods or from the context will always take priority.
	// Use the global with function when the fields being stored belong in the context and this
	// when they do not.
	With(fields ...interface{}) Logger
}

// field represents a log field.
type field struct {
	name  string
	value fieldValue
}

type Field interface {
	LogKey() string
	Value
}

type Value interface {
	LogValue() interface{}
}

type ValueFunc func() interface{}

func (fn ValueFunc) LogValue() interface{} {
	return fn()
}

type componentField string

// Component represents the component a log is being logged for.
// If there is already a component set, it will be joined by ".".
// E.g. if the component is currently "my_component" and then later
// the component "my_pkg" is set, then the final component will be
// "my_component.my_pkg".
func Component(name string) interface{} {
	return componentField(name)
}

type errorField struct {
	name string
	err  error
}

func (e errorField) LogKey() string {
	return e.name
}

func (e errorField) LogValue() interface{} {
	return e.err
}

// Error is the standard key used for logging a Go error value.
func Error(err error) Field {
	return errorField{
		name: "error",
		err:  err,
	}
}

type loggerKey struct{}

func withContext(ctx context.Context, l parsedFields) context.Context {
	return context.WithValue(ctx, loggerKey{}, l)
}

func fromContext(ctx context.Context) parsedFields {
	l, _ := ctx.Value(loggerKey{}).(parsedFields)
	return l
}

// With returns a context that contains the given fields.
// Any logs written with the provided context will contain
// the given fields.
func With(ctx context.Context, fields ...interface{}) context.Context {
	l := fromContext(ctx)
	l = l.withFields(fields)
	return withContext(ctx, l)
}
