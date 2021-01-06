package main

import (
	"encoding/json"
	"io/ioutil"
)

//TODO: Write a test for the functionality of TextFacility

type TextFacility struct {
	lookup map[string]interface{}
}

func NewTextFacility() TextFacility {
	byteConfig, _ := ioutil.ReadFile("texts.json")
	var lookup map[string]interface{}
	json.Unmarshal(byteConfig, &lookup)

	return TextFacility{lookup: lookup}
}

func (t *TextFacility) getMessageText(key string) string {
	if key == "" {
		return ""
	}

	return t.lookup[key].(string)
}

func (t *TextFacility) getKeyboardTemplate(keyboardTemplateKey string) [][]string {
	var keyboardLookup map[string]interface{}
	keyboardLookup = t.lookup["inlineKeyboards"].(map[string]interface{})
	data := keyboardLookup[keyboardTemplateKey].([]interface{})

	keyboardField := make([][]string, 0)
	keyboardRow := make([]string, 0)
	rowLen := 2
	for c, keyInterface := range data {
		keyString := keyInterface.(string)
		keyboardRow = append(keyboardRow, keyString)
		if c%rowLen == rowLen-1 {
			keyboardField = append(keyboardField, keyboardRow)
			keyboardRow = make([]string, 0)
		}
	}
	if len(keyboardRow) != 0 {
		keyboardField = append(keyboardField, keyboardRow)
	}
	return keyboardField
}
