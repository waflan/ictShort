package apifuncs

import (
	"log"
	"net/http"
	"strconv"
	"x19053/ictshort/articles"
	"x19053/ictshort/config"
	"x19053/ictshort/summarize"
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
	pageNum, err := strconv.Atoi(c.Query("page"))
	if err != nil {
		c.Status(http.StatusForbidden)
		c.Abort()
		return
	}
	articles := []articles.Article{}
	switch site {
	case "Qiita":
		articles = clientQiita.GetListArticles(keyword, pageNum)
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

	var resposeData []byte
	var err error
	contentType := ""

	if keyword != "" {
		resposeData, err = voice.GetVoiceData(keyword)
		contentType = http.DetectContentType(resposeData)
		c.Data(http.StatusOK, contentType, resposeData)
		return
	}

	switch site {
	case "Qiita":
		var articleBody string
		var summarizedText string
		articleBody, err = clientQiita.GetArticleBody(id)
		if err != nil {
			log.Println(err)
			break
		}
		summarizedText, err = summarize.SummarizeText(articleBody)
		if err != nil {
			log.Println(err)
			break
		}
		log.Println(summarizedText)
		resposeData, err = voice.GetVoiceData(summarizedText)
		contentType = http.DetectContentType(resposeData)
	default:
		c.Status(http.StatusForbidden)
		c.Abort()
		return
	}
	if err != nil || resposeData == nil {
		c.Status(http.StatusBadRequest)
		c.Abort()
		return
	}
	log.Println(contentType)

	c.Data(http.StatusOK, contentType, resposeData)
}
