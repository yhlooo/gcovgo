package gcov

import "context"

// gccVersionContextKey context.Context 中存储 gcc 版本信息的键
type gccVersionContextKey struct{}

// ContextWithGCCVersion 创建携带指定版本信息的 context.Context
func ContextWithGCCVersion(ctx context.Context, version Version) context.Context {
	return context.WithValue(ctx, gccVersionContextKey{}, version)
}

// GCCVersionFromContext 从 context.Context 获取 gcc 版本信息
func GCCVersionFromContext(ctx context.Context) Version {
	v, ok := ctx.Value(gccVersionContextKey{}).(Version)
	if ok {
		return v
	}
	return Version{}
}
