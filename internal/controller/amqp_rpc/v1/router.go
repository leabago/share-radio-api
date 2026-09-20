package v1

import (
	"github.com/leabago/share-radio/adder/pkg/jwt"
	"github.com/leabago/share-radio/adder/pkg/logger"
	"github.com/leabago/share-radio/adder/pkg/rabbitmq/rmq_rpc/server"
)

// NewRoutes -.
func NewRoutes(routes map[string]server.CallHandler, j *jwt.Manager, l logger.Interface) {

}
