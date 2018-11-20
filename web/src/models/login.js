import { routerRedux } from 'dva/router';
import { stringify, parse } from 'qs';
import { storeLogoutKey } from '../utils/utils';
import * as loginService from '../services/login';

export default {
  namespace: 'login',

  state: {
    status: undefined,
    submitting: false,
  },

  effects: {
    *submit({ payload }, { call, put }) {
      yield put({
        type: 'changeSubmitting',
        payload: true,
      });
      const response = yield call(loginService.login, payload);
      yield put({
        type: 'changeLoginStatus',
        payload: response.status,
      });
      yield put({
        type: 'changeSubmitting',
        payload: false,
      });

      if (response.status === 'ok') {
        sessionStorage.removeItem(storeLogoutKey);
        const params = parse(window.location.href.split('?')[1]);
        const { redirect } = params;
        if (redirect) {
          window.location.href = redirect;
          return;
        }
        yield put(routerRedux.replace('/'));
      }
    },
    *logout(_, { put, call }) {
      yield put({
        type: 'changeLoginStatus',
        payload: false,
      });
      const response = yield call(loginService.logout);
      if (response === 'ok') {
        if (sessionStorage.getItem(storeLogoutKey) === '1') {
          return;
        }
        sessionStorage.setItem(storeLogoutKey, '1');

        yield put(
          routerRedux.push({
            pathname: '/user/login',
            search: stringify({
              redirect: window.location.href,
            }),
          })
        );
      }
    },
  },

  reducers: {
    changeLoginStatus(state, { payload }) {
      return {
        ...state,
        status: payload,
      };
    },
    changeSubmitting(state, { payload }) {
      return {
        ...state,
        submitting: payload,
      };
    },
  },
};
