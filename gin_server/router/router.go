package router

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"gin_server/controllers"
	"gin_server/utils"

	"github.com/gin-contrib/cors"
	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/swag/example/basic/docs"
	"go.uber.org/zap"
)

func SetupRouter() *gin.Engine {
	// envファイルの初期化
	utils.Init()

	// set zap logger as default logger
	logger, err := zap.NewProduction()
	if err != nil {
		os.Exit(1)
	}

	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowMethods:    []string{"GET", "POST", "OPTIONS"},
		AllowAllOrigins: true,
		AllowWebSockets: true,
		MaxAge:          12 * time.Hour,
	}))
	r.Use(GenerateRequestBodySaveMiddleware())
	r.Use(ginzap.GinzapWithConfig(logger, &ginzap.Config{
		UTC:        true,
		TimeFormat: time.RFC3339,
		Context: ginzap.Fn(func(c *gin.Context) []zap.Field {
			start := time.Now()

			bodyBytes, _ := c.Get(keyRequestBody)
			body := bodyBytes.([]byte)

			// CloudLoggingに送信用のメッセージを定義
			heaserBytes, err := json.Marshal(c.Request.Header)
			if err != nil {
				fmt.Println("failed to marshal header")
			}
			return []zap.Field{
				zap.Int("status", c.Writer.Status()),
				zap.Int64("content_length", c.Request.ContentLength),
				zap.String("method", c.Request.Method),
				zap.String("path", c.Request.URL.Path),
				zap.String("query", c.Request.URL.RawQuery),
				zap.ByteString("body_uryu_data", body),
				zap.String("ip", c.ClientIP()),
				zap.String("user_agent", c.Request.UserAgent()),
				zap.String("errors", c.Errors.ByType(gin.ErrorTypePrivate).String()),
				zap.Duration("elapsed", time.Since(start)),
				zap.ByteString("header", heaserBytes),
				zap.Int("response_size(bytes)", c.Writer.Size()),
			}
		}),
	}))
	docs.SwaggerInfo.Title = "API Docs"
	docs.SwaggerInfo.Description = "This is a http server."
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Schemes = []string{"http", "https"}

	// swaggerのURL(HTTP APIドキュメント)
	// use ginSwagger middleware to serve the API docs
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// PubSubメッセージの受信するcontrollerのセットアップ
	ctx := context.Background()
	controller := controllers.NewFormatCon(ctx)

	// PubSubメッセージの受信するエンドポイント
	r.POST("/", controller.Create)

	return r
}
