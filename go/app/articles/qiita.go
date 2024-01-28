package articles

import (
	"encoding/json"
	"encoding/xml"
	"io"
	"net/http"
	"net/url"
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
type JSONBodyQiita struct {
	Body string `json:"body"`
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
	response, _ := http.DefaultClient.Do(request)
	bodyData, _ := io.ReadAll(response.Body)
	xmlData := new(XMLQiitaFeed)
	xml.Unmarshal(bodyData, xmlData)
	return xmlData.GetArticles()
}
func (a ApiClientQiita) GetListArticles(keyword string, index int) []Article {

	return nil
}

func (a ApiClientQiita) GetArticleBody(id string) string {
	request, _ := http.NewRequest(http.MethodGet, "https://qiita.com/popular-items/feed", nil)
	response, _ := http.DefaultClient.Do(request)
	bodyData, _ := io.ReadAll(response.Body)
	bodyStruct := JSONBodyQiita{}
	json.Unmarshal(bodyData, bodyStruct)
	return bodyStruct.Body
}
