package grpc

import (
	"github.com/leabago/share-radio/adder/pkg/logger"
	pbgrpc "google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// NewRouter -.
func NewRouter(app *pbgrpc.Server, l logger.Interface) {

	reflection.Register(app)
}
