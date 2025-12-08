package qrcode

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"strings"
	"web-tool-backend/container"
	"web-tool-backend/util"

	"github.com/skip2/go-qrcode"
)

var (
	//go:embed schema.json
	schemaJSON []byte
)

// QrcodeTask 实现了Tool接口的演示任务
type QrcodeTask struct {
	Input string `json:"input"`
}

// NewQrcodeTask 创建QrcodeTask实例
func NewQrcodeTask() *QrcodeTask {
	return &QrcodeTask{}
}

// GetDescription 返回任务的描述信息
func (d *QrcodeTask) GetDescription() *container.Description {
	return &container.Description{
		Title:     "二维码生成",
		InputForm: util.ParseJSON(string(schemaJSON)),
	}
}

// RecvInput 接收任务输入
func (d *QrcodeTask) RecvInput(bytes []byte) error {
	if err := json.Unmarshal(bytes, d); err != nil {
		return err
	}
	return nil
}

// Run 执行任务，通过sendFunc发送结果
func (d *QrcodeTask) Run(sendFunc func(string)) error {
	// 如果有输入，处理输入
	sendFunc(fmt.Sprintf("输入: %s", d.Input))
	qr, err := qrcode.New(d.Input, qrcode.High)
	if err != nil {
		return err
	}
	qrASCII := qr.ToSmallString(false)
	parts := strings.Split(qrASCII, "\n")
	for _, part := range parts {
		sendFunc(part)
	}
	return nil
}
