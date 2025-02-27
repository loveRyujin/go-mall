package domainservice

import (
	"context"
	"github.com/loveRyujin/go-mall/common/enum"
	"github.com/loveRyujin/go-mall/common/errcode"
	"github.com/loveRyujin/go-mall/common/logger"
	"github.com/loveRyujin/go-mall/common/utils"
	"github.com/loveRyujin/go-mall/dal/cache"
	"github.com/loveRyujin/go-mall/dal/dao"
	"github.com/loveRyujin/go-mall/logic/do"
	"time"
)

type UserDomainService struct {
	ctx     context.Context
	userDao *dao.UserDao
}

func NewUserDomainService(ctx context.Context) *UserDomainService {
	return &UserDomainService{
		ctx:     ctx,
		userDao: dao.NewUserDao(ctx),
	}
}

func (uds *UserDomainService) RegisterUser(userInfo *do.UserBaseInfo, password string) (*do.UserBaseInfo, error) {
	user, err := uds.userDao.FetchUserByLoginName(userInfo.LoginName)
	if err != nil {
		err = errcode.Wrap("UserDomainServiceRegisterUserError", err)
		return nil, err
	}
	if user.LoginName != userInfo.LoginName {
		return nil, errcode.ErrUserNameOccupied
	}
	hashPassword, err := utils.BcryptPassword(password)
	if err != nil {
		err = errcode.Wrap("UserDomainServiceRegisterUserError", err)
		return nil, err
	}
	user, err = uds.userDao.CreateUser(userInfo, hashPassword)
	if err != nil {
		err = errcode.Wrap("UserDomainServiceRegisterUserError", err)
		return nil, err
	}
	if err = utils.CopyProperties(userInfo, user); err != nil {
		err = errcode.Wrap("UserDomainServiceRegisterUserError", err)
		return nil, err
	}
	return userInfo, nil
}

func (uds *UserDomainService) LoginUser(loginName, password, platform string) (*do.TokenInfo, error) {
	existedUser, err := uds.userDao.FetchUserByLoginName(loginName)
	if err != nil {
		err = errcode.Wrap("UserDomainServiceLoginUserError", err)
		return nil, err
	}
	if existedUser.ID == 0 {
		return nil, errcode.ErrUserNotRight
	}
	if !utils.BcryptCompare(existedUser.Password, password) {
		return nil, errcode.ErrUserNotRight
	}
	tokenInfo, err := uds.GenAuthToken(existedUser.ID, platform, "")
	if err != nil {
		err = errcode.Wrap("UserDomainServiceLoginUserError", err)
		return nil, err
	}
	return tokenInfo, nil
}

func (uds *UserDomainService) LogoutUser(userId int64, platform string) error {
	log := logger.New(uds.ctx)
	userSession, err := cache.GetUserPlatformSession(uds.ctx, userId, platform)
	if err != nil {
		log.Error("LogoutUserError", "err", err)
		err = errcode.Wrap("UserDomainServiceLogoutUserError", err)
		return err
	}
	if err = cache.DeleteAccessToken(uds.ctx, userSession.AccessToken); err != nil {
		log.Error("LogoutUserError", "err", err)
		err = errcode.Wrap("UserDomainServiceLogoutUserError", err)
		return err
	}
	if err = cache.DeleteRefreshToken(uds.ctx, userSession.RefreshToken); err != nil {
		log.Error("LogoutUserError", "err", err)
		err = errcode.Wrap("UserDomainServiceLogoutUserError", err)
		return err
	}
	// 删除用户对应平台的session
	if err = cache.DeleteUserSessionOnPlatform(uds.ctx, userId, platform); err != nil {
		log.Error("LogoutUserError", "err", err)
		err = errcode.Wrap("UserDomainServiceLogoutUserError", err)
		return err
	}
	return nil
}

// GetUserBaseInfo 因为还没开发注册登录功能, 这里先Mock一个返回
func (uds *UserDomainService) GetUserBaseInfo(uid int64) *do.UserBaseInfo {
	return &do.UserBaseInfo{
		ID:        12345678,
		Nickname:  "Kevin",
		LoginName: "kev@gomall.com",
		Verified:  1,
		Avatar:    "",
		Slogan:    "",
		IsBlocked: enum.UserBlockStateNormal,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

// GenAuthToken 生成用户的AccessToken和RefreshToken
// 在缓存中会存储最新的token，以及与Platform对应的UserSession 同时会删除缓存中旧的token，其中RefreshToken采用的延时删除
func (uds *UserDomainService) GenAuthToken(uid int64, platform string, sessionId string) (*do.TokenInfo, error) {
	user := uds.GetUserBaseInfo(uid)
	if user == nil {
		return nil, errcode.ErrUserInvalid
	}
	// 处理参数异常情况，用户不存在、被删除、被禁用
	if user.ID == 0 || user.IsBlocked == enum.UserBlockStateBlocked {
		return nil, errcode.ErrUserInvalid
	}
	userSession := new(do.SessionInfo)
	userSession.UserId = user.ID
	userSession.Platform = platform
	if sessionId == "" {
		// 用户登录，生成sessionId
		sessionId = utils.GenSessionId(uid)
	}
	userSession.SessionId = sessionId

	accessToken, refreshToken, err := utils.GenUserAuthToken(uid)
	if err != nil {
		err = errcode.Wrap("gen user auth token error", err)
		return nil, err
	}
	userSession.AccessToken = accessToken
	userSession.RefreshToken = refreshToken
	if err := cache.SetUserAuthToken(uds.ctx, userSession); err != nil {
		err = errcode.Wrap("set user auth token error", err)
		return nil, err
	}
	if err := cache.DeleteOldSessionTokens(uds.ctx, userSession); err != nil {
		err = errcode.Wrap("delete old session tokens error", err)
		return nil, err
	}
	if err := cache.SetUserSession(uds.ctx, userSession); err != nil {
		err = errcode.Wrap("set user session error", err)
		return nil, err
	}

	tokenInfo := &do.TokenInfo{
		AccessToken:   accessToken,
		RefreshToken:  refreshToken,
		Duration:      int64(enum.AccessTokenDuration.Seconds()),
		SrvCreateTime: time.Now(),
	}
	return tokenInfo, nil
}

func (uds *UserDomainService) RefreshToken(refreshToken string) (*do.TokenInfo, error) {
	ok, err := cache.LockTokenRefresh(uds.ctx, refreshToken)
	if err != nil {
		err = errcode.Wrap("刷新Token时设置redis锁错误", err)
		return nil, err
	}
	defer cache.UnLockTokenRefresh(uds.ctx, refreshToken)

	if !ok {
		err = errcode.ErrTooManyRequests
		return nil, err
	}
	tokenSession, err := cache.GetRefreshToken(uds.ctx, refreshToken)
	if err != nil {
		logger.New(uds.ctx).Error("GetRefreshToken cache error", "err", err)
		err = errcode.ErrToken
		return nil, err
	}
	if tokenSession == nil || tokenSession.UserId == 0 {
		err = errcode.ErrToken
		return nil, err
	}
	userSession, err := cache.GetUserPlatformSession(uds.ctx, tokenSession.UserId, tokenSession.Platform)
	if err != nil {
		logger.New(uds.ctx).Error("GetUserPlatformSession cache error", "err", err)
		err = errcode.ErrToken
		return nil, err
	}
	if userSession.RefreshToken != refreshToken {
		logger.New(uds.ctx).Warn("ExpiredRefreshToken", "requestToken", refreshToken, "newToken", userSession.RefreshToken, "userId", userSession.UserId)
		err = errcode.ErrToken
		return nil, err
	}

	// 重新生成token，因为非用户主动登录，sessionId不变
	tokenInfo, err := uds.GenAuthToken(tokenSession.UserId, tokenSession.Platform, userSession.SessionId)
	if err != nil {
		err = errcode.Wrap("GenAuthTokenErr", err)
		return nil, err
	}
	return tokenInfo, nil
}

func (uds *UserDomainService) VerifyToken(accessToken string) (*do.TokenVerify, error) {
	tokenInfo, err := cache.GetAccessToken(uds.ctx, accessToken)
	if err != nil {
		logger.New(uds.ctx).Error("GetAccessToken cache error", "err", err)
		err = errcode.ErrToken
		return nil, err
	}
	tokenVerify := new(do.TokenVerify)
	if tokenInfo != nil && tokenInfo.UserId != 0 {
		tokenVerify.Approved = true
		tokenVerify.UserId = tokenInfo.UserId
		tokenVerify.SessionId = tokenInfo.SessionId
	} else {
		tokenVerify.Approved = false
	}
	return tokenVerify, nil
}
