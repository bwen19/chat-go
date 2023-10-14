package util

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

const (
	foreignKeyViolation = "23503"
	uniqueViolation     = "23505"
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
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		if pgErr.Code == uniqueViolation {
			ctx.String(http.StatusForbidden, "data already exists")
			return
		}
	}
	ctx.String(http.StatusInternalServerError, "internal database error")
}

func RecordNotFoundResponse(ctx *gin.Context, err error) {
	if errors.Is(err, pgx.ErrNoRows) {
		ctx.String(http.StatusNotFound, "record not found")
	} else {
		ctx.String(http.StatusInternalServerError, "internal database error")
	}
}
