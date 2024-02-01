package apifuncs

import (
	"bytes"
	"io"
	"log"
	"mime/multipart"
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

	var summarizedText string
	var voiceData []byte
	var titleVoiceData []byte
	responseBody := &bytes.Buffer{}
	multiwriter := multipart.NewWriter(responseBody)

	var err error

	if keyword != "" {
		voiceData, err = voice.GetVoiceData(keyword)
		contentType := http.DetectContentType(voiceData)
		c.Data(http.StatusOK, contentType, voiceData)
		return
	}

	switch site {
	case "Qiita":
		var articleContext *articles.ArticleContext
		articleContext, err = clientQiita.GetArticleContext(id)
		if err != nil {
			log.Println(err)
			break
		}
		summarizedText, err = summarize.SummarizeText(articleContext.Body)
		if err != nil {
			log.Println(err)
			break
		}
		log.Println(summarizedText)
		voiceData, err = voice.GetVoiceData(summarizedText)
		titleVoiceData, err = voice.GetVoiceData(articleContext.Title)

	default:
		c.Status(http.StatusForbidden)
		c.Abort()
		return
	}
	if err != nil || voiceData == nil || titleVoiceData == nil {
		c.Status(http.StatusBadRequest)
		c.Abort()
		return
	}

	var writer io.Writer
	writer, _ = multiwriter.CreateFormFile("summarizedText", "summarizedText.txt")
	writer.Write([]byte(summarizedText))
	writer, _ = multiwriter.CreateFormFile("titleVoiceData", "titleVoiceData.wav")
	writer.Write(titleVoiceData)
	writer, _ = multiwriter.CreateFormFile("voiceData", "voiceData.wav")
	writer.Write(voiceData)

	multiwriter.Close()
	c.Data(http.StatusOK, "multipart/form-data; boundary="+multiwriter.Boundary(), responseBody.Bytes())
}
