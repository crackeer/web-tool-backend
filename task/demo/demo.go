package demo

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"time"
	"web-tool-backend/container"
	"web-tool-backend/util"
)

var (
	//go:embed schema.json
	schemaJSON []byte
)

// DemoTask 实现了Tool接口的演示任务
type DemoTask struct {
	input string
}

// NewDemoTask 创建DemoTask实例
func NewDemoTask() *DemoTask {
	return &DemoTask{}
}

// GetDescription 返回任务的描述信息
func (d *DemoTask) GetDescription() *container.Description {
	return &container.Description{
		Title:     "DemoTask",
		InputForm: util.ParseJSON(string(schemaJSON)),
	}
}

// RecvInput 接收任务输入
func (d *DemoTask) RecvInput(bytes []byte) error {
	var input struct {
		Input string `json:"input"`
	}
	if err := json.Unmarshal(bytes, &input); err != nil {
		return err
	}
	d.input = input.Input
	return nil
}

// Run 执行任务，通过sendFunc发送结果
func (d *DemoTask) Run(sendFunc func(string)) error {
	// 发送任务开始消息
	sendFunc("Demo task started")

	// 如果有输入，处理输入
	if d.input != "" {
		sendFunc(fmt.Sprintf("Processing input: %s", d.input))
	}

	// 模拟处理步骤
	for i := 0; i < 5; i++ {
		sendFunc(fmt.Sprintf("Processing step %d/5", i+1))
		// 模拟处理延迟
		time.Sleep(time.Second)
	}

	// 发送任务完成消息
	sendFunc(fmt.Sprintf("Demo task completed successfully with input: %s", d.input))
	return nil
}
