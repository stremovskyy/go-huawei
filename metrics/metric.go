package metrics

import (
	"context"
)

type Reporter interface {
	NewRequest(name string) Request
}

type Request interface {
	EndRequest(ctx context.Context, err error, httpResp []byte, metro string)
}

type NoOpReporter struct {
}

func (n NoOpReporter) NewRequest(_ string) Request {
	return noOpRequest{}
}

type noOpRequest struct {
}

func (n noOpRequest) EndRequest(_ context.Context, _ error, _ []byte, _ string) {
}
