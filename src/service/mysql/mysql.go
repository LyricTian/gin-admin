package mysql

import (
	"bytes"
	"database/sql"
	"fmt"
	"log"
	"os"
	"reflect"
	"strings"
	"time"

	"github.com/LyricTian/retry"
	"gopkg.in/gorp.v2"
)

type (
	// M 定义字典类型
	M map[string]interface{}

	// Logger 定义日志输出
	Logger interface {
		Printf(format string, args ...interface{})
	}

	// Option 配置项
	Option func(*options)

	options struct {
		dsn          string        // 连接串
		trace        bool          // 追踪调试
		maxLifetime  time.Duration // 设置连接可以被重新使用的最大时间量
		maxOpenConns int           // 设置打开连接到数据库的最大数量
		maxIdleConns int           // 设置空闲连接池中的最大连接数
		logger       Logger        // 日志
		engine       string        // 数据库表的存储引擎
		encoding     string        // 数据库表的编码格式
	}
)

// SetDSN 设置连接串
func SetDSN(dsn string) Option {
	return func(o *options) {
		o.dsn = dsn
	}
}

// SetTrace 设置追踪调试
func SetTrace(t bool) Option {
	return func(o *options) {
		o.trace = t
	}
}

// SetLogger 设定追踪日志
func SetLogger(logger Logger) Option {
	return func(o *options) {
		o.logger = logger
	}
}

// SetMaxLifetime 设置连接可以被重新使用的最大时间量
func SetMaxLifetime(maxLifetime time.Duration) Option {
	return func(o *options) {
		o.maxLifetime = maxLifetime
	}
}

// SetMaxOpenConns 设置打开连接到数据库的最大数量
func SetMaxOpenConns(maxOpenConns int) Option {
	return func(o *options) {
		o.maxOpenConns = maxOpenConns
	}
}

// SetMaxIdleConns 设置空闲连接池中的最大连接数
func SetMaxIdleConns(maxIdleConns int) Option {
	return func(o *options) {
		o.maxIdleConns = maxIdleConns
	}
}

// SetEngine 设定数据库表的存储引擎
func SetEngine(engine string) Option {
	return func(o *options) {
		o.engine = engine
	}
}

// SetEncoding 设定数据库表的编码格式
func SetEncoding(encoding string) Option {
	return func(o *options) {
		o.encoding = encoding
	}
}

// NewDB 创建MySQL数据库实例
func NewDB(opts ...Option) (*DB, error) {
	o := &options{
		maxLifetime:  time.Hour * 2,
		maxOpenConns: 150,
		maxIdleConns: 50,
		logger:       log.New(os.Stderr, "", log.LstdFlags),
		engine:       "InnoDB",
		encoding:     "UTF8",
	}
	for _, opt := range opts {
		opt(o)
	}

	db, err := sql.Open("mysql", o.dsn)
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(o.maxOpenConns)
	db.SetMaxIdleConns(o.maxIdleConns)
	db.SetConnMaxLifetime(o.maxLifetime)

	err = retry.DoFunc(10, func() error {
		err := db.Ping()
		if err != nil {
			fmt.Printf("[mysql]发送ping值发生错误：%v\n", err)
		}
		return err
	}, func(i int) time.Duration {
		fmt.Printf("[mysql]正在尝试发送第[%d]次ping...\n", i)
		return time.Second * 2 * time.Duration(i)
	})
	if err != nil {
		return nil, err
	}

	dbMap := &DB{
		logger: o.logger,
		DbMap: &gorp.DbMap{
			Db: db,
			Dialect: gorp.MySQLDialect{
				Encoding: o.encoding,
				Engine:   o.engine,
			},
		},
	}
	if o.trace && o.logger != nil {
		dbMap.TraceOn("[mysql]", o.logger)
	}

	return dbMap, nil
}

// DB 数据库管理
type DB struct {
	*gorp.DbMap
	logger Logger
}

// Close 关闭数据库连接
func (d *DB) Close() error {
	if d.DbMap == nil {
		return nil
	}
	return d.Db.Close()
}

// CreateTableIfNotExists 创建表
func (d *DB) CreateTableIfNotExists(i interface{}, name string) {
	tableMap := d.AddTableWithName(i, name)
	query := tableMap.SqlForCreate(true)
	_, err := d.Exec(query)
	if err != nil {
		d.logger.Printf("创建表[%s]发生错误:%s", query, err.Error())
	}
}

// CreateTableIndex 创建索引
func (d *DB) CreateTableIndex(table, idx string, unique bool, columns ...string) {
	s := bytes.Buffer{}
	s.WriteString("CREATE")
	s.WriteByte(' ')
	if unique {
		s.WriteString("UNIQUE")
		s.WriteByte(' ')
	}
	s.WriteString("INDEX")
	s.WriteByte(' ')
	s.WriteString(idx)
	s.WriteByte(' ')
	s.WriteString("ON")
	s.WriteByte(' ')
	s.WriteString(table)
	s.WriteByte(' ')
	s.WriteByte('(')
	s.WriteString(strings.Join(columns, ","))
	s.WriteByte(')')
	s.WriteByte(';')

	_, err := d.Exec(s.String())
	if err != nil {
		errString := err.Error()
		if !strings.HasPrefix(errString, "Error 1061:") {
			d.logger.Printf("创建索引[%s]发生错误:%s", s.String(), errString)
		}
	}
}

// InsertSQL 获取插入SQL
func (d *DB) InsertSQL(table string, info M) (string, []interface{}) {
	q := fmt.Sprintf("INSERT INTO %s", table)

	var (
		cols []string
		vals []interface{}
	)

	for k, v := range info {
		cols = append(cols, k)
		vals = append(vals, v)
	}

	q = fmt.Sprintf("%s(%s) VALUES(%s)", q, strings.Join(cols, ","), strings.Repeat(",?", len(cols))[1:])
	return q, vals
}

// InsertM 插入数据
func (d *DB) InsertM(table string, info M) (int64, error) {
	q, vals := d.InsertSQL(table, info)
	result, err := d.Exec(q, vals...)
	if err != nil {
		return 0, err
	}
	lastInsertID, _ := result.LastInsertId()

	return lastInsertID, nil
}

// InsertMWithTran 基于事物插入数据
func (d *DB) InsertMWithTran(tran *gorp.Transaction, table string, info M) (int64, error) {
	q, vals := d.InsertSQL(table, info)
	result, err := tran.Exec(q, vals...)
	if err != nil {
		return 0, err
	}
	lastInsertID, _ := result.LastInsertId()

	return lastInsertID, nil
}

// UpdateSQL 获取更新SQL
func (d *DB) UpdateSQL(table string, pk, info M) (string, []interface{}) {
	q := fmt.Sprintf("UPDATE %s SET", table)

	var (
		cols []string
		vals []interface{}
	)

	for k, v := range info {
		cols = append(cols, fmt.Sprintf("%s=?", k))
		vals = append(vals, v)
	}

	q = fmt.Sprintf("%s %s", q, strings.Join(cols, ","))
	cols = nil

	for k, v := range pk {
		cols = append(cols, fmt.Sprintf("%s=?", k))
		vals = append(vals, v)
	}

	q = fmt.Sprintf("%s WHERE %s", q, strings.Join(cols, " and "))
	return q, vals
}

// UpdateByPK 更新数据
func (d *DB) UpdateByPK(table string, pk, info M) (int64, error) {
	q, vals := d.UpdateSQL(table, pk, info)
	result, err := d.Exec(q, vals...)
	if err != nil {
		return 0, err
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}

	return affected, nil
}

// UpdateByPKWithTran 基于事物更新数据
func (d *DB) UpdateByPKWithTran(tran *gorp.Transaction, table string, pk, info M) (int64, error) {
	q, vals := d.UpdateSQL(table, pk, info)
	result, err := tran.Exec(q, vals...)
	if err != nil {
		return 0, err
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}

	return affected, nil
}

// DeleteSQL 获取删除SQL
func (d *DB) DeleteSQL(table string, pk M) (string, []interface{}) {
	q := fmt.Sprintf("DELETE FROM %s", table)

	var (
		cols []string
		vals []interface{}
	)

	for k, v := range pk {
		cols = append(cols, fmt.Sprintf("%s=?", k))
		vals = append(vals, v)
	}

	q = fmt.Sprintf("%s WHERE %s", q, strings.Join(cols, " and "))
	return q, vals
}

// DeleteByPK 删除数据
func (d *DB) DeleteByPK(table string, pk M) (int64, error) {
	q, vals := d.DeleteSQL(table, pk)
	result, err := d.Exec(q, vals...)
	if err != nil {
		return 0, err
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}

	return affected, nil
}

// DeleteByPKWithTran 基于事物删除表数据
func (d *DB) DeleteByPKWithTran(tran *gorp.Transaction, table string, pk M) (int64, error) {
	q, vals := d.DeleteSQL(table, pk)
	result, err := tran.Exec(q, vals...)
	if err != nil {
		return 0, err
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}

	return affected, nil
}

// In 组织带有IN查询的SQL和参数
func (d *DB) In(query string, args ...interface{}) (string, []interface{}, error) {
	type argMeta struct {
		v      reflect.Value
		i      interface{}
		length int
	}

	var flatArgsCount int
	var anySlices bool

	meta := make([]argMeta, len(args))

	for i, arg := range args {
		v := reflect.ValueOf(arg)

		t := v.Type()
		if t.Kind() == reflect.Ptr {
			t = t.Elem()
		}

		if t.Kind() == reflect.Slice {
			meta[i].length = v.Len()
			meta[i].v = v

			anySlices = true
			flatArgsCount += meta[i].length

			if meta[i].length == 0 {
				return "", nil, fmt.Errorf("empty slice passed to 'in' query")
			}
		} else {
			meta[i].i = arg
			flatArgsCount++
		}
	}

	if !anySlices {
		return query, args, nil
	}

	newArgs := make([]interface{}, 0, flatArgsCount)
	buf := bytes.NewBuffer(make([]byte, 0, len(query)+len(", ?")*flatArgsCount))

	var arg, offset int

	for i := strings.IndexByte(query[offset:], '?'); i != -1; i = strings.IndexByte(query[offset:], '?') {
		if arg >= len(meta) {
			return "", nil, fmt.Errorf("number of bindVars exceeds arguments")
		}

		argMeta := meta[arg]
		arg++

		if argMeta.length == 0 {
			offset = offset + i + 1
			newArgs = append(newArgs, argMeta.i)
			continue
		}

		buf.WriteString(query[:offset+i+1])

		for si := 1; si < argMeta.length; si++ {
			buf.WriteString(", ?")
		}

		newArgs = d.appendReflectSlice(newArgs, argMeta.v, argMeta.length)

		query = query[offset+i+1:]
		offset = 0
	}

	buf.WriteString(query)

	if arg < len(meta) {
		return "", nil, fmt.Errorf("number of bindVars less than number arguments")
	}

	return buf.String(), newArgs, nil
}

func (d *DB) appendReflectSlice(args []interface{}, v reflect.Value, vlen int) []interface{} {
	switch val := v.Interface().(type) {
	case []interface{}:
		args = append(args, val...)
	case []int:
		for i := range val {
			args = append(args, val[i])
		}
	case []string:
		for i := range val {
			args = append(args, val[i])
		}
	default:
		for si := 0; si < vlen; si++ {
			args = append(args, v.Index(si).Interface())
		}
	}

	return args
}
