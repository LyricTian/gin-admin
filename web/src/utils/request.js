import axios from 'axios';
import { notification } from 'antd';
import { storeAccessTokenKey } from './utils';

// 提供API前缀
export const v1API = '/api/v1';
const tokenKey = 'access-token';

function handle(showNotify) {
  return response => {
    const { status, data, headers } = response;

    const accessToken = headers[tokenKey];
    if (accessToken) {
      sessionStorage.setItem(storeAccessTokenKey, accessToken);
    }

    if (status >= 200 && status < 300) {
      return data;
    }

    if (status === 401) {
      /* eslint-disable no-underscore-dangle */
      window.g_app._store.dispatch({ type: 'login/logout' });
      return {};
    }

    const error = {
      code: 0,
      message: '服务器发生错误',
    };
    if (status === 504) {
      error.message = '未连接到服务器';
    } else if (data) {
      const {
        error: { message, code },
      } = data;
      error.message = message;
      error.code = code;
    } else if (status >= 400 && status < 500) {
      error.message = '请求发生错误';
    }

    if (showNotify) {
      notification.error({
        message: error.message,
      });
    }

    return { error, status };
  };
}

export default async function request(url, options) {
  const defaultHeader = {};
  defaultHeader[tokenKey] = sessionStorage.getItem(storeAccessTokenKey) || '';

  let showNotify = true;
  const newOptions = {
    url,
    validateStatus() {
      return true;
    },
    ...options,
  };

  if (newOptions.notNotify) {
    showNotify = false;
  }

  if (newOptions.method === 'POST' || newOptions.method === 'PUT') {
    defaultHeader['Content-Type'] = 'application/json; charset=utf-8';
    newOptions.data = newOptions.body;
  }
  newOptions.headers = { ...defaultHeader, ...newOptions.headers };

  return axios.request(newOptions).then(handle(showNotify));
}
