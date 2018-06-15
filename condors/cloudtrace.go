package condors

import (
	"context"

	"google.golang.org/appengine/log"

	"go.opencensus.io/exporter/stackdriver"
	"go.opencensus.io/trace"
)

func span(c context.Context, text string, f func()) context.Context {
	exporter, err := stackdriver.NewExporter(stackdriver.Options{})
	if err != nil {
		log.Errorf(c, "creating StackDriver exporter: %v", err)
		return c
	}
	trace.RegisterExporter(exporter)

	cc, span := trace.StartSpan(c, text)
	f()
	span.End()
	return cc
}
