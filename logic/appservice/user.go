package appservice

import (
	"context"
	"github.com/loveRyujin/go-mall/apis/reply"
	"github.com/loveRyujin/go-mall/common/logger"
	"github.com/loveRyujin/go-mall/common/utils"
	"github.com/loveRyujin/go-mall/logic/domainservice"
)

type UserAppService struct {
	ctx              context.Context
	usrDomainService *domainservice.UserDomainService
}

func NewUserAppService(ctx context.Context) *UserAppService {
	return &UserAppService{
		ctx:              ctx,
		usrDomainService: domainservice.NewUserDomainService(ctx),
	}
}

func (uas *UserAppService) GenToken() (*reply.TokenReply, error) {
	token, err := uas.usrDomainService.GenAuthToken(12345678, "h5", "")
	if err != nil {
		return nil, err
	}
	logger.New(uas.ctx).Info("generate token success", "tokenData", token)
	tokenReply := new(reply.TokenReply)
	if err := utils.CopyProperties(tokenReply, token); err != nil {
		return nil, err
	}
	return tokenReply, err
}

func (uas *UserAppService) RefreshToken(refreshToken string) (*reply.TokenReply, error) {
	token, err := uas.usrDomainService.RefreshToken(refreshToken)
	if err != nil {
		return nil, err
	}
	logger.New(uas.ctx).Info("refresh token success", "tokenData", token)
	tokenReply := new(reply.TokenReply)
	if err := utils.CopyProperties(tokenReply, token); err != nil {
		return nil, err
	}
	return tokenReply, err
}
