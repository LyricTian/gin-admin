package schema

// Demo 测试
type Demo struct {
	ID   int64  `json:"id" db:"id,primarykey,autoincrement"`
	Code string `json:"code" db:"code,size:50" binding:"required"`
	Name string `json:"name" db:"name,size:50" binding:"required"`
}
