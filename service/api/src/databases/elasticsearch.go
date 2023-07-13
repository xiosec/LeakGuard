package databases

import (
	"LeakGuard/config"
	"bytes"
	"context"
	"encoding/json"
	"fmt"

	"github.com/elastic/go-elasticsearch/v8"
)

func Init(config config.Elastic) error {
	cfg := elasticsearch.Config{
		Addresses: config.Addresses,
		Username:  config.Username,
		Password:  config.Password,
	}
	var err error

	ElasticClient, err = elasticsearch.NewClient(cfg)

	return err
}

func IndexExists(indexName string) (bool, error) {
	response, err := ElasticClient.Indices.Exists([]string{indexName})
	if err != nil {
		return false, err
	}
	defer response.Body.Close()

	if response.IsError() {
		return false, fmt.Errorf("Elasticsearch error: %s", response.String())
	}

	return response.StatusCode == 200, nil
}

func ExistPassword(password string) (bool, error) {

	var buf bytes.Buffer
	query := map[string]interface{}{
		"query": map[string]interface{}{
			"match": map[string]interface{}{
				"message": password,
			},
		},
	}
	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		return false, err
	}

	res, err := ElasticClient.Search(
		ElasticClient.Search.WithContext(context.Background()),
		ElasticClient.Search.WithIndex(config.Conf.Elastic.Index),
		ElasticClient.Search.WithBody(&buf),
		ElasticClient.Search.WithTrackTotalHits(true),
		ElasticClient.Search.WithPretty(),
	)

	if err != nil {
		return false, err
	}
	defer res.Body.Close()

	var result map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return false, err
	}

	hits, _ := result["hits"].(map[string]interface{})
	total, _ := hits["total"].(map[string]interface{})
	value, _ := total["value"].(float64)

	if value > 0 {
		return true, nil
	}

	return false, nil
}
