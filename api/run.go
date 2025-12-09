package api

import (
	"fmt"
	"web-tool-backend/container"

	"github.com/gin-gonic/gin"
)

// GetTaskByID 根据ID获取任务详情

// RunTaskSSE 处理SSE请求
func RunTaskSSE(ctx *gin.Context) {
	// 设置SSE头部
	ctx.Header("Content-Type", "text/event-stream")
	ctx.Header("Cache-Control", "no-cache")
	ctx.Header("Connection", "keep-alive")
	ctx.Header("Transfer-Encoding", "chunked")

	var printMessage func(string) = func(msg string) {
		ctx.SSEvent("message", msg)
		ctx.Writer.Flush()
	}

	taskID := ctx.Query("task_id")
	task := container.GetTask(taskID)
	if task == nil {
		printMessage(fmt.Sprintf("task with id %s not found", taskID))
		return
	}
	toolFunc := container.GetTool(task.TaskType)
	if toolFunc == nil {
		printMessage(fmt.Sprintf("tool %s not found", task.TaskType))
		return
	}

	printMessage(fmt.Sprintf("Task started, task_id: %s", taskID))
	printMessage("")
	output, err := toolFunc(task.Input, printMessage)
	if err != nil {
		printMessage("")
		printMessage(fmt.Sprintf("failed to run task: %v", err))
		return
	}
	if len(output) > 0 {
		printMessage("output: " + output)
	}
	printMessage("")
	printMessage("Task completed successfully")
	closeEvent := fmt.Sprintf("event: close\nretry: 0\ndata: %s\n\n", "Task completed successfully")
	ctx.SSEvent("close", closeEvent)
	ctx.Writer.Flush()
}
