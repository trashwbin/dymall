package service

import (
	"context"
	"testing"
	scheduler "github.com/trashwbin/dymall/rpc_gen/kitex_gen/scheduler"
)

func TestCancelTask_Run(t *testing.T) {
	ctx := context.Background()
	s := NewCancelTaskService(ctx)
	// init req and assert value

	req := &scheduler.CancelTaskReq{}
	resp, err := s.Run(req)
	t.Logf("err: %v", err)
	t.Logf("resp: %v", resp)

	// todo: edit your unit test

}
