package main

import (
	"demo/router"

	"net/http"
	"os"

	_ "github.com/denisenkom/go-mssqldb"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/joho/godotenv/autoload"
)

//讀取環境變數
var (
	SERVER   = os.Getenv("SERVER")
	PORT     = os.Getenv("PORT")
	USERNAME = os.Getenv("USERNAME")
	PASSWORD = os.Getenv("PASSWORD")
	DATABASE = os.Getenv("DATABASE")
)

func main() {
	//建立Web API Server
	server := gin.Default()
	server.GET("/", HelloWorld)
	server.POST("/create", router.CreateShortURL)
	server.GET("/load/:key", router.LoadShortURL)
	server.Run(":8080")
}

func HelloWorld(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Hello World!"})
	return
}
