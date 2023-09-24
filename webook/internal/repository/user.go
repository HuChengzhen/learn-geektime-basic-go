package repository

import (
	"context"
	"learn-geektime-basic-go/webook/internal/domain"
	"learn-geektime-basic-go/webook/internal/repository/dao"
)

var (
	ErrUserDuplicateEmail = dao.ErrUserDuplicateEmail
	ErrUserNotFound       = dao.ErrUserNotFound
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
	err := ur.dao.Insert(ctx, u)
	return err
}

func (ur *UserRepository) FindByEmail(ctx context.Context, email string) (domain.User, error) {
	u, err := ur.dao.FindByEmail(ctx, email)

	if err != nil {
		return domain.User{}, err
	}

	return domain.User{
		Id:       u.Id,
		Email:    u.Email,
		Password: u.Password,
	}, nil
}
