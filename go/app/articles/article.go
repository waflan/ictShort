package articles

type Article struct {
	Title  string `json:"title"`
	Id     string `json:"id"`
	Url    string `json:"url"`
	Date   string `json:"date"`
	Author string `json:"author"`
}
type ArticleContext struct {
	Title string `json:"title"`
	Body  string `json:"body"`
}

type ApiClient interface {
	GetListTrendArticles() []Article
	GetListArticles() []Article
}
