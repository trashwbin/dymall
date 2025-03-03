package service

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/trashwbin/dymall/app/scheduler/biz/dal/mysql"
	scheduler "github.com/trashwbin/dymall/rpc_gen/kitex_gen/scheduler"
)

func TestExecuteTask_Run(t *testing.T) {
	ctx := context.Background()
	db := mysql.DB

	// 清理测试数据
	db.Exec("DELETE FROM scheduler_tasks WHERE target_id = ?", "test_order_1")

	// 创建一个可以立即执行的测试任务（设置为过去时间）
	s := NewCreateTaskService(ctx)
	createResp, err := s.Run(&scheduler.CreateTaskReq{
		Type:      scheduler.TaskType_TASK_TYPE_ORDER_EXPIRATION,
		TargetId:  "test_order_1",
		ExecuteAt: time.Now().Add(-time.Second).Unix(), // 设置为过去时间
		Metadata: map[string]string{
			"test": "true", // 标记为测试任务，跳过执行时间验证
		},
	})
	if !assert.NoError(t, err) {
		t.Fatalf("创建测试任务失败: %v", err)
		return
	}
	if !assert.NotNil(t, createResp) {
		t.Fatal("创建测试任务响应为空")
		return
	}
	if !assert.NotNil(t, createResp.Task) {
		t.Fatal("创建的任务为空")
		return
	}

	taskId := createResp.Task.TaskId

	tests := []struct {
		name    string
		req     *scheduler.ExecuteTaskReq
		wantErr bool
	}{
		{
			name: "执行存在的任务",
			req: &scheduler.ExecuteTaskReq{
				TaskId: taskId,
			},
			wantErr: false,
		},
		{
			name: "执行不存在的任务",
			req: &scheduler.ExecuteTaskReq{
				TaskId: "non_existent_task",
			},
			wantErr: true,
		},
		{
			name: "使用空任务ID",
			req: &scheduler.ExecuteTaskReq{
				TaskId: "",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewExecuteTaskService(ctx)
			resp, err := s.Run(tt.req)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			if !assert.NoError(t, err) {
				return
			}
			if !assert.NotNil(t, resp) {
				return
			}
			if !assert.NotNil(t, resp.Task) {
				return
			}

			// 检查任务状态为COMPLETED（因为我们是同步执行的）
			assert.Equal(t, scheduler.TaskStatus_TASK_STATUS_COMPLETED, resp.Task.Status)

			// 等待足够长的时间确保数据库更新完成
			time.Sleep(time.Second * 2)

			// 获取最新的任务状态
			updatedTask, err := s.mysqlRepo.GetTask(tt.req.TaskId)
			require.NoError(t, err)

			// 验证任务状态
			require.Equal(t, scheduler.TaskStatus_TASK_STATUS_COMPLETED, scheduler.TaskStatus(updatedTask.Status))
		})
	}

	// 清理测试数据
	db.Exec("DELETE FROM scheduler_tasks WHERE target_id = ?", "test_order_1")
}
