import React, { PureComponent } from 'react';
import { connect } from 'dva';
import { Form, Input, Modal, Select } from 'antd';

@connect(state => ({
  resource: state.resource,
}))
@Form.create()
class ResourceCard extends PureComponent {
  state = {
    methodData: ['GET', 'POST', 'PUT', 'PATCH', 'DELETE'],
  };

  onOKClick = () => {
    const { form, onSubmit } = this.props;

    form.validateFieldsAndScroll((err, values) => {
      if (!err) {
        const formData = { ...values };
        formData.status = parseInt(formData.status, 10);
        onSubmit(formData);
      }
    });
  };

  dispatch = action => {
    const { dispatch } = this.props;
    dispatch(action);
  };

  render() {
    const {
      onCancel,
      resource: { formTitle, formVisible, formData, submitting },
      form: { getFieldDecorator },
    } = this.props;

    const formItemLayout = {
      labelCol: {
        xs: { span: 24 },
        sm: { span: 6 },
      },
      wrapperCol: {
        xs: { span: 24 },
        sm: { span: 16 },
      },
    };

    const { methodData } = this.state;

    return (
      <Modal
        title={formTitle}
        width={600}
        visible={formVisible}
        maskClosable={false}
        confirmLoading={submitting}
        destroyOnClose
        onOk={this.onOKClick}
        onCancel={onCancel}
        style={{ top: 20 }}
      >
        <Form>
          <Form.Item {...formItemLayout} label="资源编号">
            {getFieldDecorator('code', {
              initialValue: formData.code,
              rules: [
                {
                  required: true,
                  message: '请输入资源编号',
                },
              ],
            })(<Input placeholder="请输入资源编号" />)}
          </Form.Item>
          <Form.Item {...formItemLayout} label="资源名称">
            {getFieldDecorator('name', {
              initialValue: formData.name,
              rules: [
                {
                  required: true,
                  message: '请输入资源名称',
                },
              ],
            })(<Input placeholder="请输入资源名称" />)}
          </Form.Item>
          <Form.Item {...formItemLayout} label="访问路径">
            {getFieldDecorator('path', {
              initialValue: formData.path,
              rules: [
                {
                  required: true,
                  message: '请输入访问路径',
                },
              ],
            })(<Input placeholder="请输入访问路径" />)}
          </Form.Item>
          <Form.Item {...formItemLayout} label="请求方式">
            {getFieldDecorator('method', {
              initialValue: formData.method,
              rules: [
                {
                  required: true,
                  message: '请选择请求方式',
                },
              ],
            })(
              <Select style={{ width: '100%' }} placeholder="请选择">
                {methodData.map(item => (
                  <Select.Option key={item} value={item}>
                    {item}
                  </Select.Option>
                ))}
              </Select>
            )}
          </Form.Item>
        </Form>
      </Modal>
    );
  }
}

export default ResourceCard;
