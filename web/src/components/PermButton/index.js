import React from 'react';
import { Button } from 'antd';
import GetGlobalContext from '@/utils/context';

export default ({ code, children, ...rest }) => {
  const GlobalContext = GetGlobalContext();
  return (
    <GlobalContext.Consumer>
      {global => {
        const { menuPaths } = global;
        if (menuPaths) {
          const item = menuPaths[window.location.pathname];
          if (item && item.actions) {
            const { actions } = item;
            for (let i = 0; i < actions.length; i += 1) {
              if (actions[i].code === code) {
                return <Button {...rest}>{children}</Button>;
              }
            }
          }
        }
        return null;
      }}
    </GlobalContext.Consumer>
  );
};
