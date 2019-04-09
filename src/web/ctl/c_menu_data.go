package ctl

// 初始化菜单数据
const menuData = `
[
  {
    "name": "首页",
    "icon": "dashboard",
    "router": "/dashboard",
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
        "router": "/system/menu",
        "sequence": 1190000,
        "actions": [
          { "code": "add", "name": "新增" },
          { "code": "edit", "name": "编辑" },
          { "code": "del", "name": "删除" },
          { "code": "query", "name": "查询" }
        ],
        "resources": [
          {
            "code": "query",
            "name": "查询菜单数据",
            "method": "GET",
            "path": "/api/v1/menus"
          },
          {
            "code": "get",
            "name": "精确查询菜单数据",
            "method": "GET",
            "path": "/api/v1/menus/:id"
          },
          {
            "code": "create",
            "name": "创建菜单数据",
            "method": "POST",
            "path": "/api/v1/menus"
          },
          {
            "code": "update",
            "name": "更新菜单数据",
            "method": "PUT",
            "path": "/api/v1/menus/:id"
          },
          {
            "code": "delete",
            "name": "删除菜单数据",
            "method": "DELETE",
            "path": "/api/v1/menus/:id"
          }
        ]
      },
      {
        "name": "角色管理",
        "icon": "audit",
        "router": "/system/role",
        "sequence": 1180000,
        "actions": [
          { "code": "add", "name": "新增" },
          { "code": "edit", "name": "编辑" },
          { "code": "del", "name": "删除" },
          { "code": "query", "name": "查询" }
        ],
        "resources": [
          {
            "code": "query",
            "name": "查询角色数据",
            "method": "GET",
            "path": "/api/v1/roles"
          },
          {
            "code": "get",
            "name": "精确查询角色数据",
            "method": "GET",
            "path": "/api/v1/roles/:id"
          },
          {
            "code": "create",
            "name": "创建角色数据",
            "method": "POST",
            "path": "/api/v1/roles"
          },
          {
            "code": "update",
            "name": "更新角色数据",
            "method": "PUT",
            "path": "/api/v1/roles/:id"
          },
          {
            "code": "delete",
            "name": "删除角色数据",
            "method": "DELETE",
            "path": "/api/v1/roles/:id"
          },
          {
            "code": "queryMenu",
            "name": "查询菜单数据",
            "method": "GET",
            "path": "/api/v1/menus"
          }
        ]
      },
      {
        "name": "用户管理",
        "icon": "user",
        "router": "/system/user",
        "sequence": 1170000,
        "actions": [
          { "code": "add", "name": "新增" },
          { "code": "edit", "name": "编辑" },
          { "code": "del", "name": "删除" },
          { "code": "query", "name": "查询" }
        ],
        "resources": [
          {
            "code": "query",
            "name": "查询用户数据",
            "method": "GET",
            "path": "/api/v1/users"
          },
          {
            "code": "get",
            "name": "精确查询用户数据",
            "method": "GET",
            "path": "/api/v1/users/:id"
          },
          {
            "code": "create",
            "name": "创建用户数据",
            "method": "POST",
            "path": "/api/v1/users"
          },
          {
            "code": "update",
            "name": "更新用户数据",
            "method": "PUT",
            "path": "/api/v1/users/:id"
          },
          {
            "code": "delete",
            "name": "删除用户数据",
            "method": "DELETE",
            "path": "/api/v1/users/:id"
          },
          {
            "code": "disable",
            "name": "禁用用户数据",
            "method": "PATCH",
            "path": "/api/v1/users/:id/disable"
          },
          {
            "code": "enable",
            "name": "启用用户数据",
            "method": "PATCH",
            "path": "/api/v1/users/:id/enable"
          },
          {
            "code": "queryRole",
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
