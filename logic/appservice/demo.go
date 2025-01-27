package appservice

import (
	"context"
	"github.com/loveRyujin/go-mall/apis/reply"
	"github.com/loveRyujin/go-mall/apis/request"
	"github.com/loveRyujin/go-mall/common/errcode"
	"github.com/loveRyujin/go-mall/common/utils"
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
		errcode.Wrap("request转换成demoOrderDo失败", err)
		return nil, err
	}

	demoOrderDo, err := das.DemoDomainService.CreateDemoOrder(demoOrderDo)
	if err != nil {
		return nil, err
	}

	// 做一些其他的创建订单成功后的外围逻辑
	// 比如异步发送创建订单创建通知

	reply := new(reply.DemoOrder)
	if err := utils.CopyProperties(reply, demoOrderDo); err != nil {
		errcode.Wrap("demoOrderDo转换成reply失败", err)
		return nil, err
	}
	return reply, nil
}
