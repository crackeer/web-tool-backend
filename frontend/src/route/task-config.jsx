import { createFileRoute } from '@tanstack/react-router'
import { useState, useEffect } from 'react'
import { Table, Form, Input, Button, Modal, message, Card, Space } from 'antd'
import { PlusOutlined, EditOutlined, DeleteOutlined, SaveOutlined, CloseOutlined } from '@ant-design/icons'

const { TextArea } = Input

const TaskConfigPage = () => {
  const [taskConfigs, setTaskConfigs] = useState([])
  const [loading, setLoading] = useState(false)
  const [form] = Form.useForm()
  const [modalVisible, setModalVisible] = useState(false)
  const [isEditing, setIsEditing] = useState(false)
  const [currentConfig, setCurrentConfig] = useState(null)
  const [modal, modalContext] = Modal.useModal()

  // 获取TaskConfig列表
  const fetchTaskConfigs = () => {
    setLoading(true)
    fetch('/api/task/config/list')
      .then(response => response.json())
      .then(data => {
        setTaskConfigs(data)
        setLoading(false)
      })
      .catch(error => {
        message.error('获取TaskConfig列表失败')
        setLoading(false)
      })
  }

  // 组件加载时获取数据
  useEffect(() => {
    fetchTaskConfigs()
  }, [])

  // 打开新增模态框
  const handleAdd = () => {
    setIsEditing(false)
    setCurrentConfig(null)
    form.resetFields()
    setModalVisible(true)
  }

  // 打开编辑模态框
  const handleEdit = (record) => {
    setIsEditing(true)
    setCurrentConfig(record)
    form.setFieldsValue({
      task_type: record.task_type,
      title: record.title,
      form: record.form,
      run_endpoint: record.run_endpoint,
      input_endpoint: record.input_endpoint
    })
    setModalVisible(true)
  }

  // 删除TaskConfig
  const handleDelete = (record) => {
    modal.confirm({
      title: '确认删除',
      content: `确定要删除TaskConfig: ${record.task_type}吗？`,
      okText: '确定',
      okType: 'danger',
      cancelText: '取消',
      onOk: () => {
        fetch('/api/task/config/delete', {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json'
          },
          body: JSON.stringify({ task_type: record.task_type })
        })
          .then(response => response.json())
          .then(data => {
            message.success('删除成功')
            fetchTaskConfigs()
          })
          .catch(error => {
            message.error('删除失败')
          })
      }
    })
  }

  // 保存TaskConfig（新增或编辑）
  const handleSave = () => {
    form.validateFields().then(values => {
      const url = isEditing ? '/api/task/config/update' : '/api/task/config/create'
      fetch(url, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json'
        },
        body: JSON.stringify(values)
      })
        .then(response => response.json())
        .then(data => {
          message.success(isEditing ? '更新成功' : '创建成功')
          setModalVisible(false)
          fetchTaskConfigs()
        })
        .catch(error => {
          message.error(isEditing ? '更新失败' : '创建失败')
        })
    }).catch(errorInfo => {
      console.log('表单验证失败:', errorInfo)
    })
  }

  // 表格列定义
  const columns = [
    {
      title: 'TaskType',
      dataIndex: 'task_type',
      key: 'task_type',
      width: 150
    },
    {
      title: 'Title',
      dataIndex: 'title',
      key: 'title',
      width: 150
    },
    {
      title: 'Form',
      dataIndex: 'form',
      key: 'form',
      ellipsis: true,
      render: (text) => (
        <span title={text}>{text.substring(0, 50)}...</span>
      )
    },
    {
      title: 'Run',
      dataIndex: 'run',
      key: 'run',
      width: 200
    },
    {
      title: 'Input',
      dataIndex: 'input',
      key: 'input',
      width: 200
    },
    {
      title: '操作',
      key: 'action',
      width: 150,
      render: (_, record) => (
        <Space size="middle">
          <Button type="link" icon={<EditOutlined />} onClick={() => handleEdit(record)}>
            编辑
          </Button>
          <Button type="link" danger icon={<DeleteOutlined />} onClick={() => handleDelete(record)}>
            删除
          </Button>
        </Space>
      )
    }
  ]

  return (
    <div style={{ padding: '16px' }}>
        {modalContext}
      <Card title="TaskConfig管理" bordered={false} extra={
        <Button type="primary" icon={<PlusOutlined />} onClick={handleAdd}>
          新增
        </Button>
      }>
        <Table
          columns={columns}
          dataSource={taskConfigs}
          rowKey="task_type"
          loading={loading}
          pagination={{ pageSize: 10 }}
          scroll={{ x: 1000 }}
        />
      </Card>

      {/* 新增/编辑模态框 */}
      <Modal
        title={isEditing ? '编辑TaskConfig' : '新增TaskConfig'}
        open={modalVisible}
        onOk={handleSave}
        onCancel={() => setModalVisible(false)}
        footer={[
          <Button key="back" onClick={() => setModalVisible(false)} icon={<CloseOutlined />}>
            取消
          </Button>,
          <Button key="submit" type="primary" onClick={handleSave} icon={<SaveOutlined />}>
            保存
          </Button>
        ]}
        width={800}
      >
        <Form
          form={form}
          layout="vertical"
          initialValues={{
            form: '{}',
            run: '',
            input: ''
          }}
        >
          <Form.Item
            name="task_type"
            label="TaskType"
            rules={[{ required: true, message: '请输入TaskType' }]}
          >
            <Input placeholder="请输入TaskType" disabled={isEditing} />
          </Form.Item>
          <Form.Item
            name="title"
            label="Title"
            rules={[{ required: true, message: '请输入Title' }]}
          >
            <Input placeholder="请输入Title" />
          </Form.Item>
          <Form.Item
            name="form"
            label="Form (JSON)"
            rules={[{ required: true, message: '请输入Form' }]}
          >
            <TextArea placeholder="请输入Form JSON" rows={6} />
          </Form.Item>
          <Form.Item
            name="run_endpoint"
            label="Run Endpoint"
          >
            <Input placeholder="请输入Run Endpoint" />
          </Form.Item>
          <Form.Item
            name="input_endpoint"
            label="Input Endpoint"
          >
            <Input placeholder="请输入Input Endpoint" />
          </Form.Item>
        </Form>
      </Modal>
    </div>
  )
}

export const Route = createFileRoute('/task-config')({
  component: TaskConfigPage,
})