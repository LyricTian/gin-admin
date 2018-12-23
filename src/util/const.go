package util

const (
	contextKeyPrefix = "github.com/LyricTian/gin-admin"
	// DebugMode 调试模式
	DebugMode = "debug"
	// TestMode 测试模式
	TestMode = "test"
	// ReleaseMode 正式模式
	ReleaseMode = "release"
	// SessionKeyUserID 存储在session中的键(用户ID)
	SessionKeyUserID = "user_id"
	// ContextKeyUserID 存储上下文中的键(用户ID)
	ContextKeyUserID = contextKeyPrefix + "/user_id"
	// ContextKeyURLMemo 存储上下文中的键(请求URL说明)
	ContextKeyURLMemo = contextKeyPrefix + "/url_memo"
	// ContextKeyTraceID 存储上下文中的键(跟踪ID)
	ContextKeyTraceID = contextKeyPrefix + "/trace_id"
	// ContextKeyResBody 存储上下文中的键(响应Body数据)
	ContextKeyResBody = contextKeyPrefix + "/res_body"
)
