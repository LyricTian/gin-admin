import React from 'react';
import { Layout, Menu, Icon, Avatar, Dropdown, Spin } from 'antd';
import DocumentTitle from 'react-document-title';
import { connect } from 'dva';
import { Link, Route, Redirect, Switch } from 'dva/router';
import { ContainerQuery } from 'react-container-query';
import classNames from 'classnames';
import Debounce from 'lodash-decorators/debounce';
import GlobalFooter from '../components/GlobalFooter';
import CopyRight from '../components/CopyRight';
import NotFound from '../routes/Exception/404';
import styles from './AdminLayout.less';
import logo from '../assets/logo.svg';

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
  defaultURL: state.global.defaultURL,
  collapsed: state.global.collapsed,
  openKeys: state.global.openKeys,
  selectedKeys: state.global.selectedKeys,
  user: state.global.user,
  menuPaths: state.global.menuPaths,
  menus: state.global.menus,
}))
class AdminLayout extends React.PureComponent {
  componentDidMount() {
    const {
      location: { pathname },
    } = this.props;

    this.dispatch({
      type: 'global/fetchUser',
    });

    this.dispatch({
      type: 'global/fetchMenus',
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
            key={item.level_code}
          >
            {this.renderNavMenuItems(item.children)}
          </SubMenu>
        );
      }

      const itemPath = item.path;
      const icon = item.icon && <Icon type={item.icon} />;
      const {
        location: { pathname },
      } = this.props;

      return (
        <Menu.Item key={item.level_code}>
          {itemPath.startsWith('http') ? (
            <a href={itemPath} target="_blank" rel="noopener noreferrer">
              {icon}
              <span>{item.name}</span>
            </a>
          ) : (
            <Link to={itemPath} replace={itemPath === pathname}>
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
      getRouteData,
      title,
    } = this.props;
    let ptitle = title;

    getRouteData('AdminLayout').forEach(item => {
      if (item.path === pathname) {
        ptitle = `${item.name} - ${title}`;
      }
    });
    return ptitle;
  }

  render() {
    const {
      user,
      collapsed,
      getRouteData,
      menus,
      copyRight,
      defaultURL,
      openKeys,
      title,
      selectedKeys,
    } = this.props;

    const menu = (
      <Menu className={styles.menu} selectedKeys={[]} onClick={this.onMenuClick}>
        <Menu.Item key="setting" disabled>
          <Icon type="user" />
          个人信息
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
            theme="light"
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
              <Switch>
                {getRouteData('AdminLayout').map(item => (
                  <Route
                    exact={item.exact}
                    key={item.path}
                    path={item.path}
                    component={item.component}
                  />
                ))}
                <Redirect exact from="/" to={defaultURL} />
                <Route component={NotFound} />
              </Switch>
            </div>
            <GlobalFooter copyright={<CopyRight title={copyRight} />} />
          </Content>
        </Layout>
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
