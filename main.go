package main

import (
	// "errors"
	"context"
	"database/sql"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	_ "github.com/denisenkom/go-mssqldb"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/joho/godotenv/autoload"
)

//資料庫結構
type Data struct {
	originalUrl  string
	shortUrl_key string
	create_date  string
	expire_date  string
	call_time    int
}

//讀取環境變數
var (
	AZURE_SQL_SERVER = os.Getenv("AZURE_SQL_SERVER")
	AZURE_SQL_PORT   = os.Getenv("AZURE_SQL_PORT")
	AZURE_USERNAME   = os.Getenv("AZURE_USERNAME")
	AZURE_PASSWORD   = os.Getenv("AZURE_PASSWORD")
	AZURE_SQL_DB     = os.Getenv("AZURE_SQL_DB")
)

func main() {
	server := SetupServer()
	server.Run()
}

func SetupServer() *gin.Engine {
	server := gin.Default()
	server.GET("/", HelloWorld)
	server.POST("/create", CreateShortURL)
	server.GET("/load/:key", LoadShortURL)
	return server
}

func HelloWorld(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Hello World!"})
	return
}

//建立 ShortUrl 資料的 Router
func CreateShortURL(c *gin.Context) {

	// 與 Azure Database Server 連線
	connString := fmt.Sprintf("server=%s;user id=%s;password=%s;port=%s;database=%s;", AZURE_SQL_SERVER, AZURE_USERNAME, AZURE_PASSWORD, AZURE_SQL_PORT, AZURE_SQL_DB)

	var db *sql.DB
	var err error
	// Create connection pool
	db, err = sql.Open("sqlserver", connString)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		log.Fatal("Error creating connection pool: ", err.Error())
		return
	}
	ctx := context.Background()
	err = db.PingContext(ctx)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		log.Fatal(err.Error())
		return
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
		c.JSON(http.StatusBadRequest, gin.H{"error": "originalUrl can't be empty"})
		return
	}

	//確認是否為可用的URL
	_, err = url.ParseRequestURI(query.OriginalUrl)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "This originalUrl is invalid."})
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

//呼叫短連結
func LoadShortURL(c *gin.Context) {

	// 與 Azure Database Server 連線
	connString := fmt.Sprintf("server=%s;user id=%s;password=%s;port=%s;database=%s;", AZURE_SQL_SERVER, AZURE_USERNAME, AZURE_PASSWORD, AZURE_SQL_PORT, AZURE_SQL_DB)
	var db *sql.DB
	var err error

	db, err = sql.Open("sqlserver", connString)
	if err != nil {
		fmt.Println(err.Error())
		log.Fatal("Error creating connection pool: ", err.Error())
	}
	ctx := context.Background()
	err = db.PingContext(ctx)
	if err != nil {
		fmt.Println(err.Error())
		log.Fatal(err.Error())
	}
	fmt.Printf("Connected!")
	defer db.Close()

	//取得URL中的Key
	key := c.Param("key")

	var exist_data Data
	fmt.Println(exist_data)
	row := db.QueryRow("SELECT * FROM DemoTable WHERE shortUrl_key	=@key", sql.Named("key", key))
	err = row.Scan(&exist_data.originalUrl, &exist_data.shortUrl_key, &exist_data.create_date, &exist_data.expire_date, &exist_data.call_time)
	//如果找不到Row
	if err == sql.ErrNoRows {
		c.JSON(http.StatusBadRequest, gin.H{"error": "undefined shortURL"})
		return

		//找到對應資料
	} else {
		//該ShortURL 呼叫次數加一
		_, err := db.Exec("UPDATE DemoTable SET call_time = @callTime WHERE shortUrl_key = @key",
			sql.Named("callTime", exist_data.call_time+1),
			sql.Named("key", key))

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		//重新導向至原連結
		c.Redirect(http.StatusMovedPermanently, exist_data.originalUrl)
		return
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
