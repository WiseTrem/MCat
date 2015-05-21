package tmdb

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
)

type jsonSearch struct {
	Results []struct {
		Id int
	}
}

type JsonInfo struct {
	Budget               int
	Tagline              string
	Imdb_id              string
	Original_title       string
	Overview             string
	Poster_path          string
	Release_date         string
	Title                string
	Vote_average         float32
	Production_countries []struct {
		Iso_3166_1 string
	}
	Genres []struct {
		Name string
	}
}

type JsonCast struct {
	Cast []struct {
		Character    string
		Name         string
		Profile_path string
	}
	Crew []struct {
		Job          string
		Name         string
		Profile_path string
	}
}

func (j *jsonSearch) jsonSearchDecode(r []byte) int {
	if err := json.Unmarshal(r, &j); err != nil {
		fmt.Println(err)
	}
	return len(j.Results)
}

func (j *JsonInfo) jsonInfoDecode(r []byte) {
	if err := json.Unmarshal(r, &j); err != nil {
		fmt.Println(err)
	}
}

func (j *JsonCast) jsonCastDecode(r []byte) {
	if err := json.Unmarshal(r, &j); err != nil {
		fmt.Println(err)
	}
}

func request(r string) []byte {
	client := &http.Client{}
	req, _ := http.NewRequest("GET", r, nil)
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	bodyByte, _ := ioutil.ReadAll(resp.Body)
	return bodyByte
}

func exist(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	fmt.Print(err)
	return false
}

func download(ji *JsonInfo) {
	path := ji.Poster_path
	url := fmt.Sprint("http://image.tmdb.org/t/p/w500", path)
	home := os.Getenv("HOME")
	os.MkdirAll((home + "/.MCat/posters"), 0777)
	if !exist((home + "/.MCat/posters") + path) {
		out, err := os.Create((home + "/.MCat/posters") + path)
		if err != nil {
			fmt.Println(err)
		}
		defer out.Close()
		resp, err := http.Get(url)
		defer resp.Body.Close()
		_, errCopy := io.Copy(out, resp.Body)
		if errCopy != nil {
			fmt.Println(errCopy)
		}
	}
}

func GetInfo(s string) (*JsonInfo, *JsonCast, error) {
	urlSearch, err := url.Parse("http://api.themoviedb.org/3/search/movie?query=template&api_key=a5c697bcbfb66710e125f672937c78c0")
	if err != nil {
		return nil, nil, err
	}
	q := urlSearch.Query()
	q.Set("query", s)
	urlSearch.RawQuery = q.Encode()

	respSearch := request(urlSearch.String())

	j := &jsonSearch{}
	if n := j.jsonSearchDecode(respSearch); n == 0 {
		err := fmt.Errorf(fmt.Sprint(s, " - no movie found"))
		return nil, nil, err
	}

	id := j.Results[0].Id

	urlGet, err := url.Parse("http://api.themoviedb.org/template?api_key=a5c697bcbfb66710e125f672937c78c0&language=ru")
	if err != nil {
		return nil, nil, err
	}
	urlGet.Path = fmt.Sprint("3/movie/", id)
	respGet := request(urlGet.String())

	ji := &JsonInfo{}
	ji.jsonInfoDecode(respGet)

	go download(ji)

	urlCast, err := url.Parse("http://api.themoviedb.org/template?api_key=a5c697bcbfb66710e125f672937c78c0&language=ru")
	if err != nil {
		return nil, nil, err
	}
	urlCast.Path = fmt.Sprint("3/movie/", id, "/credits")
	respCast := request(urlCast.String())

	jc := &JsonCast{}
	jc.jsonCastDecode(respCast)

	return ji, jc, nil
}
