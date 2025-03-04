package mysql

import (
	"time"

	"github.com/trashwbin/dymall/app/scheduler/biz/model"
	"gorm.io/gorm"
)

type TaskRepo struct {
	db *gorm.DB
}

func NewTaskRepo() *TaskRepo {
	return &TaskRepo{db: DB}
}

// CreateTask 创建任务
func (r *TaskRepo) CreateTask(task *model.Task) error {
	taskDO := &TaskDO{}
	if err := taskDO.FromModel(task); err != nil {
		return err
	}
	return r.db.Create(taskDO).Error
}

// GetTask 获取任务
func (r *TaskRepo) GetTask(taskID string) (*model.Task, error) {
	var taskDO TaskDO
	if err := r.db.First(&taskDO, "id = ?", taskID).Error; err != nil {
		return nil, err
	}
	return taskDO.ToModel()
}

// UpdateTask 更新任务
func (r *TaskRepo) UpdateTask(task *model.Task) error {
	taskDO := &TaskDO{}
	if err := taskDO.FromModel(task); err != nil {
		return err
	}
	return r.db.Save(taskDO).Error
}

// DeleteTask 删除任务
func (r *TaskRepo) DeleteTask(taskID string) error {
	return r.db.Delete(&TaskDO{}, "id = ?", taskID).Error
}

// GetPendingTasks 获取待执行的任务列表
func (r *TaskRepo) GetPendingTasks(executeTime time.Time, limit int) ([]*model.Task, error) {
	var taskDOs []TaskDO
	err := r.db.Where("status = ? AND execute_at <= ?",
		model.TaskStatusPending, executeTime).
		Order("execute_at ASC").
		Limit(limit).
		Find(&taskDOs).Error
	if err != nil {
		return nil, err
	}

	tasks := make([]*model.Task, 0, len(taskDOs))
	for _, taskDO := range taskDOs {
		task, err := taskDO.ToModel()
		if err != nil {
			continue
		}
		tasks = append(tasks, task)
	}
	return tasks, nil
}

// Transaction 事务处理
func (r *TaskRepo) Transaction(fn func(txRepo *TaskRepo) error) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		txRepo := &TaskRepo{db: tx}
		return fn(txRepo)
	})
}
