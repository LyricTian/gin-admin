import React from 'react';
import { Router, Route, Switch } from 'dva/router';
import { LocaleProvider, Spin } from 'antd';
import zhCN from 'antd/lib/locale-provider/zh_CN';
import dynamic from 'dva/dynamic';
import cloneDeep from 'lodash/cloneDeep';
import RouteData from './router.config';
import { getPlainNode } from './utils/utils';
import styles from './index.less';

dynamic.setDefaultLoadingComponent(() => <Spin size="large" className={styles.globalSpin} />);

function getRouteData(navData, path) {
  if (
    !navData.some(item => item.layout === path) ||
    !navData.filter(item => item.layout === path)[0].children
  ) {
    return null;
  }
  const route = cloneDeep(navData.filter(item => item.layout === path)[0]);
  const nodeList = getPlainNode(route.children);
  return nodeList;
}

function getLayout(data, path) {
  if (
    !data.some(item => item.layout === path) ||
    !data.filter(item => item.layout === path)[0].children
  ) {
    return null;
  }
  const route = data.filter(item => item.layout === path)[0];
  return {
    component: route.component,
    layout: route.layout,
    name: route.name,
    path: route.path,
  };
}

function RouterConfig({ history, app }) {
  const routeData = RouteData(app);
  const UserLayout = getLayout(routeData, 'UserLayout').component;
  const AdminLayout = getLayout(routeData, 'AdminLayout').component;

  const passProps = {
    app,
    getRouteData: path => getRouteData(routeData, path),
  };

  return (
    <LocaleProvider locale={zhCN}>
      <Router history={history}>
        <Switch>
          <Route path="/user" render={props => <UserLayout {...props} {...passProps} />} />
          <Route path="/" render={props => <AdminLayout {...props} {...passProps} />} />
        </Switch>
      </Router>
    </LocaleProvider>
  );
}

export default RouterConfig;
