package service

import (
	"context"
	"testing"
	order "github.com/trashwbin/dymall/rpc_gen/kitex_gen/order"
)

func TestGetOrder_Run(t *testing.T) {
	ctx := context.Background()
	s := NewGetOrderService(ctx)
	// init req and assert value

	req := &order.GetOrderReq{}
	resp, err := s.Run(req)
	t.Logf("err: %v", err)
	t.Logf("resp: %v", resp)

	// todo: edit your unit test

}
