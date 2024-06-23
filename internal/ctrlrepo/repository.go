package ctrlrepo

import "github.com/dimitargrozev5/expenses-go-1/internal/models"

type ControllerRepository interface {
	// DB status
	GetVersion() (int64, error)

	// Users
	GetMinUserVersion() (int64, error)
	GetMaxUserVersion() (int64, error)
	AddNewUser(email string, version int64) error

	// Node actions
	GetNodes() ([]models.DBNode, error)
	GetActiveNodes() ([]models.DBNode, error)
	NewNode() (int64, error)
	RegisterNode(params *models.DBNodeData) (*models.GrpcEmpty, error)
}
