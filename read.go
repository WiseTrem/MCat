package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"regexp"
	"strings"

	"github.com/WiseTrem/MCat/base"
	"github.com/WiseTrem/MCat/tmdb"
)

func main() {
	str, err := ioutil.ReadDir("/mnt/DS/Video/Films")
	if err != nil {
		log.Fatal(err)
	}

	for i := range str {
		filestr := str[i].Name()
		s := regexp.MustCompile("((20[0-9][0-9])|(x264)|(.mkv)|(1080)|(19[0-9][0-9])|(\\(.*)|(\\[.*)).*")
		//filter mkv and m4v
		f := regexp.MustCompile(".*((.mkv)|(.m4v))")
		strFiltered := f.FindString(filestr)
		var searchString string
		if strFiltered != "" {
			cut := s.FindString(strFiltered)
			//fmt.Printf("%s\n", cut)
			searchString = strings.Join((strings.Split(strings.TrimSuffix(strFiltered, cut), ".")), " ")
			if searchString == "" {
				fmt.Printf("Cant parse %s\n", strFiltered)
				searchString = strFiltered
			}
			//fmt.Printf("%s\n", searchString)

			info, cast, err := tmdb.GetInfo(searchString)
			if err == nil {
				base.Save(searchString, info, cast)
			} else {
				fmt.Println(err)
			}
		}
	}
}
