package persist

import (
	"github.com/olivere/elastic/v7"
	"go_crawler/engine"
	"go_crawler/persist"
	"log"
)

type ItemSaverService struct {
	Client *elastic.Client
	Index  string
}

func (s *ItemSaverService) Save(item engine.Item,
	result *string) error {
	err := persist.Save(s.Client, item, s.Index)
	log.Printf("Item %v saved.", item)
	if err == nil {
		*result = "ok"
	} else {
		log.Printf("Error saving %v: %v ", item, err)
	}
	return err
}
