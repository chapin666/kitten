package goworkflow

import (
	"context"

	"goworkflow/pkg/expression"
)

type (
	expKey  struct{}
	flagKey struct{}
)

// FromExpContext 获取表达式的上下文
func FromExpContext(ctx context.Context) (expression.ExpContext, bool) {
	exp, ok := ctx.Value(expKey{}).(expression.ExpContext)
	return exp, ok
}

// FromFlagContext 获取flag的上下文
func FromFlagContext(ctx context.Context) (string, bool) {
	flag, ok := ctx.Value(flagKey{}).(string)
	return flag, ok
}
