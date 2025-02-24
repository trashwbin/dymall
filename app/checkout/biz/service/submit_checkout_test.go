package service

import (
	"context"
	"testing"
	checkout "github.com/trashwbin/dymall/rpc_gen/kitex_gen/checkout"
)

func TestSubmitCheckout_Run(t *testing.T) {
	ctx := context.Background()
	s := NewSubmitCheckoutService(ctx)
	// init req and assert value

	req := &checkout.SubmitCheckoutReq{}
	resp, err := s.Run(req)
	t.Logf("err: %v", err)
	t.Logf("resp: %v", resp)

	// todo: edit your unit test

}
