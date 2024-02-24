package db

import (
	"fmt"
	"os"
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/mysqldialect"
	"github.com/uptrace/bun/extra/bundebug"
)

var DB *bun.DB

func Init() {
  // api.envに定義したDB関係の環境変数を取得
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
  // tcp（）の中にdocker-composeで定義したDB用コンテナのサービス名を入れれば、
  // 自動的にホストとポートを読み取ってくれる
	dsn := fmt.Sprintf(
		"%s:%s@tcp(db)/%s?charset=utf8mb4&parseTime=true&loc=Local",
		dbUser,
		dbPassword,
		dbName,
	)
	sqldb, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal(err)
	}
	//defer sqldb.Close()

	DB = bun.NewDB(sqldb, mysqldialect.New())
	DB.AddQueryHook(bundebug.NewQueryHook(
		bundebug.WithVerbose(true),
		bundebug.FromEnv("BUNDEBUG"),
	))
}