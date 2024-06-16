package dbnoderpc

import (
	"context"
	"fmt"

	"github.com/dimitargrozev5/expenses-go-1/internal/models"
)

func (m *DatabaseServer) GetCategories(ctx context.Context, params *models.GrpcEmpty) (*models.GetCategoriesReturns, error) {
	// Get db
	db, ok := m.GetDB(ctx)
	if !ok {
		return nil, fmt.Errorf("can't find user db connection")
	}

	ret, err := db.GetCategories(params)
	if err != nil {
		return nil, err
	}

	return ret, nil
}

func (m *DatabaseServer) GetCategoriesOverview(ctx context.Context, params *models.GrpcEmpty) (*models.GetCategoriesOverviewReturns, error) {
	// Get db
	db, ok := m.GetDB(ctx)
	if !ok {
		return nil, fmt.Errorf("can't find user db connection")
	}

	ret, err := db.GetCategoriesOverview(params)
	if err != nil {
		return nil, err
	}

	return ret, nil
}

func (m *DatabaseServer) AddCategory(ctx context.Context, params *models.AddCategoryParams) (*models.GrpcEmpty, error) {
	// Get db
	db, ok := m.GetDB(ctx)
	if !ok {
		return nil, fmt.Errorf("can't find user db connection")
	}

	ret, err := db.AddCategory(params)
	if err != nil {
		return nil, err
	}

	return ret, nil
}

func (m *DatabaseServer) ReorderCategory(ctx context.Context, params *models.ReorderCategoryParams) (*models.GrpcEmpty, error) {
	// Get db
	db, ok := m.GetDB(ctx)
	if !ok {
		return nil, fmt.Errorf("can't find user db connection")
	}

	ret, err := db.ReorderCategory(params)
	if err != nil {
		return nil, err
	}

	return ret, nil
}

func (m *DatabaseServer) DeleteCategory(ctx context.Context, params *models.DeleteCategoryParams) (*models.GrpcEmpty, error) {
	// Get db
	db, ok := m.GetDB(ctx)
	if !ok {
		return nil, fmt.Errorf("can't find user db connection")
	}

	ret, err := db.DeleteCategory(params)
	if err != nil {
		return nil, err
	}

	return ret, nil
}

func (m *DatabaseServer) ResetCategories(ctx context.Context, params *models.ResetCategoriesParams) (*models.GrpcEmpty, error) {
	// Get db
	db, ok := m.GetDB(ctx)
	if !ok {
		return nil, fmt.Errorf("can't find user db connection")
	}

	ret, err := db.ResetCategories(params)
	if err != nil {
		return nil, err
	}

	return ret, nil
}
