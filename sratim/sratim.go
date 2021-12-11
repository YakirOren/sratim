package sratim

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"path/filepath"
	"sratim/progress"
	"strings"

	"golang.org/x/net/publicsuffix"
)

const (
	URL     = "https://Sratim.tv"
	API_URL = "https://api.Sratim.tv"
)

type Sratim struct {
	client *http.Client
	token  string
	url    string
	apiUrl string
}

func New(url string, apiUrl string) (*Sratim, error) {
	jar, err := cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})
	if err != nil {
		log.Fatal(err)
	}
	client := &http.Client{
		Jar: jar,
	}

	a := &Sratim{client: client, url: url, apiUrl: apiUrl}

	if a.init() != nil {
		return nil, err
	}

	return a, nil
}

func (sr Sratim) Search(term string) ([]SearchResult, error) {
	var result SearchResponse

	form := url.Values{}
	form.Add("term", term)

	response, err := http.PostForm(sr.apiUrl+"/movie/search", form)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(data, &result)
	if err != nil {
		return nil, err
	}

	return result.Results, nil
}

func (sr *Sratim) init() error {
	resp, err := sr.client.Get(sr.url)
	if err != nil {
		return err
	}

	u, err := url.Parse(sr.apiUrl)
	if err != nil {
		return err
	}

	sr.client.Jar.SetCookies(u, resp.Cookies())

	resp, err = sr.client.Get(fmt.Sprintf("%s/movie/preWatch", sr.apiUrl))
	if err != nil {
		return err
	}

	token, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	sr.token = string(token)
	return nil
}

func (sr Sratim) GetMovie(id string) (*Response, error) {
	progress.Loader(30)
	resp, err := sr.client.Get(fmt.Sprintf("%s/movie/watch/id/%s/token/%s", sr.apiUrl, id, sr.token))
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(resp.Body)

	var response Response

	if err := json.Unmarshal(body, &response); err != nil {
		return nil, err
	}

	return &response, nil
}

func (sr Sratim) download(movieURL string, writer io.Writer) error {
	resp, err := sr.client.Get(movieURL)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	counter := &progress.WriteCounter{}

	_, err = io.Copy(writer, io.TeeReader(resp.Body, counter))
	if err != nil {
		return err
	}

	return nil
}

func (sr Sratim) DownloadMovie(id string) error {
	movieURL, err := sr.GetMovieURL(id)
	if err != nil {
		return err
	}

	fileName := filepath.Base(movieURL.Path)
	// the ID could be used here
	// but this assures the file is saved in the correct format.
	file, err := os.Create(fileName)
	if err != nil {
		log.Fatal(err)
	}

	err = sr.download(movieURL.String(), file)
	if err != nil {
		return err
	}

	return nil
}

func (sr Sratim) GetMovieURL(id string) (*url.URL, error) {
	response, err := sr.GetMovie(id)
	if err != nil {
		return nil, err
	}
	if !response.Success {
		return nil, errors.New(strings.Join(response.Errors, ","))
	}
	u, err := url.Parse(response.Watch.URL)
	if err != nil {
		return nil, err
	}
	u.Scheme = "https"

	return u, nil
}
