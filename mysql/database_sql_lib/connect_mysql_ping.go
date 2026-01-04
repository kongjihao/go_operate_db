// package main

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql" // 注意这里要导入mysql驱动，init()函数会自动注册mysql驱动
)

var db *sql.DB

func initMysql() (err error) {
	// DB DSN (Data Source Name)
	dsn := "root:12345678@tcp(127.0.0.1:3306)/sql_demo?charset=utf8mb4&parseTime=True&loc=Local"

	// 此处db为全局变量，不要用 := 初始化来赋值
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		panic(err)
	}

	err = db.Ping()
	if err != nil {
		fmt.Println("connect mysql failed, err:%v\n", err)
		return err
	}

	// 设置数据库连接池参数,根据应用需求进行调整
	db.SetConnMaxLifetime(0) // 设置连接的最大可复用时间，0表示无限制
	db.SetMaxOpenConns(200)  // 设置与数据库建立的最大连接数
	db.SetMaxIdleConns(10)   // 设置连接池中空闲连接的最大数量

	return nil
}

func main() {
	if err := initMysql(); err != nil {
		fmt.Println("init mysql failed, err:%v\n", err)
	}

	defer db.Close() // 注意这个defer关闭数据库连接相关资源，要写main函数中，只有当main函数结束了才释放数据库连接资源
	fmt.Println("connet mysql success")
	// db.xxx()去使用数据库操作
}
