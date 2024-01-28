package apifuncs

import (
	"net/http"
	"strconv"
	"x19053/ictshort/articles"
	"x19053/ictshort/config"
	"x19053/ictshort/voice"

	"github.com/gin-gonic/gin"
)

var clientQiita articles.ApiClientQiita

func init() {
	clientQiita.Key = config.MainConfig.AppKeyQiita
}

func GetTrendArticleApi(c *gin.Context) {
	site := c.Query("site")
	articles := []articles.Article{}
	switch site {
	case "Qiita":
		articles = clientQiita.GetListTrendArticles()
	default:
		c.Status(http.StatusForbidden)
		c.Abort()
		return
	}
	c.JSON(http.StatusOK, articles)
}

func GetArticlesApi(c *gin.Context) {
	site := c.Query("site")
	keyword := c.Query("keyword")
	articleIndex, err := strconv.Atoi(c.Query("index"))
	if err != nil {
		c.Status(http.StatusForbidden)
		c.Abort()
		return
	}
	articles := []articles.Article{}
	switch site {
	case "Qiita":
		articles = clientQiita.GetListArticles(keyword, articleIndex)
	default:
		c.Status(http.StatusForbidden)
		c.Abort()
		return
	}
	c.JSON(http.StatusOK, articles)
}

func GetVoiceApi(c *gin.Context) {
	site := c.Query("site")
	keyword := c.Query("keyword")
	id := c.Query("id")

	resposeData := []byte{}

	switch site {
	case "Qiita":
		clientQiita.GetArticleBody(id)
		resposeData = voice.GetVoiceData(keyword)
	default:
		c.Status(http.StatusForbidden)
		c.Abort()
		return
	}
	c.Data(http.StatusOK, http.DetectContentType(resposeData), resposeData)
}
