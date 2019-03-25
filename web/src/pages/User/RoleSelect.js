import React, { PureComponent } from 'react';
import { Select } from 'antd';
import { query } from '@/services/role';

function parseValue(value) {
  if (!value) {
    return [];
  }
  return value.map(v => v.role_id);
}

export default class RoleSelect extends PureComponent {
  constructor(props) {
    super(props);

    this.state = {
      value: parseValue(props.value),
      data: [],
    };
  }

  componentDidMount() {
    query({ q: 'select' }).then(data => {
      this.setState({ data: data.list || [] });
    });
  }

  static getDerivedStateFromProps(nextProps, state) {
    if ('value' in nextProps) {
      return { ...state, value: parseValue(nextProps.value) };
    }
    return state;
  }

  handleChange = value => {
    this.setState({ value });
    this.triggerChange(value);
  };

  triggerChange = data => {
    const { onChange } = this.props;
    if (onChange) {
      const newData = data.map(v => ({ role_id: v }));
      onChange(newData);
    }
  };

  render() {
    const { value, data } = this.state;

    return (
      <Select
        mode="tags"
        value={value}
        onChange={this.handleChange}
        placeholder="请选择角色"
        style={{ width: '100%' }}
      >
        {data.map(item => (
          <Select.Option key={item.record_id} value={item.record_id}>
            {item.name}
          </Select.Option>
        ))}
      </Select>
    );
  }
}
