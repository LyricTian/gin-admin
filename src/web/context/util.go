package context

import (
	"fmt"
	"regexp"
	"strings"
	"sync"
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

// SetRouterTitle 设定路由标题
func SetRouterTitle(method, router, title string) {
	routerTitle.Store(fmt.Sprintf("%s-%s", method, router), title)
}

// GetRouterTitleAndKey 获取路由标题和键
func GetRouterTitleAndKey(method, router string) (string, string) {
	key := fmt.Sprintf("%s-%s", method, router)
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
	Total    int64 `json:"total,omitempty"`
	Current  uint  `json:"current,omitempty"`
	PageSize uint  `json:"pageSize,omitempty"`
}
