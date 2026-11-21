package progress

import "context"

type reporterContextKey struct{}

func WithReporter(ctx context.Context, reporter Reporter) context.Context {
	if reporter == nil {
		return ctx
	}

	return context.WithValue(ctx, reporterContextKey{}, reporter)
}

func ReporterFrom(ctx context.Context) Reporter {
	if ctx == nil {
		return nil
	}

	reporter, _ := ctx.Value(reporterContextKey{}).(Reporter)
	return reporter
}
