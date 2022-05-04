package mock

import "log"

// Scrapper implements a mock for the Scrapper from "github.com/kubeshop/testkube/pkg/executor/content"
type Scrapper struct {
	ScrapeFn func(id string, directories []string) error
}

func (s Scrapper) Scrape(id string, directories []string) error {
	if s.ScrapeFn == nil {
		log.Fatal("not implemented")
	}
	return s.ScrapeFn(id, directories)
}
