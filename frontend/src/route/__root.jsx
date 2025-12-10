import { createRootRoute, Link, Outlet, useLocation } from '@tanstack/react-router'
import { Layout, Menu } from 'antd'
import 'antd/dist/reset.css'
import { HomeOutlined, SettingOutlined, TableOutlined } from '@ant-design/icons'

const RootLayout = () => {
  const { Header, Content } = Layout
  const location = useLocation()

  return (
    <Layout style={{ minHeight: '100vh' }}>
      <Header style={{ padding: 0, background: '#fff', boxShadow: '0 2px 8px rgba(0, 0, 0, 0.1)' }}>
        <div style={{ display: 'flex', alignItems: 'center', justifyContent: 'space-between', padding: '0 16px' }}>
          <div style={{ fontSize: '20px', fontWeight: 'bold', color: '#1890ff', display: 'flex', alignItems: 'center' }}>
            任务控制台
          </div>
          <Menu
            mode="horizontal"
            selectedKeys={[location.pathname]}
            style={{ borderBottom: 0 }}
            items={[
              {
                key: '/',
                icon: <HomeOutlined />,
                label: <Link to="/">任务执行</Link>,
              },
              {
                key: '/tasks',
                icon: <TableOutlined />,
                label: <Link to="/tasks">任务记录</Link>,
              }
            ]}
          />
        </div>
      </Header>
      <Content style={{ margin: '16px' }}>
        <Outlet />
      </Content>
    </Layout>
  )
}

export const Route = createRootRoute({ component: RootLayout })