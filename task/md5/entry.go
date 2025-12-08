package md5

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

// Md5Task 实现了Tool接口的演示任务
type Md5Task struct {
	Input string `json:"input"`
}

// NewMd5Task 创建Md5Task实例
func NewMd5Task() *Md5Task {
	return &Md5Task{}
}

// GetDescription 返回任务的描述信息
func (d *Md5Task) GetDescription() *container.Description {
	return &container.Description{
		Title:     "MD5计算",
		InputForm: util.ParseJSON(string(schemaJSON)),
	}
}

// RecvInput 接收任务输入
func (d *Md5Task) RecvInput(bytes []byte) error {
	if err := json.Unmarshal(bytes, d); err != nil {
		return err
	}
	return nil
}

// Run 执行任务，通过sendFunc发送结果
func (d *Md5Task) Run(sendFunc func(string)) error {
	// 如果有输入，处理输入
	sendFunc(fmt.Sprintf("输入: %s", d.Input))
	// 计算MD5值
	md5Value := util.CalculateMD5(d.Input)
	sendFunc(fmt.Sprintf("MD5值: %s", md5Value))
	return nil
}
