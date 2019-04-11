package main

import (
	"gopkg.in/olivere/elastic.v6"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
	"suggestion/es"
	"suggestion/utils/log"
	log "github.com/sirupsen/logrus"
	"net/url"
)

var esSuggestor es.EsSuggestor

func main() {
	initSuggestor()
	utils.InitLogger()

	r := gin.Default()
	r.LoadHTMLFiles("templates/corporation.html", "templates/school.html")
	r.Static("/static", "./static")
	r.GET("/corporations", corporation)
	r.GET("/search/corporations", searchCorporations)
	r.GET("/schools", school)
	r.GET("/search/schools", searchSchools)
	if err := r.Run(":62500"); err != nil {
		log.Panicf("init web server error")
	}
}

func initSuggestor() {
	var err error
	var elasticClient *elastic.Client
	for {
		elasticClient, err = elastic.NewClient(
			elastic.SetURL("http://172.17.31.225:9200"))
		if err != nil {
			log.Println(err)
			time.Sleep(3 * time.Second)
		} else {
			break
		}
	}
	esSuggestor = es.EsSuggestor{
		Client: elasticClient,
		SuggestName: "suggest",
		Field: "name_suggest",
		Pretty: true,
	}
}


func corporation(c *gin.Context) {
	c.HTML(200, "corporation.html", gin.H{"title": "corporation suggest",})
}

func school(c *gin.Context) {
	c.HTML(200, "school.html", gin.H{"title": "school suggest"})
}

func searchCorporations(c *gin.Context) {
	queryRaw := c.Query("query")
	query, _ := url.QueryUnescape(queryRaw)
	print(query)
	if query == "12334" {
		errorResponse(c, http.StatusUnavailableForLegalReasons, query)
	}
	suggestions := esSuggestor.Search(query, "corporations")

	c.Writer.Header().Set("Content-Type", "application/json")
	c.JSON(http.StatusOK, gin.H{"corporations": suggestions})
}

func searchSchools(c *gin.Context) {
	query := c.Query("query")
	print(query)
	suggestions := esSuggestor.Search(query, "schools")

	c.Writer.Header().Set("Content-Type", "application/json")
	c.JSON(http.StatusOK, gin.H{"schools": suggestions})
}

func errorResponse(c *gin.Context, code int, err string) {
	c.JSON(code, gin.H{"error": err})
}

