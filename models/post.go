package models

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"time"
)

type Post struct {
	Title string    `json:"title"`
	Body  string    `json:"body"`
	Url   string    `json:"url"`
	Thumb string    `json:"thumbnail"`
	Date  time.Time `json:"date"`
}

func (p *Post) ToJSON() ([]byte, error) {
	ret, err := json.Marshal(p)
	if err != nil {
		return nil, err
	}
	return ret, nil
}

func (p *Post) DownloadThumb(rootDir string) error {
	thumbUrl, err := url.Parse(p.Thumb)
	if err != nil {
		return err
	}
	// We follow the Hugo path structure for saving images in the static directory.

	// Check if the images directory exists, if not create
	imgDir := filepath.Join(rootDir, "static", "images")
	err = os.MkdirAll(imgDir, 0750)
	if err != nil && !os.IsExist(err) {
		return err
	}
	filePath := filepath.Join(imgDir, filepath.Base(thumbUrl.Path))
	imgPath := filepath.Join("/images", filepath.Base(thumbUrl.Path))
	file, err := os.Create(filePath)
	if err != nil {
		log.Printf("err: %s creating file\n", err.Error())
		return err
	}
	response, err := http.Get(thumbUrl.String())
	if err != nil {
		log.Printf("err: %s making http request\n", err.Error())
		return err
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		log.Println("err: non OK http response")
		return fmt.Errorf("%s non OK http response", response.Status)
	}
	_, err = io.Copy(file, response.Body)
	if err != nil {
		log.Printf("err: %s writing file\n", err.Error())
		return err
	}
	// Since the download went well, we assume that the file was downloaded ok
	// I guess we could also check the file size, but for now this will do
	p.Thumb = imgPath
	return nil
}

func (p *Post) SaveThumb(rootDir string) (string, error) {
	thumbUrl, err := url.Parse(p.Thumb)
	if err != nil {
		return "", err
	}
	response, err := http.Get(thumbUrl.String())
	if err != nil {
		log.Printf("err: %s making http request\n", err.Error())
		return "", err
	}
	defer response.Body.Close()
	contentType := response.Header.Get("Content-Type")
	if response.StatusCode != http.StatusOK {
		log.Println("err: non OK http response")
		return "", fmt.Errorf("%s non OK http response", response.Status)
	}
	var b bytes.Buffer
	encoder := base64.NewEncoder(base64.StdEncoding, &b)

	_, err = io.Copy(encoder, response.Body)
	if err != nil {
		log.Printf("err: %s writing file\n", err.Error())
		return "", err
	}
	final := fmt.Sprintf("data:%s;base64,%s", contentType, b.String())
	p.Thumb = final
	return final, nil

}

func (p Post) String() string {
	b, err := p.ToJSON()
	if err != nil {
		log.Println(err.Error())
		return ""
	}
	return string(b)
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
