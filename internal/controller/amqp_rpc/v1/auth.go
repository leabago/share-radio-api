package v1

import (
	"fmt"

	"github.com/goccy/go-json"
	"github.com/leabago/share-radio/adder/internal/controller/amqp_rpc/v1/request"
	"github.com/leabago/share-radio/adder/pkg/jwt"
	amqp "github.com/rabbitmq/amqp091-go"
)

func extractUserID(d *amqp.Delivery, jwtManager *jwt.Manager) (userID string, data json.RawMessage, err error) {
	var req request.AuthenticatedRequest

	err = json.Unmarshal(d.Body, &req)
	if err != nil {
		return "", nil, fmt.Errorf("invalid request format: %w", err)
	}

	userID, err = jwtManager.ParseToken(req.Token)
	if err != nil {
		return "", nil, fmt.Errorf("invalid or expired token: %w", err)
	}

	return userID, req.Data, nil
}
