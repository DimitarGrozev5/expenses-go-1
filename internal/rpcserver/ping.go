package rpcserver

import (
	"context"

	"github.com/dimitargrozev5/expenses-go-1/internal/models"
)

func (s *DatabaseServer) Ping(ctx context.Context, msg *models.SimpleMessage) (*models.SimpleMessage, error) {
	if msg.Msg == "Ping" {
		return &models.SimpleMessage{Msg: "Pong"}, nil
	}

	return &models.SimpleMessage{Msg: "No message"}, nil
}
