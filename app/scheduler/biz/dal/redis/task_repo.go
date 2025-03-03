package redis

import (
	"context"

	"github.com/trashwbin/dymall/app/scheduler/biz/model"
)

type TaskRepo struct{}

func NewTaskRepo() *TaskRepo {
	return &TaskRepo{}
}

// GetTask 获取任务缓存
func (r *TaskRepo) GetTask(ctx context.Context, taskID string) (*model.Task, error) {
	cache := &TaskCache{ID: taskID}
	data, err := RedisClient.Get(ctx, cache.GetKey()).Bytes()
	if err != nil {
		return nil, err
	}

	if err := cache.Unmarshal(data); err != nil {
		return nil, err
	}

	return cache.ToModel(), nil
}

// SetTask 设置任务缓存
func (r *TaskRepo) SetTask(ctx context.Context, task *model.Task) error {
	cache := &TaskCache{}
	cache.FromModel(task)

	data, err := cache.Marshal()
	if err != nil {
		return err
	}

	return RedisClient.Set(ctx, cache.GetKey(), data, TaskExpiration).Err()
}

// DeleteTask 删除任务缓存
func (r *TaskRepo) DeleteTask(ctx context.Context, taskID string) error {
	cache := &TaskCache{ID: taskID}
	return RedisClient.Del(ctx, cache.GetKey()).Err()
}

// AcquireLock 获取任务锁
func (r *TaskRepo) AcquireLock(ctx context.Context, taskID string) (bool, error) {
	cache := &TaskCache{ID: taskID}
	return RedisClient.SetNX(ctx, cache.GetLockKey(), 1, TaskLockExpiration).Result()
}

// ReleaseLock 释放任务锁
func (r *TaskRepo) ReleaseLock(ctx context.Context, taskID string) error {
	cache := &TaskCache{ID: taskID}
	return RedisClient.Del(ctx, cache.GetLockKey()).Err()
}

// ExtendLock 延长任务锁过期时间
func (r *TaskRepo) ExtendLock(ctx context.Context, taskID string) error {
	cache := &TaskCache{ID: taskID}
	return RedisClient.Expire(ctx, cache.GetLockKey(), TaskLockExpiration).Err()
}
