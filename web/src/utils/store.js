const accessTokenKey = 'access_token';

export default class store {
  // 设定访问令牌
  static setAccessToken(token) {
    sessionStorage.setItem(accessTokenKey, JSON.stringify(token));
  }

  // 获取访问令牌
  static getAccessToken() {
    const token = sessionStorage.getItem(accessTokenKey);
    if (!token || token === '') {
      return null;
    }
    return JSON.parse(token);
  }

  // 清空访问令牌
  static clearAccessToken() {
    sessionStorage.removeItem(accessTokenKey);
  }
}
