package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

type Keys struct {
	XMLName  xml.Name `xml:"MediaContainer"`
	VideoKey []Key    `xml:"Video"`
}

type Key struct {
	Key string `xml:"key,attr"`
}

type Metadata struct {
	XMLName xml.Name `xml:"MediaContainer"`
	Video   Video    `xml:"Video"`
}

type Video struct {
	Imdb_id string `xml:"guid,attr"`
	Media   Media  `xml:"Media"`
}

type Media struct {
	Movie_width  string `xml:"width,attr"`
	Movie_height string `xml:"height,attr"`
	Video_codac  string `xml:"videoCodec,attr"`
	Audio_codac  string `xml:"audioCodec,attr"`
	Container    string `xml:"container,attr"`
	Frame_rate   string `xml:"videoFrameRate,attr"`
	Aspect_ratio string `xml:"aspectRatio,attr"`
}

func plexMovieDataGrabber(api_base string) {
	xml_url_base := "http://192.168.0.8:32400"
	xml_key_url := xml_url_base + "/library/sections/1/all"

	response, err := http.Get(xml_key_url)
	if err != nil {
		fmt.Println(err)
		return
	}

	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println(err)
		return
	}

	keys := Keys{}
	err = xml.Unmarshal(body, &keys)
	if err != nil {
		fmt.Println(err)
		return
	}

	for _, url_extention := range keys.VideoKey {
		xml_metadata_url := xml_url_base + url_extention.Key

		metadata_response, err := http.Get(xml_metadata_url)
		if err != nil {
			fmt.Println(err)
			return
		}

		defer metadata_response.Body.Close()
		metadata_body, err := ioutil.ReadAll(metadata_response.Body)
		if err != nil {
			fmt.Println(err)
			return
		}

		metadata := Metadata{}
		err = xml.Unmarshal(metadata_body, &metadata)
		if err != nil {
			fmt.Println(err)
			return
		}
		metadata.Video.Imdb_id = strings.Split(strings.Split(metadata.Video.Imdb_id, "//")[1], "?")[0]
		fmt.Println(metadata)
		addMovie(api_base, "66IliKuYo5wZNlyXdlsLCWoUBliSWKu8Rms7wbTfmk3yz5Pn3sThrdrx914UivQF",
			metadata.Video.Media.Movie_width,
			metadata.Video.Media.Movie_height,
			metadata.Video.Media.Video_codac,
			metadata.Video.Media.Audio_codac,
			metadata.Video.Imdb_id,
			metadata.Video.Media.Container,
			metadata.Video.Media.Frame_rate,
			metadata.Video.Media.Aspect_ratio)
	}
}
