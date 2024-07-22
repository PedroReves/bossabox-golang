package controllers

import (
	"context"
	"github.com/PedroReves/bossabox-golang/db"
	"github.com/PedroReves/bossabox-golang/model"
	"github.com/gin-gonic/gin"
	_ "github.com/joho/godotenv/autoload"
	"net/http"
)

func GetTools(g *gin.Context) {
	conn := db.InitConn(g)

	defer conn.Close(context.Background())

	rows, err := conn.Query(context.Background(), "SELECT * FROM TOOLS")

	if err != nil {
		g.JSON(http.StatusInternalServerError, gin.H{"Internal Server Error": "Unable to Finish Query!"})
	}

	var tools []model.Tool
	for rows.Next() {
		var tool model.Tool
		if err := rows.Scan(&tool.Id, &tool.Title, &tool.Link, &tool.Description, &tool.Tags); err != nil {
			g.JSON(http.StatusInternalServerError, gin.H{"Internal Server Error": "Unable to list Tools!"})
		}

		tools = append(tools, tool)
	}

	g.JSON(http.StatusOK, tools)
}

func CreateTool(g *gin.Context) {
	var tool model.Tool

	if err := g.ShouldBindJSON(&tool); err != nil {
		g.JSON(http.StatusBadRequest, gin.H{"Bad Request": "There is an error with the request, Try Again!"})
	}

	conn := db.InitConn(g)

	defer conn.Close(context.Background())

	_, err := conn.Exec(context.Background(), "INSERT INTO tools (title, link, description, tags) VALUES ($1, $2, $3, $4)", &tool.Title, &tool.Link, &tool.Description, &tool.Tags)

	if err != nil {
		g.JSON(http.StatusInternalServerError, gin.H{"Internal Server Error": "Unable to Finish Query!"})
	}

	g.JSON(http.StatusCreated, gin.H{"Tool": tool})
}

func DeleteTool(g *gin.Context) {
	id := g.Param("id")

	conn := db.InitConn(g)

	defer conn.Close(context.Background())

	var count int

	err := conn.QueryRow(context.Background(), "SELECT COUNT(*) FROM tools WHERE id = $1", id).Scan(&count)

	if err != nil {
		g.JSON(http.StatusInternalServerError, gin.H{"Internal Server Error": "Unable to Finish Query!"})
	}

	if count < 1 {
		g.JSON(http.StatusNotFound, gin.H{"Message": "This Tool was not Found in the Database!"})
		return
	}

	_, err = conn.Exec(context.Background(), "DELETE FROM tools WHERE id = $1", id)

	if err != nil {
		g.JSON(http.StatusInternalServerError, gin.H{"Internal Server Error": "Unable to Finish Query"})
	}

	g.JSON(http.StatusNoContent, nil)

}

func GetFilteredTool(g *gin.Context) {
	name := g.Query("name")

	conn := db.InitConn(g)

	defer conn.Close(context.Background())

	rows, err := conn.Query(context.Background(), "SELECT * FROM tools WHERE $1 = ANY(tags)", name)

	if err != nil {
		g.JSON(http.StatusInternalServerError, gin.H{"Internal Server Error": "Unable to Finish Query"})
	}

	var tools []model.Tool
	for rows.Next() {
		var tool model.Tool
		if err := rows.Scan(&tool.Id, &tool.Link, &tool.Description, &tool.Title, &tool.Tags); err != nil {
			g.JSON(http.StatusInternalServerError, gin.H{"Internal Server Error": "Unable to List Filtered Tools"})
		}
		tools = append(tools, tool)
	}

	g.JSON(http.StatusOK, tools)

}
