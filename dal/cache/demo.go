package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/loveRyujin/go-mall/common/enum"
	"github.com/loveRyujin/go-mall/common/logger"
	"github.com/loveRyujin/go-mall/logic/do"
)

func SetDemoOrder(ctx context.Context, demoOrder *do.DemoOrder) error {
	jsonData, _ := json.Marshal(demoOrder)
	redisKey := fmt.Sprintf(enum.REDIS_KEY_DEMO_ORDER_DETAIL, demoOrder.OrderNo)
	if err := Redis().Set(ctx, redisKey, jsonData, 0).Err(); err != nil {
		logger.New(ctx).Error("set demo order redis error", "err", err)
		return err
	}

	return nil
}

func GetDemoOrder(ctx context.Context, orderNo string) (*do.DemoOrder, error) {
	redisKey := fmt.Sprintf(enum.REDIS_KEY_DEMO_ORDER_DETAIL, orderNo)
	data, err := Redis().Get(ctx, redisKey).Result()
	if err != nil {
		logger.New(ctx).Error("get demo order redis error", "err", err)
		return nil, err
	}

	demoOrder := new(do.DemoOrder)
	if err := json.Unmarshal([]byte(data), demoOrder); err != nil {
		logger.New(ctx).Error("unmarshal demo order redis error", "err", err)
		return nil, err
	}

	return demoOrder, nil
}
