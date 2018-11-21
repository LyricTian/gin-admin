import moment from 'moment';
import md5 from 'md5';

// 登出key
export const storeLogoutKey = 'is_logout';

// 扩展节点
export function getPlainNode(nodeList, parentPath = '') {
  const arr = [];
  nodeList.forEach(node => {
    const item = node;
    let itemPath = item.path;
    if (itemPath !== '') {
      itemPath = `/${itemPath}`;
    }
    item.path = `${parentPath}${itemPath}`;

    item.exact = true;
    if (item.children && !item.component) {
      arr.push(...getPlainNode(item.children, item.path));
    } else {
      if (item.children && item.component) {
        item.exact = false;
      }
      arr.push(item);
    }
  });
  return arr;
}

// 获取分级码
export function getLevelCode(pathname, menuPaths) {
  const node = menuPaths[pathname];
  if (node) {
    return node.level_code;
  }

  // 如果匹配不到全部路径，则优先匹配前缀
  const pp = pathname.split('/');
  if (pp.length > 2) {
    for (let i = pp.length - 1; i > 1; i -= 1) {
      const n = menuPaths[pp.slice(0, i).join('/')];
      if (n) {
        return n.level_code;
      }
    }
  }

  return '';
}

// 获取菜单key
export function getMenuKeys(pathname, menuPaths, menus) {
  const levelCode = getLevelCode(pathname, menuPaths);
  if (levelCode.length === 0) {
    return [];
  }

  let prefix = '';
  for (let i = 0; i < menus.length; i += 1) {
    const { level_code: nodeLevelCode } = menus[i];
    const l = levelCode
      .split('')
      .slice(0, nodeLevelCode.length)
      .join('');

    if (nodeLevelCode === l) {
      prefix = nodeLevelCode;
      break;
    }
  }

  const keys = [prefix];
  levelCode
    .split('')
    .slice(prefix.length)
    .forEach((s, i) => {
      prefix += s;
      if ((i + 1) % 2 === 0) {
        keys.push(prefix);
      }
    });

  return keys;
}

// 格式化时间戳
export function formatTimestamp(val, format) {
  let f = 'YYYY-MM-DD HH:mm:ss';
  if (format) {
    f = format;
  }
  return moment.unix(val).format(f);
}

// 解析时间戳
export function parseTimestamp(val) {
  return moment.unix(val);
}

// 解析日期
export function parseDate(val) {
  return moment(val);
}

// md5加密
export function md5Hash(value) {
  return md5(value);
}
