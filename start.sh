#!/bin/bash

# 设置crd目录和输出文件路径
crdDir="./crd"
injectFile="./container/inject.go"

# 创建并清空输出文件
> "$injectFile"

# 写入包声明
cat << EOF > "$injectFile"
package container

import (
EOF

# 收集所有需要导入的包和工具注册语句
for file in "$crdDir"/*.json; do
    if [ -f "$file" ]; then
        # 获取文件名（不含扩展名）作为工具名
        name=$(basename "$file" .json)
        
        # 使用jq从json文件中提取包名
        package=$(jq -r '.package' "$file")
        
        # 生成导入语句，使用别名以避免冲突
        alias="${name}Task"
        echo "    $alias \"$package\"" >> "$injectFile"
    fi
done

# 关闭import块并添加init函数
cat << EOF >> "$injectFile"
)

// init 函数在包初始化时自动执行，注册所有工具
// 这些工具与crd目录下的json文件一一对应
func init() {
EOF

# 写入所有工具注册语句
for file in "$crdDir"/*.json; do
    if [ -f "$file" ]; then
        # 获取文件名（不含扩展名）作为工具名
        name=$(basename "$file" .json)
        alias="${name}Task"
        echo "    RegisterTool(\"$name\", $alias.Run)" >> "$injectFile"
    fi
done

# 关闭init函数
cat << EOF >> "$injectFile"
}
EOF

echo "已生成 $injectFile 文件"

case $1 in
    "build")
        go build -o web-tool-backend .
        ;;
    "run")
        go run main.go
        ;;
    *)
        echo "Usage: $0 {build|run}"
        exit 1
        ;;
esac  

