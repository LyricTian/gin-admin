import { getLevelCode, getMenuKeys } from '../utils/utils';
import * as loginService from '../services/login';

export default {
  namespace: 'global',

  state: {
    title: 'RBAC权限管理平台',
    copyRight: '2018 LyricTian',
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
      const p = pathname;
      if (p === '/') {
        return;
      }

      const menuPaths = yield select(state => state.global.menuPaths);
      const menus = yield select(state => state.global.menus);
      const keys = getMenuKeys(p, menuPaths, menus);

      if (keys.length > 0) {
        yield put({
          type: 'changeOpenKeys',
          payload: keys.slice(0, keys.length - 1),
        });
      }

      const levelCode = getLevelCode(p, menuPaths);
      yield put({
        type: 'changeSelectedKeys',
        payload: [levelCode],
      });
    },
    *fetchUser(_, { call, put }) {
      const response = yield call(loginService.getCurrentUser);
      yield put({
        type: 'saveUser',
        payload: response,
      });
    },
    *fetchMenus({ pathname }, { call, put }) {
      const response = yield call(loginService.queryCurrentMenus);

      const menuData = response.list || [];
      yield put({
        type: 'saveMenus',
        payload: menuData,
      });

      const menuPaths = {};
      function findPath(data) {
        for (let i = 0; i < data.length; i += 1) {
          if (data[i].path !== '') {
            menuPaths[data[i].path] = data[i];
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
