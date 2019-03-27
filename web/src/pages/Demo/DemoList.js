import React, { PureComponent } from 'react';
import { connect } from 'dva';
import { Row, Col, Card, Form, Input, Button, Table, Modal, Badge, Select } from 'antd';
import PageHeaderLayout from '@/layouts/PageHeaderLayout';
import PButton from '@/components/PermButton';
import { formatDate } from '@/utils/utils';
import DemoCard from './DemoCard';

import styles from './DemoList.less';

@connect(state => ({
  loading: state.loading.models.demo,
  demo: state.demo,
}))
@Form.create()
class DemoList extends PureComponent {
  state = {
    selectedRowKeys: [],
    selectedRows: [],
  };

  componentDidMount() {
    this.dispatch({
      type: 'demo/fetch',
      search: {},
      pagination: {},
    });
  }

  onItemDisableClick = item => {
    this.dispatch({
      type: 'demo/changeStatus',
      payload: { record_id: item.record_id, status: 2 },
    });
  };

  onItemEnableClick = item => {
    this.dispatch({
      type: 'demo/changeStatus',
      payload: { record_id: item.record_id, status: 1 },
    });
  };

  onItemEditClick = item => {
    this.dispatch({
      type: 'demo/loadForm',
      payload: {
        type: 'E',
        id: item.record_id,
      },
    });
  };

  onAddClick = () => {
    this.dispatch({
      type: 'demo/loadForm',
      payload: {
        type: 'A',
      },
    });
  };

  onDelOKClick(id) {
    this.dispatch({
      type: 'demo/del',
      payload: { record_id: id },
    });
    this.clearSelectRows();
  }

  clearSelectRows = () => {
    const { selectedRowKeys } = this.state;
    if (selectedRowKeys.length === 0) {
      return;
    }
    this.setState({ selectedRowKeys: [], selectedRows: [] });
  };

  onItemDelClick = item => {
    Modal.confirm({
      title: `确定删除【基础示例数据：${item.name}】？`,
      okText: '确认',
      okType: 'danger',
      cancelText: '取消',
      onOk: this.onDelOKClick.bind(this, item.record_id),
    });
  };

  onTableSelectRow = (keys, rows) => {
    this.setState({
      selectedRowKeys: keys,
      selectedRows: rows,
    });
  };

  onTableChange = pagination => {
    this.dispatch({
      type: 'demo/fetch',
      pagination: {
        current: pagination.current,
        pageSize: pagination.pageSize,
      },
    });
    this.clearSelectRows();
  };

  onResetFormClick = () => {
    const { form } = this.props;
    form.resetFields();
    this.dispatch({
      type: 'demo/fetch',
      search: {},
      pagination: {},
    });
  };

  onSearchFormSubmit = e => {
    if (e) {
      e.preventDefault();
    }
    const { form } = this.props;
    form.validateFields({ force: true }, (err, values) => {
      if (err) {
        return;
      }
      this.dispatch({
        type: 'demo/fetch',
        search: values,
        pagination: {},
      });
      this.clearSelectRows();
    });
  };

  onDataFormSubmit = data => {
    this.dispatch({
      type: 'demo/submit',
      payload: data,
    });
    this.clearSelectRows();
  };

  onDataFormCancel = () => {
    this.dispatch({
      type: 'demo/changeFormVisible',
      payload: false,
    });
  };

  dispatch = action => {
    const { dispatch } = this.props;
    dispatch(action);
  };

  renderDataForm() {
    return <DemoCard onCancel={this.onDataFormCancel} onSubmit={this.onDataFormSubmit} />;
  }

  renderSearchForm() {
    const {
      form: { getFieldDecorator },
    } = this.props;
    return (
      <Form onSubmit={this.onSearchFormSubmit} layout="inline">
        <Row gutter={16}>
          <Col md={8} sm={24}>
            <Form.Item label="编号">
              {getFieldDecorator('code')(<Input placeholder="请输入" />)}
            </Form.Item>
          </Col>
          <Col md={8} sm={24}>
            <Form.Item label="名称">
              {getFieldDecorator('name')(<Input placeholder="请输入" />)}
            </Form.Item>
          </Col>
          <Col md={8} sm={24}>
            <Form.Item label="状态">
              {getFieldDecorator('status')(
                <Select placeholder="请选择" style={{ width: '100%' }}>
                  <Select.Option value="1">正常</Select.Option>
                  <Select.Option value="2">停用</Select.Option>
                </Select>
              )}
            </Form.Item>
          </Col>
        </Row>
        <div style={{ overflow: 'hidden' }}>
          <span style={{ float: 'right', marginBottom: 24 }}>
            <Button type="primary" htmlType="submit">
              查询
            </Button>
            <Button style={{ marginLeft: 8 }} onClick={this.onResetFormClick}>
              重置
            </Button>
          </span>
        </div>
      </Form>
    );
  }

  render() {
    const {
      loading,
      demo: {
        data: { list, pagination },
      },
    } = this.props;

    const { selectedRows, selectedRowKeys } = this.state;

    const columns = [
      {
        title: '编号',
        dataIndex: 'code',
      },
      {
        title: '名称',
        dataIndex: 'name',
      },
      {
        title: '备注',
        dataIndex: 'memo',
      },
      {
        title: '状态',
        dataIndex: 'status',
        render: val => {
          if (val === 1) {
            return <Badge status="success" text="启用" />;
          }
          return <Badge status="error" text="停用" />;
        },
      },
      {
        title: '创建时间',
        dataIndex: 'created_at',
        render: val => <span>{formatDate(val, 'YYYY-MM-DD HH:mm')}</span>,
      },
    ];

    const paginationProps = {
      showSizeChanger: true,
      showQuickJumper: true,
      showTotal: total => <span>共{total}条</span>,
      ...pagination,
    };

    const breadcrumbList = [{ title: '演示用例' }, { title: '基础示例', href: '/example/demo' }];

    return (
      <PageHeaderLayout title="基础示例" breadcrumbList={breadcrumbList}>
        <Card bordered={false}>
          <div className={styles.tableList}>
            <div className={styles.tableListForm}>{this.renderSearchForm()}</div>
            <div className={styles.tableListOperator}>
              <PButton code="add" icon="plus" type="primary" onClick={() => this.onAddClick()}>
                新建
              </PButton>
              {selectedRows.length === 1 && [
                <PButton
                  key="edit"
                  code="edit"
                  icon="edit"
                  onClick={() => this.onItemEditClick(selectedRows[0])}
                >
                  编辑
                </PButton>,
                <PButton
                  key="del"
                  code="del"
                  icon="delete"
                  type="danger"
                  onClick={() => this.onItemDelClick(selectedRows[0])}
                >
                  删除
                </PButton>,
                selectedRows[0].status === 2 && (
                  <PButton
                    key="enable"
                    code="enable"
                    icon="check"
                    onClick={() => this.onItemEnableClick(selectedRows[0])}
                  >
                    启用
                  </PButton>
                ),
                selectedRows[0].status === 1 && (
                  <PButton
                    key="disable"
                    code="disable"
                    icon="stop"
                    type="danger"
                    onClick={() => this.onItemDisableClick(selectedRows[0])}
                  >
                    禁用
                  </PButton>
                ),
              ]}
            </div>
            <div>
              <Table
                rowSelection={{
                  selectedRowKeys,
                  onChange: this.onTableSelectRow,
                }}
                loading={loading}
                rowKey={record => record.record_id}
                dataSource={list}
                columns={columns}
                pagination={paginationProps}
                onChange={this.onTableChange}
                size="small"
              />
            </div>
          </div>
        </Card>
        {this.renderDataForm()}
      </PageHeaderLayout>
    );
  }
}

export default DemoList;
