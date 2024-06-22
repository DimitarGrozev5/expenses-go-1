package ctrlrepo

import "github.com/dimitargrozev5/expenses-go-1/internal/models"

type ControllerRepository interface {
	// Node actions
	RegisterNode(params *models.DBNodeData) (*models.GrpcEmpty, error)
}
