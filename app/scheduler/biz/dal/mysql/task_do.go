package mysql

import (
	"encoding/json"
	"time"

	"github.com/trashwbin/dymall/app/scheduler/biz/model"
	"gorm.io/gorm"
)

// TaskDO 任务数据对象
type TaskDO struct {
	ID           string         `gorm:"primaryKey;type:varchar(36)"`
	Type         int32          `gorm:"type:int;not null;comment:任务类型"`
	Status       int32          `gorm:"type:int;not null;default:1;comment:任务状态"`
	TargetID     string         `gorm:"type:varchar(32);index:idx_target_id;not null;comment:目标ID"`
	ExecuteAt    time.Time      `gorm:"index:idx_execute_at;not null;comment:执行时间"`
	CreatedAt    time.Time      `gorm:"not null;comment:创建时间"`
	UpdatedAt    time.Time      `gorm:"not null;comment:更新时间"`
	ErrorMessage string         `gorm:"type:text;comment:执行失败原因"`
	Metadata     string         `gorm:"type:text;comment:额外元数据"`
	DeletedAt    gorm.DeletedAt `gorm:"index;comment:删除时间"`
}

// TableName 设置表名
func (TaskDO) TableName() string {
	return "scheduler_tasks"
}

// ToModel 转换为领域模型
func (t *TaskDO) ToModel() (*model.Task, error) {
	metadata := make(map[string]string)
	if t.Metadata != "" {
		if err := json.Unmarshal([]byte(t.Metadata), &metadata); err != nil {
			return nil, err
		}
	}

	return &model.Task{
		ID:           t.ID,
		Type:         model.TaskType(t.Type),
		Status:       model.TaskStatus(t.Status),
		TargetID:     t.TargetID,
		ExecuteAt:    t.ExecuteAt,
		CreatedAt:    t.CreatedAt,
		UpdatedAt:    t.UpdatedAt,
		ErrorMessage: t.ErrorMessage,
		Metadata:     metadata,
	}, nil
}

// FromModel 从领域模型转换
func (t *TaskDO) FromModel(m *model.Task) error {
	metadata, err := json.Marshal(m.Metadata)
	if err != nil {
		return err
	}

	t.ID = m.ID
	t.Type = int32(m.Type)
	t.Status = int32(m.Status)
	t.TargetID = m.TargetID
	t.ExecuteAt = m.ExecuteAt
	t.CreatedAt = m.CreatedAt
	t.UpdatedAt = m.UpdatedAt
	t.ErrorMessage = m.ErrorMessage
	t.Metadata = string(metadata)
	return nil
}
