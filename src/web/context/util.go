package context

import (
	"fmt"
	"regexp"
	"strings"
	"sync"
)

// 定义上下文中的键
const (
	contextKeyPrefix = "github.com/LyricTian/gin-admin"
	// ContextKeyUserID 存储上下文中的键(用户ID)
	ContextKeyUserID = contextKeyPrefix + "/user_id"
	// ContextKeyURLTitle 存储上下文中的键(请求URL说明)
	ContextKeyURLTitle = contextKeyPrefix + "/url_title"
	// ContextKeyTraceID 存储上下文中的键(跟踪ID)
	ContextKeyTraceID = contextKeyPrefix + "/trace_id"
	// ContextKeyResBody 存储上下文中的键(响应Body数据)
	ContextKeyResBody = contextKeyPrefix + "/res_body"
)

// 定义响应状态数据
const (
	StatusOK    = "OK"
	StatusError = "error"
	StatusFail  = "fail"
)

// 路由关联的标题数据
var (
	routerTitle = &sync.Map{}
	routerRe    = regexp.MustCompile(`(.*):[^/]+(.*)`)
)

// GetRouter 获取路由
func GetRouter(method, path string) string {
	if len(path) > 0 && path[0] != '/' {
		path = "/" + path
	}
	return fmt.Sprintf("%s%s", method, path)
}

// SetRouterTitle 设定路由标题
func SetRouterTitle(method, path, title string) {
	routerTitle.Store(GetRouter(method, path), title)
}

// GetRouterTitleAndKey 获取路由标题和键
func GetRouterTitleAndKey(method, path string) (string, string) {
	key := GetRouter(method, path)
	vv, ok := routerTitle.Load(key)
	if ok {
		return vv.(string), key
	}

	var title string
	routerTitle.Range(func(vk, vv interface{}) bool {
		vkey := vk.(string)
		if !strings.Contains(vkey, "/:") {
			return true
		}

		rkey := "^" + routerRe.ReplaceAllString(vkey, "$1[^/]+$2") + "$"
		b, _ := regexp.MatchString(rkey, key)
		if b {
			title = vv.(string)
			key = vkey
		}
		return !b
	})

	return title, key
}

// HTTPError HTTP响应错误
type HTTPError struct {
	Error HTTPErrorItem `json:"error,omitempty"`
}

// HTTPErrorItem HTTP响应错误项
type HTTPErrorItem struct {
	Code    int    `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
}

// HTTPStatus HTTP响应状态
type HTTPStatus struct {
	Status string `json:"status,omitempty"`
}

// HTTPList HTTP响应列表数据
type HTTPList struct {
	List       interface{}     `json:"list,omitempty"`
	Pagination *HTTPPagination `json:"pagination,omitempty"`
}

// HTTPPagination HTTP分页数据
type HTTPPagination struct {
	Total    int  `json:"total,omitempty"`
	Current  uint `json:"current,omitempty"`
	PageSize uint `json:"pageSize,omitempty"`
}

// HTTPNewItem HTTP响应创建成功后的记录ID
type HTTPNewItem struct {
	RecordID string `json:"record_id"`
}
