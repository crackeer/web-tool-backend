import { useState, useEffect, useRef } from 'react'
import { createFileRoute } from '@tanstack/react-router'
import { Terminal } from '@xterm/xterm'
import { FitAddon } from '@xterm/addon-fit'
import { WebLinksAddon } from '@xterm/addon-web-links'
import { ImageAddon } from 'xterm-addon-image';
import '@xterm/xterm/css/xterm.css'
import { Typography, Button, Card, Space, message, Radio, Form } from 'antd'
import 'antd/dist/reset.css'
import FormRender, { useForm } from 'form-render'

export const Route = createFileRoute('/')({
  component: Index,
})

function Index() {
  const [isConnected, setIsConnected] = useState(false)
  const eventSourceRef = useRef(null)
  const terminalRef = useRef(null)
  const terminalInstanceRef = useRef(null)
  const fitAddonRef = useRef(null)
  const form = useForm()
  const [tools, setTools] = useState([])
  const [selectedTool, setSelectedTool] = useState('')
  const [toolConfig, setToolConfig] = useState(null)
  const [formData, setFormData] = useState({})
  const location = window.location

  // 获取工具列表
  useEffect(() => {
    fetch('/api/crd')
      .then(response => response.json())
      .then(data => {
        setTools(data)
        if (data.length > 0) {
          // 尝试从localStorage读取最后选择的任务类型
          const lastSelectedTool = localStorage.getItem('lastSelectedTool')
          // 如果有存储的任务类型且在工具列表中存在，则使用它；否则使用第一个工具
          const defaultTool = lastSelectedTool && data.some(tool => tool.name === lastSelectedTool) 
            ? lastSelectedTool 
            : data[data.length - 1].name // 使用最后一个工具作为默认值
          
          setSelectedTool(defaultTool)
          const selectedToolConfig = data.find(tool => tool.name === defaultTool)
          if (selectedToolConfig) {
            setToolConfig(selectedToolConfig.form)
          }
        }
      })
      .catch(error => {
        console.error('Failed to fetch tools:', error)
        message.error('获取工具列表失败')
      })
  }, [])

  // 处理重跑任务参数
  useEffect(() => {
    // 使用URLSearchParams解析查询参数
    const params = new URLSearchParams(location.search)
    const rerunTaskId = params.get('rerun_task_id')
    if (rerunTaskId) {
      // 获取任务详情
      fetch(`/api/task/detail?task_id=${rerunTaskId}`)
        .then(response => response.json())
        .then(task => {
          if (task) {
            // 设置选中的工具类型
            setSelectedTool(task.task_type)
            // 查找对应的工具配置
            const tool = tools.find(t => t.name === task.task_type)
            if (tool) {
              setToolConfig(tool.form)
            }
            // 解析并设置表单数据
            try {
              const parsedInput = JSON.parse(task.input)
              console.log(parsedInput, task.input, form)

              Object.keys(parsedInput).forEach(key => {
                form.setFieldValue(key, parsedInput[key])
              })
            } catch (error) {
              console.error('解析任务输入参数失败:', error)
              message.error('解析任务参数失败')
            }
          }
        })
        .catch(error => {
          console.error('获取任务详情失败:', error)
          message.error('获取任务详情失败')
        })
    }
  }, [location.search, tools])

  // 处理工具选择变化
  const handleToolChange = (e) => {
    const value = e.target.value;
    const tool = tools.find(t => t.name === value)
    setSelectedTool(value)
    setToolConfig(tool.form)
    setFormData({})
    // 将选择的工具保存到localStorage
    localStorage.setItem('lastSelectedTool', value)
    // 重置表单 - 不需要调用resetFields，通过设置formData为空对象即可重置
  }

  // 使用task_id运行任务
  const startSSEWithInputID = (taskID) => {
    if (!selectedTool) {
      message.error('请选择一个工具')
      return
    }

    // 关闭现有的连接
    if (eventSourceRef.current) {
      eventSourceRef.current.close()
    }

    // 初始化终端
    if (terminalRef.current && !terminalInstanceRef.current) {
      const terminal = new Terminal({
        theme: {
          background: '#1e1e1e',
          foreground: '#d4d4d4'
        },
        fontSize: 14,
        fontFamily: 'Consolas, "Courier New", monospace',
        cursorBlink: true
      })
      const fitAddon = new FitAddon()
      const webLinksAddon = new WebLinksAddon()
      const imageAddon = new ImageAddon();
      terminal.loadAddon(fitAddon)
      terminal.loadAddon(webLinksAddon)
      terminal.loadAddon(imageAddon)
      terminal.open(terminalRef.current)
      fitAddon.fit()

      terminalInstanceRef.current = terminal
      fitAddonRef.current = fitAddon
    }

    // 创建新的SSE连接，传递task_id
    const eventSource = new EventSource(`/api/run?task_id=${taskID}`)
    eventSourceRef.current = eventSource
    setIsConnected(true)
    terminalInstanceRef.current.clear()

    if (terminalInstanceRef.current) {
      console.log('SSE connection established')
    }

    // 处理message事件
    eventSource.onmessage = (event) => {
      console.log('Received message:', event.data)
      if (terminalInstanceRef.current) {
        let msg = event.data.replace(/{WINDOW_HOSTNAME}/g, window.location.host)
        terminalInstanceRef.current.writeln(msg)
      }
    }

    // 处理close事件
    eventSource.addEventListener('close', (event) => {
      if (terminalInstanceRef.current) {
        console.log('SSE connection closed')
      }
      setIsConnected(false)
      eventSource.close()
    })

    // 处理error事件
    eventSource.onerror = (error) => {
      console.error('SSE error:', error, eventSource)
      if (terminalInstanceRef.current) {
        terminalInstanceRef.current.writeln(`\x1b[31mSSE error: ${error.message}\x1b[0m`)
      }
      setIsConnected(false)
      eventSource.close()
    }
  }

  // 保留原有的startSSE函数，用于向后兼容
  const startSSE = () => {
    if (!selectedTool) {
      message.error('请选择一个工具')
      return
    }

    form.validateFields().then((formData) => {
      // 调用task接口创建任务
      fetch('/api/task/create', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json'
        },
        body: JSON.stringify({
          task_type: selectedTool,
          input: JSON.stringify(formData)
        })
      })
        .then(response => response.json())
        .then(data => {
          const taskID = data.task_id
          console.log('任务已创建，task_id:', taskID)
          // 使用task_id运行任务
          startSSEWithInputID(taskID)
        })
        .catch(error => {
          console.error('创建任务失败:', error)
          message.error('创建任务失败，请重试')
        })
    }).catch(() => {
      // 表单验证失败，提示用户
      message.error('请填写完整的表单')
    })
  }

  const stopSSE = () => {
    if (eventSourceRef.current) {
      eventSourceRef.current.close()
      eventSourceRef.current = null
      setIsConnected(false)
      if (terminalInstanceRef.current) {
        terminalInstanceRef.current.writeln('')
        terminalInstanceRef.current.writeln('SSE connection manually closed')
      }
    }
  }

  // 处理窗口大小变化
  useEffect(() => {
    const handleResize = () => {
      if (terminalInstanceRef.current && fitAddonRef.current) {
        fitAddonRef.current.fit()
      }
    }

    window.addEventListener('resize', handleResize)
    return () => {
      window.removeEventListener('resize', handleResize)
    }
  }, [])

  // 组件卸载时清理资源
  useEffect(() => {
    return () => {
      if (eventSourceRef.current) {
        eventSourceRef.current.close()
      }
      if (terminalInstanceRef.current) {
        terminalInstanceRef.current.dispose()
      }
    }
  }, [])

  const { Title, Text } = Typography

  return (
    <div>
      <Card title={<Text strong>选择任务执行</Text>} style={{ marginBottom: '24px' }}>
        <Space direction="vertical" size="large" style={{ width: '100%' }}>
          <div style={{ marginTop: '8px' }}>
            <Typography.Text strong>任务类型:</Typography.Text>
            <Radio.Group 
              value={selectedTool} 
              onChange={handleToolChange}
              style={{ marginLeft: '16px', display: 'flex', flexWrap: 'wrap', gap: '16px' }}
            >
              {tools.map(tool => (
                <Radio key={tool.name} value={tool.name}>
                  {tool.title}
                </Radio>
              ))}
            </Radio.Group>
          </div>
          {toolConfig && (
            <div>
              <FormRender
                schema={toolConfig}
                form={form}
                onFinish={() => {
                  form.validate().then((values) => {
                    console.log(values)
                  })
                }}
                footer={false}
              />
            </div>
          )}

          <Space>
            <Button
              type="primary"
              onClick={startSSE}
              disabled={isConnected}
              loading={isConnected}
            >
              执行
            </Button>
            <Button
              danger
              onClick={stopSSE}
              disabled={!isConnected}
            >
              停止
            </Button>
          </Space>
        </Space>
      </Card>

      <Card title={<Text strong>任务输出</Text>} bordered={false}>
        <div
          ref={terminalRef}
          style={{ width: '100%', height: '600px', borderRadius: '1px' }}
        ></div>
      </Card>
    </div>
  )
}
