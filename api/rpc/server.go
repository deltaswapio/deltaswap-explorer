package rpc

import (
	"github.com/deltaswapio/deltaswap/node/pkg/common"
	publicrpcv1 "github.com/deltaswapio/deltaswap/node/pkg/proto/publicrpc/v1"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type Server struct {
	Srv *grpc.Server
}

// NewServer creates a GRPC server.
func NewServer(h *Handler, logger *zap.Logger) *grpc.Server {
	grpcServer := common.NewInstrumentedGRPCServer(logger, common.GrpcLogDetailMinimal)
	publicrpcv1.RegisterPublicRPCServiceServer(grpcServer, h)
	return grpcServer
}
