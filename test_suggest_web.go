package main
//
//import (
//	"gopkg.in/olivere/elastic.v6"
//	"log"
//	"github.com/gin-gonic/gin"
//	"net/http"
//	"time"
//	"encoding/json"
//	"context"
//	"fmt"
//)
//
//type SuggestResult struct {
//	Name  string  `json:"name"`
//	Score float64 `json:"score"`
//}
//
//type SuggestResponse struct {
//	Status  int             `json:"status"`
//	Suggest []SuggestResult `json:"suggest"`
//}
//
//func (response *SuggestResponse) showSuggestResponse() {
//	fmt.Println("show suggest response:")
//	fmt.Println("\tstatus:\t", response.Status)
//	for _, suggest := range response.Suggest {
//		fmt.Println("\t\tname:\t", suggest.Name)
//		fmt.Println("\t\tscore:\t", suggest.Score)
//	}
//}
//
//var (
//	elasticClient *elastic.Client
//)
//
//func main() {
//	var err error
//	for {
//		elasticClient, err = elastic.NewClient(
//			elastic.SetURL("http://172.17.31.225:9200"))
//		if err != nil {
//			log.Println(err)
//			time.Sleep(3 * time.Second)
//		} else {
//			break
//		}
//	}
//
//	r := gin.Default()
//	r.LoadHTMLFiles("templates/corporation.html", "templates/school.html")
//	r.Static("/static", "./static")
//	r.GET("/search", searchCorporations)
//	r.GET("/corporation", corporation)
//	r.GET("school", school)
//	if err := r.Run(":8080"); err != nil {
//		log.Fatal(err)
//	}
//}
//
//
//func corporation(c *gin.Context) {
//	c.HTML(200, "corporation.html", gin.H{"title": "corporation suggest",})
//}
//
//func school(c *gin.Context) {
//	c.HTML(200, "school.html", gin.H{"title": "school suggest"})
//}
//
//func searchCorporations(c *gin.Context) {
//	query := c.Query("query")
//	print(query)
//	if query == "12334" {
//		errorResponse(c, http.StatusUnavailableForLegalReasons, query)
//	}
//	s := elastic.NewCompletionSuggester("corporation_suggest").Prefix(query).Field("name_suggest")
//	src, err := s.Source(true)
//	if err != nil {
//		log.Fatal(err)
//	}
//	data, err := json.Marshal(src)
//	if err != nil {
//		log.Fatal(err)
//	}
//	got := string(data)
//	println(got)
//	searchResult, err := elasticClient.Search().
//		Index("corporations").
//		Query(elastic.NewMatchAllQuery()).
//		Suggester(s).Pretty(true).
//		Do(context.TODO())
//	if err != nil {
//		log.Fatal(err)
//	}
//	if searchResult.Suggest == nil {
//		log.Fatal("suggest is nil")
//	}
//	result, fond := searchResult.Suggest["corporation_suggest"]
//	println(fond)
//	println(len(result))
//	resultData := SuggestResponse{
//		Status:  200,
//		Suggest: make([]SuggestResult, 0),
//	}
//	for _, suggests := range result {
//		for _, option := range suggests.Options {
//			resultData.Suggest = append(resultData.Suggest, SuggestResult{
//				Name:  option.Text,
//				Score: option.ScoreUnderscore})
//		}
//	}
//
//	// resultData.showSuggestResponse()
//	//b, err := json.Marshal(resultData)
//	//if err != nil {
//	//	log.Fatal(err)
//	//	errorResponse(c, http.StatusUnavailableForLegalReasons, query)
//	//}
//	// println(string(b))
//	c.Writer.Header().Set("Content-Type", "application/json")
//	c.JSON(http.StatusOK, resultData.Suggest)
//}
//
//func errorResponse(c *gin.Context, code int, err string) {
//	c.JSON(code, gin.H{"error": err})
//}
//
