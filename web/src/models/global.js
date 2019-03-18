import * as loginService from '@/services/login';

export default {
  namespace: 'global',

  state: {
    title: '权限管理脚手架',
    copyRight: '2019 LyricTian',
    defaultURL: '/dashboard',
    collapsed: false,
    openKeys: [],
    selectedKeys: [],
    user: {
      user_name: 'admin',
      real_name: '管理员',
      role_names: [],
    },
    menuPaths: {},
    menus: [],
  },

  effects: {
    *menuEvent({ pathname }, { put, select }) {
      let p = pathname;
      if (p === '/') {
        p = yield select(state => state.global.defaultURL);
      }

      const menuPaths = yield select(state => state.global.menuPaths);
      const item = menuPaths[p];
      if (!item) {
        return;
      }

      if (item.parent_path && item.parent_path !== '') {
        yield put({
          type: 'changeOpenKeys',
          payload: item.parent_path.split('/'),
        });
      }

      yield put({
        type: 'changeSelectedKeys',
        payload: [item.record_id],
      });
    },
    *fetchUser(_, { call, put }) {
      const response = yield call(loginService.getCurrentUser);
      yield put({
        type: 'saveUser',
        payload: response,
      });
    },
    *fetchMenuTree({ pathname }, { call, put }) {
      const response = yield call(loginService.queryMenuTree);

      const menuData = response.list || [];
      yield put({
        type: 'saveMenus',
        payload: menuData,
      });

      const menuPaths = {};
      function findPath(data) {
        for (let i = 0; i < data.length; i += 1) {
          if (data[i].router !== '') {
            menuPaths[data[i].router] = data[i];
          }
          if (data[i].children && data[i].children.length > 0) {
            findPath(data[i].children);
          }
        }
      }
      findPath(menuData);

      yield put({
        type: 'saveMenuPaths',
        payload: menuPaths,
      });

      yield put({
        type: 'menuEvent',
        pathname,
      });
    },
  },

  reducers: {
    changeLayoutCollapsed(state, { payload }) {
      return {
        ...state,
        collapsed: payload,
      };
    },
    changeOpenKeys(state, { payload }) {
      return {
        ...state,
        openKeys: payload,
      };
    },
    changeSelectedKeys(state, { payload }) {
      return {
        ...state,
        selectedKeys: payload,
      };
    },
    saveUser(state, { payload }) {
      return { ...state, user: payload };
    },
    saveMenuPaths(state, { payload }) {
      return { ...state, menuPaths: payload };
    },
    saveMenus(state, { payload }) {
      return { ...state, menus: payload };
    },
  },
  subscriptions: {
    setup({ dispatch, history }) {
      history.listen(({ pathname }) => {
        dispatch({
          type: 'menuEvent',
          pathname,
        });
      });
    },
  },
};
