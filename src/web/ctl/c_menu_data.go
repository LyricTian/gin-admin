package ctl

// 初始化菜单数据
const menuData = `
[
  {
    "name": "首页",
    "icon": "dashboard",
    "path": "/dashboard",
    "sequence": 1900000
  },
  {
    "name": "系统管理",
    "icon": "setting",
    "sequence": 1100000,
    "children": [
      {
        "name": "菜单管理",
        "icon": "solution",
        "path": "/system/menu",
        "sequence": 1190000,
        "resources": [
          {
            "name": "查询菜单数据",
            "method": "GET",
            "path": "/api/v1/menus"
          },
          {
            "name": "精确查询菜单数据",
            "method": "GET",
            "path": "/api/v1/menus/:id"
          },
          {
            "name": "创建菜单数据",
            "method": "POST",
            "path": "/api/v1/menus"
          },
          {
            "name": "更新菜单数据",
            "method": "PUT",
            "path": "/api/v1/menus/:id"
          },
          {
            "name": "删除菜单数据",
            "method": "DELETE",
            "path": "/api/v1/menus/:id"
          }
        ]
      },
      {
        "name": "角色管理",
        "icon": "audit",
        "path": "/system/role",
        "sequence": 1180000,
        "resources": [
          {
            "name": "查询角色数据",
            "method": "GET",
            "path": "/api/v1/roles"
          },
          {
            "name": "精确查询角色数据",
            "method": "GET",
            "path": "/api/v1/roles/:id"
          },
          {
            "name": "创建角色数据",
            "method": "POST",
            "path": "/api/v1/roles"
          },
          {
            "name": "更新角色数据",
            "method": "PUT",
            "path": "/api/v1/roles/:id"
          },
          {
            "name": "删除角色数据",
            "method": "DELETE",
            "path": "/api/v1/roles/:id"
          },
          {
            "name": "查询菜单数据",
            "method": "GET",
            "path": "/api/v1/menus"
          }
        ]
      },
      {
        "name": "用户管理",
        "icon": "user",
        "path": "/system/user",
        "sequence": 1170000,
        "resources": [
          {
            "name": "查询用户数据",
            "method": "GET",
            "path": "/api/v1/users"
          },
          {
            "name": "精确查询用户数据",
            "method": "GET",
            "path": "/api/v1/users/:id"
          },
          {
            "name": "创建用户数据",
            "method": "POST",
            "path": "/api/v1/users"
          },
          {
            "name": "更新用户数据",
            "method": "PUT",
            "path": "/api/v1/users/:id"
          },
          {
            "name": "删除用户数据",
            "method": "DELETE",
            "path": "/api/v1/users/:id"
          },
          {
            "name": "禁用用户数据",
            "method": "PATCH",
            "path": "/api/v1/users/:id/disable"
          },
          {
            "name": "启用用户数据",
            "method": "PATCH",
            "path": "/api/v1/users/:id/enable"
          },
          {
            "name": "查询角色数据",
            "method": "GET",
            "path": "/api/v1/roles"
          }
        ]
      }
    ]
  }
]
`
