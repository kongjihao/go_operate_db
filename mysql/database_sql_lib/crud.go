// package main

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql" // 注意这里要导入mysql驱动，init()函数会自动注册mysql驱动
)

var db *sql.DB

type user struct {
	id   int64
	name string
	age  uint32
}

func initMysql() (err error) {
	// DB DSN (Data Source Name)
	dsn := "root:12345678@tcp(127.0.0.1:3306)/sql_demo?charset=utf8mb4&parseTime=True&loc=Local"

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
	// db.SetConnMaxLifetime(0) // 设置连接的最大可复用时间，0表示无限制
	db.SetMaxOpenConns(20) // 设置与数据库建立的最大连接数
	db.SetMaxIdleConns(2)  // 设置连接池中空闲连接的最大数量

	return nil
}

// 查询单条记录示例
func queryRowDemo() {
	sqlStr := "select id, name, age from user where id = ?"
	var u user
	row := db.QueryRow(sqlStr, 1)
	err := row.Scan(&u.id, &u.name, &u.age) // 把值scan出来，只有scan后db的连接资源才会被释放，所以当使用database/sql库的时候query后一定要scan
	if err != nil {
		fmt.Printf("query row failed, err:%v\n", err)
		return
	}

	fmt.Printf("query row success: id:%d name:%s age:%d\n", u.id, u.name, u.age)
}

// 查询多行
func queryMultiRowDemo() {
	sqlStr := "select id, name, age from user where id > ?"
	rows, err := db.Query(sqlStr, 0)
	if err != nil {
		fmt.Printf("query multi row failed, err:%v\n", err)
		return
	}
	// 注意关闭rows释放持有的数据库连接，避免连接泄漏
	defer rows.Close() // 非常重要，因为不一定能保证下面的 for 循环一定执行完，防止连接泄露

	// 循环读取结果集中的数据
	for rows.Next() {
		var u user
		err := rows.Scan(&u.id, &u.name, &u.age)
		if err != nil {
			fmt.Printf("scan row failed, err:%v\n", err)
			return
		}
		fmt.Printf("query rows success: id:%d name:%s age:%d\n", u.id, u.name, u.age)
	}
}

// 插入数据
func insertRowDemo() {
	sqlStr := "insert into user(name, age) values (?, ?)"
	ret, err := db.Exec(sqlStr, "小王子", 18)
	if err != nil {
		fmt.Println("insert failed, err:%v\n", err)
		return
	}

	theID, err := ret.LastInsertId() // 返回新插入数据的ID
	if err != nil {
		fmt.Println("get lastinsert ID failed, err:%v\n", err)
		return
	}
	fmt.Println("insert success, the id is: ", theID)
}

// 更新数据
func updateRowDemo() {
	sqlStr := "update user set name = ? where id = ?"
	ret, err := db.Exec(sqlStr, "xiaowangba", 3)
	if err != nil {
		fmt.Println("update failed, err:%v\n", err)
		return
	}

	n, err := ret.RowsAffected() // 操作影响的行数
	if err != nil {
		fmt.Println("get RowsAffected failed, err:%v\n", err)
		return
	}

	fmt.Println("update success, affected rows:", n)
}

// 删除数据
func deleteRowDemo() {
	sqlStr := "delete from user where id = ?"
	ret, err := db.Exec(sqlStr, 3)
	if err != nil {
		fmt.Println("delete failed, err:%v\n", err)
		return
	}

	n, err := ret.RowsAffected() // 操作影响的行数
	if err != nil {
		fmt.Println("get RowsAffected failed, err:%v\n", err)
		return
	}

	fmt.Println("delete success, affected rows:", n)
}

func main() {
	if err := initMysql(); err != nil {
		fmt.Println("init mysql failed, err:%v\n", err)
	}

	defer db.Close() // 注意这个defer关闭数据库连接相关资源，要写main函数中，只有当main函数结束了才释放数据库连接资源
	fmt.Println("connet mysql success")

	// updateRowDemo()
	// insertRowDemo()
	// queryRowDemo()
	// deleteRowDemo()
	queryMultiRowDemo()
}
