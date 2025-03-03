package service

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	scheduler "github.com/trashwbin/dymall/rpc_gen/kitex_gen/scheduler"
)

func TestCreateTask_Run(t *testing.T) {
	tests := []struct {
		name    string
		req     *scheduler.CreateTaskReq
		wantErr bool
	}{
		{
			name: "创建有效的订单过期任务",
			req: &scheduler.CreateTaskReq{
				Type:      scheduler.TaskType_TASK_TYPE_ORDER_EXPIRATION,
				TargetId:  "test_order_1",
				ExecuteAt: time.Now().Add(time.Hour).Unix(),
			},
			wantErr: false,
		},
		{
			name: "创建过期时间在过去的任务",
			req: &scheduler.CreateTaskReq{
				Type:      scheduler.TaskType_TASK_TYPE_ORDER_EXPIRATION,
				TargetId:  "test_order_1",
				ExecuteAt: time.Now().Add(-time.Hour).Unix(),
			},
			wantErr: true,
		},
		{
			name: "创建无效类型的任务",
			req: &scheduler.CreateTaskReq{
				Type:      scheduler.TaskType_TASK_TYPE_UNSPECIFIED,
				TargetId:  "test_order_1",
				ExecuteAt: time.Now().Add(time.Hour).Unix(),
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			s := NewCreateTaskService(ctx)
			resp, err := s.Run(tt.req)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.NotNil(t, resp)
			assert.NotEmpty(t, resp.Task.TaskId)
			assert.Equal(t, tt.req.Type, resp.Task.Type)
			assert.Equal(t, tt.req.TargetId, resp.Task.TargetId)
			assert.Equal(t, tt.req.ExecuteAt, resp.Task.ExecuteAt)
			assert.Equal(t, scheduler.TaskStatus_TASK_STATUS_PENDING, resp.Task.Status)
		})
	}
}
