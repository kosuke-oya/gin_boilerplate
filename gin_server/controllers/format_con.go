//go:generate mockgen -source=$GOFILE -package=mock_$GOPACKAGE -destination=./mock/$GOFILE
package controllers

import (
	"context"
	"gin_server/dto"
	"gin_server/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type IFormatCon interface {
	Create(c *gin.Context)
}

type FormatCon struct {
	FormatService services.IFormatService
}

var _ IFormatCon = (*FormatCon)(nil)

func NewFormatCon(c context.Context) *FormatCon {
	return &FormatCon{
		FormatService: services.NewFormatService(c),
	}
}

func (f *FormatCon) Create(c *gin.Context) {
	var req dto.FormatDataRequest
	var res dto.FormatDataResponse
	if c.ShouldBindJSON(&req) != nil {
		res.Error = "Invalid request"
		c.JSON(http.StatusBadRequest, res)
		return
	}

	data, err := f.FormatService.Create(c, req)
	if err != nil {
		res.Error = err.Error()
		c.JSON(http.StatusInternalServerError, res)
		return
	}

	res.Data = data

	c.JSON(http.StatusOK, res)
}
