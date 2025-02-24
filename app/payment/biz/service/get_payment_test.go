package service

import (
	"context"
	"testing"
	payment "github.com/trashwbin/dymall/rpc_gen/kitex_gen/payment"
)

func TestGetPayment_Run(t *testing.T) {
	ctx := context.Background()
	s := NewGetPaymentService(ctx)
	// init req and assert value

	req := &payment.GetPaymentReq{}
	resp, err := s.Run(req)
	t.Logf("err: %v", err)
	t.Logf("resp: %v", resp)

	// todo: edit your unit test

}
