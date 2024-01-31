package voice

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
)

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
	bodyData, err := io.ReadAll(response.Body)
	response.Body.Close()

	return bodyData, nil
}
