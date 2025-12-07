package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"web-tool-backend/container"
	"web-tool-backend/task/demo"
	"web-tool-backend/task/demo2"

	"github.com/gin-gonic/gin"

	"github.com/caarlos0/env/v11"
	_ "github.com/joho/godotenv/autoload"
)

type AppConfig struct {
	Port        string `env:"PORT" envDefault:"8080"`
	FrontendDir string `env:"FRONTEND_DIR" envDefault:"./frontend"`
}

var (
	cfg *AppConfig
)

func enableCORS(ctx *gin.Context) {
	ctx.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	ctx.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
	ctx.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	ctx.Writer.Header().Set("Access-Control-Expose-Headers", "Content-Length")
	ctx.Writer.Header().Set("Access-Control-Allow-Credentials", "true")

	if ctx.Request.Method == http.MethodOptions {
		ctx.AbortWithStatus(http.StatusNoContent)
		return
	}
	ctx.Next()
}

func main() {
	// 从环境变量中解析配置
	cfg = &AppConfig{}
	if err := env.Parse(cfg); err != nil {
		log.Fatalf("Failed to parse environment variables: %v", err)
	}

	// 注册DemoTask
	container.RegisterTool("demo", demo.NewDemoTask())
	container.RegisterTool("demo2", demo2.NewDemoTask())

	// 创建 Gin 实例
	router := gin.Default()
	router.Use(enableCORS)

	apiGroup := router.Group("/api")
	{
		apiGroup.GET("/task/run", RunTaskSSE)
		apiGroup.GET("/tools", GetTools)
		apiGroup.POST("/task/create", CreateTask)
		apiGroup.GET("/task/list", GetTasks)
		apiGroup.POST("/task/delete", DeleteTask)
		apiGroup.GET("/task/detail", GetTaskByID)
	}

	router.NoRoute(func(ctx *gin.Context) {
		fileServer := http.StripPrefix("", http.FileServer(http.Dir(cfg.FrontendDir)))
		fileServer.ServeHTTP(ctx.Writer, ctx.Request)
	})

	// 启动服务器
	router.Run(":" + cfg.Port)
}

// CreateTask 创建任务
func CreateTask(ctx *gin.Context) {
	var task container.Task
	if err := ctx.ShouldBindJSON(&task); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	inputID, err := container.CreateInput(task.TaskType, task.Input)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"task_id": inputID})
}

// GetTasks 获取任务列表
func GetTasks(ctx *gin.Context) {
	// 获取查询参数
	taskType := ctx.Query("task_type")

	// 手动解析整数参数
	page := 1
	pageSize := 10

	if pageStr := ctx.Query("page"); pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}

	if pageSizeStr := ctx.Query("page_size"); pageSizeStr != "" {
		if ps, err := strconv.Atoi(pageSizeStr); err == nil && ps > 0 && ps <= 100 {
			pageSize = ps
		}
	}

	// 查询任务列表
	tasks, total, err := container.ListTasks(taskType, page, pageSize)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 返回结果
	ctx.JSON(http.StatusOK, gin.H{
		"tasks":       tasks,
		"total":       total,
		"page":        page,
		"page_size":   pageSize,
		"total_pages": (total + int64(pageSize) - 1) / int64(pageSize),
	})
}

func GetTools(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, container.GetToolConfig())
}

// DeleteTask 删除任务
func DeleteTask(ctx *gin.Context) {
	// 获取路径参数中的任务ID
	idStr := ctx.Query("task_id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "无效的任务ID"})
		return
	}

	// 调用删除任务函数
	if err := container.DeleteTask(uint(id)); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "任务删除成功"})
}

// GetTaskByID 根据ID获取任务详情
func GetTaskByID(ctx *gin.Context) {
	// 获取路径参数中的任务ID
	idStr := ctx.Query("task_id")
	if idStr == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "任务ID不能为空"})
		return
	}

	// 调用获取任务函数
	task := container.GetTask(idStr)
	if task == nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "任务不存在"})
		return
	}

	ctx.JSON(http.StatusOK, task)
}

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

	printMessage("Task started")

	taskID := ctx.Query("task_id")
	task := container.GetTask(taskID)
	if task == nil {
		printMessage(fmt.Sprintf("task with id %s not found", taskID))
		return
	}
	fmt.Printf("task: %v\n", task)

	tool := container.GetTool(task.TaskType)
	if tool == nil {
		printMessage(fmt.Sprintf("tool %s not found", task.TaskType))
		return
	}

	var input []byte = []byte(task.Input)

	if err := tool.RecvInput(input); err != nil {
		printMessage(fmt.Sprintf("failed to recv input: %v", err))
		return
	}

	if err := tool.Run(printMessage); err != nil {
		printMessage(fmt.Sprintf("failed to run task: %v", err))
		return
	}
	printMessage("Task completed successfully")
	closeEvent := fmt.Sprintf("event: close\nretry: 0\ndata: %s\n\n", "Task completed successfully")
	ctx.SSEvent("close", closeEvent)
	ctx.Writer.Flush()
}
