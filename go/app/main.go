package main

import (
	"encoding/json"
	"net/http"
	"x19053/ictshort/apifuncs"
	"x19053/ictshort/articles"

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

	router.GET("/", rootHandler)

	apiGroup := router.Group("/api")
	{
		apiGroup.GET("/listtrend", apifuncs.GetTrendArticleApi)
		apiGroup.GET("/list", apifuncs.GetArticlesApi)
	}

	http.ListenAndServe(":80", router)

}

func rootHandler(c *gin.Context) {
	hoge := articles.ApiClientQiita{}
	articles := hoge.GetListTrendArticles()
	jsonData, _ := json.Marshal(articles)
	c.String(200, string(jsonData))
}
