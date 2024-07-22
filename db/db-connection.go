package db

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"net/http"
	"os"
)

func InitConn(g *gin.Context) *pgx.Conn {
	conn, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))

	if err != nil {
		g.JSON(http.StatusInternalServerError, gin.H{"Internal Server Error": "Unable to Connect to db!"})
	}

	return conn
}
