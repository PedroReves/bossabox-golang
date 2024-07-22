package db

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
)

func InitConn(g *gin.Context) *pgx.Conn {
	err := godotenv.Load()

	if err != nil {
		fmt.Printf("Unable to parse the env files %v", err)
	}
	conn, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))

	if err != nil {
		g.JSON(http.StatusInternalServerError, gin.H{"Internal Server Error": "Unable to Connect to db!"})
	}

	return conn
}
