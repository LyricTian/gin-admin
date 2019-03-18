import axios from 'axios';
import { notification } from 'antd';
import moment from 'moment';
import store from './store';

function checkAccessTokenExpires(expiresAt) {
  const now = moment().unix();
  if (expiresAt - now <= 0) {
    return -1;
  }
  if (expiresAt - now <= 600) {
    return 0;
  }
  return 1;
}

async function getAccessToken() {
  const tokenInfo = store.getAccessToken();
  if (!tokenInfo) {
    return '';
  }

  if (checkAccessTokenExpires(tokenInfo.expires_at) === 0) {
    return axios
      .request({
        url: '/api/v1/refresh_token',
        method: 'POST',
        headers: {
          Authorization: `${tokenInfo.token_type} ${tokenInfo.access_token}`,
        },
      })
      .then(response => {
        const { status, data } = response;
        if (status === 200) {
          store.setAccessToken(data);
          return `${data.token_type} ${data.access_token}`;
        }
        return '';
      });
  }
  return `${tokenInfo.token_type} ${tokenInfo.access_token}`;
}

export default async function request(url, options) {
  let showNotify = true;
  const opts = {
    baseURL: '/api',
    url,
    validateStatus() {
      return true;
    },
    ...options,
  };
  if (opts.notNotify) {
    showNotify = false;
  }

  const defaultHeader = {
    Authorization: await getAccessToken(),
  };
  if (opts.method === 'POST' || opts.method === 'PUT') {
    defaultHeader['Content-Type'] = 'application/json; charset=utf-8';
    opts.data = opts.body;
  }
  opts.headers = { ...defaultHeader, ...opts.headers };

  return axios.request(opts).then(response => {
    const { status, data } = response;
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
  });
}
