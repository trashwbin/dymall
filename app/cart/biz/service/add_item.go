package service

import (
	"context"
	"fmt"
	cart "github.com/trashwbin/dymall/rpc_gen/kitex_gen/cart"
)

type AddItemService struct {
	ctx context.Context
} // NewAddItemService new AddItemService
func NewAddItemService(ctx context.Context) *AddItemService {
	return &AddItemService{ctx: ctx}
}

// Run create note info
func (s *AddItemService) Run(req *cart.AddItemReq) (resp *cart.AddItemResp, err error) {
	// Finish your business logic.
	fmt.Println("AddItemService")
	return &cart.AddItemResp{}, nil
}
