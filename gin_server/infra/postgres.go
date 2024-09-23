//go:generate mockgen -source=$GOFILE -package=mock_$GOPACKAGE -destination=./mock/$GOFILE
package infra

import (
	"context"
	"database/sql"
	"encoding/base64"
	"fmt"
	"log"
	"net"
	"os"

	"gin_server/models"
	"gin_server/utils"

	"cloud.google.com/go/cloudsqlconn"
	_ "github.com/GoogleCloudPlatform/cloudsql-proxy/proxy/dialers/postgres"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/stdlib"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type IMyPostgres interface {
	Create(ctx context.Context, st models.FormatData) error
}

type MyPostgres struct {
	DB *gorm.DB
}

var _ IMyPostgres = (*MyPostgres)(nil)

// 環境変数の違いで読み込む.envファイルを変更
func init() {
	utils.Init()
}

var (
	dbHost      = os.Getenv("DB_HOST")
	dbUser      = os.Getenv("DB_USER")
	dbPassword  = os.Getenv("DB_PASSWORD")
	dbName      = os.Getenv("DB_NAME")
	port        = os.Getenv("DB_PORT")
	cloudSqlKey = os.Getenv("CLOUD_SQL_KEY")
)

func NewDB(ctx context.Context) *MyPostgres {
	var dsn string
	var db *gorm.DB
	var err error
	if os.Getenv("ENV") == "prod" {
		sqlDb, err := connectWithConnector(ctx)
		if err != nil {
			panic(err)
		}

		// sql.DBからGORMのDBインスタンスを作成
		db, err = gorm.Open(postgres.New(postgres.Config{
			Conn: sqlDb,
		}), &gorm.Config{})
		if err != nil {
			panic(err)
		}

	} else {
		// localのpostgresへの接続
		dsn = fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Tokyo",
			dbHost, dbUser, dbPassword, dbName, port)
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err != nil {
			panic(err)
		}
	}

	// 追加: PostgreSQL に拡張機能 "uuid-ossp" を追加する
	db.Exec(`CREATE EXTENSION IF NOT EXISTS "uuid-ossp";`)
	p := new(MyPostgres)
	p.DB = db

	return p
}

func (d *MyPostgres) Create(ctx context.Context, u models.FormatData) error {
	return d.DB.Create(&u).Error
}
func connectWithConnector(ctx context.Context) (*sql.DB, error) {
	mustGetenv := func(k string) string {
		v := os.Getenv(k)
		if v == "" {
			log.Fatalf("Fatal Error in connect_connector.go: %s environment variable not set.\n", k)
		}
		return v
	}
	// Note: Saving credentials in environment variables is convenient, but not
	// secure - consider a more secure solution such as
	// Cloud Secret Manager (https://cloud.google.com/secret-manager) to help
	// keep passwords and other secrets safe.
	var (
		dbUser                 = mustGetenv("DB_USER")                  // e.g. 'my-db-user'
		dbPassword             = mustGetenv("DB_PASSWORD")              // e.g. 'my-db-password'
		dbName                 = mustGetenv("DB_NAME")                  // e.g. 'my-database'
		instanceConnectionName = mustGetenv("INSTANCE_CONNECTION_NAME") // e.g. 'project:region:instance'
		usePrivate             = os.Getenv("PRIVATE_IP")
	)

	dsn := fmt.Sprintf("user=%s password=%s database=%s", dbUser, dbPassword, dbName)
	config, err := pgx.ParseConfig(dsn)
	if err != nil {
		return nil, err
	}
	var opts []cloudsqlconn.Option
	if usePrivate != "" {
		opts = append(opts, cloudsqlconn.WithDefaultDialOptions(cloudsqlconn.WithPrivateIP()))
	}
	// Create a new dialer with any options
	key_bytes, err := base64.StdEncoding.DecodeString(cloudSqlKey)
	if err != nil {
		panic(err)
	}
	op := cloudsqlconn.WithCredentialsJSON(key_bytes)
	opts = append(opts, cloudsqlconn.WithIAMAuthN())
	opts = append(opts, op)

	d, err := cloudsqlconn.NewDialer(
		ctx,
		opts...,
	)
	if err != nil {
		return nil, err
	}
	// Use the Cloud SQL connector to handle connecting to the instance.
	// This approach does *NOT* require the Cloud SQL proxy.
	config.DialFunc = func(ctx context.Context, network, instance string) (net.Conn, error) {
		return d.Dial(ctx, instanceConnectionName)
	}
	dbURI := stdlib.RegisterConnConfig(config)
	dbPool, err := sql.Open("pgx", dbURI)
	if err != nil {
		return nil, fmt.Errorf("sql.Open: %w", err)
	}
	return dbPool, nil
}
