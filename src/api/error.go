package api

import (
	"errors"
	"gochat/src/db"
	"net/http"

	"github.com/gin-gonic/gin"
)

var (
	ErrUnmarshal = errors.New("json parsing error")
	ErrInternal  = errors.New("internal server error")
)

func InvalidArgumentResponse(ctx *gin.Context) {
	ctx.String(http.StatusBadRequest, "invalid argument")
}

func InternalErrorResponse(ctx *gin.Context) {
	ctx.String(http.StatusInternalServerError, "internal server error")
}

func PermissionDeniedResponse(ctx *gin.Context) {
	ctx.String(http.StatusForbidden, "permission denied")
}

func UniqueViolationResponse(ctx *gin.Context, err error) {
	if db.HasViolation(err) {
		ctx.String(http.StatusForbidden, "already exists")
		return
	}
	ctx.String(http.StatusInternalServerError, "internal database error")
}

func RecordNotFoundResponse(ctx *gin.Context, err error) {
	if errors.Is(err, db.ErrRecordNotFound) {
		ctx.String(http.StatusNotFound, "not found")
		return
	}
	ctx.String(http.StatusInternalServerError, "internal database error")
}
