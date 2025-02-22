package service

import (
	"context"
	"testing"
	auth "github.com/trashwbin/dymall/rpc_gen/kitex_gen/auth"
)

func TestAddPolicy_Run(t *testing.T) {
	ctx := context.Background()
	s := NewAddPolicyService(ctx)
	// init req and assert value

	req := &auth.PolicyReq{}
	resp, err := s.Run(req)
	t.Logf("err: %v", err)
	t.Logf("resp: %v", resp)

	// todo: edit your unit test

}
