package utils

import (
	"encoding/json"
	"errors"
	"strings"
)

type CampoCollection struct {
	Chave string `json:"chave"`
}

func ExtractField(data []byte) ([]CampoCollection, error) {
	var part struct {
		Variable []struct {
			Key string `json:"key"`
		} `json:"variable"`
	}

	err := json.Unmarshal(data, &part)
	if err != nil {
		return nil, errors.New("Erro inesperado")
	}

	field := make([]CampoCollection, 0, len(part.Variable))
	for _, v := range part.Variable {
		field = append(field, CampoCollection{Chave: v.Key})
	}
	return field, nil
}

func CompletePlaceholders(data []byte, value map[string]string) []byte {
	text := string(data)

	for key, value := range value {
		if value == "" {
			continue
		}
		variable := "{{ " + key + "}}"
		text = strings.ReplaceAll(text, variable, key)
	}
	return []byte(text)
}
