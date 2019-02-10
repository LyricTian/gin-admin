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

@connect(({ menu, loading }) => ({
  menu,
  loading: loading.models.menu,
}))
@Form.create()
class MenuList extends PureComponent {
  state = {
    selectedRows: [],
    treeSelectedKeys: [],
    selectedType: 0,
  };

  componentDidMount() {
    this.dispatch({
      type: 'menu/fetchTree',
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
            ...values,
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
            <Form.Item label="菜单编号">
              {getFieldDecorator('code')(<Input placeholder="请输入" />)}
            </Form.Item>
          </Col>
          <Col md={8} sm={24}>
            <Form.Item label="菜单名称">
              {getFieldDecorator('name')(<Input placeholder="请输入" />)}
            </Form.Item>
          </Col>
          <Col md={8} sm={24}>
            <Form.Item label="菜单类型">
              {getFieldDecorator('type')(
                <Select placeholder="请选择" style={{ width: '100%' }}>
                  <Select.Option value="1">模块</Select.Option>
                  <Select.Option value="2">功能</Select.Option>
                  <Select.Option value="3">资源</Select.Option>
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
        treeData,
        expandedKeys,
      },
    } = this.props;

    const { selectedRows, selectedType } = this.state;

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
        width: 180,
      },
      {
        title: '菜单编号',
        dataIndex: 'code',
        width: 150,
      },
      {
        title: '菜单类型',
        dataIndex: 'type',
        width: 100,
        render: val => {
          let v = '';
          switch (val) {
            case 1:
              v = '模块';
              break;
            case 2:
              v = '功能';
              break;
            case 3:
              v = '资源';
              break;
            default:
              v = '-';
              break;
          }
          return <span>{v}</span>;
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
        title: '访问路径',
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

    const breadcrumbList = [{ title: '系统管理' }, { title: '菜单管理', href: '/system/menu' }];

    return (
      <PageHeaderLayout title="菜单管理" breadcrumbList={breadcrumbList}>
        <Layout>
          <Layout.Sider
            width={200}
            style={{ background: '#fff', borderRight: '1px solid lightGray' }}
          >
            <Tree
              expandedKeys={expandedKeys}
              onSelect={(keys, { selectedNodes }) => {
                this.setState({
                  treeSelectedKeys: keys,
                  selectedType: selectedNodes.length > 0 ? selectedNodes[0].props.dataRef.type : 0,
                });

                const {
                  menu: { search },
                } = this.props;

                const item = {
                  parent_id: '',
                };

                if (keys.length > 0) {
                  [item.parent_id] = keys;
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
              {this.renderTreeNodes(treeData)}
            </Tree>
          </Layout.Sider>
          <Layout.Content>
            <Card bordered={false}>
              <div className={styles.tableList}>
                <div className={styles.tableListForm}>{this.renderSearchForm()}</div>
                <div className={styles.tableListOperator}>
                  {selectedType === 2 && (
                    <Button icon="select" type="primary" onClick={() => this.onBatchDelClick()}>
                      选择资源
                    </Button>
                  )}
                  <Button icon="plus" type="primary" onClick={() => this.onAddClick()}>
                    新建
                  </Button>
                  {selectedRows.length > 0 && (
                    <Button icon="delete" type="danger" onClick={() => this.onBatchDelClick()}>
                      删除
                    </Button>
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
                  scroll={{ x: 890 }}
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
export default MenuList;
