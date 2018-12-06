import React from 'react';
import { connect } from 'dva';
import GlobalFooter from '@/components/GlobalFooter';
import CopyRight from '@/components/CopyRight';
import styles from './UserLayout.less';
import logo from '../assets/logo.svg';

@connect(state => ({
  global: state.global,
}))
class UserLayout extends React.PureComponent {
  render() {
    const {
      children,
      global: { title, copyRight },
    } = this.props;

    return (
      <div className={styles.container}>
        <div className={styles.top}>
          <div className={styles.header}>
            <img alt="" className={styles.logo} src={logo} />
            <span className={styles.title}>{title}</span>
          </div>
          <div className={styles.desc} />
        </div>
        {children}
        <GlobalFooter className={styles.footer} copyright={<CopyRight title={copyRight} />} />
      </div>
    );
  }
}

export default UserLayout;
