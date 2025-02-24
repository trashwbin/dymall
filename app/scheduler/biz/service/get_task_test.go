package service

import (
	"context"
	"testing"
	scheduler "github.com/trashwbin/dymall/rpc_gen/kitex_gen/scheduler"
)

func TestGetTask_Run(t *testing.T) {
	ctx := context.Background()
	s := NewGetTaskService(ctx)
	// init req and assert value

	req := &scheduler.GetTaskReq{}
	resp, err := s.Run(req)
	t.Logf("err: %v", err)
	t.Logf("resp: %v", resp)

	// todo: edit your unit test

}
