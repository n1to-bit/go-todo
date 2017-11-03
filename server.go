package main

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"os"
)

type DB struct {
	*sql.DB
}

type Tx struct {
	*sql.Tx
}

func Open(dataSourceName string) (*DB, error) {
	db, err := sql.Open("mysql", dataSourceName)
	if err != nil {
		return nil, err
	}
	return &DB{db}, nil
}

func (db *DB) Begin() (*Tx, error) {
	tx, err := db.DB.Begin()
	if err != nil {
		return nil, err
	}
	return &Tx{tx}, nil
}

func main() {
	var (
		dbUser     = "root"
		dbPassword = ""
		dbHost     = "127.0.0.1"
		dbName     = "todo_with_gin"
	)
	db, err := Open(dbUser + ":" + dbPassword + "@tcp(" + dbHost + ":3306)/" + dbName)
	defer db.Close()

	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	todoResource := &TodoResource{db: db}

	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.String(200, "pong")
	})
	r.GET("/todos", todoResource.GetAllTodos)
	r.GET("/todos/:id", todoResource.GetTodo)
	r.POST("/todos", todoResource.CreateTodo)
	r.Run(":8080")
}
