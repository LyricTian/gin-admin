import React, { PureComponent } from 'react';
import { connect } from 'dva';
import { Form, Input, Modal, message, Card, Row, Col } from 'antd';

import MenuTree from '../Menu/MenuTree';

@connect(state => ({
  role: state.role,
}))
@Form.create()
class RoleCard extends PureComponent {
  onOKClick = () => {
    const {
      form,
      role: { menuKeys },
      onSubmit,
    } = this.props;

    form.validateFieldsAndScroll((err, values) => {
      if (!err) {
        if (menuKeys.length === 0) {
          message.warning('请选择权限菜单！');
          return;
        }

        const formData = { ...values, menu_ids: menuKeys };
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
      role: { formTitle, formVisible, formData, submitting, menuKeys },
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
      },
    };

    const menuTreeProps = {
      checkedKeys: menuKeys,
      onCheck: keys => {
        this.dispatch({
          type: 'role/saveMenuKeys',
          payload: keys,
        });
      },
    };

    return (
      <Modal
        title={formTitle}
        width={450}
        visible={formVisible}
        maskClosable={false}
        confirmLoading={submitting}
        destroyOnClose
        onOk={this.onOKClick}
        onCancel={onCancel}
        style={{ top: 20 }}
        bodyStyle={{ height: 550, overflowY: 'scroll' }}
      >
        <Form>
          <Row>
            <Col md={24} sm={24}>
              <Form.Item {...formItemLayout} label="角色名称">
                {getFieldDecorator('name', {
                  initialValue: formData.name,
                  rules: [
                    {
                      required: true,
                      message: '请输入角色名称',
                    },
                  ],
                })(<Input placeholder="请输入角色名称" />)}
              </Form.Item>
            </Col>
          </Row>
          <Row>
            <Col md={24} sm={24}>
              <Form.Item {...formItemLayout} label="备注">
                {getFieldDecorator('memo', {
                  initialValue: formData.memo,
                })(<Input.TextArea rows={2} placeholder="请输入备注" />)}
              </Form.Item>
            </Col>
          </Row>
        </Form>
        <Card title="选择菜单权限">
          <div style={{ paddingLeft: 20 }}>
            <MenuTree treeProps={menuTreeProps} />
          </div>
        </Card>
      </Modal>
    );
  }
}

export default RoleCard;
