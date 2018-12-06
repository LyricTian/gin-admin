import { message } from 'antd';
import * as menuService from '../services/menu';

export default {
  namespace: 'menu',
  state: {
    search: {},
    pagination: {},
    data: {
      list: [],
      pagination: {},
    },
    submitting: false,
    formType: '',
    formTitle: '',
    formID: '',
    formVisible: false,
    formData: {},
    searchTreeData: [],
    treeData: [],
    expandedKeys: [],
  },
  effects: {
    *fetch({ search, pagination }, { call, put, select }) {
      let params = {
        type: 'page',
      };

      if (search) {
        params = { ...params, ...search };
        yield put({
          type: 'saveSearch',
          payload: search,
        });
      } else {
        const s = yield select(state => state.menu.search);
        if (s) {
          params = { ...params, ...s };
        }
      }

      if (pagination) {
        params = { ...params, ...pagination };
        yield put({
          type: 'savePagination',
          payload: pagination,
        });
      } else {
        const p = yield select(state => state.menu.pagination);
        if (p) {
          params = { ...params, ...p };
        }
      }

      const response = yield call(menuService.query, params);
      yield put({
        type: 'saveData',
        payload: response,
      });
    },
    *fetchSearchTree({ payload }, { call, put, select }) {
      let params = {
        type: 'tree',
        is_menu: 1,
      };
      if (payload) {
        params = { ...params, ...payload };
      }
      const response = yield call(menuService.query, params);
      const list = response.list || [];
      yield put({
        type: 'saveSearchTreeData',
        payload: list,
      });

      const expandedKeys = yield select(state => state.menu.expandedKeys);
      if (expandedKeys.length === 0 && list.length > 0) {
        yield put({
          type: 'saveExpandedKeys',
          payload: [list[0].record_id],
        });
      }
    },
    *loadForm({ payload }, { call, put, select }) {
      yield put({
        type: 'changeFormVisible',
        payload: true,
      });

      yield [
        put({
          type: 'saveFormType',
          payload: payload.type,
        }),
        put({
          type: 'saveFormTitle',
          payload: '新建菜单',
        }),
        put({
          type: 'saveFormID',
          payload: '',
        }),
        put({
          type: 'saveFormData',
          payload: {},
        }),
        put({
          type: 'fetchTree',
          payload: { is_menu: 1 },
        }),
      ];

      if (payload.type === 'A') {
        const search = yield select(state => state.menu.search);
        if (search.parent_id) {
          const { parent_id: parentID } = search;
          const response = yield call(menuService.get, { record_id: parentID });
          yield put({
            type: 'saveFormData',
            payload: {
              parent_id: parentID,
              type: response.type + 10,
            },
          });
        }
      }

      if (payload.type === 'E') {
        yield [
          put({
            type: 'saveFormTitle',
            payload: '编辑菜单',
          }),
          put({
            type: 'saveFormID',
            payload: payload.id,
          }),
          put({
            type: 'fetchForm',
            payload: { record_id: payload.id },
          }),
        ];
      }
    },
    *fetchForm({ payload }, { call, put }) {
      const response = yield call(menuService.get, payload);
      yield put({
        type: 'saveFormData',
        payload: response,
      });
    },
    *submit({ payload }, { call, put, select }) {
      yield put({
        type: 'changeSubmitting',
        payload: true,
      });

      const params = { ...payload };
      const { type } = params;

      let valid = true;
      if (type === 10) {
        if (params.parent_id !== '') {
          message.error('系统不能包含上级菜单');
          valid = false;
        }
      } else if (params.parent_id === '') {
        message.error('请选择上级菜单');
        valid = false;
      } else {
        const parent = yield call(menuService.get, { record_id: params.parent_id });
        const { type: ptype } = parent;
        if (type === 20 && !(ptype === 10 || ptype === 20)) {
          message.error('模块依赖于系统或模块');
          valid = false;
        } else if (type === 30 && !(ptype === 10 || ptype === 20)) {
          message.error('功能依赖于系统或模块');
          valid = false;
        } else if (type === 40 && ptype !== 30) {
          message.error('资源依赖于功能');
          valid = false;
        }
      }

      if (!valid) {
        yield put({
          type: 'changeSubmitting',
          payload: false,
        });
        return;
      }

      const formType = yield select(state => state.menu.formType);
      const formID = yield select(state => state.menu.formID);
      let success = false;
      if (formType === 'E') {
        params.record_id = formID;
        const response = yield call(menuService.update, params);
        if (response.status === 'OK') {
          success = true;
        }
      } else {
        const response = yield call(menuService.create, params);
        if (response.record_id && response.record_id !== '') {
          success = true;
        }
      }

      yield put({
        type: 'changeSubmitting',
        payload: false,
      });

      if (success) {
        message.success('保存成功');
        yield put({
          type: 'changeFormVisible',
          payload: false,
        });

        yield put({ type: 'fetchSearchTree' });
        yield put({ type: 'fetch' });
      }
    },
    *del({ payload }, { call, put }) {
      const response = yield call(menuService.del, payload);
      if (response.status === 'OK') {
        message.success('删除成功');
        yield put({ type: 'fetchSearchTree' });
        yield put({ type: 'fetch' });
      }
    },
    *delMany({ payload }, { call, put }) {
      const response = yield call(menuService.delMany, payload);
      if (response.status === 'OK') {
        message.success('删除成功');
        yield put({ type: 'fetchSearchTree' });
        yield put({ type: 'fetch' });
      }
    },
    *changeStatus({ payload }, { call, put, select }) {
      let response;
      if (payload.status === 1) {
        response = yield call(menuService.enable, payload);
      } else {
        response = yield call(menuService.disable, payload);
      }

      if (response.status === 'OK') {
        let msg = '启用成功';
        if (payload.status === 2) {
          msg = '停用成功';
        }
        message.success(msg);
        const data = yield select(state => state.menu.data);
        const newData = { list: [], pagination: data.pagination };

        for (let i = 0; i < data.list.length; i += 1) {
          const item = data.list[i];
          if (item.record_id === payload.record_id) {
            item.status = payload.status;
          }
          newData.list.push(item);
        }

        yield put({
          type: 'saveData',
          payload: newData,
        });
      }
    },
    *fetchTree({ payload }, { call, put }) {
      let params = {
        type: 'tree',
        status: 1,
      };
      if (payload) {
        params = { ...params, ...payload };
      }
      const response = yield call(menuService.query, params);
      yield put({
        type: 'saveTreeData',
        payload: response.list || [],
      });
    },
  },
  reducers: {
    saveData(state, { payload }) {
      return { ...state, data: payload };
    },
    saveSearch(state, { payload }) {
      return { ...state, search: payload };
    },
    savePagination(state, { payload }) {
      return { ...state, pagination: payload };
    },
    changeFormVisible(state, { payload }) {
      return { ...state, formVisible: payload };
    },
    saveFormType(state, { payload }) {
      return { ...state, formType: payload };
    },
    saveFormTitle(state, { payload }) {
      return { ...state, formTitle: payload };
    },
    saveFormID(state, { payload }) {
      return { ...state, formID: payload };
    },
    saveFormData(state, { payload }) {
      return { ...state, formData: payload };
    },
    changeSubmitting(state, { payload }) {
      return { ...state, submitting: payload };
    },
    saveSearchTreeData(state, { payload }) {
      return { ...state, searchTreeData: payload };
    },
    saveTreeData(state, { payload }) {
      return { ...state, treeData: payload };
    },
    saveExpandedKeys(state, { payload }) {
      return { ...state, expandedKeys: payload };
    },
  },
};
