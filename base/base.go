package base

import (
	"fmt"
	"github.com/boltdb/bolt"
	"log"

	"MCat/tmdb"
)

func Save(s string, ji *tmdb.JsonInfo, jc *tmdb.JsonCast) {
	db, err := bolt.Open("/tmp/TEST.db", 0666, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	db.Update(func(tx *bolt.Tx) error {
		tx.CreateBucketIfNotExists([]byte(s))
		b := tx.Bucket([]byte(s))
		fmt.Printf("Saving %s\n", ji.Title)
		//Saving JsonInfo struct meta
		b.Put([]byte("Budget"), []byte(fmt.Sprint(ji.Budget)))
		b.Put([]byte("Tagline"), []byte(ji.Tagline))
		b.Put([]byte("Imdb_id"), []byte(ji.Imdb_id))
		b.Put([]byte("Original_title"), []byte(ji.Original_title))
		b.Put([]byte("Overview"), []byte(ji.Overview))
		b.Put([]byte("Poster_path"), []byte(ji.Poster_path))
		b.Put([]byte("Release_date"), []byte(ji.Release_date))
		b.Put([]byte("Title"), []byte(ji.Title))
		b.Put([]byte("Vote_average"), []byte(fmt.Sprint(ji.Vote_average)))
		for i := 0; i < len(ji.Production_countries); i++ {
			b.Put([]byte(fmt.Sprint("Country", i)), []byte(ji.Production_countries[i].Iso_3166_1))
		}
		for i := 0; i < len(ji.Genres); i++ {
			b.Put([]byte(fmt.Sprint("Genres", i)), []byte(ji.Genres[i].Name))
		}
		//Saving JsonCast struct meta
		for i := 0; i < len(jc.Cast); i++ {
			b.Put([]byte(jc.Cast[i].Name), []byte(jc.Cast[i].Character))
			b.Put([]byte(fmt.Sprint(jc.Cast[i].Name, "&Photo")), []byte(jc.Cast[i].Character))
		}
		for i := 0; i < len(jc.Crew); i++ {
			if jc.Crew[i].Job == "Director" || jc.Crew[i].Job == "Screenplay" || jc.Crew[i].Job == "Original Music Composer" {
				b.Put([]byte(fmt.Sprint(jc.Crew[i].Job, i)), []byte(jc.Crew[i].Name))
				b.Put([]byte(fmt.Sprint(jc.Crew[i].Name, "&Photo")), []byte(jc.Crew[i].Profile_path))
			}
		}
		return nil
	})
}
