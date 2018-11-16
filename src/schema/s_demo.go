package schema

// Demo 示例程序
type Demo struct {
	ID       int64  `json:"id" db:"id,primarykey,autoincrement" structs:"id"`         // 唯一标识(自增ID)
	RecordID string `json:"record_id" db:"record_id,size:36" structs:"record_id"`     // 记录内码(uuid)
	Code     string `json:"code" db:"code,size:50" structs:"code" binding:"required"` // 编号
	Name     string `json:"name" db:"name,size:50" structs:"name" binding:"required"` // 名称
	Creator  string `json:"creator" db:"creator,size:36" structs:"creator"`           // 创建者
	Created  int64  `json:"created" db:"created" structs:"created"`                   // 创建时间戳
	Updated  int64  `json:"updated" db:"updated" structs:"updated"`                   // 更新时间戳
	Deleted  int64  `json:"deleted" db:"deleted" structs:"deleted"`                   // 删除时间戳
}

// DemoQueryParam 示例查询条件
type DemoQueryParam struct {
	Code string // 编号
	Name string // 名称
}

// DemoQueryResult 示例查询结果
type DemoQueryResult struct {
	ID       int64  `json:"id" db:"id"`               // 唯一标识(自增ID)
	RecordID string `json:"record_id" db:"record_id"` // 记录内码(uuid)
	Code     string `json:"code" db:"code"`           // 编号
	Name     string `json:"name" db:"name"`           // 名称
}
