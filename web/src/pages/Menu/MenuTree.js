import React, { PureComponent } from 'react';
import { connect } from 'dva';
import { Tree } from 'antd';

@connect(state => ({
  menu: state.menu,
}))
class MenuTree extends PureComponent {
  componentDidMount() {
    this.dispatch({
      type: 'menu/fetchTree',
      payload: { include_resource: 1 },
    });
  }

  dispatch = action => {
    const { dispatch } = this.props;
    dispatch(action);
  };

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

  render() {
    const {
      treeProps,
      menu: { treeData },
    } = this.props;

    const props = treeProps || {};

    return (
      <div>
        {treeData.length > 0 && (
          <Tree checkable showLine defaultExpandAll {...props}>
            {this.renderTreeNodes(treeData)}
          </Tree>
        )}
      </div>
    );
  }
}

export default MenuTree;
