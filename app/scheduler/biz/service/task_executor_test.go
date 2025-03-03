package service

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/trashwbin/dymall/app/scheduler/biz/dal/mysql"
	"github.com/trashwbin/dymall/app/scheduler/biz/model"
	scheduler "github.com/trashwbin/dymall/rpc_gen/kitex_gen/scheduler"
)

func TestTaskExecutor_ExecuteBatch(t *testing.T) {
	ctx := context.Background()
	executor := NewTaskExecutor()
	db := mysql.DB

	// 清理测试数据
	db.Exec("DELETE FROM scheduler_tasks WHERE target_id LIKE ?", "test_order_%")

	// 创建多个测试任务
	tasks := []struct {
		targetID  string
		executeAt time.Time
	}{
		{
			targetID:  "test_order_1",
			executeAt: time.Now().Add(-time.Second), // 已过期，应该被执行
		},
		{
			targetID:  "test_order_2",
			executeAt: time.Now().Add(time.Hour), // 未过期，不应该被执行
		},
	}

	for _, task := range tasks {
		s := NewCreateTaskService(ctx)
		resp, err := s.Run(&scheduler.CreateTaskReq{
			Type:      scheduler.TaskType_TASK_TYPE_ORDER_EXPIRATION,
			TargetId:  task.targetID,
			ExecuteAt: task.executeAt.Unix(),
			Metadata: map[string]string{
				"test": "true", // 标记为测试任务
			},
		})
		require.NoError(t, err)
		require.NotNil(t, resp)
	}

	// 执行批量任务
	err := executor.executeBatch(ctx)
	require.NoError(t, err)

	// 等待任务执行完成
	time.Sleep(time.Second * 2)

	// 验证任务状态
	for _, task := range tasks {
		// 查询任务
		var taskDO mysql.TaskDO
		err := db.Where("target_id = ?", task.targetID).First(&taskDO).Error
		require.NoError(t, err)

		if task.executeAt.Before(time.Now()) {
			// 已过期的任务应该被执行
			assert.Equal(t, int32(model.TaskStatusCompleted), taskDO.Status)
		} else {
			// 未过期的任务应该保持待执行状态
			assert.Equal(t, int32(model.TaskStatusPending), taskDO.Status)
		}
	}

	// 清理测试数据
	db.Exec("DELETE FROM scheduler_tasks WHERE target_id LIKE ?", "test_order_%")
}

func TestTaskExecutor_HandleOrderExpiration(t *testing.T) {
	ctx := context.Background()
	executor := NewTaskExecutor()
	db := mysql.DB

	tests := []struct {
		name           string
		targetID       string
		expectedStatus model.TaskStatus
	}{
		{
			name:           "处理待支付订单",
			targetID:       "test_order_1", // mock客户端会返回待支付状态
			expectedStatus: model.TaskStatusCompleted,
		},
		{
			name:           "处理已支付订单",
			targetID:       "test_order_2", // mock客户端会返回已支付状态
			expectedStatus: model.TaskStatusCompleted,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 清理测试数据
			db.Exec("DELETE FROM scheduler_tasks WHERE target_id = ?", tt.targetID)

			// 创建测试任务
			s := NewCreateTaskService(ctx)
			resp, err := s.Run(&scheduler.CreateTaskReq{
				Type:      scheduler.TaskType_TASK_TYPE_ORDER_EXPIRATION,
				TargetId:  tt.targetID,
				ExecuteAt: time.Now().Add(-time.Second).Unix(),
				Metadata: map[string]string{
					"test": "true",
				},
			})
			require.NoError(t, err)
			require.NotNil(t, resp)

			// 获取任务模型
			task, err := executor.mysqlRepo.GetTask(resp.Task.TaskId)
			require.NoError(t, err)

			// 标记任务开始执行
			task.MarkRunning()
			err = executor.mysqlRepo.UpdateTask(task)
			require.NoError(t, err)

			// 执行订单过期处理
			err = executor.handleOrderExpiration(ctx, task)
			require.NoError(t, err)

			// 标记任务完成
			task.MarkCompleted()
			err = executor.mysqlRepo.UpdateTask(task)
			require.NoError(t, err)

			// 验证任务状态
			updatedTask, err := executor.mysqlRepo.GetTask(resp.Task.TaskId)
			require.NoError(t, err)
			assert.Equal(t, tt.expectedStatus, updatedTask.Status)

			// 清理测试数据
			db.Exec("DELETE FROM scheduler_tasks WHERE target_id = ?", tt.targetID)
		})
	}
}

func TestTaskExecutor_Start(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	executor := NewTaskExecutor()
	db := mysql.DB

	// 清理测试数据
	db.Exec("DELETE FROM scheduler_tasks WHERE target_id = ?", "test_order_1")

	// 创建一个测试任务
	s := NewCreateTaskService(ctx)
	resp, err := s.Run(&scheduler.CreateTaskReq{
		Type:      scheduler.TaskType_TASK_TYPE_ORDER_EXPIRATION,
		TargetId:  "test_order_1",
		ExecuteAt: time.Now().Add(-time.Second).Unix(),
		Metadata: map[string]string{
			"test": "true",
		},
	})
	require.NoError(t, err)
	require.NotNil(t, resp)

	// 启动执行器
	go executor.Start(ctx)

	// 等待任务执行完成
	time.Sleep(time.Second * 2)

	// 验证任务状态
	var taskDO mysql.TaskDO
	err = db.Where("target_id = ?", "test_order_1").First(&taskDO).Error
	require.NoError(t, err)
	assert.Equal(t, int32(model.TaskStatusCompleted), taskDO.Status)

	// 清理测试数据
	db.Exec("DELETE FROM scheduler_tasks WHERE target_id = ?", "test_order_1")
}
