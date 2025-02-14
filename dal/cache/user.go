package cache

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/loveRyujin/go-mall/common/enum"
	"github.com/loveRyujin/go-mall/common/logger"
	"github.com/loveRyujin/go-mall/logic/do"
	"github.com/redis/go-redis/v9"
	"time"
)

func SetUserAuthToken(ctx context.Context, session *do.SessionInfo) error {
	if err := setAccessToken(ctx, session); err != nil {
		logger.New(ctx).Error("set access token error", "err", err)
		return err
	}
	if err := setRefreshToken(ctx, session); err != nil {
		logger.New(ctx).Error("set refresh token error", "err", err)
		return err
	}
	return nil
}

func SetUserSession(ctx context.Context, session *do.SessionInfo) error {
	key := fmt.Sprintf(enum.REDIS_KEY_USER_SESSION, session.UserId)
	sessionDataBytes, err := json.Marshal(session)
	if err != nil {
		logger.New(ctx).Error("marshal session data error", "err", err)
		return err
	}
	if err := Redis().HSet(ctx, key, session.Platform, sessionDataBytes).Err(); err != nil {
		logger.New(ctx).Error("set user session error", "err", err)
		return err
	}
	return nil
}

func DeleteOldSessionTokens(ctx context.Context, session *do.SessionInfo) error {
	oldSession, err := GetUserPlatformSession(ctx, session.UserId, session.Platform)
	if err != nil {
		logger.New(ctx).Error("get user platform session error", "err", err)
		return err
	}
	if oldSession == nil {
		return nil
	}
	if err := DeleteAccessToken(ctx, oldSession.AccessToken); err != nil {
		logger.New(ctx).Error("delete access token error", "err", err)
		return err
	}
	if err := DeleteRefreshToken(ctx, oldSession.RefreshToken); err != nil {
		logger.New(ctx).Error("delete refresh token error", "err", err)
		return err
	}
	return nil
}

func GetUserPlatformSession(ctx context.Context, userId int64, platform string) (*do.SessionInfo, error) {
	key := fmt.Sprintf(enum.REDIS_KEY_USER_SESSION, userId)
	res, err := Redis().HGet(ctx, key, platform).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, nil
		}
		return nil, err
	}
	session := new(do.SessionInfo)
	if err := json.Unmarshal([]byte(res), session); err != nil {
		logger.New(ctx).Error("unmarshal session data error", "err", err)
		return nil, err
	}
	return session, nil
}

func setAccessToken(ctx context.Context, session *do.SessionInfo) error {
	key := fmt.Sprintf(enum.REDIS_KEY_ACCESS_TOKEN, session.AccessToken)
	sessionDataBytes, err := json.Marshal(session)
	if err != nil {
		logger.New(ctx).Error("marshal session data error", "err", err)
		return err
	}
	res, err := Redis().Set(ctx, key, sessionDataBytes, enum.AccessTokenDuration).Result()
	if err != nil {
		return err
	}
	logger.New(ctx).Debug("set access token success", "res", res)
	return nil
}

func setRefreshToken(ctx context.Context, session *do.SessionInfo) error {
	key := fmt.Sprintf(enum.REDIS_KEY_REFRESH_TOKEN, session.RefreshToken)
	sessionDataBytes, err := json.Marshal(session)
	if err != nil {
		logger.New(ctx).Error("marshal session data error", "err", err)
		return err
	}
	if err := Redis().Set(ctx, key, sessionDataBytes, enum.RefreshTokenDuration).Err(); err != nil {
		return err
	}
	return nil
}

func DeleteAccessToken(ctx context.Context, accessToken string) error {
	key := fmt.Sprintf(enum.REDIS_KEY_ACCESS_TOKEN, accessToken)
	if err := Redis().Del(ctx, key).Err(); err != nil {
		return err
	}
	return nil
}

// DelayDeleteRefreshToken 刷新Token时让旧的RefreshToken 保留一段时间自己过期
func DelayDeleteRefreshToken(ctx context.Context, refreshToken string) error {
	key := fmt.Sprintf(enum.REDIS_KEY_REFRESH_TOKEN, refreshToken)
	return Redis().Expire(ctx, key, enum.OldRefreshTokenHoldingDuration).Err()
}

// DeleteRefreshToken 直接删除RefreshToken缓存 修改密码、退出登录时使用
func DeleteRefreshToken(ctx context.Context, refreshToken string) error {
	key := fmt.Sprintf(enum.REDIS_KEY_REFRESH_TOKEN, refreshToken)
	if err := Redis().Del(ctx, key).Err(); err != nil {
		return err
	}
	return nil
}

func LockTokenRefresh(ctx context.Context, refreshToken string) (bool, error) {
	key := fmt.Sprintf(enum.REDISKEY_TOKEN_REFRESH_LOCK, refreshToken)
	return Redis().SetNX(ctx, key, "locked", 10*time.Second).Result()
}

func UnLockTokenRefresh(ctx context.Context, refreshToken string) error {
	key := fmt.Sprintf(enum.REDISKEY_TOKEN_REFRESH_LOCK, refreshToken)
	return Redis().Del(ctx, key).Err()
}

func GetRefreshToken(ctx context.Context, refreshToken string) (*do.SessionInfo, error) {
	key := fmt.Sprintf(enum.REDIS_KEY_REFRESH_TOKEN, refreshToken)
	session := new(do.SessionInfo)
	result, err := Redis().Get(ctx, key).Result()
	if err != nil {
		if !errors.Is(err, redis.Nil) {
			return nil, err
		}
		return session, nil
	}
	if err := json.Unmarshal([]byte(result), session); err != nil {
		return nil, err
	}
	return session, nil
}

func GetAccessToken(ctx context.Context, accessToken string) (*do.SessionInfo, error) {
	key := fmt.Sprintf(enum.REDIS_KEY_ACCESS_TOKEN, accessToken)
	session := new(do.SessionInfo)
	result, err := Redis().Get(ctx, key).Result()
	if err != nil {
		if !errors.Is(err, redis.Nil) {
			return nil, err
		}
		return session, nil
	}
	if err := json.Unmarshal([]byte(result), session); err != nil {
		return nil, err
	}
	return session, nil
}
