const path = require('path');
const config = require('./src/config');

export default {
  entry: 'src/index.js',
  extraBabelPlugins: [['import', { libraryName: 'antd', libraryDirectory: 'es', style: true }]],
  env: {
    development: {
      extraBabelPlugins: ['dva-hmr'],
    },
  },
  alias: {
    components: path.resolve(__dirname, 'src/components/'),
  },
  ignoreMomentLocale: true,
  theme: './src/theme.js',
  html: {
    title: config.title,
    template: './src/index.ejs',
  },
  lessLoaderOptions: {
    javascriptEnabled: true,
  },
  disableDynamicImport: false,
  publicPath: '/',
  hash: true,
  proxy: {
    '/api/': {
      target: 'http://127.0.0.1:8086/',
      changeOrigin: true,
    },
  },
};
