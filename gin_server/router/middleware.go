package router

import (
	"bytes"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

const (
	keyRequestBody = "requestBodyCopy"
)

func GenerateRequestBodySaveMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		var buf bytes.Buffer
		// io.TeeReaderを使い、読み出し先とは別のbufferにもセットされるようにする
		tee := io.TeeReader(c.Request.Body, &buf)
		body, err := io.ReadAll(tee)
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		// リクエストボディは1度読み出すと空になるので、再度読み出せるようにする
		c.Request.Body = io.NopCloser(&buf)
		c.Set(keyRequestBody, body)
		c.Next()
	}
}
