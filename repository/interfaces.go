// This file contains the interfaces for the repository layer.
// The repository layer is responsible for interacting with the database.
// For testing purpose we will generate mock implementations of these
// interfaces using mockgen. See the Makefile for more information.
package repository

import (
	"context"
)

type RepositoryInterface interface {
	GetTestById(ctx context.Context, input GetTestByIdInput) (output GetTestByIdOutput, err error)
	RegisterUser(ctx context.Context, input RegisterUser) (string, error)
	GetUserByPhone(ctx context.Context, phone string) (User, error)
	IncrSuccessLogin(ctx context.Context, phone string) error
	UpdateUser(ctx context.Context, input UpdateUser, identifier string) error
}
