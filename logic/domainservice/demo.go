package domainservice

import (
	"context"
	"github.com/loveRyujin/go-mall/common/errcode"
	"github.com/loveRyujin/go-mall/common/utils"
	"github.com/loveRyujin/go-mall/dal/dao"
	"github.com/loveRyujin/go-mall/logic/do"
)

type DemoDomainService struct {
	ctx     context.Context
	DemoDao *dao.DemoDao
}

func NewDemoDomainService(ctx context.Context) *DemoDomainService {
	return &DemoDomainService{
		ctx:     ctx,
		DemoDao: dao.NewDemoDao(ctx),
	}
}

func (dds *DemoDomainService) GetDemos() ([]*do.DemoOrder, error) {
	demos, err := dds.DemoDao.GetAllDemos()
	if err != nil {
		err = errcode.Wrap("query entity error", err)
		return nil, err
	}

	demoOrders := make([]*do.DemoOrder, 0, len(demos))
	for _, demo := range demos {
		demoOrder := new(do.DemoOrder)
		if err := utils.CopyProperties(demoOrder, demo); err != nil {
			err = errcode.Wrap("CopyProperties失败", err)
			return nil, err
		}
		demoOrders = append(demoOrders, demoOrder)
	}

	return demoOrders, nil
}

func (dds *DemoDomainService) CreateDemoOrder(demoOrder *do.DemoOrder) (*do.DemoOrder, error) {
	// 生成订单号  先随便写个
	demoOrder.OrderNo = "20240627596615375920904456"

	demoOrderModel, err := dds.DemoDao.CreateDemoOrder(demoOrder)
	if err != nil {
		err = errcode.Wrap("创建DemoOrder失败", err)
		return nil, err
	}

	// TODO1: 写订单快照
	// 这里一般要在事务里写订单商品快照表, 这个等后面做需求时再演示
	if err := utils.CopyProperties(demoOrder, demoOrderModel); err != nil {
		err = errcode.Wrap("CopyProperties失败", err)
		return nil, err
	}
	// 返回领域对象
	return demoOrder, nil
}
