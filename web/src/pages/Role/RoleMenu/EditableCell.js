import React, { PureComponent } from 'react';
import { Checkbox, Row, Col } from 'antd';

export default class EditableCell extends PureComponent {
  findItem = () => {
    const { data, record } = this.props;
    for (let i = 0; i < data.length; i += 1) {
      if (data[i].menu_id === record.record_id) {
        return data[i];
      }
    }
    return null;
  };

  handleChange = values => {
    const { record, dataIndex, handleSave } = this.props;
    handleSave(record, dataIndex, values);
  };

  renderAction = () => {
    const { record } = this.props;
    if (!record.actions || record.actions.length === 0) {
      return null;
    }

    const item = this.findItem();
    return (
      <Checkbox.Group
        disabled={!item}
        value={item ? item.actions : []}
        onChange={this.handleChange}
      >
        <Row>
          {record.actions.map(v => (
            <Col key={v.code}>
              <Checkbox value={v.code}>{v.name}</Checkbox>
            </Col>
          ))}
        </Row>
      </Checkbox.Group>
    );
  };

  renderResource = () => {
    const { record } = this.props;
    if (!record.resources || record.resources.length === 0) {
      return null;
    }

    const item = this.findItem();
    return (
      <Checkbox.Group
        disabled={!item}
        value={item ? item.resources : []}
        onChange={this.handleChange}
      >
        <Row>
          {record.resources.map(v => (
            <Col key={v.code}>
              <Checkbox value={v.code}>{v.name}</Checkbox>
            </Col>
          ))}
        </Row>
      </Checkbox.Group>
    );
  };

  render() {
    const { dataIndex, record, menuData, handleSave, ...restProps } = this.props;
    return (
      <td {...restProps}>
        {dataIndex === 'actions' && this.renderAction()}
        {dataIndex === 'resources' && this.renderResource()}
        {!(dataIndex === 'actions' || dataIndex === 'resources') && restProps.children}
      </td>
    );
  }
}
