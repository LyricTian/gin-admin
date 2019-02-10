import React, { PureComponent } from 'react';
import { connect } from 'dva';
import { Row, Col, Card, Form, Input, Button, Table, Modal, Icon, Dropdown, Menu } from 'antd';
import PageHeaderLayout from '@/layouts/PageHeaderLayout';
import ResourceCard from './ResourceCard';

import styles from './ResourceList.less';

@connect(state => ({
  loading: state.loading.models.resource,
  resource: state.resource,
}))
@Form.create()
class ResourceList extends PureComponent {
  state = {
    selectedRows: [],
  };

  componentDidMount() {
    this.dispatch({
      type: 'resource/fetch',
      search: {},
      pagination: {},
    });
  }

  onDelBatchOKClick = () => {
    const { selectedRows } = this.state;
    if (selectedRows.length === 0) {
      return;
    }
    this.setState({
      selectedRows: [],
    });
    this.dispatch({
      type: 'resource/delMany',
      payload: { batch: selectedRows.join(',') },
    });
  };

  onBatchDelClick = () => {
    Modal.confirm({
      title: '确认删除选中的数据吗？',
      okText: '确认',
      okType: 'danger',
      cancelText: '取消',
      onOk: this.onDelBatchOKClick.bind(this),
    });
  };

  onItemEditClick = id => {
    this.dispatch({
      type: 'resource/loadForm',
      payload: {
        type: 'E',
        id,
      },
    });
  };

  onAddClick = () => {
    this.dispatch({
      type: 'resource/loadForm',
      payload: {
        type: 'A',
      },
    });
  };

  onDelOKClick(id) {
    this.dispatch({
      type: 'resource/del',
      payload: { record_id: id },
    });
  }

  onItemDelClick = item => {
    Modal.confirm({
      title: `确定删除【资源数据：${item.name}】？`,
      okText: '确认',
      okType: 'danger',
      cancelText: '取消',
      onOk: this.onDelOKClick.bind(this, item.record_id),
    });
  };

  onTableSelectRow = rows => {
    this.setState({
      selectedRows: rows,
    });
  };

  onTableChange = pagination => {
    this.dispatch({
      type: 'resource/fetch',
      pagination: {
        current: pagination.current,
        pageSize: pagination.pageSize,
      },
    });
  };

  onResetFormClick = () => {
    const { form } = this.props;
    form.resetFields();
    this.dispatch({
      type: 'resource/fetch',
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
      if (!err) {
        this.dispatch({
          type: 'resource/fetch',
          search: values,
          pagination: {},
        });
      }
    });
  };

  onDataFormSubmit = data => {
    this.dispatch({
      type: 'resource/submit',
      payload: data,
    });
  };

  onDataFormCancel = () => {
    this.dispatch({
      type: 'resource/changeFormVisible',
      payload: false,
    });
  };

  dispatch = action => {
    const { dispatch } = this.props;
    dispatch(action);
  };

  renderDataForm() {
    return <ResourceCard onCancel={this.onDataFormCancel} onSubmit={this.onDataFormSubmit} />;
  }

  renderSearchForm() {
    const {
      form: { getFieldDecorator },
    } = this.props;
    return (
      <Form onSubmit={this.onSearchFormSubmit} layout="inline">
        <Row gutter={16}>
          <Col md={8} sm={24}>
            <Form.Item label="名称">
              {getFieldDecorator('name')(<Input placeholder="请输入" />)}
            </Form.Item>
          </Col>
          <Col md={8} sm={24}>
            <Form.Item label="访问路径">
              {getFieldDecorator('path')(<Input placeholder="请输入" />)}
            </Form.Item>
          </Col>
          <Col md={8} sm={24}>
            <div style={{ overflow: 'hidden' }}>
              <span style={{ marginBottom: 24 }}>
                <Button type="primary" htmlType="submit">
                  查询
                </Button>
                <Button style={{ marginLeft: 8 }} onClick={this.onResetFormClick}>
                  重置
                </Button>
              </span>
            </div>
          </Col>
        </Row>
      </Form>
    );
  }

  render() {
    const {
      loading,
      resource: {
        data: { list, pagination },
      },
    } = this.props;

    const { selectedRows } = this.state;

    const columns = [
      {
        dataIndex: 'record_id',
        width: 80,
        render: (val, record) => (
          <div>
            {
              <Dropdown
                overlay={
                  <Menu>
                    <Menu.Item>
                      <a
                        onClick={() => {
                          this.onItemEditClick(val);
                        }}
                      >
                        编辑
                      </a>
                    </Menu.Item>
                    <Menu.Item>
                      <a
                        onClick={() => {
                          this.onItemDelClick(record);
                        }}
                      >
                        删除
                      </a>
                    </Menu.Item>
                  </Menu>
                }
              >
                <a>
                  操作 <Icon type="down" />
                </a>
              </Dropdown>
            }
          </div>
        ),
      },
      {
        title: '资源编号',
        dataIndex: 'code',
      },
      {
        title: '资源名称',
        dataIndex: 'name',
      },
      {
        title: '访问路径',
        dataIndex: 'path',
      },
      {
        title: '请求方式',
        dataIndex: 'method',
      },
    ];

    const paginationProps = {
      showSizeChanger: true,
      showQuickJumper: true,
      showTotal: total => <span>共{total}条</span>,
      ...pagination,
    };

    const breadcrumbList = [{ title: '系统管理' }, { title: '资源管理', href: '/system/resource' }];

    return (
      <PageHeaderLayout title="资源管理" breadcrumbList={breadcrumbList}>
        <Card bordered={false}>
          <div className={styles.tableList}>
            <div className={styles.tableListForm}>{this.renderSearchForm()}</div>
            <div className={styles.tableListOperator}>
              <Button icon="plus" type="primary" onClick={() => this.onAddClick()}>
                新建
              </Button>
              {selectedRows.length > 0 && (
                <span>
                  <Button icon="delete" type="danger" onClick={() => this.onBatchDelClick()}>
                    删除
                  </Button>
                </span>
              )}
            </div>
            <div>
              <Table
                rowSelection={{
                  selectedRowKeys: selectedRows,
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

export default ResourceList;
