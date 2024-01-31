package summarize

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strings"
	"x19053/ictshort/config"
)

var reqHeader http.Header
var summarizePrompts []*Message

const gpt_max_tokens = 4097

type Request struct {
	Model     string     `json:"model"`
	Messages  []*Message `json:"messages"`
	MaxTokens int        `json:"max_tokens"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type Response struct {
	ID      string    `json:"id"`
	Object  string    `json:"object"`
	Created int       `json:"created"`
	Model   string    `json:"model"`
	Usage   *Usage    `json:"usage"`
	Choices []*Choice `json:"choices"`
	Error   struct {
		Code string `json:"code"`
	} `json:"error"`
}

type Usage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

type Choice struct {
	Message      *Message `json:"message"`
	FinishReason string   `json:"finish_reason"`
	Index        int      `json:"index"`
}

func init() {
	reqHeader = make(http.Header)
	reqHeader.Set("Content-Type", "application/json")
	reqHeader.Set("Authorization", "Bearer "+config.MainConfig.AppKeyGPT)
	summarizePrompts = []*Message{
		{Role: "system", Content: "Summarize the following article in Japanese in about 5 lines:"},
		{Role: "system", Content: "If it contains English words, convert them to Japanese pronunciation."},
	}
}

func SummarizeText(text string) (string, error) {
	var response *Response
	for cutPhase := 0; cutPhase <= 2; cutPhase++ {
		log.Println(len(text))
		var err error
		response, err = _SummarizeText(text)
		if err != nil {
			return "", err
		}
		if response.Error.Code == "context_length_exceeded" {
			switch cutPhase {
			case 0:
				tmp := ""
				for i, str := range strings.Split(text, "```") {
					if i%2 == 0 {
						tmp += str
					} else {
						tmp += "(※省略済コード文)"
					}
				}
				text = tmp
			case 1:
				if len(text) > gpt_max_tokens {
					text = text[:gpt_max_tokens]
				}
			default:
				text = "要約に失敗しました"
				break
			}

		} else {
			break
		}
	}

	result := ""
	for _, choice := range response.Choices {
		result += choice.Message.Content
	}
	return result, nil
}

func _SummarizeText(text string) (*Response, error) {
	// log.Println(strings.Trim(text, "\n"))
	requestData, err := json.Marshal(Request{
		Model:     "gpt-3.5-turbo",
		Messages:  append(summarizePrompts, &Message{Role: "user", Content: text}),
		MaxTokens: 300,
	})
	if err != nil {
		return nil, err
	}
	request, _ := http.NewRequest(http.MethodPost, "https://api.openai.com/v1/chat/completions", bytes.NewReader(requestData))
	request.Header = reqHeader
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	bodyData, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	responseStruct := new(Response)
	err = json.Unmarshal(bodyData, responseStruct)
	if err != nil {
		return nil, err
	}

	// log.Println(string(bodyData))

	return responseStruct, nil
}
