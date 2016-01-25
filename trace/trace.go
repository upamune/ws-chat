package trace

import (
	"fmt"
	"io"
)

// Tracer is interface
type Tracer interface {
	Trace(...interface{})
}

type tracer struct {
	out io.Writer
}

type nilTracer struct{}

func (t *nilTracer) Trace(a ...interface{}) {}

// Off return ignoring calling method Tracer
func Off() Tracer {
	return &nilTracer{}
}

func (t *tracer) Trace(a ...interface{}) {
	if t.out == nil {
		return
	}

	fmt.Fprintln(t.out, a...)
}

// New return Tracer
func New(w io.Writer) Tracer {
	return &tracer{out: w}
}
