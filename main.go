package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"log/slog"
	"os"
	"strings"
	"time"

	"github.com/blevesearch/bleve"
	"github.com/blevesearch/bleve/mapping"
	"github.com/gammazero/workerpool"
	"github.com/vbauerster/mpb/v7"
	"github.com/vbauerster/mpb/v7/decor"
)

type LegalEntity struct {
	FullName    string
	OKPO        string
	OKATO       string
	OKTMO       string
	OKOGU       string
	OKFS        string
	OKOPF       string
	INN         string
	OGRN_OGRNIP string
}

func defineIndexMapping() mapping.IndexMapping {
	indexMapping := bleve.NewIndexMapping()

	legalEntityMapping := bleve.NewDocumentMapping()

	legalEntityMapping.AddFieldMappingsAt("FullName", bleve.NewTextFieldMapping())
	legalEntityMapping.AddFieldMappingsAt("OKPO", bleve.NewTextFieldMapping())
	legalEntityMapping.AddFieldMappingsAt("OKATO", bleve.NewTextFieldMapping())
	legalEntityMapping.AddFieldMappingsAt("OKTMO", bleve.NewTextFieldMapping())
	legalEntityMapping.AddFieldMappingsAt("OKOGU", bleve.NewTextFieldMapping())
	legalEntityMapping.AddFieldMappingsAt("OKFS", bleve.NewTextFieldMapping())
	legalEntityMapping.AddFieldMappingsAt("OKOPF", bleve.NewTextFieldMapping())
	legalEntityMapping.AddFieldMappingsAt("INN", bleve.NewTextFieldMapping())
	legalEntityMapping.AddFieldMappingsAt("OGRN_OGRNIP", bleve.NewTextFieldMapping())

	indexMapping.AddDocumentMapping("legalEntity", legalEntityMapping)

	return indexMapping
}

func main() {
	file, err := os.Open("data-20220314-structure-20220314.csv")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.Comma = ';'

	lines, err := reader.ReadAll()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	lineCount := len(lines)

	indexMapping := defineIndexMapping()
	index, err := bleve.New("index", indexMapping)
	if err != nil {
		slog.Warn("try start bleve index err", err, "try open exists bleve index")
		index, err = bleve.Open("index")
		if err != nil {
			log.Fatalf("open bleve index: %v", err)
			return
		}
	}

	p := mpb.New(mpb.WithWidth(64))

	indexBar := p.AddBar(int64(lineCount),
		mpb.PrependDecorators(
			decor.Name("Indexing: "),
			decor.Percentage(),
		),
		mpb.AppendDecorators(
			decor.CountersNoUnit("%d / %d"),
		),
	)

	pool := workerpool.New(5)

	startTime := time.Now()
	for _, record := range lines {
		entity := LegalEntity{
			FullName:    strings.ToLower(record[0]),
			OKPO:        record[1],
			OKATO:       record[2],
			OKTMO:       record[3],
			OKOGU:       record[4],
			OKFS:        record[5],
			OKOPF:       record[6],
			INN:         record[7],
			OGRN_OGRNIP: record[8],
		}

		pool.Submit(func() {
			index.Index(entity.FullName, entity)
			indexBar.IncrBy(1)
		})
	}

	pool.StopWait()

	duration := time.Since(startTime)
	fmt.Println("Indexing took:", duration.String())

	query := bleve.NewQueryStringQuery(strings.ToLower("автомобильные дороги"))
	searchRequest := bleve.NewSearchRequest(query)
	searchResult, err := index.Search(searchRequest)
	if err != nil {
		log.Fatalf("search: %s", err)
	}

	fmt.Println("Результаты поиска для запроса:", query)
	for _, hit := range searchResult.Hits {
		indexedEntity, ok := hit.Fields["FullName"].(string)
		if !ok {
			fmt.Println("Error: unable to retrieve indexed entity")
			continue
		}
		fmt.Printf("%v", indexedEntity)
	}
}
