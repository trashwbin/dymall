package service

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	scheduler "github.com/trashwbin/dymall/rpc_gen/kitex_gen/scheduler"
)

func TestCancelTask_Run(t *testing.T) {
	// 先创建一个测试任务
	taskId := createTestTask(t)

	tests := []struct {
		name    string
		req     *scheduler.CancelTaskReq
		wantErr bool
	}{
		{
			name: "取消存在的任务",
			req: &scheduler.CancelTaskReq{
				TaskId: taskId,
			},
			wantErr: false,
		},
		{
			name: "取消不存在的任务",
			req: &scheduler.CancelTaskReq{
				TaskId: "non_existent_task",
			},
			wantErr: true,
		},
		{
			name: "使用空任务ID",
			req: &scheduler.CancelTaskReq{
				TaskId: "",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			s := NewCancelTaskService(ctx)
			resp, err := s.Run(tt.req)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.NotNil(t, resp)

			// 验证任务确实被取消了
			getResp, err := NewGetTaskService(ctx).Run(&scheduler.GetTaskReq{TaskId: tt.req.TaskId})
			assert.NoError(t, err)
			assert.Equal(t, scheduler.TaskStatus_TASK_STATUS_CANCELLED, getResp.Task.Status)
		})
	}
}
