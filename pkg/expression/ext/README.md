# sql 扩展模块使用

``` go

// db 为已经初始化后的数据库连接
ext.Reg(db)

pwd, _ := os.Getwd()
exp := expression.CreateExecer(pwd)

// sql 表示成当前目录下开始导入 sql 脚本，如果
exp.ScriptImportAlias("ext", "sql")

// 这样就可以在脚本中使用：
    sql.query(query,args...) sql.one(query,args...)  四个函数
// 也可以使用：
    sqlctx.QueryDB(ctx,db,query,args)  sqlctx.OneDB(ctx,db,query,args)
// 也可以使用：
    sqlctx.Query(query,args) sqlctx.One(query,args)
```