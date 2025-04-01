package client

import "context"

func contextWithFileName(ctx context.Context, filename string) context.Context {
	return context.WithValue(ctx, filenameKey{}, filename)
}

func filenameFromContextOrPanic(ctx context.Context) string {
	filename, ok := filenameFromContext(ctx)
	if !ok {
		panic("could not get filename from context")
	}
	return filename
}

func filenameFromContext(ctx context.Context) (string, bool) {
	res, ok := ctx.Value(filenameKey{}).(string)
	return res, ok
}

func contextWithPattern(ctx context.Context, pattern string) context.Context {
	return context.WithValue(ctx, patternKey{}, pattern)
}

func patternFromContextOrPanic(ctx context.Context) string {
	pattern, ok := patternFromContext(ctx)
	if !ok {
		panic("could not get pattern from context")
	}
	return pattern
}

func patternFromContext(ctx context.Context) (string, bool) {
	res, ok := ctx.Value(patternKey{}).(string)
	return res, ok
}

type filenameKey struct{}
type patternKey struct{}
