package redis

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/trashwbin/dymall/app/scheduler/biz/model"
)

const (
	// TaskKeyPrefix 任务缓存key前缀
	TaskKeyPrefix = "scheduler:task:"
	// TaskLockKeyPrefix 任务锁key前缀
	TaskLockKeyPrefix = "scheduler:lock:task:"
	// TaskExpiration 任务缓存过期时间
	TaskExpiration = 24 * time.Hour
	// TaskLockExpiration 任务锁过期时间
	TaskLockExpiration = 5 * time.Minute
)

// TaskCache Redis任务缓存数据对象
type TaskCache struct {
	ID           string            `json:"id"`
	Type         int32             `json:"type"`
	Status       int32             `json:"status"`
	TargetID     string            `json:"target_id"`
	ExecuteAt    time.Time         `json:"execute_at"`
	CreatedAt    time.Time         `json:"created_at"`
	UpdatedAt    time.Time         `json:"updated_at"`
	ErrorMessage string            `json:"error_message"`
	Metadata     map[string]string `json:"metadata"`
}

// ToModel 转换为领域模型
func (c *TaskCache) ToModel() *model.Task {
	return &model.Task{
		ID:           c.ID,
		Type:         model.TaskType(c.Type),
		Status:       model.TaskStatus(c.Status),
		TargetID:     c.TargetID,
		ExecuteAt:    c.ExecuteAt,
		CreatedAt:    c.CreatedAt,
		UpdatedAt:    c.UpdatedAt,
		ErrorMessage: c.ErrorMessage,
		Metadata:     c.Metadata,
	}
}

// FromModel 从领域模型转换
func (c *TaskCache) FromModel(m *model.Task) {
	c.ID = m.ID
	c.Type = int32(m.Type)
	c.Status = int32(m.Status)
	c.TargetID = m.TargetID
	c.ExecuteAt = m.ExecuteAt
	c.CreatedAt = m.CreatedAt
	c.UpdatedAt = m.UpdatedAt
	c.ErrorMessage = m.ErrorMessage
	c.Metadata = m.Metadata
}

// GetKey 获取任务缓存key
func (c *TaskCache) GetKey() string {
	return fmt.Sprintf("%s%s", TaskKeyPrefix, c.ID)
}

// GetLockKey 获取任务锁key
func (c *TaskCache) GetLockKey() string {
	return fmt.Sprintf("%s%s", TaskLockKeyPrefix, c.ID)
}

// Marshal 序列化
func (c *TaskCache) Marshal() ([]byte, error) {
	return json.Marshal(c)
}

// Unmarshal 反序列化
func (c *TaskCache) Unmarshal(data []byte) error {
	return json.Unmarshal(data, c)
}
