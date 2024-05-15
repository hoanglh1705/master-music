package elasticsearch

import (
	"fmt"
	"music-master/config"

	"github.com/olivere/elastic/v7"
)

func NewESClient(cfg *config.Configuration) (*elastic.Client, error) {
	client, err := elastic.NewClient(elastic.SetURL("http://localhost:9200"),
		elastic.SetSniff(false),
		elastic.SetHealthcheck(false))

	fmt.Println("ES initialized...")

	return client, err
}
