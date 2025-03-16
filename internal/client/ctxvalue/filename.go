package ctxvalue

import "context"

func ContextWithFileName(ctx context.Context, filename string) context.Context {
	return context.WithValue(ctx, filenameKey{}, filename)
}

func FilenameFromContext(ctx context.Context) (string, bool) {
	res, ok := ctx.Value(filenameKey{}).(string)
	return res, ok
}

func ContextWithPattern(ctx context.Context, pattern string) context.Context {
	return context.WithValue(ctx, patternKey{}, pattern)
}

func PatternFromContext(ctx context.Context) (string, bool) {
	res, ok := ctx.Value(patternKey{}).(string)
	return res, ok
}

type filenameKey struct{}
type patternKey struct{}
