package router

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

//呼叫短連結
func LoadShortURL(c *gin.Context) {

	// 與 Azure Database Server 連線
	connString := fmt.Sprintf("server=%s;user id=%s;password=%s;port=%s;database=%s;", SERVER, USERNAME, PASSWORD, PORT, DATABASE)
	var db *sql.DB
	var err error

	db, err = sql.Open("sqlserver", connString)
	if err != nil {
		log.Fatal("Error creating connection pool: ", err.Error())
	}
	ctx := context.Background()
	err = db.PingContext(ctx)
	if err != nil {
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
