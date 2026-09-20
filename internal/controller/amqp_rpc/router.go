package v1

import (
	v1 "github.com/leabago/share-radio/adder/internal/controller/amqp_rpc/v1"
	"github.com/leabago/share-radio/adder/pkg/jwt"
	"github.com/leabago/share-radio/adder/pkg/logger"
	"github.com/leabago/share-radio/adder/pkg/rabbitmq/rmq_rpc/server"
)

// NewRouter -.
func NewRouter(j *jwt.Manager, l logger.Interface) map[string]server.CallHandler {
	routes := make(map[string]server.CallHandler)

	{
		v1.NewRoutes(routes, j, l)
	}

	return routes
}
