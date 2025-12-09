package container

import (
	"fmt"
	"time"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

var db *gorm.DB

// Task 定义Task结构体，对应task表
type Task struct {
	ID         uint      `gorm:"autoIncrement;column:id" json:"id"`
	TaskType   string    `gorm:"column:task_type" json:"task_type"`
	Input      string    `gorm:"column:input" json:"input"`
	CreateTime time.Time `gorm:"column:create_time" json:"create_time"`
}

// 设置表名
func (Task) TableName() string {
	return "task"
}

func InitDB() error {
	var err error
	db, err = gorm.Open(sqlite.Open(GetConfig().SQLLiteDB), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("failed to connect database: %w", err)
	}

	// 自动迁移创建表
	db.AutoMigrate(&Task{})
	return nil
}

// CreateInput 创建任务并保存到数据库
func CreateInput(taskType string, input string) (string, error) {
	// 生成当前时间
	now := time.Now()

	// 创建Task实例，ID由SQLite自动生成
	task := Task{
		TaskType:   taskType,
		Input:      input,
		CreateTime: now,
	}

	// 保存到数据库
	if result := db.Create(&task); result.Error != nil {
		return "", result.Error
	}

	// 将uint类型的ID转换为string返回
	return fmt.Sprintf("%d", task.ID), nil
}

// GetTask 根据ID从数据库获取任务
func GetTask(inputID string) *Task {
	var task Task
	result := db.Where("id = ?", inputID).First(&task)
	if result.Error != nil {
		return nil
	}
	return &task
}

// UpdateTaskType 更新任务类型
func UpdateTaskType(inputID string, taskType string) error {
	result := db.Model(&Task{}).Where("id = ?", inputID).Update("task_type", taskType)
	return result.Error
}

// ListTasks 查询任务列表，支持分页和任务类型过滤
func ListTasks(taskType string, page int, pageSize int) ([]*Task, int64, error) {
	var tasks []*Task
	var total int64

	// 构建查询
	query := db.Model(&Task{})

	// 如果提供了任务类型，则添加过滤条件
	if taskType != "" {
		query = query.Where("task_type = ?", taskType)
	}

	// 获取总记录数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 计算偏移量
	offset := (page - 1) * pageSize

	// 执行分页查询，按创建时间倒序排列
	if err := query.Offset(offset).Limit(pageSize).Order("create_time DESC").Find(&tasks).Error; err != nil {
		return nil, 0, err
	}

	return tasks, total, nil
}

// DeleteTask 根据ID删除任务
func DeleteTask(taskID uint) error {
	result := db.Delete(&Task{}, taskID)
	return result.Error
}
