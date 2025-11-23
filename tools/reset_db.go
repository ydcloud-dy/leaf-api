package main

import (
	"fmt"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	// 连接到MySQL服务器（不指定数据库）
	dsn := "root:123456@tcp(127.0.0.1:3306)/?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to MySQL: ", err)
	}

	// 删除数据库
	fmt.Println("Dropping database leaf_admin...")
	if err := db.Exec("DROP DATABASE IF EXISTS leaf_admin").Error; err != nil {
		log.Fatal("Failed to drop database: ", err)
	}

	// 创建数据库
	fmt.Println("Creating database leaf_admin...")
	if err := db.Exec("CREATE DATABASE leaf_admin CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci").Error; err != nil {
		log.Fatal("Failed to create database: ", err)
	}

	fmt.Println("Database reset successfully!")
	fmt.Println("Now you can start the application and use admin/admin123 to login.")
}
