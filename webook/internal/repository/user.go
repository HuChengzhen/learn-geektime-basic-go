package repository

import (
	"context"
	"learn-geektime-basic-go/webook/internal/domain"
	"learn-geektime-basic-go/webook/internal/repository/dao"
)

type UserRepository struct {
	dao *dao.UserDAO
}

func NewUserRepository(dao *dao.UserDAO) *UserRepository {
	return &UserRepository{
		dao: dao,
	}
}

func (ur *UserRepository) FindById(id int64) {

}

func (ur *UserRepository) Create(ctx context.Context, user domain.User) error {
	u := dao.User{
		Email:    user.Email,
		Password: user.Password,
	}
	return ur.dao.Insert(ctx, u)
}
