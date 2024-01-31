package main

import (
	"net/http"
	"x19053/ictshort/apifuncs"
	"x19053/ictshort/summarize"

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
		apiGroup.GET("/summurize", test)
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
