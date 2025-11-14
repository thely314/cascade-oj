package main

import (
	"cascade-oj/ent"
	"context"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	// 连接数据库
	client, err := ent.Open("mysql", "cascade_oj:123456@tcp(localhost:3306)/cascade_oj_db?parseTime=True")
	if err != nil {
		log.Fatalf("failed opening connection to mysql: %v", err)
	}
	defer client.Close()

	// 创建所有表（根据 schema）并确保迁移生效
	ctx := context.Background()
	if err := client.Schema.Create(ctx); err != nil {
		log.Fatalf("failed creating schema resources: %v", err)
	}

	log.Println("Database schema created successfully")
}
