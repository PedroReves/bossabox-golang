package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
)

type Tool struct {
	Id          string   `json:"id"`
	Title       string   `json:"title"`
	Link        string   `json:"link"`
	Description string   `json:"description"`
	Tags        []string `json:"tags"`
}

func main() {
	r := gin.Default()

	r.GET("/tools", getTools)
	r.GET("/tool", getFilteredTool)
	r.POST("/tools", createTool)
	r.DELETE("/tools/:id", deleteTool)

	if err := r.Run(); err != nil {
		fmt.Println("Unable to start server")
	}
}

func getTools(g *gin.Context) {
	conn, err := pgx.Connect(context.Background(), "postgres://docker:root@localhost:5432/tools")

	if err != nil {
		g.JSON(http.StatusInternalServerError, gin.H{"Internal Server Error": err.Error()})
	}

	defer conn.Close(context.Background())

	rows, err := conn.Query(context.Background(), "SELECT * FROM TOOLS")

	if err != nil {
		g.JSON(http.StatusInternalServerError, gin.H{"Internal Server Error": err.Error()})
	}

	var tools []Tool
	for rows.Next() {
		var id, title, link, description string
		var tags []string
		if err := rows.Scan(&id, &title, &link, &description, &tags); err != nil {
			g.JSON(http.StatusInternalServerError, gin.H{"Internal Server Error": err.Error()})
		}

		tools = append(tools, Tool{
			Id:          id,
			Title:       title,
			Link:        link,
			Description: description,
			Tags:        tags,
		})
	}

	g.JSON(http.StatusOK, tools)
}

func createTool(g *gin.Context) {
	var tool Tool

	if err := g.ShouldBindJSON(&tool); err != nil {
		g.JSON(http.StatusBadRequest, gin.H{"Bad Request": err.Error()})
	}

	conn, err := pgx.Connect(context.Background(), "postgres://docker:root@localhost:5432/tools")

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

func deleteTool(g *gin.Context) {
	id := g.Param("id")

	conn, err := pgx.Connect(context.Background(), "postgres://docker:root@localhost:5432/tools")

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

func getFilteredTool(g *gin.Context) {
	name := g.Query("name")

	conn, err := pgx.Connect(context.Background(), "postgres://docker:root@localhost:5432/tools")

	if err != nil {
		g.JSON(http.StatusInternalServerError, gin.H{"Internal Server Error": err.Error()})
	}

	defer conn.Close(context.Background())

	rows, err := conn.Query(context.Background(), "SELECT * FROM tools WHERE $1 = ANY(tags)", name)

	if err != nil {
		g.JSON(http.StatusInternalServerError, gin.H{"Internal Server Error": err.Error()})
	}

	var tools []Tool
	for rows.Next() {

		var id, link, title, description string
		var tags []string
		if err := rows.Scan(&id, &link, &description, &title, &tags); err != nil {
			g.JSON(http.StatusInternalServerError, gin.H{"Internal Server Error": err.Error()})
		}
		tools = append(tools, Tool{
			Id:          id,
			Title:       title,
			Link:        link,
			Description: description,
			Tags:        tags,
		})
	}

	g.JSON(http.StatusOK, tools)

}
