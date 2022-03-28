package main

import (
	// "errors"
	"net/http"
	"github.com/gin-gonic/gin"
	"math/rand"
	"fmt"
	"strings"
)

func main() {
	server := gin.Default()
	server.GET("/", HelloWorld)
	server.GET("/create", CreateShortURL)
	// server.POST("/load", LoginAuth)
	server.Run(":8888")
}

func HelloWorld(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Hello World!"})
	return 
}

//建立 ShortUrl 資料的 Router
func CreateShortURL(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{ "status": "create succsss"})
	fmt.Println(CreateBase62Key(6))
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
func CreateBase62Key(keyLen int) string{

	//Base62字符
	base62 := "abcdefghijklmnopqrstuvwxyz"+"ABCDEFGHIJKLMNOPQRSTUVWXYZ"+"0123456789"

	//儲存Key的容器
	var keyBuilder strings.Builder
	// keyBuilder.Grow(keyLen)

	//迴圈方式產生keyLen個字元的Key
	for i := 0; i < keyLen; i++{
		//隨機選擇一個字元並加入
		base62_index := rand.Intn(len(base62))
		keyBuilder.WriteByte(base62[base62_index])
	}

	return keyBuilder.String()
}
