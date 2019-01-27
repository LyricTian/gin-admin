package context

import (
	"fmt"
	"regexp"
	"strings"
	"sync"
)

// 定义上下文中的键
const (
	prefix = "github.com/LyricTian/gin-admin"
	// UserIDKey 存储上下文中的键(用户ID)
	UserIDKey = prefix + "/user_id"
	// TraceIDKey 存储上下文中的键(跟踪ID)
	TraceIDKey = prefix + "/trace_id"
	// ResBodyKey 存储上下文中的键(响应Body数据)
	ResBodyKey = prefix + "/res_body"
)

// 定义响应状态数据
const (
	StatusOK    = "OK"
	StatusError = "error"
	StatusFail  = "fail"
)

var (
	routerData = &sync.Map{}
	routerRe   = regexp.MustCompile(`(.*):[^/]+(.*)`)
)

// RouterItem 路由项
type RouterItem struct {
	Code string // 路由编号
	Name string // 路由名称
}

// JoinRouter 拼接路由
func JoinRouter(method, path string) string {
	if len(path) > 0 && path[0] != '/' {
		path = "/" + path
	}
	return fmt.Sprintf("%s%s", strings.ToUpper(method), path)
}

// SetRouterItem 存储路由项
func SetRouterItem(key string, item RouterItem) {
	routerData.Store(key, item)
}

// GetRouterItem 获取路由项
func GetRouterItem(key string) RouterItem {
	vv, ok := routerData.Load(key)
	if ok {
		return vv.(RouterItem)
	}

	var item RouterItem
	routerData.Range(func(vk, vv interface{}) bool {
		vkey := vk.(string)
		if !strings.Contains(vkey, "/:") {
			return true
		}

		rkey := "^" + routerRe.ReplaceAllString(vkey, "$1[^/]+$2") + "$"
		b, _ := regexp.MatchString(rkey, key)
		if b {
			item = vv.(RouterItem)
		}
		return !b
	})
	return item
}

// HTTPError HTTP响应错误
type HTTPError struct {
	Error HTTPErrorItem `json:"error"`
}

// HTTPErrorItem HTTP响应错误项
type HTTPErrorItem struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// HTTPStatus HTTP响应状态
type HTTPStatus struct {
	Status string `json:"status"`
}

// HTTPList HTTP响应列表数据
type HTTPList struct {
	List       interface{}     `json:"list"`
	Pagination *HTTPPagination `json:"pagination,omitempty"`
}

// HTTPPagination HTTP分页数据
type HTTPPagination struct {
	Total    int  `json:"total"`
	Current  uint `json:"current"`
	PageSize uint `json:"pageSize"`
}

// HTTPNewItem HTTP响应创建成功后的记录ID
type HTTPNewItem struct {
	RecordID string `json:"record_id"`
}
