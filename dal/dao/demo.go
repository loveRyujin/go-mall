package dao

import (
	"context"
	"github.com/loveRyujin/go-mall/common/utils"
	"github.com/loveRyujin/go-mall/dal/model"
	"github.com/loveRyujin/go-mall/logic/do"
)

type DemoDao struct {
	ctx context.Context
}

func NewDemoDao(ctx context.Context) *DemoDao {
	return &DemoDao{ctx: ctx}
}

func (demo *DemoDao) GetAllDemos() (demos []*model.DemoOrder, err error) {
	if err = DB().WithContext(demo.ctx).Find(&demos).Error; err != nil {
		return nil, err
	}
	return demos, err
}

func (demo *DemoDao) CreateDemoOrder(demoOrder *do.DemoOrder) (*model.DemoOrder, error) {
	model := new(model.DemoOrder)
	if err := utils.CopyProperties(model, demoOrder); err != nil {
		return nil, err
	}
	if err := DB().WithContext(demo.ctx).Create(model).Error; err != nil {
		return nil, err
	}
	return model, nil
}
