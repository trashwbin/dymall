package converter

import (
	"github.com/trashwbin/dymall/app/product/biz/model"
	product "github.com/trashwbin/dymall/rpc_gen/kitex_gen/product"
)

// ToProductProto 将领域模型转换为proto模型
func ToProductProto(p *model.Product) *product.Product {
	if p == nil {
		return nil
	}
	return &product.Product{
		Id:          p.ID,
		Name:        p.Name,
		Description: p.Description,
		Picture:     p.Picture,
		Price:       p.Price,
		Categories:  p.Categories,
	}
}

// ToProductProtos 批量将领域模型转换为proto模型
func ToProductProtos(products []*model.Product) []*product.Product {
	if products == nil {
		return nil
	}
	result := make([]*product.Product, len(products))
	for i, p := range products {
		result[i] = ToProductProto(p)
	}
	return result
}
