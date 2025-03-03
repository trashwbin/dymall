package service

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	scheduler "github.com/trashwbin/dymall/rpc_gen/kitex_gen/scheduler"
)

func TestGetTask_Run(t *testing.T) {
	// 先创建一个测试任务
	taskId := createTestTask(t)

	tests := []struct {
		name    string
		req     *scheduler.GetTaskReq
		wantErr bool
	}{
		{
			name: "获取存在的任务",
			req: &scheduler.GetTaskReq{
				TaskId: taskId,
			},
			wantErr: false,
		},
		{
			name: "获取不存在的任务",
			req: &scheduler.GetTaskReq{
				TaskId: "non_existent_task",
			},
			wantErr: true,
		},
		{
			name: "使用空任务ID",
			req: &scheduler.GetTaskReq{
				TaskId: "",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			s := NewGetTaskService(ctx)
			resp, err := s.Run(tt.req)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.NotNil(t, resp)
			assert.NotNil(t, resp.Task)
			assert.Equal(t, tt.req.TaskId, resp.Task.TaskId)
			assert.NotEmpty(t, resp.Task.Type)
			assert.NotEmpty(t, resp.Task.TargetId)
			assert.NotEmpty(t, resp.Task.ExecuteAt)
			assert.NotEmpty(t, resp.Task.Status)
		})
	}
}
