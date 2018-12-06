import './polyfill';
import './global.less';
import 'moment/locale/zh-cn';
import { storeLogoutKey } from './utils/utils';

// 移除登出Key
sessionStorage.removeItem(storeLogoutKey);
