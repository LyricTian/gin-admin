import React, { PureComponent } from 'react';
import { connect } from 'dva';
import { Form, Input, Button, Icon, Alert } from 'antd';
import { md5Hash } from '../../utils/utils';

import styles from './Index.less';

@connect(({ login }) => ({
  login,
}))
@Form.create()
class Login extends PureComponent {
  handleSubmit = e => {
    e.preventDefault();

    const { form, dispatch } = this.props;
    form.validateFields({ force: true }, (err, values) => {
      if (!err) {
        dispatch({
          type: 'login/submit',
          payload: {
            user_name: values.user_name,
            password: md5Hash(values.password),
          },
        });
      }
    });
  };

  renderMessage = (type, message) => (
    <Alert style={{ marginBottom: 24 }} message={message} type={type} showIcon />
  );

  render() {
    const {
      form: { getFieldDecorator },
      login,
    } = this.props;

    return (
      <div className={styles.main}>
        <Form onSubmit={this.handleSubmit}>
          {login.status === 'fail' &&
            login.submitting === false &&
            this.renderMessage('warning', '用户名或密码错误，请重新输入！')}

          {login.status === 'error' &&
            login.submitting === false &&
            this.renderMessage('error', '服务器发生错误，请联系管理员！')}

          <Form.Item>
            {getFieldDecorator('user_name', {
              rules: [
                {
                  required: true,
                  message: '请输入账户名！',
                },
              ],
            })(
              <Input
                size="large"
                prefix={<Icon type="user" className={styles.prefixIcon} />}
                placeholder="用户名：demo"
              />
            )}
          </Form.Item>
          <Form.Item>
            {getFieldDecorator('password', {
              rules: [
                {
                  required: true,
                  message: '请输入密码！',
                },
              ],
            })(
              <Input
                size="large"
                prefix={<Icon type="lock" className={styles.prefixIcon} />}
                type="password"
                placeholder="密码：888888"
              />
            )}
          </Form.Item>

          <Form.Item className={styles.additional}>
            <Button
              size="large"
              loading={login.submitting}
              className={styles.submit}
              type="primary"
              htmlType="submit"
            >
              登录
            </Button>
          </Form.Item>
        </Form>
      </div>
    );
  }
}

export default Login;
