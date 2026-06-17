package store

import (
	"github.com/yc446833448/VisuTask/internal/model"
	"gorm.io/gorm"
)

// ─── Script Repository ───

func (db *DB) ListScripts() ([]model.Script, error) {
	var scripts []model.Script
	err := db.Order("created_at DESC").Find(&scripts).Error
	return scripts, err
}

func (db *DB) GetScript(id string) (*model.Script, error) {
	var script model.Script
	err := db.First(&script, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &script, nil
}

func (db *DB) CreateScript(script *model.Script) error {
	return db.Create(script).Error
}

func (db *DB) UpdateScript(script *model.Script) error {
	return db.Save(script).Error
}

func (db *DB) DeleteScript(id string) error {
	// Check if any tasks reference this script
	var count int64
	db.Model(&model.Task{}).Where("script_id = ?", id).Count(&count)
	if count > 0 {
		return ErrScriptInUse
	}
	return db.Delete(&model.Script{}, "id = ?", id).Error
}

func (db *DB) CountTasksByScript(scriptID string) (int64, error) {
	var count int64
	err := db.Model(&model.Task{}).Where("script_id = ?", scriptID).Count(&count).Error
	return count, err
}

// ─── Task Repository ───

func (db *DB) ListTasks() ([]model.Task, error) {
	var tasks []model.Task
	err := db.Order("created_at DESC").Find(&tasks).Error
	return tasks, err
}

func (db *DB) GetTask(id string) (*model.Task, error) {
	var task model.Task
	err := db.First(&task, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &task, nil
}

func (db *DB) CreateTask(task *model.Task) error {
	return db.Create(task).Error
}

func (db *DB) UpdateTask(task *model.Task) error {
	return db.Save(task).Error
}

func (db *DB) UpdateTaskStatus(id string, status model.TaskStatus) error {
	return db.Model(&model.Task{}).Where("id = ?", id).Update("status", status).Error
}

func (db *DB) DeleteTask(id string) error {
	return db.Delete(&model.Task{}, "id = ?", id).Error
}

func (db *DB) ListRunningTasks() ([]model.Task, error) {
	var tasks []model.Task
	err := db.Where("status = ?", model.TaskStatusRunning).Find(&tasks).Error
	return tasks, err
}

// ─── Execution Repository ───

func (db *DB) ListExecutions(taskID string, limit int) ([]model.Execution, error) {
	var executions []model.Execution
	q := db.Order("started_at DESC")
	if taskID != "" {
		q = q.Where("task_id = ?", taskID)
	}
	if limit > 0 {
		q = q.Limit(limit)
	}
	err := q.Find(&executions).Error
	return executions, err
}

func (db *DB) GetExecution(id string) (*model.Execution, error) {
	var execution model.Execution
	err := db.First(&execution, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &execution, nil
}

func (db *DB) CreateExecution(execution *model.Execution) error {
	return db.Create(execution).Error
}

func (db *DB) UpdateExecution(execution *model.Execution) error {
	return db.Save(execution).Error
}

// ─── User Repository ───

func (db *DB) GetUser() (*model.User, error) {
	var user model.User
	err := db.First(&user).Error
	if err == gorm.ErrRecordNotFound {
		// Create default user
		user = model.User{
			ID:            "default",
			VIPLevel:      0,
			MaxConcurrent: 3,
			Balance:       0,
		}
		if err := db.Create(&user).Error; err != nil {
			return nil, err
		}
		return &user, nil
	}
	return &user, err
}

func (db *DB) UpdateUser(user *model.User) error {
	return db.Save(user).Error
}

// ─── Errors ───

type StoreError string

func (e StoreError) Error() string { return string(e) }

const (
	ErrScriptInUse StoreError = "script is in use by tasks, delete tasks first"
	ErrNotFound    StoreError = "record not found"
)
