import axios from 'axios';
import { notification } from 'antd';
import store from '../index';
import { storeAccessTokenKey } from './utils';

// 默认前缀
export const baseURLV1 = '/api/v1';

function handle(url) {
  return response => {
    const { status, data, headers } = response;

    if (url === '/api/v1/login') {
      const { 'access-token': accessToken } = headers;
      if (accessToken) {
        localStorage.setItem(storeAccessTokenKey, accessToken);
      }
    }

    if (status >= 200 && status < 300) {
      return data;
    }

    if (status === 401) {
      store.dispatch({ type: 'login/logout' });
      return {};
    }

    let message = '服务器发生错误';
    if (status === 504) {
      message = '未连接到服务器';
    } else if (data) {
      const {
        error: { message: msg },
      } = data;
      message = msg;
    } else if (status >= 400 && status < 500) {
      message = '请求发生错误';
    }

    notification.error({
      message,
    });

    return {};
  };
}

export default async function request(url, options) {
  const defaultHeader = {
    'access-token': localStorage.getItem(storeAccessTokenKey) || '',
  };

  const newOptions = {
    url,
    validateStatus() {
      return true;
    },
    ...options,
  };
  if (newOptions.method === 'POST' || newOptions.method === 'PUT') {
    defaultHeader['Content-Type'] = 'application/json; charset=utf-8';
    newOptions.data = newOptions.body;
  }
  newOptions.headers = { ...defaultHeader, ...newOptions.headers };

  return axios.request(newOptions).then(handle(url));
}
