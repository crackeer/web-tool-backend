package json2csv

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"os"
	"web-tool-backend/container"
	"web-tool-backend/util"
)

var (
	//go:embed schema.json
	schemaJSON []byte
)

// Json2CsvTask 实现了Tool接口的演示任务
type Json2CsvTask struct {
	Input   string `json:"input"`
	TempDir string `json:"tempDir"`
}

// NewJson2CsvTask 创建Json2CsvTask实例
func NewJson2CsvTask(tempDir string) *Json2CsvTask {
	return &Json2CsvTask{
		TempDir: tempDir,
	}
}

// GetDescription 返回任务的描述信息
func (d *Json2CsvTask) GetDescription() *container.Description {
	return &container.Description{
		Title:     "JSON转CSV",
		InputForm: util.ParseJSON(string(schemaJSON)),
	}
}

// RecvInput 接收任务输入
func (d *Json2CsvTask) RecvInput(bytes []byte) error {
	if err := json.Unmarshal(bytes, d); err != nil {
		return err
	}
	return nil
}

// Run 执行任务，通过sendFunc发送结果
func (d *Json2CsvTask) Run(sendFunc func(string)) error {
	// 转换JSON为CSV
	sendFunc("正在转换JSON为CSV...")
	csvFilePath := fmt.Sprintf("%s/output.csv", d.TempDir)
	if err := os.MkdirAll(d.TempDir, 0755); err != nil {
		return fmt.Errorf("failed to create temp dir: %v", err)
	}
	if err := util.Json2CsvFile([]byte(d.Input), csvFilePath); err != nil {
		return fmt.Errorf("failed to convert JSON to CSV: %v", err)
	}

	sendFunc(fmt.Sprintf("CSV文件已生成: %s", csvFilePath))
	downloadURL := util.GenerationDownloadURL(csvFilePath)
	sendFunc("下载CSV：" + downloadURL)

	return nil
}
