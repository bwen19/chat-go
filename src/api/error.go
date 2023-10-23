package api

import (
	"errors"
	"gochat/src/db"
	"net/http"

	"github.com/gin-gonic/gin"
)

var (
	errArgument = errors.New("invalid argument")
	errDenied   = errors.New("permission denied")
	errInternal = errors.New("internal server error")
	errNotFound = errors.New("not found")
)

func InvalidArgumentResponse(ctx *gin.Context) {
	ctx.String(http.StatusBadRequest, errArgument.Error())
}

func InternalErrorResponse(ctx *gin.Context) {
	ctx.String(http.StatusInternalServerError, errInternal.Error())
}

func PermissionDeniedResponse(ctx *gin.Context) {
	ctx.String(http.StatusForbidden, errDenied.Error())
}

func ViolationResponse(ctx *gin.Context, err error) {
	if db.HasViolation(err) {
		ctx.String(http.StatusForbidden, "already exists")
		return
	}
	ctx.String(http.StatusInternalServerError, errInternal.Error())
}

func NotFoundResponse(ctx *gin.Context, err error) {
	if errors.Is(err, db.ErrRecordNotFound) {
		ctx.String(http.StatusNotFound, errNotFound.Error())
		return
	}
	ctx.String(http.StatusInternalServerError, errInternal.Error())
}
