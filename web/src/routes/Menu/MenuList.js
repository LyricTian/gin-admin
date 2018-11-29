import React, { PureComponent } from 'react';
import { connect } from 'dva';
import {
  Row,
  Col,
  Card,
  Form,
  Input,
  Button,
  Table,
  Badge,
  Select,
  Modal,
  Icon,
  Dropdown,
  Menu,
  Layout,
  Tree,
} from 'antd';
import PageHeaderLayout from '../../layouts/PageHeaderLayout';
import MenuCard from './MenuCard';

import styles from './MenuList.less';

@connect(state => ({
  menu: state.menu,
  loading: state.loading.models.menu,
}))
@Form.create()
export default class MenuList extends PureComponent {
  state = {
    selectedRows: [],
    treeSelectedKeys: [],
  };

  componentDidMount() {
    this.dispatch({
      type: 'menu/fetchSearchTree',
    });

    this.dispatch({
      type: 'menu/fetch',
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
      type: 'menu/delMany',
      payload: { batch: selectedRows.join(',') },
    });
  };

  onBatchDelClick = () => {
    Modal.confirm({
      title: '确认删除选中的菜单吗？',
      okText: '确认',
      okType: 'danger',
      cancelText: '取消',
      onOk: this.onDelBatchOKClick.bind(this),
    });
  };

  onItemDisableClick = id => {
    this.dispatch({
      type: 'menu/changeStatus',
      payload: { record_id: id, status: 2 },
    });
  };

  onItemEnableClick = id => {
    this.dispatch({
      type: 'menu/changeStatus',
      payload: { record_id: id, status: 1 },
    });
  };

  onItemEditClick = id => {
    this.dispatch({
      type: 'menu/loadForm',
      payload: {
        type: 'E',
        id,
      },
    });
  };

  onAddClick = () => {
    this.dispatch({
      type: 'menu/loadForm',
      payload: {
        type: 'A',
      },
    });
  };

  onDelOKClick(id) {
    this.dispatch({
      type: 'menu/del',
      payload: { record_id: id },
    });
  }

  onItemDelClick = item => {
    Modal.confirm({
      title: `确定删除【菜单数据：${item.name}】？`,
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
      type: 'menu/fetch',
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
      type: 'menu/fetch',
      search: { parent_id: this.getParentID() },
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
          type: 'menu/fetch',
          search: {
            name: values.name,
            status: values.status,
            mtype: values.type,
            parent_id: this.getParentID(),
          },
          pagination: {},
        });
      }
    });
  };

  onDataFormSubmit = data => {
    this.dispatch({
      type: 'menu/submit',
      payload: data,
    });
  };

  onDataFormCancel = () => {
    this.dispatch({
      type: 'menu/changeFormVisible',
      payload: false,
    });
  };

  getParentID = () => {
    const { treeSelectedKeys } = this.state;
    let parentID = '';
    if (treeSelectedKeys.length > 0) {
      [parentID] = treeSelectedKeys;
    }
    return parentID;
  };

  dispatch = action => {
    const { dispatch } = this.props;
    dispatch(action);
  };

  renderDataForm() {
    return <MenuCard onCancel={this.onDataFormCancel} onSubmit={this.onDataFormSubmit} />;
  }

  renderTreeNodes = data =>
    data.map(item => {
      if (item.children) {
        return (
          <Tree.TreeNode title={item.name} key={item.record_id} dataRef={item}>
            {this.renderTreeNodes(item.children)}
          </Tree.TreeNode>
        );
      }
      return <Tree.TreeNode title={item.name} key={item.record_id} dataRef={item} />;
    });

  renderSearchForm() {
    const {
      form: { getFieldDecorator },
    } = this.props;
    return (
      <Form onSubmit={this.onSearchFormSubmit} layout="inline">
        <Row gutter={16}>
          <Col md={8} sm={24}>
            <Form.Item label="菜单名称">
              {getFieldDecorator('name')(<Input placeholder="请输入" />)}
            </Form.Item>
          </Col>
          <Col md={8} sm={24}>
            <Form.Item label="菜单类型">
              {getFieldDecorator('type')(
                <Select placeholder="请选择" style={{ width: '100%' }}>
                  <Select.Option value="10">系统</Select.Option>
                  <Select.Option value="20">模块</Select.Option>
                  <Select.Option value="30">功能</Select.Option>
                  <Select.Option value="40">动作</Select.Option>
                </Select>
              )}
            </Form.Item>
          </Col>
          <Col md={8} sm={24}>
            <Form.Item label="菜单状态">
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
      menu: {
        data: { list, pagination },
        searchTreeData,
        expandedKeys,
      },
    } = this.props;

    const { selectedRows } = this.state;

    const columns = [
      {
        dataIndex: 'record_id',
        width: 80,
        fixed: 'left',
        render: (val, record) => (
          <div>
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
                    {record.status === 1 ? (
                      <a
                        onClick={() => {
                          this.onItemDisableClick(val);
                        }}
                      >
                        设置为停用
                      </a>
                    ) : (
                      <a
                        onClick={() => {
                          this.onItemEnableClick(val);
                        }}
                      >
                        设置为启用
                      </a>
                    )}
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
          </div>
        ),
      },
      {
        title: '菜单名称',
        dataIndex: 'name',
        width: 150,
      },
      {
        title: '菜单编号',
        dataIndex: 'code',
        width: 100,
      },
      {
        title: '菜单类型',
        dataIndex: 'type',
        width: 100,
        render: val => {
          let v = '';
          switch (val) {
            case 10:
              v = '系统';
              break;
            case 20:
              v = '模块';
              break;
            case 30:
              v = '功能';
              break;
            case 40:
              v = '动作';
              break;
            default:
              v = '-';
              break;
          }
          return <span>{v}</span>;
        },
      },
      {
        title: '菜单状态',
        dataIndex: 'status',
        width: 100,
        render: val => {
          if (val === 1) {
            return <Badge status="success" text="正常" />;
          }
          return <Badge status="error" text="停用" />;
        },
      },
      {
        title: '排序值',
        dataIndex: 'sequence',
        width: 80,
      },
      {
        title: '菜单图标',
        dataIndex: 'icon',
        width: 100,
      },
      {
        title: '跳转链接',
        dataIndex: 'path',
        width: 200,
      },
    ];

    const paginationProps = {
      showSizeChanger: true,
      showQuickJumper: true,
      showTotal: total => <span>共{total}条</span>,
      ...pagination,
    };

    return (
      <PageHeaderLayout title="菜单管理">
        <Layout>
          <Layout.Sider
            width={200}
            style={{ background: '#fff', borderRight: '1px solid lightGray' }}
          >
            <div>
              <Input.Search
                placeholder="请输入"
                onChange={e => {
                  const { value } = e.target;
                  this.dispatch({
                    type: 'menu/fetchSearchTree',
                    payload: { name: value },
                  });
                }}
              />
            </div>
            <Tree
              expandedKeys={expandedKeys}
              onSelect={(keys, { selectedNodes }) => {
                this.setState({ treeSelectedKeys: keys });

                const {
                  menu: { search },
                } = this.props;

                const item = {
                  parent_id: '',
                  parent_type: '',
                };

                if (keys.length > 0) {
                  [item.parent_id] = keys;
                  item.parent_type = selectedNodes[0].props.dataRef.type;
                }

                this.dispatch({
                  type: 'menu/fetch',
                  search: { ...search, ...item },
                  pagination: {},
                });
              }}
              onExpand={keys => {
                this.dispatch({
                  type: 'menu/saveExpandedKeys',
                  payload: keys,
                });
              }}
            >
              {this.renderTreeNodes(searchTreeData)}
            </Tree>
          </Layout.Sider>
          <Layout.Content>
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
                  scroll={{ x: 910 }}
                  size="small"
                />
              </div>
            </Card>
          </Layout.Content>
        </Layout>
        {this.renderDataForm()}
      </PageHeaderLayout>
    );
  }
}
