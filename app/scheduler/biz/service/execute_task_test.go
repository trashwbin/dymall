package service

import (
	"context"
	"testing"
	scheduler "github.com/trashwbin/dymall/rpc_gen/kitex_gen/scheduler"
)

func TestExecuteTask_Run(t *testing.T) {
	ctx := context.Background()
	s := NewExecuteTaskService(ctx)
	// init req and assert value

	req := &scheduler.ExecuteTaskReq{}
	resp, err := s.Run(req)
	t.Logf("err: %v", err)
	t.Logf("resp: %v", resp)

	// todo: edit your unit test

}
