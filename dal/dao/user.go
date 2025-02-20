package dao

import (
	"context"
	"errors"
	"github.com/loveRyujin/go-mall/common/errcode"
	"github.com/loveRyujin/go-mall/common/utils"
	"github.com/loveRyujin/go-mall/dal/model"
	"github.com/loveRyujin/go-mall/logic/do"
	"gorm.io/gorm"
)

type UserDao struct {
	ctx context.Context
}

func NewUserDao(ctx context.Context) *UserDao {
	return &UserDao{ctx: ctx}
}

func (ud *UserDao) CreateUser(userInfo *do.UserBaseInfo, hashPassword string) (*model.User, error) {
	user := new(model.User)
	if err := utils.CopyProperties(user, userInfo); err != nil {
		err = errcode.Wrap("UserDaoCreateUserError", err)
		return nil, err
	}
	user.Password = hashPassword

	if err := DBMaster().WithContext(ud.ctx).Create(user).Error; err != nil {
		err = errcode.Wrap("UserDaoCreateUserError", err)
		return nil, err
	}
	return user, nil
}

func (ud *UserDao) FetchUserByLoginName(loginName string) (*model.User, error) {
	user := new(model.User)
	if err := DB().Where(model.User{LoginName: loginName}).First(user).Error; err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		err = errcode.Wrap("UserDaoFetchUserError", err)
		return nil, err
	}
	return user, nil
}
