package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

func main() {
	api_base_url := "http://localhost:8080"
	testApi(api_base_url)
	plexMovieDataGrabber(api_base_url)
}

func testApi(api_base_url string) {
	api_url := api_base_url + "/api/test"

	client := &http.Client{}
	r, _ := http.NewRequest("POST", api_url, nil)
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	response, err := client.Do(r)
	if err != nil {
		fmt.Println(err)
	}

	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(string(body))
}

func getAllMovies(api_base_url string) {
	api_url := api_base_url + "/api/test"

	client := &http.Client{}
	r, _ := http.NewRequest("POST", api_url, nil)
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	response, err := client.Do(r)
	if err != nil {
		fmt.Println(err)
	}

	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(string(body))
}

func addMovie(api_base_url, user_key, movie_width, movie_height, video_codac, audio_codac, imdb_id, container, frame_rate, aspect_ratio string) {
	api_url := api_base_url + "/api/addmovie"

	data := url.Values{}
	data.Add("user_key", user_key)
	data.Add("movie_width", movie_width)
	data.Add("movie_height", movie_height)
	data.Add("video_codac", video_codac)
	data.Add("audio_codac", audio_codac)
	data.Add("imdb_id", imdb_id)
	data.Add("container", container)
	data.Add("frame_rate", frame_rate)
	data.Add("aspect_ratio", aspect_ratio)

	client := &http.Client{}
	r, _ := http.NewRequest("POST", api_url, bytes.NewBufferString(data.Encode()))
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	response, err := client.Do(r)
	if err != nil {
		fmt.Println(err)
	}

	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(string(body))
}
