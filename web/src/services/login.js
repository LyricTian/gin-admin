import request, { v1API } from '../utils/request';

// 验证码ID
export async function captchaID() {
  return request(`${v1API}/login/captchaid`);
}

// 图形验证码
export function captcha(id) {
  return `${v1API}/login/captcha?id=${id}`;
}

// 登录
export async function login(params) {
  return request(`${v1API}/login`, {
    method: 'POST',
    body: params,
    notNotify: true,
  });
}

// 退出
export async function logout() {
  return request(`${v1API}/login/exit`, {
    method: 'POST',
  });
}

// 获取当前用户信息
export async function getCurrentUser() {
  return request(`${v1API}/current/user`);
}

// 查询当前用户菜单树
export async function queryMenuTree() {
  return request(`${v1API}/current/menutree`);
}
