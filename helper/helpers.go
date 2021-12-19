package helper

import (
	"errors"
	"io"
	"net/http"
	"sratim/progress"
	"strconv"

	"github.com/dustin/go-humanize"
)

func SaveFile(client *http.Client, fileURL string, writer io.Writer) error {
	resp, err := client.Get(fileURL)
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

func GetFileSize(client *http.Client, fileURL string) (string, error) {
	resp, err := client.Head(fileURL)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != http.StatusOK {
		return "", errors.New(resp.Status)
	}

	// the Header "Content-Length" will let us know
	// the total file size to download
	size, err := strconv.ParseUint(resp.Header.Get("Content-Length"), 10, 64)
	if err != nil {
		return "", err
	}

	return humanize.Bytes(size), nil
}
