package goworkflow

import (
	"context"
)

type (
	flagKey struct{}
)

// FromFlagContext 获取flag的上下文
func FromFlagContext(ctx context.Context) (string, bool) {
	flag, ok := ctx.Value(flagKey{}).(string)
	return flag, ok
}
