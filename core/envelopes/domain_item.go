package envelopes

import (
	"context"
	"github.com/segmentio/ksuid"
	"github.com/tietang/dbx"
	"imooc.com/resk/infra/base"
	"imooc.com/resk/services"
)

type itemDomain struct {
	RedEnvelopeItem
}

// 生成itemNo
func (domain *itemDomain) createItemNo() {
	domain.ItemNo = ksuid.New().Next().String()
}

// 创建item
func (domain *itemDomain) Create(dto *services.RedEnvelopeItemDTO) {
	domain.RedEnvelopeItem.FromDTO(dto)
	domain.RecvUsername.Valid = true
	domain.createItemNo()
}

// 保存item
func (domain *itemDomain) Save(ctx context.Context) (id int64, err error) {
	err = base.ExecuteContext(ctx, func(runner *dbx.TxRunner) error {
		dao := &RedEnvelopeItemDao{runner: runner}
		id, err = dao.Insert(&domain.RedEnvelopeItem)
		if err != nil {
			return err
		}
		return nil
	})
	return
}

// 通过itemNo查询抢红包明细数据
func (domain *itemDomain) GetOne(ctx context.Context, itemNo string) (dto *services.RedEnvelopeItemDTO) {
	err := base.ExecuteContext(ctx, func(runner *dbx.TxRunner) error {
		dao := &RedEnvelopeItemDao{runner: runner}
		data := dao.GetOne(itemNo)
		if data != nil {
			dto = data.ToDTO()
		}
		return nil
	})
	if err != nil {
		return nil
	}
	return dto
}

// 通过envelopeNo查询已抢红包列表
func (domain *itemDomain) FindItems(envelopeNo string) (itemDTOs []*services.RedEnvelopeItemDTO) {
	var items []*RedEnvelopeItem
	err := base.Tx(func(runner *dbx.TxRunner) error {
		dao := &RedEnvelopeItemDao{runner: runner}
		items = dao.FindItems(envelopeNo)
		return nil
	})
	if err != nil {
		return nil
	}
	itemDTOs = make([]*services.RedEnvelopeItemDTO, 0)
	for _, item := range items {
		itemDTOs = append(itemDTOs, item.ToDTO())
	}
	return itemDTOs
}
