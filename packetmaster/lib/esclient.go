package lib

import (
	"encoding/json"
	"errors"
	"gopkg.in/olivere/elastic.v2"
	"log"
	"time"
)

type ESClient struct {
	*elastic.Client
}

func Getclient(server string) (*ESClient, error) {
	client, err := elastic.NewClient(
		elastic.SetURL(server))
	if err != nil {
		return nil, err

	}
	esclient := &ESClient{Client: client}
	return esclient, nil
}

//the input info should be the json byte[]
func (client *ESClient) Push(info []byte, indexstr string, instancetype string) error {
	infostr := string(info)
	id := time.Now().String()
	returnmessage, err := client.Index().Index(indexstr).Type(instancetype).Id(id).BodyJson(infostr).Do()
	if err != nil {
		log.Println(returnmessage)
		return err
	}
	log.Println("upload ok")
	return nil
}

//the input info should be the json byte[]
func (client *ESClient) Aggregationterm_direct(indexstr string, instancetype string, factterm string, queryname string, queryvalue string) (*elastic.SearchResult, error) {
	Agg_str := `{
  "query": {
    "term": {
      "` + queryname + `": {
        "value":"` + queryvalue + `"
      }
    }
  },
  "facets": {
    "tags": {
      "terms": {
        "field":"` + factterm + `"
      }
    }
  }
}`
	//PerformRequest(method, path string, params url.Values, body interface{})
	path := "/" + indexstr + "/" + instancetype + "/" + "_search"
	response, err := client.PerformRequest("POST", path, nil, Agg_str)
	if err != nil {
		return nil, errors.New("fial to aggregation:" + err.Error())
	}
	ret := new(elastic.SearchResult)
	if err := json.Unmarshal(response.Body, ret); err != nil {
		return nil, err
	}

	log.Println("aggregation ok")
	return ret, nil
}

//already have the structure better to extrac the information
func (client *ESClient) Aggregationterm_indirect(indexstr string, instancetype string, factterm string, querypara string) (*elastic.SearchResult, error) {
	//problem here how to build complex combination?
	Facetquery := elastic.NewQueryFacet().Query(elastic.NewTermQuery("Srcport", "48486")).FacetFilter(elastic.NewTermsFacet().Field(factterm))
	log.Println(Facetquery.Source())
	Facetdestip := elastic.NewTermsFacet().Field(factterm)
	//Facetdestport := elastic.NewTermsFacet().Field(termb)
	log.Println(querypara)
	//searchResult, err := client.Search().Index(indexstr).Type(instancetype).Facet("tags", Facetdestip).Facet("tags", Facetdestport).Do()
	searchResult, err := client.Search().Index(indexstr).Type(instancetype).Facet("tags", Facetdestip).Do()
	if err != nil {
		return nil, errors.New("fial to aggregation:" + err.Error())
	}
	log.Println("aggregation ok")
	return searchResult, nil
}
