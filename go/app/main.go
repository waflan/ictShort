package main

import (
	"net/http"
	"x19053/ictshort/apifuncs"

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

	router.GET("/", rootHandler)

	apiGroup := router.Group("/api")
	{
		apiGroup.GET("/listtrend", apifuncs.GetTrendArticleApi)
		apiGroup.GET("/list", apifuncs.GetArticlesApi)
		apiGroup.GET("/voice", apifuncs.GetVoiceApi)
	}

	http.ListenAndServe(":80", router)

}

func rootHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", gin.H{})
}
