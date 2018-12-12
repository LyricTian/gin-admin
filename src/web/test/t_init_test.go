package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/LyricTian/gin-admin/src"
	"github.com/LyricTian/gin-admin/src/util"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"

	_ "github.com/go-sql-driver/mysql"
)

const (
	configFile = "../../../config/config.toml"
	apiPrefix  = "/api/v1/"
)

var engine *gin.Engine

func init() {
	viper.SetConfigFile(configFile)
	if err := viper.ReadInConfig(); err != nil {
		panic("加载配置文件发生错误：" + err.Error())
	}
	viper.Set("run_mode", "debug")
	viper.Set("casbin_model", "../../../config/model.conf")

	engine, _ = src.Init("1.0.0", util.MustUUID())
}

func toReader(v interface{}) io.Reader {
	buf := new(bytes.Buffer)
	json.NewEncoder(buf).Encode(v)
	return buf
}

func parseReader(r io.Reader, v interface{}) error {
	return json.NewDecoder(r).Decode(v)
}

func releaseReader(r io.Reader) {
	ioutil.ReadAll(r)
}

func newPageParam(extra map[string]string) map[string]string {
	data := map[string]string{
		"current":  "1",
		"pageSize": "1",
	}
	for k, v := range extra {
		data[k] = v
	}
	return data
}

type PaginationResult struct {
	Total    int64 `json:"total"`
	Current  int   `json:"current"`
	PageSize int   `json:"pageSize"`
}

type PageResult struct {
	List       interface{}       `json:"list"`
	Pagination *PaginationResult `json:"pagination"`
}

func parsePageReader(r io.Reader, v interface{}) error {
	result := &PageResult{List: v}
	return parseReader(r, result)
}

func newPostRequest(router string, v interface{}) *http.Request {
	req, _ := http.NewRequest("POST", apiPrefix+router, toReader(v))
	return req
}

func newPutRequest(formatRouter string, v interface{}, args ...interface{}) *http.Request {
	req, _ := http.NewRequest("PUT", apiPrefix+fmt.Sprintf(formatRouter, args...), toReader(v))
	return req
}

func newPatchRequest(formatRouter string, args ...interface{}) *http.Request {
	req, _ := http.NewRequest("PATCH", apiPrefix+fmt.Sprintf(formatRouter, args...), nil)
	return req
}

func newDeleteRequest(formatRouter string, args ...interface{}) *http.Request {
	req, _ := http.NewRequest("DELETE", apiPrefix+fmt.Sprintf(formatRouter, args...), nil)
	return req
}

func newGetRequest(formatRouter string, params map[string]string, args ...interface{}) *http.Request {
	values := make(url.Values)
	for k, v := range params {
		values.Set(k, v)
	}

	urlStr := apiPrefix + fmt.Sprintf(formatRouter, args...)
	if len(values) > 0 {
		urlStr += "?" + values.Encode()
	}

	req, _ := http.NewRequest("GET", urlStr, nil)
	return req
}
