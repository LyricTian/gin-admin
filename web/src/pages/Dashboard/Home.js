import React, { PureComponent } from 'react';
import { connect } from 'dva';
import { Card } from 'antd';
import moment from 'moment';
import PageHeaderLayout from '../../layouts/PageHeaderLayout';
import styles from './Home.less';

@connect(state => ({
  global: state.global,
}))
class Home extends PureComponent {
  state = {
    currentTime: moment().format('HH:mm:ss'),
  };

  componentDidMount() {
    this.interval = setInterval(() => {
      this.setState({ currentTime: moment().format('HH:mm:ss') });
    }, 1000);
  }

  componentWillUnmount() {
    clearInterval(this.interval);
  }

  render() {
    const {
      global: { user },
    } = this.props;

    const { currentTime } = this.state;

    const breadcrumbList = [{ title: '首页' }];

    return (
      <PageHeaderLayout
        title={`您好，${user.real_name}，祝你开心每一天！`}
        breadcrumbList={breadcrumbList}
        content={user.role_names && user.role_names.length > 0 ? user.role_names.join('|') : null}
      >
        <Card bordered={false}>
          <div className={styles.showTime}>{currentTime}</div>
        </Card>
      </PageHeaderLayout>
    );
  }
}

export default Home;
