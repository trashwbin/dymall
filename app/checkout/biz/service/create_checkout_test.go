package service

import (
	"context"
	"testing"
	checkout "github.com/trashwbin/dymall/rpc_gen/kitex_gen/checkout"
)

func TestCreateCheckout_Run(t *testing.T) {
	ctx := context.Background()
	s := NewCreateCheckoutService(ctx)
	// init req and assert value

	req := &checkout.CreateCheckoutReq{}
	resp, err := s.Run(req)
	t.Logf("err: %v", err)
	t.Logf("resp: %v", resp)

	// todo: edit your unit test

}
