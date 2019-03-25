import React, { PureComponent } from 'react';
import { Table } from 'antd';
import EditableCell from './EditableCell';

import * as menuService from '@/services/menu';

export default class RoleMenu extends PureComponent {
  constructor(props) {
    super(props);

    this.columns = [
      {
        title: '菜单名称',
        dataIndex: 'name',
        width: '30%',
      },
      {
        title: '动作权限',
        dataIndex: 'actions',
        editable: true,
        width: '30%',
      },
      {
        title: '资源权限',
        dataIndex: 'resources',
        editable: true,
        width: '40%',
      },
    ];

    this.state = {
      menuData: [],
      dataSource: props.value || [],
    };
  }

  componentDidMount() {
    menuService.query({ q: 'tree', include_actions: '1', include_resources: '1' }).then(data => {
      const list = data.list || [];
      this.setState({ menuData: this.fillData(list) });
    });
  }

  static getDerivedStateFromProps(nextProps, state) {
    if ('value' in nextProps) {
      return {
        ...state,
        dataSource: nextProps.value || [],
      };
    }
    return state;
  }

  fillData = data => {
    const newData = [...data];
    for (let i = 0; i < newData.length; i += 1) {
      const { children } = newData[i];
      const item = { ...newData[i], hasChild: children && children.length > 0 };
      if (item.hasChild) {
        item.children = this.fillData(children);
      }
      newData[i] = item;
    }
    return newData;
  };

  handleSave = (record, dataIndex, values) => {
    const { dataSource } = this.state;
    const data = [...dataSource];
    const index = data.findIndex(item => item.menu_id === record.record_id);
    let item = data[index];
    if (!item) {
      item = {
        menu_id: record.record_id,
        dataIndex: values,
      };
    } else {
      item[dataIndex] = values;
    }
    data.splice(index, 1, {
      ...item,
    });
    this.setState({ dataSource: data }, () => {
      this.triggerChange(data);
    });
  };

  triggerChange = data => {
    const { onChange } = this.props;
    if (onChange) {
      onChange(data);
    }
  };

  handleSelectedRow = (_, rows) => {
    const { dataSource } = this.state;
    const list = [];

    for (let i = 0; i < rows.length; i += 1) {
      let exists = false;
      for (let j = 0; j < dataSource.length; j += 1) {
        if (dataSource[j].menu_id === rows[i].record_id) {
          exists = true;
          list.push({ ...dataSource[j] });
          break;
        }
      }

      if (!exists) {
        list.push({
          menu_id: rows[i].record_id,
          actions: rows[i].actions ? rows[i].actions.map(v => v.code) : [],
          resources: rows[i].resources ? rows[i].resources.map(v => v.code) : [],
        });
      }
    }

    this.setState({ dataSource: list }, () => {
      this.triggerChange(list);
    });
  };

  render() {
    const { dataSource, menuData } = this.state;
    const components = {
      body: {
        cell: EditableCell,
      },
    };
    const columns = this.columns.map(col => {
      if (!col.editable) {
        return col;
      }
      return {
        ...col,
        onCell: record => ({
          record,
          data: dataSource,
          dataIndex: col.dataIndex,
          handleSave: this.handleSave,
        }),
      };
    });

    return (
      menuData.length > 0 && (
        <Table
          bordered
          defaultExpandAllRows
          rowSelection={{
            selectedRowKeys: dataSource.map(v => v.menu_id),
            onChange: this.handleSelectedRow,
            getCheckboxProps: record => ({
              disabled: record.hasChild,
            }),
          }}
          rowKey={record => record.record_id}
          components={components}
          dataSource={menuData}
          columns={columns}
          pagination={false}
        />
      )
    );
  }
}
