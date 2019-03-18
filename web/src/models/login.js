import { routerRedux } from 'dva/router';
import { stringify, parse } from 'qs';
import store from '@/utils/store';
import * as loginService from '@/services/login';

let isLogout = false;

export default {
  namespace: 'login',

  state: {
    status: '',
    tip: '',
    submitting: false,
    captchaID: '',
    captcha: '',
  },

  effects: {
    *loadCaptcha(_, { call, put }) {
      const response = yield call(loginService.captchaID);
      const { captcha_id: captchaID } = response;

      yield put({
        type: 'saveCaptchaID',
        payload: captchaID,
      });
      yield put({
        type: 'saveCaptcha',
        payload: loginService.captcha(captchaID),
      });
    },
    *reloadCaptcha(_, { put, select }) {
      const captchaID = yield select(state => state.login.captchaID);
      yield put({
        type: 'saveCaptcha',
        payload: `${loginService.captcha(captchaID)}&reload=${Math.random()}`,
      });
    },
    *submit({ payload }, { call, put }) {
      yield put({
        type: 'changeSubmitting',
        payload: true,
      });
      const response = yield call(loginService.login, payload);
      if (response.error) {
        const { message } = response.error;
        yield [
          put({
            type: 'saveTip',
            payload: message,
          }),
          put({
            type: 'saveStatus',
            payload: response.status >= 500 ? 'ERROR' : 'FAIL',
          }),
        ];
        yield put({
          type: 'changeSubmitting',
          payload: false,
        });
        yield put({
          type: 'loadCaptcha',
        });
        return;
      }

      // 保存访问令牌
      store.setAccessToken(response);
      yield put({
        type: 'changeSubmitting',
        payload: false,
      });

      isLogout = false;
      const params = parse(window.location.href.split('?')[1]);
      const { redirect } = params;
      if (redirect) {
        window.location.href = redirect;
        return;
      }
      yield put(routerRedux.replace('/'));
    },
    *logout(_, { put, call }) {
      if (isLogout) {
        return;
      }
      isLogout = true;

      const response = yield call(loginService.logout);
      if (response.status === 'OK') {
        yield put(
          routerRedux.push({
            pathname: '/user/login',
            search: stringify({
              redirect: window.location.href,
            }),
          })
        );
      }
      store.clearAccessToken();
    },
  },

  reducers: {
    saveCaptchaID(state, { payload }) {
      return {
        ...state,
        captchaID: payload,
      };
    },
    saveCaptcha(state, { payload }) {
      return {
        ...state,
        captcha: payload,
      };
    },
    saveStatus(state, { payload }) {
      return {
        ...state,
        status: payload,
      };
    },
    saveTip(state, { payload }) {
      return {
        ...state,
        tip: payload,
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
