package appservice

import (
	"context"
	"github.com/loveRyujin/go-mall/apis/reply"
	"github.com/loveRyujin/go-mall/apis/request"
	"github.com/loveRyujin/go-mall/common/errcode"
	"github.com/loveRyujin/go-mall/common/logger"
	"github.com/loveRyujin/go-mall/common/utils"
	"github.com/loveRyujin/go-mall/dal/cache"
	"github.com/loveRyujin/go-mall/logic/do"
	"github.com/loveRyujin/go-mall/logic/domainservice"
)

type DemoAppService struct {
	ctx               context.Context
	DemoDomainService *domainservice.DemoDomainService
}

func NewDemoAppService(ctx context.Context) *DemoAppService {
	return &DemoAppService{
		ctx:               ctx,
		DemoDomainService: domainservice.NewDemoDomainService(ctx),
	}
}

func (das *DemoAppService) GetDemoIds() ([]int64, error) {
	demos, err := das.DemoDomainService.GetDemos()
	if err != nil {
		return nil, err
	}

	ids := make([]int64, 0, len(demos))
	for _, demo := range demos {
		ids = append(ids, demo.Id)
	}
	return ids, nil
}

func (das *DemoAppService) CreateDemoOrder(orderRequest *request.DemoOrderCreate) (*reply.DemoOrder, error) {
	demoOrderDo := new(do.DemoOrder)
	if err := utils.CopyProperties(demoOrderDo, orderRequest); err != nil {
		err = errcode.Wrap("request转换成demoOrderDo失败", err)
		return nil, err
	}

	// 测试redis的使用，后续删除
	if err := cache.SetDemoOrder(das.ctx, demoOrderDo); err != nil {
		err = errcode.Wrap("SetDemoOrder失败", err)
		return nil, err
	}
	cacheData, err := cache.GetDemoOrder(das.ctx, demoOrderDo.OrderNo)
	if err != nil {
		err = errcode.Wrap("GetDemoOrder失败", err)
		return nil, err
	}
	logger.New(das.ctx).Info("redis cache data", "cacheData", cacheData)

	demoOrder, err := das.DemoDomainService.CreateDemoOrder(demoOrderDo)
	if err != nil {
		return nil, err
	}

	// 做一些其他的创建订单成功后的外围逻辑
	// 比如异步发送创建订单创建通知

	reply := new(reply.DemoOrder)
	if err := utils.CopyProperties(reply, demoOrder); err != nil {
		err = errcode.Wrap("demoOrderDo转换成reply失败", err)
		return nil, err
	}
	return reply, nil
}
