package ext

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/xushiwei/qlang"
)


// Reg 注册数据库DB
// 有默认数据库操作
// 也支持多数据库
func Reg(db *sql.DB) {
	sqlCtx := context.Background()
	qlang.Import("sqlctx", map[string]interface{}{
		"QueryDB": QueryDB,
		"Query": func(query string, args ...interface{}) []map[string]interface{} {
			return QueryDB(sqlCtx, db, query, args...)
		},
		"OneDB": QueryOneDB,
		"One": func(query string, args ...interface{}) map[string]interface{} {
			return QueryOneDB(sqlCtx, db, query, args...)
		},
	})

}

// QueryDB 查询sql返回的所有行
func QueryDB(ctx context.Context, db *sql.DB, query string, args ...interface{}) (out []map[string]interface{}) {
	// var rows *ext.Rows
	// var err error
	rows, err := db.QueryContext(ctx, query, args...)

	if err != nil {
		panic(fmt.Sprintf("查询失败:%s  %v ==> %v", query, args, err))
	}
	cols, err := rows.Columns() // Remember to check err afterwards
	if err != nil {
		panic(fmt.Sprintf("查询失败:%s  %v ==> %v", query, args, err))
	}
	vals := make([]interface{}, len(cols))
	for i := range cols {
		vals[i] = new(sql.RawBytes)
	}
	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(vals...)
		if err != nil {
			panic(fmt.Sprintf("提取数据失败:%s  %v ==> %v", query, args, err))
		}
		vmap := make(map[string]interface{})
		for i, col := range cols {
			var s string
			rb, ok := vals[i].(*sql.RawBytes)
			if ok {
				s = string(*rb)
			}
			vmap[col] = s
		}
		out = append(out, vmap)
	}

	return
}

// QueryOneDB 查询sql返回的第一条记录
func QueryOneDB(ctx context.Context, db *sql.DB, query string, args ...interface{}) (out map[string]interface{}) {
	// var rows *ext.Rows
	// var err error
	if strings.Index(strings.ToLower(query), "limit") < 0 {
		query = "SELECT * FROM (" + query + ") AS __t__ LIMIT 1"
	}
	rows, err := db.QueryContext(ctx, query, args...)

	if err != nil {
		panic(fmt.Sprintf("查询失败:%s  %v ==> %v", query, args, err))
	}
	cols, err := rows.Columns() // Remember to check err afterwards
	if err != nil {
		panic(fmt.Sprintf("查询失败:%s  %v ==> %v", query, args, err))
	}
	vals := make([]interface{}, len(cols))
	for i := range cols {
		vals[i] = new(sql.RawBytes)
	}
	defer rows.Close()
	if rows.Next() {
		err = rows.Scan(vals...)
		if err != nil {
			panic(fmt.Sprintf("提取数据失败:%s  %v ==> %v", query, args, err))
		}
		out = make(map[string]interface{})
		for i, col := range cols {
			out[col] = vals[i]

		}
	}

	return
}
