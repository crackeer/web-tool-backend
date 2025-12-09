package container

import (
    base64Task "web-tool-backend/task/base64"
    demoTask "web-tool-backend/task/demo"
    downloadworkTask "web-tool-backend/task/downloadwork"
    json2csvTask "web-tool-backend/task/json2csv"
    md5Task "web-tool-backend/task/md5"
    qrcodeTask "web-tool-backend/task/qrcode"
)

// init 函数在包初始化时自动执行，注册所有工具
// 这些工具与crd目录下的json文件一一对应
func init() {
    RegisterTool("base64", base64Task.Run)
    RegisterTool("demo", demoTask.Run)
    RegisterTool("downloadwork", downloadworkTask.Run)
    RegisterTool("json2csv", json2csvTask.Run)
    RegisterTool("md5", md5Task.Run)
    RegisterTool("qrcode", qrcodeTask.Run)
}
