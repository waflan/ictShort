package apifuncs

import (
	"encoding/json"
	"io"
	"net/http"
	"x19053/ictshort/voice"

	"github.com/gin-gonic/gin"
)

func GetDictionaryApi(c *gin.Context) {
	allFlag := c.Query("all")
	if allFlag == "true" {
		c.Data(http.StatusOK, "application/json", voice.GetDictionaryData())
	} else {
		c.JSON(http.StatusOK, voice.DictionaryData)
	}
}
func SetDictionaryApi(c *gin.Context) {
	id := c.Query("id")
	word := &voice.Word{}
	body, _ := io.ReadAll(c.Request.Body)
	json.Unmarshal(body, word)
	if _, ok := voice.DictionaryData[id]; ok {
		if voice.UpdateDictionary(id, word) {
			voice.DictionaryData[id] = *word
			c.Status(http.StatusOK)
		} else {
			c.Status(http.StatusBadRequest)
		}
	} else {
		id = voice.AddDictionary(word)
		if id == "" {
			c.Status(http.StatusBadRequest)
		} else {
			voice.DictionaryData[id] = *word
			c.String(http.StatusOK, id)
		}
	}
}
func DeleteDictionaryApi(c *gin.Context) {
	id := c.Query("id")
	if _, ok := voice.DictionaryData[id]; ok {
		if voice.DeleteDictionary(id) {
			delete(voice.DictionaryData, id)
			c.Status(http.StatusOK)
		} else {
			c.Status(http.StatusBadRequest)
		}
	} else {
		c.Status(http.StatusForbidden)
	}
}

func ImportDictionaryApi(c *gin.Context) {
	body, _ := io.ReadAll(c.Request.Body)
	// log.Println(string(body))
	if voice.ImportDictionary(body) {
		c.Status(http.StatusOK)
	} else {
		c.Status(http.StatusBadRequest)
	}
}
