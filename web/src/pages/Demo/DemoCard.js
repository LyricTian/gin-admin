import React, { PureComponent } from 'react';
import { connect } from 'dva';
import { Form, Input, Modal, Radio } from 'antd';

@connect(state => ({
  demo: state.demo,
}))
@Form.create()
class DemoCard extends PureComponent {
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
      demo: { formTitle, formVisible, formData, submitting },
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
        bodyStyle={{ maxHeight: 'calc( 100vh - 158px )', overflowY: 'auto' }}
      >
        <Form>
          <Form.Item {...formItemLayout} label="编号">
            {getFieldDecorator('code', {
              initialValue: formData.code,
              rules: [
                {
                  required: true,
                  message: '请输入编号',
                },
              ],
            })(<Input placeholder="请输入编号" />)}
          </Form.Item>
          <Form.Item {...formItemLayout} label="名称">
            {getFieldDecorator('name', {
              initialValue: formData.name,
              rules: [
                {
                  required: true,
                  message: '请输入名称',
                },
              ],
            })(<Input placeholder="请输入名称" />)}
          </Form.Item>
          <Form.Item {...formItemLayout} label="备注">
            {getFieldDecorator('memo', {
              initialValue: formData.memo,
            })(<Input.TextArea rows={2} placeholder="请输入备注" />)}
          </Form.Item>
          <Form.Item {...formItemLayout} label="状态">
            {getFieldDecorator('status', {
              initialValue: formData.status ? formData.status.toString() : '1',
            })(
              <Radio.Group>
                <Radio value="1">正常</Radio>
                <Radio value="2">停用</Radio>
              </Radio.Group>
            )}
          </Form.Item>
        </Form>
      </Modal>
    );
  }
}

export default DemoCard;
