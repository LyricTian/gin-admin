import React, { PureComponent } from 'react';
import { connect } from 'dva';
import { Form, Input, Card, Radio, Modal, TreeSelect, InputNumber, Row, Col } from 'antd';

@connect(state => ({
  menu: state.menu,
}))
@Form.create()
export default class MenuCard extends PureComponent {
  onOKClick = () => {
    const { form, onSubmit } = this.props;
    form.validateFieldsAndScroll((err, values) => {
      if (!err) {
        const formData = { ...values };
        formData.status = parseInt(formData.status, 10);
        formData.type = parseInt(formData.type, 10);
        formData.sequence = parseInt(formData.sequence, 10);
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
      menu: { formVisible, formTitle, formData, submitting, treeData },
      form: { getFieldDecorator },
      onCancel,
    } = this.props;

    const formItemLayout = {
      labelCol: {
        xs: { span: 24 },
        sm: { span: 6 },
      },
      wrapperCol: {
        xs: { span: 24 },
        sm: { span: 18 },
        md: { span: 18 },
      },
    };

    return (
      <Modal
        title={formTitle}
        width={850}
        visible={formVisible}
        maskClosable={false}
        confirmLoading={submitting}
        destroyOnClose
        onOk={this.onOKClick}
        onCancel={onCancel}
      >
        <Card bordered={false}>
          <Form>
            <Row>
              <Col md={12} sm={24}>
                <Form.Item {...formItemLayout} label="菜单名称">
                  {getFieldDecorator('name', {
                    initialValue: formData.name,
                    rules: [
                      {
                        required: true,
                        message: '请输入菜单名称',
                      },
                    ],
                  })(<Input placeholder="请输入" />)}
                </Form.Item>
              </Col>
              <Col md={12} sm={24}>
                <Form.Item {...formItemLayout} label="菜单编号">
                  {getFieldDecorator('code', {
                    initialValue: formData.code,
                  })(<Input placeholder="请输入" />)}
                </Form.Item>
              </Col>
            </Row>
            <Row>
              <Col md={12} sm={24}>
                <Form.Item {...formItemLayout} label="菜单上级">
                  {getFieldDecorator('parent_id', {
                    initialValue: formData.parent_id,
                  })(
                    <TreeSelect
                      showSearch
                      treeNodeFilterProp="title"
                      treeDefaultExpandedKeys={[treeData.length > 0 && treeData[0].record_id]}
                      style={{ minWidth: 200 }}
                      dropdownStyle={{ maxHeight: 400, overflow: 'auto' }}
                      treeData={treeData}
                      placeholder="请选择"
                    />
                  )}
                </Form.Item>
              </Col>
              <Col md={12} sm={24}>
                <Form.Item {...formItemLayout} label="菜单类型">
                  {getFieldDecorator('type', {
                    initialValue: formData.type ? formData.type.toString() : '',
                  })(
                    <Radio.Group>
                      <Radio value="10">系统</Radio>
                      <Radio value="20">模块</Radio>
                      <Radio value="30">功能</Radio>
                      <Radio value="40">动作</Radio>
                    </Radio.Group>
                  )}
                </Form.Item>
              </Col>
            </Row>
            <Row>
              <Col md={12} sm={24}>
                <Form.Item {...formItemLayout} label="排序值">
                  {getFieldDecorator('sequence', {
                    initialValue: formData.sequence ? formData.sequence.toString() : '0',
                  })(<InputNumber min={1} style={{ width: '100%' }} />)}
                </Form.Item>
              </Col>
              <Col md={12} sm={24}>
                <Form.Item {...formItemLayout} label="菜单图标">
                  {getFieldDecorator('icon', {
                    initialValue: formData.icon,
                  })(<Input placeholder="请输入" />)}
                </Form.Item>
              </Col>
            </Row>
            <Row>
              <Col md={12} sm={24}>
                <Form.Item {...formItemLayout} label="跳转链接">
                  {getFieldDecorator('path', {
                    initialValue: formData.path,
                  })(<Input placeholder="请输入" />)}
                </Form.Item>
              </Col>
              <Col md={12} sm={24}>
                <Form.Item {...formItemLayout} label="菜单状态">
                  {getFieldDecorator('status', {
                    initialValue: formData.status ? formData.status.toString() : '1',
                  })(
                    <Radio.Group>
                      <Radio value="1">正常</Radio>
                      <Radio value="2">停用</Radio>
                    </Radio.Group>
                  )}
                </Form.Item>
              </Col>
            </Row>
          </Form>
        </Card>
      </Modal>
    );
  }
}
