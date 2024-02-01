package voice

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
)

type Word struct {
	Surface       string `json:"surface"`
	Pronunciation string `json:"pronunciation"`
	Accent        int    `json:"accent_type"`
	Type          string `json:"word_type"`
}
type WordAdded struct {
	Surface       string `json:"surface"`
	Pronunciation string `json:"pronunciation"`
	Accent        int    `json:"accent_type"`
	Type0         string `json:"part_of_speech"`
	Type1         string `json:"part_of_speech_detail_1"`
}

type Dictionary map[string]Word

var DictionaryData Dictionary

func init() {
	tmp := GetDictionary()
	DictionaryData = *tmp
}

func (w WordAdded) ConvertToWord() *Word {
	wordType := ""
	switch w.Type0 {
	case "名詞":
		switch w.Type1 {
		case "固有名詞":
			wordType = "PROPER_NOUN"
		case "一般":
			wordType = "COMMON_NOUN"
		case "接尾":
			wordType = "SUFFIX"
		}
	case "動詞":
		wordType = "VERB"
	case "形容詞":
		wordType = "ADJECTIVE"
	}
	return &Word{
		Surface:       w.Surface,
		Pronunciation: w.Pronunciation,
		Accent:        w.Accent,
		Type:          wordType,
	}
}

func GetVoiceData(text string) ([]byte, error) {
	if text == "" {
		return nil, nil
	}
	// log.Println(text)
	request, _ := http.NewRequest(http.MethodPost, "http://ictshort_voicevox_engine:50021/audio_query?speaker=1&text="+url.QueryEscape(text), nil)
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	tmp, _ := io.ReadAll(response.Body)
	// log.Println(string(tmp))
	// if err != nil {
	// 	log.Println(err)
	// 	return []byte{}
	// }
	response.Body.Close()
	request, _ = http.NewRequest(http.MethodPost, "http://ictshort_voicevox_engine:50021/synthesis?speaker=1", bytes.NewReader(tmp))
	response, err = http.DefaultClient.Do(request)
	if err != nil {
		log.Println(err)
		log.Println(strings.ReplaceAll(string(tmp), "\n", ""))
		return nil, err
	}
	bodyData, _ := io.ReadAll(response.Body)
	response.Body.Close()

	return bodyData, nil
}

func GetDictionary() *map[string]Word {
	dictData := GetDictionaryData()
	words := new(map[string]WordAdded)
	result := map[string]Word{}
	json.Unmarshal(dictData, words)

	for key, value := range *words {
		result[key] = *value.ConvertToWord()
	}
	return &result
}
func GetDictionaryData() []byte {
	request, _ := http.NewRequest(http.MethodGet, "http://ictshort_voicevox_engine:50021/user_dict", nil)
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		log.Println(err)
		return nil
	}
	bodyData, _ := io.ReadAll(response.Body)
	return bodyData
}
func AddDictionary(word *Word) string {
	request, _ := http.NewRequest(http.MethodPost, fmt.Sprintf("http://ictshort_voicevox_engine:50021/user_dict_word?surface=%s&pronunciation=%s&accent_type=%d&word_type=%s", url.QueryEscape(word.Surface), url.QueryEscape(word.Pronunciation), word.Accent, word.Type), nil)
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		log.Println(err)
		return ""
	}
	if response.StatusCode != http.StatusOK {
		log.Println(response.Status)
		return ""
	}
	bodyData, _ := io.ReadAll(response.Body)

	return string(bodyData)
}
func UpdateDictionary(id string, word *Word) bool {
	request, _ := http.NewRequest(http.MethodPut, fmt.Sprintf("http://ictshort_voicevox_engine:50021/user_dict_word/%s?surface=%s&pronunciation=%s&accent_type=%d&word_type=%s", url.QueryEscape(id), url.QueryEscape(word.Surface), url.QueryEscape(word.Pronunciation), word.Accent, word.Type), nil)
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		log.Println(err)
		return false
	}
	if response.StatusCode != http.StatusNoContent {
		log.Println(response.Status)
		return false
	}
	return true
}
func DeleteDictionary(id string) bool {
	request, _ := http.NewRequest(http.MethodDelete, "http://ictshort_voicevox_engine:50021/user_dict_word/"+url.QueryEscape(id), nil)
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		log.Println(err)
		return false
	}
	if response.StatusCode != http.StatusNoContent {
		log.Println(response.Status)
		return false
	}
	return true
}

func ImportDictionary(data []byte) bool {
	request, _ := http.NewRequest(http.MethodPost, "http://ictshort_voicevox_engine:50021/import_user_dict?override=true", bytes.NewReader(data))
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		log.Println(err)
		return false
	}
	if response.StatusCode != http.StatusNoContent {
		log.Println(response.Status)
		bodyData, _ := io.ReadAll(response.Body)
		log.Println(string(bodyData))
		return false
	}
	return true
}
