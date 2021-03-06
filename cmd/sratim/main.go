package main

import (
	"fmt"
	"sratim/sratim"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	log "github.com/sirupsen/logrus"
)

func main() {
	client, err := sratim.New(sratim.URL, sratim.API_URL)
	if err != nil {
		log.Fatal(err)
	}

	if client == nil{
		log.Fatalf("No connection to %s", sratim.URL)
	}

	movieName := ""
	prompt := &survey.Input{Message: "Search movie: "}

	err = survey.AskOne(prompt, &movieName)
	if err != nil {
		log.Fatal(err)
	}

	results, err := client.Search(movieName)
	if err != nil {
		log.Fatal(err)
	}

	values := map[string]string{}
	var movieNames []string
	for _, result := range results {
		movieNames = append(movieNames, result.Name)
		values[result.Name] = result.Id
	}

	selectedMovie := ""
	prompt2 := &survey.Select{
		Message: "Pick movie:",
		Options: movieNames,
	}

	err = survey.AskOne(prompt2, &selectedMovie)
	if err != nil {
		log.Fatal(err)
	}

	engMovieName := strings.TrimSpace(strings.Split(selectedMovie, "/")[1])

	err = client.DownloadMovie(values[selectedMovie], engMovieName)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("\nDONE")
}
