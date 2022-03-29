package main

import (
	// "errors"
	"net/http"
	"github.com/gin-gonic/gin"
	"math/rand"
	"fmt"
	"strings"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

//資料庫連線資料
const (
	USERNAME = "Remon"
	PASSWORD = "andy0709"
	NETWORK  = "tcp"
	SERVER   = "127.0.0.1"
	PORT     = 3306
	DATABASE = "dcard-backend-shorturl"
)

//資料庫結構
type Data struct {
	originalUrl string
	shortUrl_key string
	create_date string
	expire_date string
	call_time int
}

func main() {

	//建立Web API Server
	server := gin.Default()
	server.GET("/", HelloWorld)
	server.POST("/create", CreateShortURL)
	//server.POST("/load", LoginAuth)
	server.Run(":8888")
}

func HelloWorld(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Hello World!"})
	return 
}

//建立 ShortUrl 資料的 Router
func CreateShortURL(c *gin.Context) {

	//建立資料庫連線
	conn := fmt.Sprintf("%s:%s@%s(%s:%d)/%s", USERNAME, PASSWORD, NETWORK, SERVER, PORT, DATABASE)
	db, err := sql.Open("mysql", conn)
	if err != nil {
		fmt.Println("開啟SQL資料庫連線錯誤：", err)
		return
	}
	if err := db.Ping(); err != nil {
		fmt.Println("資料庫連線錯誤：", err.Error())
		return
	}
	//離開此函式時，關閉資料庫
	defer db.Close()
	

	//取得前端傳來的資訊
	type Query_Json struct {
		OriginalUrl string `json:"originalUrl"`
	}
	var query Query_Json

	if err := c.ShouldBindJSON(&quest_json); err != nil {
		fmt.Println("err:", err)
	}

	
	//如果傳入資訊缺少originalUrl，回傳錯誤訊息
	if query.OriginalUrl == "" {
		c.JSON(http.StatusBadRequest, gin.H{ "status": "undefined originalUrl"}) 
	}


	//檢查此網址是否已經建立，若已建立則回傳該Key，並重置過期時間
	row := db.QueryRow("select * from datas where original_url=?", query.originalUrl)
	if err := row.Scan(); err != nil {c.JSON(http.StatusBadRequest, gin.H{ "status": err})}
		

	//產生新的Key並檢查是否用過


	
	// fmt.Println()
	
	

	

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
