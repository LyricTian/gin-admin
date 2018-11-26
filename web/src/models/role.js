import { message } from 'antd';
import * as roleService from '../services/role';

export default {
  namespace: 'role',
  state: {
    search: {},
    pagination: {},
    data: {
      list: [],
      pagination: {},
    },
    submitting: false,
    formTitle: '',
    formID: '',
    formVisible: false,
    formData: {},
    selectData: [],
    menuKeys: [],
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
        const s = yield select(state => state.role.search);
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
        const p = yield select(state => state.role.pagination);
        if (p) {
          params = { ...params, ...p };
        }
      }

      const response = yield call(roleService.query, params);
      yield put({
        type: 'saveData',
        payload: response,
      });
    },
    *loadForm({ payload }, { put }) {
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
          payload: '新建角色',
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
          type: 'saveMenuKeys',
          payload: [],
        }),
      ];

      if (payload.type === 'E') {
        yield [
          put({
            type: 'saveFormTitle',
            payload: '编辑角色',
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
      const response = yield call(roleService.get, payload);
      yield [
        put({
          type: 'saveFormData',
          payload: response,
        }),
        put({
          type: 'saveMenuKeys',
          payload: response.menu_ids || [],
        }),
      ];
    },
    *submit({ payload }, { call, put, select }) {
      yield put({
        type: 'changeSubmitting',
        payload: true,
      });

      const params = { ...payload };
      const formType = yield select(state => state.role.formType);
      const formID = yield select(state => state.role.formID);
      let success = false;
      if (formType === 'E') {
        params.record_id = formID;
        const response = yield call(roleService.update, params);
        if (response.status === 'OK') {
          success = true;
        }
      } else {
        const response = yield call(roleService.create, params);
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
        yield put({
          type: 'fetch',
        });
      }
    },
    *del({ payload }, { call, put }) {
      const response = yield call(roleService.del, payload);
      if (response.status === 'OK') {
        message.success('删除成功');
        yield put({ type: 'fetch' });
      }
    },
    *delMany({ payload }, { call, put }) {
      const response = yield call(roleService.delMany, payload);
      if (response.status === 'OK') {
        message.success('删除成功');
        yield put({ type: 'fetch' });
      }
    },
    *changeStatus({ payload }, { call, put, select }) {
      let response;
      if (payload.status === 1) {
        response = yield call(roleService.enable, payload);
      } else {
        response = yield call(roleService.disable, payload);
      }

      if (response.status === 'OK') {
        let msg = '启用成功';
        if (payload.status === 2) {
          msg = '停用成功';
        }
        message.success(msg);
        const data = yield select(state => state.role.data);
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
    *fetchSelect(_, { call, put }) {
      const params = {
        type: 'select',
        status: 1,
      };
      const response = yield call(roleService.query, params);
      yield put({
        type: 'saveSelectData',
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
    saveFormTitle(state, { payload }) {
      return { ...state, formTitle: payload };
    },
    saveFormType(state, { payload }) {
      return { ...state, formType: payload };
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
    saveSelectData(state, { payload }) {
      return { ...state, selectData: payload };
    },
    saveMenuKeys(state, { payload }) {
      return { ...state, menuKeys: payload };
    },
  },
};
