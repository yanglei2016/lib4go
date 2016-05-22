package db

import (
	"database/sql"
	"os"
	"strings"
	//_ "github.com/mattn/go-oci8"
	_ "github.com/mattn/go-sqlite3"
)

const (
	//SQLITE3 Sqlite3数据库
	SQLITE3 = "sqlite3"
)

//DB 数据库实体
type DB struct {
	provider   string
	connString string
	db         *sql.DB
	maxIdle    int
	maxOpen    int
	lang       string
}

//NewDB 创建DB实例
func NewDB(provider string, connString string) (obj *DB, err error) {
	obj = &DB{provider: provider, connString: connString, maxIdle: 3, maxOpen: 10, lang: "AMERICAN_AMERICA.AL32UTF8"}
	obj.db, err = sql.Open(provider, connString)
	return
}

//SetPoolSize 设置连接池大小
func (db *DB) SetPoolSize(maxIdle int, maxOpen int) {
	db.db.SetMaxIdleConns(maxIdle)
	db.db.SetMaxOpenConns(maxOpen)
}

//QuerySchema 根据包含@名称占位符的查询语句执行查询语句
func (db *DB) QuerySchema(query string, data map[string]interface{}) (dataRows []map[string]interface{}, err error) {
	query, args := GetSchema(db.provider, query, data)
	return db.Query(query, args...)
}



//Query 执行SQL查询语句
func (db *DB) Query(query string, args ...interface{}) (dataRows []map[string]interface{}, err error) {
	rows, err := db.db.Query(query, args...)
	if err != nil {
		return
	}
	defer rows.Close()
	columns, err := rows.Columns()
	if err != nil {
		return
	}
	for rows.Next() {
		row := make(map[string]interface{})
		dataRows = append(dataRows, row)
		for index := 0; index < len(columns); index++ {
			var value interface{}
			err = rows.Scan(&value)
			if err != nil {
				return
			}
			key := columns[index]
			row[key] = value
		}
	}
	return
}

//ExecuteSchema 根据包含@名称占位符的语句执行查询语句
func (db *DB) ExecuteSchema(query string, data map[string]interface{}) (affectedRow int64, err error) {
	query, args := GetSchema(db.provider, query, data)
	return db.Execute(query, args...)
}

//Execute 执行SQL操作语句
func (db *DB) Execute(query string, args ...interface{}) (affectedRow int64, err error) {
	result, err := db.db.Exec(query, args...)
	if err != nil {
		return
	}
	
	affectedRow, err = result.RowsAffected()
	return
}

//SetLang 设置语言环境变量
func (db *DB) SetLang() {
	nlsLang := os.Getenv("NLS_LANG")
	if strings.EqualFold(nlsLang, db.lang) {
		os.Setenv("NLS_LANG", db.lang)
	}
}