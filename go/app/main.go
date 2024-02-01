package main

import (
	"bytes"
	"mime/multipart"
	"net/http"
	"x19053/ictshort/apifuncs"
	"x19053/ictshort/summarize"
	"x19053/ictshort/voice"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	gin.SetMode("release")
	router := gin.Default()
	router.Use(cors.New(cors.Config{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"GET", "POST"},
		AllowHeaders:     []string{"*"},
		AllowCredentials: true,
	}))

	router.LoadHTMLGlob("html/templates/*")

	router.Static("/css", "html/css")
	router.Static("/contents", "html/contents")

	router.GET("/", rootHandler)

	apiGroup := router.Group("/api")
	{
		apiGroup.GET("/listtrend", apifuncs.GetTrendArticleApi)
		apiGroup.GET("/list", apifuncs.GetArticlesApi)
		apiGroup.GET("/voice", apifuncs.GetVoiceApi)
		apiGroup.GET("/dict", apifuncs.GetDictionaryApi)
		apiGroup.POST("/dict", apifuncs.SetDictionaryApi)
		apiGroup.DELETE("/dict", apifuncs.DeleteDictionaryApi)
		apiGroup.POST("/dict/import", apifuncs.ImportDictionaryApi)
		apiGroup.GET("/summurize", test)
		apiGroup.GET("/multiple", test2)
	}

	http.ListenAndServe(":80", router)

}

func rootHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", gin.H{})
}

func test(c *gin.Context) {
	text := c.Query("text")
	context, err := summarize.SummarizeText(text)
	if err != nil {
		c.Abort()
		c.Status(http.StatusBadRequest)
	}
	c.String(http.StatusOK, context)
}
func test2(c *gin.Context) {
	buffer := &bytes.Buffer{}
	multiwriter := multipart.NewWriter(buffer)

	writer1, _ := multiwriter.CreateFormFile("voice1", "voice1.wav")
	voice1, _ := voice.GetVoiceData("ほげ")
	// voice1 := []byte("hoge")
	writer1.Write(voice1)

	writer2, _ := multiwriter.CreateFormFile("voice2", "voice2.wav")
	voice2, _ := voice.GetVoiceData("ふが")
	// voice2 := []byte("hoge")
	writer2.Write(voice2)

	multiwriter.Close()

	c.Data(200, "multipart/form-data; boundary="+multiwriter.Boundary(), buffer.Bytes())
}
