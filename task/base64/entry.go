package base64

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"web-tool-backend/container"
	"web-tool-backend/util"
)

var (
	//go:embed schema.json
	schemaJSON []byte
)

// Base64Task 实现了Tool接口的演示任务
type Base64Task struct {
	Input string `json:"input"`
	Type  string `json:"type"`
}

// NewBase64Task 创建Base64Task实例
func NewBase64Task() *Base64Task {
	return &Base64Task{}
}

// GetDescription 返回任务的描述信息
func (d *Base64Task) GetDescription() *container.Description {
	return &container.Description{
		Title:     "Base64编码/解码",
		InputForm: util.ParseJSON(string(schemaJSON)),
	}
}

// RecvInput 接收任务输入
func (d *Base64Task) RecvInput(bytes []byte) error {
	if err := json.Unmarshal(bytes, d); err != nil {
		return err
	}
	return nil
}

// Run 执行任务，通过sendFunc发送结果
func (d *Base64Task) Run(sendFunc func(string)) error {
	// 如果有输入，处理输入
	sendFunc(fmt.Sprintf("输入: %s", d.Input))
	if d.Type == "encode" {
		sendFunc(fmt.Sprintf("编码结果: %s", util.Base64Encode(d.Input)))
	} else {
		sendFunc(fmt.Sprintf("解码结果: %s", util.Base64Decode(d.Input)))
	}
	return nil
}
