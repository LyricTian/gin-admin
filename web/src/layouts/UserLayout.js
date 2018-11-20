import React from 'react';
import { connect } from 'dva';
import { Route } from 'dva/router';
import DocumentTitle from 'react-document-title';
import GlobalFooter from '../components/GlobalFooter';
import CopyRight from '../components/CopyRight';
import styles from './UserLayout.less';
import logo from '../assets/logo.svg';

@connect(state => ({
  global: state.global,
}))
class UserLayout extends React.PureComponent {
  getPageTitle() {
    const {
      getRouteData,
      location: { pathname },
      global: { title },
    } = this.props;

    let ptitle = title;
    getRouteData('UserLayout').forEach(item => {
      if (item.path === pathname) {
        ptitle = `${item.name} - ${title}`;
      }
    });
    return ptitle;
  }

  render() {
    const {
      getRouteData,
      global: { title, copyRight },
    } = this.props;

    return (
      <DocumentTitle title={this.getPageTitle()}>
        <div className={styles.container}>
          <div className={styles.top}>
            <div className={styles.header}>
              <img alt="" className={styles.logo} src={logo} />
              <span className={styles.title}>{title}</span>
            </div>
            <div className={styles.desc} />
          </div>
          {getRouteData('UserLayout').map(item => (
            <Route exact={item.exact} key={item.path} path={item.path} component={item.component} />
          ))}
          <GlobalFooter className={styles.footer} copyright={<CopyRight title={copyRight} />} />
        </div>
      </DocumentTitle>
    );
  }
}

export default UserLayout;
