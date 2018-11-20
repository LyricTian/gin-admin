import request, { baseURLV1 } from '../utils/request';

// 登录
export async function login(params) {
  return request(`${baseURLV1}/login`, {
    method: 'POST',
    body: params,
  });
}

// 登出
export async function logout() {
  return request(`${baseURLV1}/logout`, {
    method: 'POST',
  });
}

// 查询当前用户信息
export async function getCurrentUser() {
  return request(`${baseURLV1}/current/user`);
}

// 查询当前用户菜单
export async function queryCurrentMenus() {
  return request(`${baseURLV1}/current/menus`);
}
