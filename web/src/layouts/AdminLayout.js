import React from 'react';
import { Layout, Menu, Icon, Avatar, Dropdown, Spin } from 'antd';
import DocumentTitle from 'react-document-title';
import { connect } from 'dva';
import Link from 'umi/link';
import { ContainerQuery } from 'react-container-query';
import classNames from 'classnames';
import Debounce from 'lodash-decorators/debounce';
import GlobalFooter from '@/components/GlobalFooter';
import CopyRight from '@/components/CopyRight';
import UpdatePasswordDialog from '@/components/UpdatePasswordDialog';
import styles from './AdminLayout.less';
import logo from '../assets/logo.svg';
import GetGlobalContext from '@/utils/context';

const { Header, Sider, Content } = Layout;
const { SubMenu } = Menu;

const query = {
  'screen-xs': {
    maxWidth: 575,
  },
  'screen-sm': {
    minWidth: 576,
    maxWidth: 767,
  },
  'screen-md': {
    minWidth: 768,
    maxWidth: 991,
  },
  'screen-lg': {
    minWidth: 992,
    maxWidth: 1199,
  },
  'screen-xl': {
    minWidth: 1200,
  },
};

@connect(state => ({
  title: state.global.title,
  copyRight: state.global.copyRight,
  collapsed: state.global.collapsed,
  openKeys: state.global.openKeys,
  selectedKeys: state.global.selectedKeys,
  user: state.global.user,
  menuPaths: state.global.menuPaths,
  menus: state.global.menus,
  global: state.global,
}))
class AdminLayout extends React.PureComponent {
  state = {
    updatePwdVisible: false,
  };

  componentDidMount() {
    const {
      location: { pathname },
    } = this.props;

    this.dispatch({
      type: 'global/fetchUser',
    });

    this.dispatch({
      type: 'global/fetchMenuTree',
      pathname,
    });
  }

  dispatch = action => {
    const { dispatch } = this.props;
    dispatch(action);
  };

  onCollapse = collapsed => {
    this.dispatch({
      type: 'global/changeLayoutCollapsed',
      payload: collapsed,
    });
  };

  onMenuClick = ({ key }) => {
    if (key === 'logout') {
      this.dispatch({
        type: 'login/logout',
      });
    } else if (key === 'updatepwd') {
      this.setState({ updatePwdVisible: true });
    }
  };

  onMenuOpenChange = openKeys => {
    if (openKeys.length > 1) {
      const lastKey = openKeys[openKeys.length - 1];
      if (lastKey.length === openKeys[0].length) {
        this.dispatch({
          type: 'global/changeOpenKeys',
          payload: [lastKey],
        });
        return;
      }
    }
    this.dispatch({
      type: 'global/changeOpenKeys',
      payload: [...openKeys],
    });
  };

  onToggleClick = () => {
    const { collapsed } = this.props;
    this.dispatch({
      type: 'global/changeLayoutCollapsed',
      payload: !collapsed,
    });
    this.onTriggerResizeEvent();
  };

  @Debounce(600)
  onTriggerResizeEvent = () => {
    const event = document.createEvent('HTMLEvents');
    event.initEvent('resize', true, false);
    window.dispatchEvent(event);
  };

  handleUpdatePwdCancel = () => {
    this.setState({ updatePwdVisible: false });
  };

  renderNavMenuItems(menusData) {
    if (!menusData) {
      return [];
    }

    return menusData.map(item => {
      if (!item.name) {
        return null;
      }
      if (item.children && item.children.some(child => child.name)) {
        return (
          <SubMenu
            title={
              item.icon ? (
                <span>
                  <Icon type={item.icon} />
                  <span>{item.name}</span>
                </span>
              ) : (
                item.name
              )
            }
            key={item.record_id}
          >
            {this.renderNavMenuItems(item.children)}
          </SubMenu>
        );
      }

      const { router } = item;
      const icon = item.icon && <Icon type={item.icon} />;
      const {
        location: { pathname },
      } = this.props;

      return (
        <Menu.Item key={item.record_id}>
          {router.startsWith('http') ? (
            <a href={router} target="_blank" rel="noopener noreferrer">
              {icon}
              <span>{item.name}</span>
            </a>
          ) : (
            <Link to={router} replace={router === pathname}>
              {icon}
              <span>{item.name}</span>
            </Link>
          )}
        </Menu.Item>
      );
    });
  }

  renderPageTitle() {
    const {
      location: { pathname },
      menuPaths,
      title,
    } = this.props;

    let ptitle = title;
    const item = menuPaths[pathname];
    if (item) {
      ptitle = `${item.name} - ${title}`;
    }
    return ptitle;
  }

  render() {
    const {
      children,
      user,
      collapsed,
      menus,
      copyRight,
      openKeys,
      title,
      selectedKeys,
      global,
    } = this.props;

    const { updatePwdVisible } = this.state;
    const GlobalContext = GetGlobalContext();

    const menu = (
      <Menu className={styles.menu} selectedKeys={[]} onClick={this.onMenuClick}>
        <Menu.Item key="updatepwd">
          <Icon type="lock" />
          修改密码
        </Menu.Item>
        <Menu.Divider />
        <Menu.Item key="logout">
          <Icon type="logout" />
          退出登录
        </Menu.Item>
      </Menu>
    );

    // Don't show popup menu when it is been collapsed
    const menuProps = collapsed ? {} : { openKeys };

    const layout = (
      <Layout>
        <Sider
          trigger={null}
          collapsible
          collapsed={collapsed}
          breakpoint="lg"
          onCollapse={this.onCollapse}
          width={256}
          className={styles.sider}
        >
          <div className={styles.logo}>
            <Link to="/">
              <img src={logo} alt="logo" />
              <h1>{title}</h1>
            </Link>
          </div>
          <Menu
            mode="inline"
            {...menuProps}
            onOpenChange={this.onMenuOpenChange}
            selectedKeys={selectedKeys}
            style={{ margin: '16px 0', width: '100%' }}
          >
            {this.renderNavMenuItems(menus)}
          </Menu>
        </Sider>
        <Layout>
          <Header className={styles.header}>
            <Icon
              className={styles.trigger}
              type={collapsed ? 'menu-unfold' : 'menu-fold'}
              onClick={this.onToggleClick}
            />
            <div className={styles.right}>
              {user.user_name ? (
                <Dropdown overlay={menu}>
                  <span className={`${styles.action} ${styles.account}`}>
                    <Avatar size="small" className={styles.avatar} icon="user" />
                    {user.real_name !== ''
                      ? `${user.user_name}(${user.real_name})`
                      : user.user_name}
                  </span>
                </Dropdown>
              ) : (
                <Spin size="small" style={{ marginLeft: 8 }} />
              )}
            </div>
          </Header>
          <Content style={{ margin: '24px 24px 0', height: '100%' }}>
            <div style={{ minHeight: 'calc(100vh - 260px)' }}>
              <GlobalContext.Provider value={global}>{children}</GlobalContext.Provider>
            </div>
            <GlobalFooter copyright={<CopyRight title={copyRight} />} />
          </Content>
        </Layout>
        <UpdatePasswordDialog visible={updatePwdVisible} onCancel={this.handleUpdatePwdCancel} />
      </Layout>
    );

    return (
      <DocumentTitle title={this.renderPageTitle()}>
        <ContainerQuery query={query}>
          {params => <div className={classNames(params)}>{layout}</div>}
        </ContainerQuery>
      </DocumentTitle>
    );
  }
}

export default AdminLayout;
