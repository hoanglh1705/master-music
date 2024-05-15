package elasticsearch

import (
	"context"
	"encoding/json"
	"fmt"
	"music-master/internal/model"

	elastic "github.com/olivere/elastic/v7"
)

type MusicTrackES struct {
	db *elastic.Client
}

func NewMusicTrackCollection(db *elastic.Client) *MusicTrackES {
	return &MusicTrackES{
		db: db,
	}
}

func (es MusicTrackES) Search(ctx context.Context) {
	var students []*model.MusicTrack

	searchSource := elastic.NewSearchSource()
	searchSource.Query(elastic.NewMatchQuery("album", "qua"))

	/* this block will basically print out the es query */
	queryStr, err1 := searchSource.Source()
	queryJs, err2 := json.Marshal(queryStr)

	if err1 != nil || err2 != nil {
		fmt.Println("[esclient][GetResponse]err during query marshal=", err1, err2)
	}
	fmt.Println("[esclient]Final ESQuery=\n", string(queryJs))
	/* until this block */

	searchService := es.db.Search().SearchSource(searchSource)

	searchResult, err := searchService.Do(ctx)
	if err != nil {
		fmt.Println("[ProductsES][GetPIds]Error=", err)
		return
	}

	fmt.Println(searchResult.Hits.Hits)
	for _, hit := range searchResult.Hits.Hits {
		var musicTrack *model.MusicTrack
		err := json.Unmarshal(hit.Source, &musicTrack)
		if err != nil {
			fmt.Println("[Getting Students][Unmarshal] Err=", err)
		}

		students = append(students, musicTrack)
	}

	for _, s := range students {
		fmt.Printf("Student found Title: %s, Album: %v, Artist: %v \n", s.Title, s.Album, s.Artist)
	}

}
