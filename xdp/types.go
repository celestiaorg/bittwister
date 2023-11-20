package xdp

import (
	"context"

	"go.uber.org/zap"
)

type XdpLoader interface {
	Start(ctx context.Context, logger *zap.Logger)
	Ready() bool
}
