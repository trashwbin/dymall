package service

import (
	"context"
	product "github.com/trashwbin/dymall/rpc_gen/kitex_gen/product"
)

type UpdateProductService struct {
	ctx context.Context
} // NewUpdateProductService new UpdateProductService
func NewUpdateProductService(ctx context.Context) *UpdateProductService {
	return &UpdateProductService{ctx: ctx}
}

// Run create note info
func (s *UpdateProductService) Run(req *product.UpdateProductReq) (resp *product.UpdateProductResp, err error) {
	// Finish your business logic.

	return
}
