import fetch from 'dva/fetch';
import { notification } from 'antd';
import store from '../index';

// 定义路由前缀
export const baseURLV1 = '/api/v1';

function checkStatus(response) {
  const { status } = response;
  if (status >= 200 && status < 300) {
    return response;
  }

  if (status === 401) {
    store.dispatch({ type: 'login/logout' });
    return response;
  }

  if (status === 504) {
    throw new Error('未连接到服务器');
  }

  return response.json().then(body => {
    const {
      error: { message },
    } = body;
    const error = new Error(message);
    error.response = response;
    throw error;
  });
}

/**
 * Requests a URL, returning a promise.
 *
 * @param  {string} url       The URL we want to request
 * @param  {object} [options] The options we want to pass to "fetch"
 * @return {object}           An object containing either "data" or "err"
 */
export default function request(url, options, nonotification) {
  const defaultOptions = {
    credentials: 'include',
  };
  const newOptions = { ...defaultOptions, ...options };
  if (newOptions.method === 'POST' || newOptions.method === 'PUT') {
    newOptions.headers = {
      Accept: 'application/json',
      'Content-Type': 'application/json; charset=utf-8',
      ...newOptions.headers,
    };
    newOptions.body = JSON.stringify(newOptions.body);
  }

  return fetch(url, newOptions)
    .then(checkStatus)
    .then(response => response.json())
    .then(response => {
      if (!response) {
        return {};
      }
      return response;
    })
    .catch(error => {
      if (!nonotification && 'stack' in error && 'message' in error) {
        notification.error({
          message: error.message,
          // message: `请求错误: ${url}`,
          // description: error.message,
        });
      }
      return error;
    });
}
