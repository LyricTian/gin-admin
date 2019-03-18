import request from '../utils/request';

// 验证码ID
export async function captchaID() {
  return request(`/v1/login/captchaid`);
}

// 图形验证码
export function captcha(id) {
  return `/api/v1/login/captcha?id=${id}`;
}

// 登录
export async function login(params) {
  return request(`/v1/login`, {
    method: 'POST',
    body: params,
    notNotify: true,
  });
}

// 退出
export async function logout() {
  return request(`/v1/login/exit`, {
    method: 'POST',
  });
}

// 获取当前用户信息
export async function getCurrentUser() {
  return request(`/v1/current/user`);
}

// 查询当前用户菜单树
export async function queryMenuTree() {
  return request(`/v1/current/menutree`);
}
