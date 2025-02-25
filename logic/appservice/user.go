package appservice

import (
	"context"
	"errors"
	"github.com/loveRyujin/go-mall/apis/reply"
	"github.com/loveRyujin/go-mall/apis/request"
	"github.com/loveRyujin/go-mall/common/errcode"
	"github.com/loveRyujin/go-mall/common/logger"
	"github.com/loveRyujin/go-mall/common/utils"
	"github.com/loveRyujin/go-mall/logic/do"
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

func (uas *UserAppService) RegisterUser(registerInfo *request.UserRegister) error {
	userInfo := new(do.UserBaseInfo)
	if err := utils.CopyProperties(userInfo, registerInfo); err != nil {
		return err
	}
	_, err := uas.usrDomainService.RegisterUser(userInfo, registerInfo.Password)
	if err != nil {
		if errors.Is(err, errcode.ErrUserNameOccupied) {
			// 重名导致的注册不成功不额外处理
			return err
		}
		// 此处可以发通知告知用户注册失败 | 记录日志，告警监控，提示有用户注册失败
		logger.New(uas.ctx).Error("failed to register user, err: ", "error", err)
		return err
	}
	// 此处可以写注册成功后的外围逻辑，比如注册成功后给用户发确认短信|邮件
	// 如果产品的逻辑是注册成功后立即登录，此处紧跟登录的逻辑
	return nil
}

func (uas *UserAppService) LoginUser(userLoginReq *request.UserLogin) (*reply.TokenReply, error) {
	tokenInfo, err := uas.usrDomainService.LoginUser(userLoginReq.Body.LoginName, userLoginReq.Body.Password, userLoginReq.Header.Platform)
	if err != nil {
		return nil, err
	}
	tokenReply := new(reply.TokenReply)
	if err = utils.CopyProperties(tokenReply, tokenInfo); err != nil {
		return nil, err
	}

	// 执行用户登录成功后发送消息通知之类的外围辅助型逻辑
	return tokenReply, nil
}

func (uas *UserAppService) LogoutUser(userId int64, platform string) error {
	return uas.usrDomainService.LogoutUser(userId, platform)
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
