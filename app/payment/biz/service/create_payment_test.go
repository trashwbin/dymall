package service

import (
	"context"
	"testing"
	payment "github.com/trashwbin/dymall/rpc_gen/kitex_gen/payment"
)

func TestCreatePayment_Run(t *testing.T) {
	ctx := context.Background()
	s := NewCreatePaymentService(ctx)
	// init req and assert value

	req := &payment.CreatePaymentReq{}
	resp, err := s.Run(req)
	t.Logf("err: %v", err)
	t.Logf("resp: %v", resp)

	// todo: edit your unit test

}
