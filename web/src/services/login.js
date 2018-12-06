import request, { v1API } from '../utils/request';

// 登录
export async function login(params) {
  return request(`${v1API}/login`, {
    method: 'POST',
    body: params,
  });
}

// 登出
export async function logout() {
  return request(`${v1API}/logout`, {
    method: 'POST',
  });
}

// 查询当前用户信息
export async function getCurrentUser() {
  return request(`${v1API}/current/user`);
}

// 查询当前用户菜单
export async function queryCurrentMenus() {
  return request(`${v1API}/current/menus`);
}
