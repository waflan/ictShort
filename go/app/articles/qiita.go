package articles

import (
	"encoding/json"
	"encoding/xml"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

type ApiClientQiita struct {
	Key string
}

type XMLQiitaFeed struct {
	Articles []struct {
		Title string `xml:"title"`
		Link  struct {
			XMLName string `xml:"link"`
			Href    string `xml:"href,attr"`
		}
		Date   string `xml:"updated"`
		Author string `xml:"author>name"`
	} `xml:"entry"`
}
type JSONArticleQiita struct {
	Title string `json:"title"`
	Id    string `json:"id"`
	Url   string `json:"url"`
	Date  string `json:"updated_at"`
	User  struct {
		Id string `json:"id"`
	} `json:"user"`
}
type JSONContextQiita struct {
	Title string `json:"title"`
	Body  string `json:"body"`
}

func (feed *XMLQiitaFeed) GetArticles() []Article {
	result := []Article{}
	for _, articleData := range feed.Articles {
		urlData, _ := url.Parse(articleData.Link.Href)
		id := strings.Split(urlData.Path, "/")[3]
		result = append(result, Article{
			Title:  articleData.Title,
			Id:     id,
			Url:    articleData.Link.Href,
			Date:   articleData.Date,
			Author: articleData.Author,
		})
	}
	return result
}

// トレンド取得のAPIがない代わりにトレンド記事のxmlを返す機能があったのでそれを使う
func (a ApiClientQiita) GetListTrendArticles() []Article {
	request, _ := http.NewRequest(http.MethodGet, "https://qiita.com/popular-items/feed", nil)
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return nil
	}
	defer response.Body.Close()
	bodyData, _ := io.ReadAll(response.Body)
	xmlData := new(XMLQiitaFeed)
	xml.Unmarshal(bodyData, xmlData)
	return xmlData.GetArticles()
}
func (a ApiClientQiita) GetListArticles(keyword string, pageNum int) []Article {
	request, _ := http.NewRequest(http.MethodGet, "https://qiita.com/api/v2/items?per_page=100&query="+url.QueryEscape(keyword)+"&page="+strconv.Itoa(pageNum), nil)
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return nil
	}
	defer response.Body.Close()
	bodyData, _ := io.ReadAll(response.Body)
	bodyArray := new([]JSONArticleQiita)
	json.Unmarshal(bodyData, bodyArray)

	var result []Article
	for _, q := range *bodyArray {
		result = append(result, *q.ConvertToArticle())
	}
	return result
}

func (a ApiClientQiita) GetArticleContext(id string) (*ArticleContext, error) {
	request, _ := http.NewRequest(http.MethodGet, "https://qiita.com/api/v2/items/"+url.QueryEscape(id), nil)
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return nil, nil
	}
	defer response.Body.Close()
	bodyData, _ := io.ReadAll(response.Body)
	contextStruct := new(JSONContextQiita)
	json.Unmarshal(bodyData, contextStruct)
	return &ArticleContext{Title: contextStruct.Title, Body: contextStruct.Body}, nil
}

func (j *JSONArticleQiita) ConvertToArticle() *Article {
	return &Article{
		Title:  j.Title,
		Id:     j.Id,
		Url:    j.Url,
		Date:   j.Date,
		Author: j.User.Id,
	}
}
