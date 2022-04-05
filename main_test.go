package main

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

//測試伺服器開啟，嘗試呼叫IndexRouter
func Test_ServerIndexRouter(t *testing.T) {
	server := SetupServer()

	req, _ := http.NewRequest("GET", "/", nil) // 建立一個請求
	w := httptest.NewRecorder()                // 建立一個ResponseRecorder其實作http.ResponseWriter，用來記錄response狀態
	server.ServeHTTP(w, req)                   // gin.Engine.ServerHttp實作http.Handler介面，用來處理HTTP請求及回應。

	expectedStatus := http.StatusOK
	expectedContent := "Hello World"

	assert.Equal(t, expectedStatus, w.Code)
	assert.Contains(t, w.Body.String(), expectedContent)
}

//測試新增短網址
func Test_CreateShortURL_Success(t *testing.T) {
	server := SetupServer()

	//Json紀錄原網址內容
	var jsonData = []byte(`{"originalUrl": "https://www.google.com.tw"}`)

	req, _ := http.NewRequest("POST", "/create", bytes.NewBuffer(jsonData)) // 建立一個請求
	w := httptest.NewRecorder()                                             // 建立一個ResponseRecorder其實作http.ResponseWriter，用來記錄response狀態
	server.ServeHTTP(w, req)                                                // gin.Engine.ServerHttp實作http.Handler介面，用來處理HTTP請求及回應。

	expectedStatus := http.StatusOK
	expectedContent := "created successfully"

	assert.Equal(t, expectedStatus, w.Code)
	assert.Contains(t, w.Body.String(), expectedContent)
}

//測試新增短網址，但內容為空白
func Test_CreateShortURL_EmptyUrl(t *testing.T) {
	server := SetupServer()

	//Json紀錄原網址內容
	var jsonData = []byte(`{"originalUrl": ""}`)

	req, _ := http.NewRequest("POST", "/create", bytes.NewBuffer(jsonData)) // 建立一個請求
	w := httptest.NewRecorder()                                             // 建立一個ResponseRecorder其實作http.ResponseWriter，用來記錄response狀態
	server.ServeHTTP(w, req)                                                // gin.Engine.ServerHttp實作http.Handler介面，用來處理HTTP請求及回應。

	expectedStatus := http.StatusBadRequest
	expectedContent := "can't be empty"

	assert.Equal(t, expectedStatus, w.Code)
	assert.Contains(t, w.Body.String(), expectedContent)
}

//測試新增短網址，但內容為空白
func Test_CreateShortURL_InvalidUrl(t *testing.T) {
	server := SetupServer()

	//Json紀錄原網址內容
	var jsonData = []byte(`{"originalUrl": "123"}`)

	req, _ := http.NewRequest("POST", "/create", bytes.NewBuffer(jsonData)) // 建立一個請求
	w := httptest.NewRecorder()                                             // 建立一個ResponseRecorder其實作http.ResponseWriter，用來記錄response狀態
	server.ServeHTTP(w, req)                                                // gin.Engine.ServerHttp實作http.Handler介面，用來處理HTTP請求及回應。

	expectedStatus := http.StatusBadRequest
	expectedContent := "invalid"

	assert.Equal(t, expectedStatus, w.Code)
	assert.Contains(t, w.Body.String(), expectedContent)
}

//測試呼叫短網址，成功
func Test_LoadShortURL_Success(t *testing.T) {
	server := SetupServer()

	shortUrl_key := "Dsc2WD"

	req, _ := http.NewRequest("GET", "/load/"+shortUrl_key, nil) // 建立一個請求
	w := httptest.NewRecorder()                                  // 建立一個ResponseRecorder其實作http.ResponseWriter，用來記錄response狀態
	server.ServeHTTP(w, req)                                     // gin.Engine.ServerHttp實作http.Handler介面，用來處理HTTP請求及回應。

	expectedStatus := http.StatusMovedPermanently
	expectedContent := "https://www.google.com.tw"

	assert.Equal(t, expectedStatus, w.Code)
	assert.Contains(t, w.Body.String(), expectedContent)
}

//測試呼叫短網址，無效的短網址
func Test_LoadShortURL_UndefinedKey(t *testing.T) {
	server := SetupServer()

	shortUrl_key := "123456"

	req, _ := http.NewRequest("GET", "/load/"+shortUrl_key, nil) // 建立一個請求
	w := httptest.NewRecorder()                                  // 建立一個ResponseRecorder其實作http.ResponseWriter，用來記錄response狀態
	server.ServeHTTP(w, req)                                     // gin.Engine.ServerHttp實作http.Handler介面，用來處理HTTP請求及回應。

	expectedStatus := http.StatusBadRequest
	expectedContent := "undefined shortURL"

	assert.Equal(t, expectedStatus, w.Code)
	assert.Contains(t, w.Body.String(), expectedContent)
}
