import React, { PureComponent } from 'react';
import { connect } from 'dva';
import { Form, Input, Card, Radio, Modal, TreeSelect, Select, InputNumber, Row, Col } from 'antd';

@connect(({ menu }) => ({
  menu,
}))
@Form.create()
class MenuCard extends PureComponent {
  state = {
    methodData: ['GET', 'POST', 'PUT', 'PATCH', 'DELETE'],
  };

  onOKClick = () => {
    const { form, onSubmit } = this.props;
    form.validateFieldsAndScroll((err, values) => {
      if (!err) {
        const formData = { ...values };
        formData.is_hide = parseInt(formData.is_hide, 10);
        formData.type = parseInt(formData.type, 10);
        formData.sequence = parseInt(formData.sequence, 10);
        onSubmit(formData);
      }
    });
  };

  onTypeChange = e => {
    const {
      menu: { formData },
    } = this.props;
    const newFormData = { ...formData, type: parseInt(e.target.value, 10) };
    this.dispatch({
      type: 'menu/saveFormData',
      payload: newFormData,
    });
  };

  dispatch = action => {
    const { dispatch } = this.props;
    dispatch(action);
  };

  toTreeSelect = data => {
    if (!data) {
      return [];
    }
    const newData = [];
    for (let i = 0; i < data.length; i += 1) {
      const item = { ...data[i], title: data[i].name, value: data[i].record_id };
      if (item.children && item.children.length > 0) {
        item.children = this.toTreeSelect(item.children);
      }
      newData.push(item);
    }
    return newData;
  };

  render() {
    const {
      menu: { formVisible, formTitle, formData, submitting, treeData },
      form: { getFieldDecorator },
      onCancel,
    } = this.props;

    const { methodData } = this.state;

    const formItemLayout = {
      labelCol: {
        span: 6,
      },
      wrapperCol: {
        span: 18,
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
              <Col md={12}>
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
              <Col md={12}>
                <Form.Item {...formItemLayout} label="菜单编号">
                  {getFieldDecorator('code', {
                    initialValue: formData.code,
                    rules: [
                      {
                        required: true,
                        message: '请输入菜单编号',
                      },
                    ],
                  })(<Input placeholder="请输入" />)}
                </Form.Item>
              </Col>
            </Row>
            <Row>
              <Col md={12}>
                <Form.Item {...formItemLayout} label="菜单上级">
                  {getFieldDecorator('parent_id', {
                    initialValue: formData.parent_id,
                  })(
                    <TreeSelect
                      showSearch
                      treeNodeFilterProp="title"
                      style={{ minWidth: 200 }}
                      dropdownStyle={{ maxHeight: 400, overflow: 'auto' }}
                      treeData={this.toTreeSelect(treeData)}
                      placeholder="请选择"
                    />
                  )}
                </Form.Item>
              </Col>
              <Col md={12}>
                <Form.Item {...formItemLayout} label="菜单类型">
                  {getFieldDecorator('type', {
                    initialValue: formData.type ? formData.type.toString() : '1',
                  })(
                    <Radio.Group onChange={this.onTypeChange}>
                      <Radio value="1">模块</Radio>
                      <Radio value="2">功能</Radio>
                      <Radio value="3">资源</Radio>
                    </Radio.Group>
                  )}
                </Form.Item>
              </Col>
            </Row>
            {formData.type !== 3 && (
              <Row>
                <Col md={12}>
                  <Form.Item {...formItemLayout} label="菜单图标">
                    {getFieldDecorator('icon', {
                      initialValue: formData.icon,
                    })(<Input placeholder="请输入" />)}
                  </Form.Item>
                </Col>
                <Col md={12}>
                  <Form.Item {...formItemLayout} label="访问路径">
                    {getFieldDecorator('path', {
                      initialValue: formData.path,
                    })(<Input placeholder="请输入" />)}
                  </Form.Item>
                </Col>
              </Row>
            )}
            {formData.type === 3 && (
              <Row>
                <Col md={12}>
                  <Form.Item {...formItemLayout} label="资源路径">
                    {getFieldDecorator('path', {
                      initialValue: formData.path,
                      rules: [
                        {
                          required: true,
                          message: '请输入资源路径',
                        },
                      ],
                    })(<Input placeholder="请输入" />)}
                  </Form.Item>
                </Col>
                <Col md={12}>
                  <Form.Item {...formItemLayout} label="请求方式">
                    {getFieldDecorator('method', {
                      initialValue: formData.method,
                      rules: [
                        {
                          required: true,
                          message: '请输入请求方式',
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
                </Col>
              </Row>
            )}
            <Row>
              <Col md={12}>
                <Form.Item {...formItemLayout} label="排序值">
                  {getFieldDecorator('sequence', {
                    initialValue: formData.sequence ? formData.sequence.toString() : '0',
                  })(<InputNumber min={1} style={{ width: '100%' }} />)}
                </Form.Item>
              </Col>
              <Col md={12}>
                <Form.Item {...formItemLayout} label="隐藏菜单">
                  {getFieldDecorator('is_hide', {
                    initialValue: formData.is_hide ? formData.is_hide.toString() : '2',
                  })(
                    <Radio.Group>
                      <Radio value="1">是</Radio>
                      <Radio value="2">否</Radio>
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

export default MenuCard;
