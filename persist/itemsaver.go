package persist

import (
	"context"
	"github.com/olivere/elastic/v7"
	"go_crawler/engine"
	"log"
)

func ItemSaver(index string) (chan engine.Item, error) {
	client, err := elastic.NewClient(
		elastic.SetSniff(false),
	)

	if err != nil {
		return nil, err
	}
	out := make(chan engine.Item)
	go func() {
		itemCount := 0
		for {
			item := <-out
			log.Printf("Item saver: got item"+"#%d : %v", itemCount, item)
			itemCount++
			err := Save(client, item, index)
			if err != nil {
				log.Printf("Item saver: error"+"saving item %v: %v", item, err)
			}
		}
	}()
	return out, nil
}

func Save(client *elastic.Client, item engine.Item, index string) error {

	indexService := client.Index().Index(index).BodyJson(item)

	if item.Id != "" {
		indexService.Id(item.Id)
	}

	_, err := indexService.Do(context.Background())

	if err != nil {
		return err
	}
	return nil
}
