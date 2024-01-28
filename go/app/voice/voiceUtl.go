package voice

import (
	"io"
	"log"
	"net/http"
	"net/url"
)

func GetVoiceData(text string) []byte {
	request, _ := http.NewRequest(http.MethodPost, "http://ictshort_voicevox_engine:50021/audio_query?speaker=1&text="+url.QueryEscape(text), nil)
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		log.Println(err)
		return []byte{}
	}
	// bodyData, err := io.ReadAll(response.Body)
	// if err != nil {
	// 	log.Println(err)
	// 	return []byte{}
	// }
	request, _ = http.NewRequest(http.MethodPost, "http://ictshort_voicevox_engine:50021/synthesis?speaker=1", response.Body)
	response, err = http.DefaultClient.Do(request)
	if err != nil {
		log.Println(err)
		return []byte{}
	}
	bodyData, err := io.ReadAll(response.Body)
	return bodyData
}
