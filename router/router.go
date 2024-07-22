package router

import (
	"fmt"
	"github.com/PedroReves/bossabox-golang/controllers"
	"github.com/gin-gonic/gin"
)

func Initialize() {
	r := gin.Default()

	r.GET("/tools", controllers.GetTools)
	r.GET("/tool", controllers.GetFilteredTool)
	r.POST("/tools", controllers.CreateTool)
	r.DELETE("/tools/:id", controllers.DeleteTool)

	if err := r.Run(); err != nil {
		fmt.Println("Unable to start server")
	}

}
