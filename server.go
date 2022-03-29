package main

import (
	// "errors"
	"net/http"
	"github.com/gin-gonic/gin"
	"math/rand"
	"fmt"
	"strings"
	"database/sql"
	"time"
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

var (

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
		c.JSON(http.StatusBadRequest, gin.H{ "error": err.Error()}) 
		return
	}
	//離開此函式時，關閉資料庫
	defer db.Close()
	

	//取得前端傳來的資訊
	type Query_Json struct {
		OriginalUrl string `json:"originalUrl"`
	}
	var query Query_Json

	if err := c.ShouldBindJSON(&query); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{ "error": err.Error()}) 
		return
	}

	
	//如果傳入資訊缺少originalUrl，回傳錯誤訊息
	if query.OriginalUrl == "" {
		c.JSON(http.StatusBadRequest, gin.H{ "error": "undefined originalUrl"}) 
		return
	}


	//檢查此網址是否已經建立，若已建立則回傳該Key，並重置過期時間
	var exist_data Data
	row := db.QueryRow("SELECT * FROM datas WHERE original_url=?", query.OriginalUrl)
	err = row.Scan(&exist_data.originalUrl, &exist_data.shortUrl_key, &exist_data.create_date, &exist_data.expire_date, &exist_data.call_time); 
	//如果此網址已在資料庫中
	if err != sql.ErrNoRows {
		
		todayStr := time.Now().Format("2006-01-02")
		expiredateStr := time.Now().AddDate(3, 0, 0).Format("2006-01-02") //增加三年

		//更新設定日期及期限
		_, err := db.Exec("UPDATE datas SET create_date = ?, expire_date = ? WHERE shortUrl_key = ?", todayStr, expiredateStr, exist_data.shortUrl_key)

		if err != nil{
			c.JSON(http.StatusBadRequest, gin.H{ "error": err.Error()}) 
			return
		}else{
			c.JSON(http.StatusOK, gin.H{ 
				"message": "short url created successfully",
				"shortURL": c.Request.Host + "/load/" + exist_data.shortUrl_key,
				"expire_date": expiredateStr,
			});
			return
		}

	//如果是尚未登陸的網址
	}else{

		var newKey string

		//在GO裡面沒有While迴圈概念，只能用for執行
		for{
			//產生新的Key
			newKey = CreateBase62Key(6)
			//檢查是否使用過
			err := db.QueryRow("SELECT * FROM datas WHERE shortUrl_key=?", newKey).Scan()
			//如果找不到Row代表沒用過，跳脫For迴圈
			if err == sql.ErrNoRows{
				break;
			}
		}
		
		_, err := db.Exec(
			"INSERT INTO datas (original_url, shortUrl_key, create_date, expire_date, call_time) VALUES (?, ?, ?, ?, ?)",
			query.OriginalUrl,
			newKey,
			time.Now().Format("2006-01-02"),
			time.Now().AddDate(3, 0, 0).Format("2006-01-02"),
			0,
		)

		if err != nil{
			c.JSON(http.StatusBadRequest, gin.H{ "error": err.Error()}) 
			return
		}else{
			c.JSON(http.StatusOK, gin.H{ 
				"message": "short url created successfully",
				"shortURL": c.Request.Host + "/load/" + newKey,
				"expire_date": time.Now().AddDate(3, 0, 0).Format("2006-01-02"),
			});
			return
		}
	}
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
