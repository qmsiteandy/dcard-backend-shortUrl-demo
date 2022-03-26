package main

import (
	// "errors"
	"net/http"
	"github.com/gin-gonic/gin"
	// "fmt"
)

func main() {
	server := gin.Default()
	//設定靜態資源的讀取
	server.GET("/create", CreateShortURL)
	// server.POST("/load", LoginAuth)
	server.Run(":8888")
}

//建立 ShortUrl 資料的 Router
func CreateShortURL(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{ "status": "create succsss"})
	return 
}

// func LoginAuth(c *gin.Context) {
// 	var (
// 		username string
// 		password string
// 	)
// 	if in, isExist := c.GetPostForm("username"); isExist && in != "" {
// 		username = in
// 	} else {
// 		c.HTML(http.StatusBadRequest, "login.html", gin.H{
// 			"error": errors.New("必須輸入使用者名稱"),
// 		})
// 		return
// 	}
// 	if in, isExist := c.GetPostForm("password"); isExist && in != "" {
// 		password = in
// 	} else {
// 		c.HTML(http.StatusBadRequest, "login.html", gin.H{
// 			"error": errors.New("必須輸入密碼名稱"),
// 		})
// 		return
// 	}
// 	if err := Auth(username, password); err == nil {
// 		c.HTML(http.StatusOK, "login.html", gin.H{
// 			"success": "登入成功",
// 		})
// 		return
// 	} else {
// 		c.HTML(http.StatusUnauthorized, "login.html", gin.H{
// 			"error": err,
// 		})
// 		return
// 	}
// }

//以Hash方式產生Key
func CreateKey() string{
	return "123"
}
