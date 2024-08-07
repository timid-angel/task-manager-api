package usecase

import (
	"context"
	domain "task_manager_api/Domain"
	"time"
)

type UserUsecase struct {
	UserRespository domain.UserRepositoryInterface
	Timeout         time.Duration
}

func (uC *UserUsecase) CreateUser(c context.Context, user domain.User) domain.CodedError {
	ctx, cancel := context.WithTimeout(c, uC.Timeout)
	defer cancel()
	return uC.UserRespository.CreateUser(ctx, user)
}

func (uC *UserUsecase) ValidateAndGetToken(c context.Context, user domain.User) (string, domain.CodedError) {
	ctx, cancel := context.WithTimeout(c, uC.Timeout)
	defer cancel()
	return uC.UserRespository.ValidateAndGetToken(ctx, user)
}
