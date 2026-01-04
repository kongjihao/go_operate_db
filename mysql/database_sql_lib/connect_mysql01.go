// package main

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql" // 注意这里要导入mysql驱动，init()函数会自动注册mysql驱动
)

func main() {
	// DB DSN (Data Source Name)
	dsn := "root:12345678@tcp(127.0.0.1:3306)/sql_demo?charset=utf8mb4&parseTime=True&loc=Local"

	// Use the DSN to connect to the database
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		panic(err)
	}

	// 做完错误检查之后，确保db不为nil，当连接数据库出错时也能正常关闭数据库的连接资源
	defer db.Close() // 注意这个defer关闭数据库连接相关资源，要写在err判断的下面，防止出现第三方库打开db为nil的情况

	fmt.Println("connet mysql success ???")
}
