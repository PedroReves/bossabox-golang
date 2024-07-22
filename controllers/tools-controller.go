package controllers

import (
	"context"
	"net/http"
	"os"

	"github.com/PedroReves/bossabox-golang/model"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	_ "github.com/joho/godotenv/autoload"
)

func GetTools(g *gin.Context) {
	conn, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))

	if err != nil {
		g.JSON(http.StatusInternalServerError, gin.H{"Internal Server Error": err.Error()})
	}

	defer conn.Close(context.Background())

	rows, err := conn.Query(context.Background(), "SELECT * FROM TOOLS")

	if err != nil {
		g.JSON(http.StatusInternalServerError, gin.H{"Internal Server Error": err.Error()})
	}

	var tools []model.Tool
	for rows.Next() {
		var id, title, link, description string
		var tags []string
		if err := rows.Scan(&id, &title, &link, &description, &tags); err != nil {
			g.JSON(http.StatusInternalServerError, gin.H{"Internal Server Error": err.Error()})
		}

		tools = append(tools, model.Tool{
			Id:          id,
			Title:       title,
			Link:        link,
			Description: description,
			Tags:        tags,
		})
	}

	g.JSON(http.StatusOK, tools)
}

func CreateTool(g *gin.Context) {
	var tool model.Tool

	if err := g.ShouldBindJSON(&tool); err != nil {
		g.JSON(http.StatusBadRequest, gin.H{"Bad Request": err.Error()})
	}

	conn, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))

	if err != nil {
		g.JSON(http.StatusInternalServerError, gin.H{"Internal Server Error": err.Error()})
	}

	defer conn.Close(context.Background())

	_, err = conn.Exec(context.Background(), "INSERT INTO tools (title, link, description, tags) VALUES ($1, $2, $3, $4)", &tool.Title, &tool.Link, &tool.Description, &tool.Tags)

	if err != nil {
		g.JSON(http.StatusInternalServerError, gin.H{"Internal Server Error": err.Error()})
	}

	g.JSON(http.StatusCreated, gin.H{"Tool": tool})
}

func DeleteTool(g *gin.Context) {
	id := g.Param("id")

	conn, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))

	if err != nil {
		g.JSON(http.StatusInternalServerError, gin.H{"Internal Server Error": err.Error()})
	}

	defer conn.Close(context.Background())

	var count int

	err = conn.QueryRow(context.Background(), "SELECT COUNT(*) FROM tools WHERE id = $1", id).Scan(&count)

	if err != nil {
		g.JSON(http.StatusInternalServerError, gin.H{"Internal Server Error": err.Error()})
	}

	if count < 1 {
		g.JSON(http.StatusInternalServerError, gin.H{"Message": "This tool doesnt exist"})
		return
	}

	_, err = conn.Exec(context.Background(), "DELETE FROM tools WHERE id = $1", id)

	if err != nil {
		g.JSON(http.StatusInternalServerError, gin.H{"Internal Server Error": err.Error()})
	}

	g.JSON(http.StatusNoContent, nil)

}

func GetFilteredTool(g *gin.Context) {
	name := g.Query("name")

	conn, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))

	if err != nil {
		g.JSON(http.StatusInternalServerError, gin.H{"Internal Server Error": err.Error()})
	}

	defer conn.Close(context.Background())

	rows, err := conn.Query(context.Background(), "SELECT * FROM tools WHERE $1 = ANY(tags)", name)

	if err != nil {
		g.JSON(http.StatusInternalServerError, gin.H{"Internal Server Error": err.Error()})
	}

	var tools []model.Tool
	for rows.Next() {

		var id, link, title, description string
		var tags []string
		if err := rows.Scan(&id, &link, &description, &title, &tags); err != nil {
			g.JSON(http.StatusInternalServerError, gin.H{"Internal Server Error": err.Error()})
		}
		tools = append(tools, model.Tool{
			Id:          id,
			Title:       title,
			Link:        link,
			Description: description,
			Tags:        tags,
		})
	}

	g.JSON(http.StatusOK, tools)

}
