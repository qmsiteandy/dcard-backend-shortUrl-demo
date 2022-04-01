package router

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

//讀取環境變數
var (
	SERVER   = os.Getenv("SERVER")
	PORT     = os.Getenv("PORT")
	USERNAME = os.Getenv("USERNAME")
	PASSWORD = os.Getenv("PASSWORD")
	DATABASE = os.Getenv("DATABASE")
)

//資料庫結構
type Data struct {
	originalUrl  string
	shortUrl_key string
	create_date  string
	expire_date  string
	call_time    int
}

//建立 ShortUrl 資料的 Router
func CreateShortURL(c *gin.Context) {

	// 與 Azure Database Server 連線
	connString := fmt.Sprintf("server=%s;user id=%s;password=%s;port=%s;database=%s;", SERVER, USERNAME, PASSWORD, PORT, DATABASE)

	var db *sql.DB
	var err error
	// Create connection pool
	db, err = sql.Open("sqlserver", connString)
	if err != nil {
		fmt.Printf(err.Error())
		log.Fatal("Error creating connection pool: ", err.Error())
	}
	ctx := context.Background()
	err = db.PingContext(ctx)
	if err != nil {
		fmt.Printf(err.Error())
		log.Fatal(err.Error())
	}
	fmt.Println("Connected!")
	defer db.Close()

	//取得前端傳來的資訊
	type Query_Json struct {
		OriginalUrl string `json:"originalUrl"`
	}
	var query Query_Json

	if err := c.ShouldBindJSON(&query); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	//如果傳入資訊缺少originalUrl，回傳錯誤訊息
	if query.OriginalUrl == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "undefined originalUrl"})
		return
	}

	//檢查此網址是否已經建立，若已建立則回傳該Key，並重置過期時間
	var exist_data Data
	row := db.QueryRow("SELECT * FROM DemoTable WHERE original_url = @url", sql.Named("url", query.OriginalUrl))
	err = row.Scan(&exist_data.originalUrl, &exist_data.shortUrl_key, &exist_data.create_date, &exist_data.expire_date, &exist_data.call_time)

	//如果此網址已在資料庫中
	if err != sql.ErrNoRows {

		todayStr := time.Now().Format("2006-01-02")
		expiredateStr := time.Now().AddDate(3, 0, 0).Format("2006-01-02") //增加三年

		//更新設定日期及期限
		_, err := db.Exec("UPDATE DemoTable SET create_date = @createDate, expire_date = @expireDate WHERE shortUrl_key = @key",
			sql.Named("createDate", todayStr), sql.Named("expireDate", expiredateStr), sql.Named("key", exist_data.shortUrl_key))

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return

		} else {
			c.JSON(http.StatusOK, gin.H{
				"message":     "short url created successfully",
				"shortURL":    c.Request.Host + "/load/" + exist_data.shortUrl_key,
				"expire_date": expiredateStr,
			})
			return
		}

		//如果是尚未登陸的網址
	} else {

		var newKey string

		//在GO裡面沒有While迴圈概念，只能用for執行
		for {
			//產生新的Key
			newKey = CreateBase62Key(6)
			//檢查是否使用過
			err := db.QueryRow("SELECT * FROM DemoTable WHERE shortUrl_key = @key", sql.Named("key", newKey)).Scan()
			//如果找不到Row代表沒用過，跳脫For迴圈
			if err == sql.ErrNoRows {
				break
			}
		}

		todayStr := time.Now().Format("2006-01-02")
		expiredateStr := time.Now().AddDate(3, 0, 0).Format("2006-01-02") //增加三年

		//插入新資料
		_, err := db.Exec(
			"INSERT INTO DemoTable (original_url, shortUrl_key, create_date, expire_date, call_time) VALUES (@url, @key, @createDate, @expireDate, @callTime)",
			sql.Named("url", query.OriginalUrl),
			sql.Named("key", newKey),
			sql.Named("createDate", todayStr),
			sql.Named("expireDate", expiredateStr),
			sql.Named("callTime", 0),
		)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		} else {
			c.JSON(http.StatusOK, gin.H{
				"message":     "short url created successfully",
				"shortURL":    c.Request.Host + "/load/" + newKey,
				"expire_date": expiredateStr,
			})
			return
		}
	}
}

//以隨機方式產生Base62的Key
func CreateBase62Key(keyLen int) string {

	//Base62字符
	base62 := "abcdefghijklmnopqrstuvwxyz" + "ABCDEFGHIJKLMNOPQRSTUVWXYZ" + "0123456789"

	//儲存Key的容器
	var keyBuilder strings.Builder
	// keyBuilder.Grow(keyLen)

	//迴圈方式產生keyLen個字元的Key
	for i := 0; i < keyLen; i++ {
		//隨機選擇一個字元並加入
		base62_index := rand.Intn(len(base62))
		keyBuilder.WriteByte(base62[base62_index])
	}

	return keyBuilder.String()
}
