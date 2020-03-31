package api

// Role 角色管理
type Role struct {
}

// Query 查询数据
// @Tags 角色管理
// @Summary 查询数据
// @Param Authorization header string false "Bearer 用户令牌"
// @Param current query int true "分页索引" default(1)
// @Param pageSize query int true "分页大小" default(10)
// @Param likeName query string false "角色名称(模糊查询)"
// @Param status query int false "状态(1:启用 2:禁用)"
// @Success 200 {array} schema.Role "查询结果：{list:列表数据,pagination:{current:页索引,pageSize:页大小,total:总数量}}"
// @Failure 401 {object} schema.HTTPError "{error:{code:0,message:未授权}}"
// @Failure 500 {object} schema.HTTPError "{error:{code:0,message:服务器错误}}"
// @Router /api/v1/roles [get]
func (a *Role) Query() {
}

// QuerySelect 查询选择数据
// @Tags 角色管理
// @Summary 查询选择数据
// @Param Authorization header string false "Bearer 用户令牌"
// @Param likeName query string false "角色名称(模糊查询)"
// @Param status query int false "状态(1:启用 2:禁用)"
// @Success 200 {array} schema.Role "查询结果：{list:角色列表}"
// @Failure 400 {object} schema.HTTPError "{error:{code:0,message:未知的查询类型}}"
// @Failure 401 {object} schema.HTTPError "{error:{code:0,message:未授权}}"
// @Failure 500 {object} schema.HTTPError "{error:{code:0,message:服务器错误}}"
// @Router /api/v1/roles.select [get]
func (a *Role) QuerySelect() {
}

// Get 查询指定数据
// @Tags 角色管理
// @Summary 查询指定数据
// @Param Authorization header string false "Bearer 用户令牌"
// @Param id path string true "记录ID"
// @Success 200 {object} schema.Role
// @Failure 401 {object} schema.HTTPError "{error:{code:0,message:未授权}}"
// @Failure 404 {object} schema.HTTPError "{error:{code:0,message:资源不存在}}"
// @Failure 500 {object} schema.HTTPError "{error:{code:0,message:服务器错误}}"
// @Router /api/v1/roles/{id} [get]
func (a *Role) Get() {
}

// Create 创建数据
// @Tags 角色管理
// @Summary 创建数据
// @Param Authorization header string false "Bearer 用户令牌"
// @Param body body schema.Role true "创建数据"
// @Success 200 {object} schema.HTTPRecordID
// @Failure 400 {object} schema.HTTPError "{error:{code:0,message:无效的请求参数}}"
// @Failure 401 {object} schema.HTTPError "{error:{code:0,message:未授权}}"
// @Failure 500 {object} schema.HTTPError "{error:{code:0,message:服务器错误}}"
// @Router /api/v1/roles [post]
func (a *Role) Create() {
}

// Update 更新数据
// @Tags 角色管理
// @Summary 更新数据
// @Param Authorization header string false "Bearer 用户令牌"
// @Param id path string true "记录ID"
// @Param body body schema.Role true "更新数据"
// @Success 200 {object} schema.Role
// @Failure 400 {object} schema.HTTPError "{error:{code:0,message:无效的请求参数}}"
// @Failure 401 {object} schema.HTTPError "{error:{code:0,message:未授权}}"
// @Failure 500 {object} schema.HTTPError "{error:{code:0,message:服务器错误}}"
// @Router /api/v1/roles/{id} [put]
func (a *Role) Update() {
}

// Delete 删除数据
// @Tags 角色管理
// @Summary 删除数据
// @Param Authorization header string false "Bearer 用户令牌"
// @Param id path string true "记录ID"
// @Success 200 {object} schema.HTTPStatus "{status:OK}"
// @Failure 401 {object} schema.HTTPError "{error:{code:0,message:未授权}}"
// @Failure 500 {object} schema.HTTPError "{error:{code:0,message:服务器错误}}"
// @Router /api/v1/roles/{id} [delete]
func (a *Role) Delete() {
}

// Enable 启用数据
// @Tags 角色管理
// @Summary 启用数据
// @Param Authorization header string false "Bearer 用户令牌"
// @Param id path string true "记录ID"
// @Success 200 {object} schema.HTTPStatus "{status:OK}"
// @Failure 401 {object} schema.HTTPError "{error:{code:0,message:未授权}}"
// @Failure 500 {object} schema.HTTPError "{error:{code:0,message:服务器错误}}"
// @Router /api/v1/roles/{id}/enable [patch]
func (a *Role) Enable() {
}

// Disable 禁用数据
// @Tags 角色管理
// @Summary 禁用数据
// @Param Authorization header string false "Bearer 用户令牌"
// @Param id path string true "记录ID"
// @Success 200 {object} schema.HTTPStatus "{status:OK}"
// @Failure 401 {object} schema.HTTPError "{error:{code:0,message:未授权}}"
// @Failure 500 {object} schema.HTTPError "{error:{code:0,message:服务器错误}}"
// @Router /api/v1/roles/{id}/disable [patch]
func (a *Role) Disable() {
}
