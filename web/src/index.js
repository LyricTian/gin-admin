import './polyfill';
import './index.less';
import dva from 'dva';
import 'moment/locale/zh-cn';
import browserHistory from 'history/createBrowserHistory';
import router from './router';
import { storeLogoutKey } from './utils/utils';

// 1. Initialize
const app = dva({
  history: browserHistory(),
});

// 2. Plugins
// app.use({});

// 3. Register global model
app.model(require('./models/global').default);

// 4. Router
app.router(router);

// 5. Start
app.start('#root');

// 移除重定向
sessionStorage.removeItem(storeLogoutKey);

export default app._store; // eslint-disable-line
