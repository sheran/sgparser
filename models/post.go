package models

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"encoding/json"
	"log"
)

type Post struct {
	Title string `json:"title"`
	Body  string `json:"body"`
	Url   string `json:"url"`
	Thumb string `json:"thumbnail"`
}

func (p *Post) ToJSON() ([]byte, error) {
	ret, err := json.Marshal(p)
	if err != nil {
		return nil, err
	}
	return ret, nil
}

func (p *Post) ToBase64() (string, error) {
	var b bytes.Buffer
	gz := gzip.NewWriter(&b)

	ret, err := p.ToJSON()
	if err != nil {
		log.Println(err.Error())
		return "", err
	}

	_, err = gz.Write(ret)
	if err != nil {
		log.Println(err.Error())
		return "", err
	}

	if err := gz.Close(); err != nil {
		log.Println(err.Error())
		return "", err
	}
	return base64.StdEncoding.EncodeToString(b.Bytes()), nil

}
