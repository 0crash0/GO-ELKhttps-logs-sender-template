package main

import (
	"bufio"
	"fmt"
	"github.com/Nerzal/gocloak/v13"
	"github.com/elastic/elastic-transport-go/v8/elastictransport"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/gofiber/fiber/v2"
	"go.elastic.co/ecszap"
	"go.uber.org/zap"
	"io"
	"log"
	"os"
	"strings"
)

const elastic_index_name = "my_index"

func elastic_init() *elasticsearch.Client {
	cfg := elasticsearch.Config{
		Addresses: []string{
			"https://192.168.137.128:9200",
		},
		APIKey:                 "TUNzX0FwQUIxN1RTUjRKeDRmN3U6ZS1WZkVXYVRTUGVFNjhPamt3ZWRYdw==",
		CertificateFingerprint: "3B214124C99672377B142724261CD0ECE2AC9CA253439029AC795463B09FE24D",
		//Username:               "elastic",
		//Password:               "wXk51bclsiDOPsjG8h_x",
	}

	es, err := elasticsearch.NewClient(cfg)
	if err != nil {
		log.Fatal(err)
	}

	return es
}

type keycloak struct {
	gocloak      gocloak.GoCloak // keycloak client
	clientId     string          // clientId specified in Keycloak
	clientSecret string          // client secret specified in Keycloak
	realm        string          // realm specified in Keycloak
}

func main() {

	es := elastic_init()
	res, err := es.Info()
	if err != nil {
		log.Fatalf("Error getting response: %s", err)
	}

	defer res.Body.Close()
	//log.Println(res)
	log.Print(es.Transport.(*elastictransport.Client).URLs())

	// redirect stdout to elasticsearch log via pipe to https
	r, w, _ := os.Pipe()
	os.Stdout = w

	go func(reader io.Reader) {

		scanner := bufio.NewScanner(reader)
		for scanner.Scan() {
			res, err = es.Index(
				elastic_index_name,                // Index name
				strings.NewReader(scanner.Text()), // Document body
				//es.Index.WithDocumentID("7"), // Document ID
				es.Index.WithRefresh("true"), // Refresh
			)
			if err != nil {
				log.Fatalf("ERROR: %s", err)
			}
		}
		if err := scanner.Err(); err != nil {
			fmt.Fprintln(os.Stderr, "There was an error with the scanner", err)
		}
	}(r)

	encoderConfig := ecszap.NewDefaultEncoderConfig()
	core := ecszap.NewCore(encoderConfig, os.Stdout, zap.DebugLevel)
	zapLogger := zap.New(core, zap.AddCaller())

	zapLogger.Info("this is info")
	zapLogger.Debug("this is debug")
	zapLogger.Warn("this is warn")

	zapLogger.Info("Server started",
		zap.String("logger", "ZAP"),
		zap.String("host", "localhost"),
		zap.String("port", "3000"),
	)

	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		zapLogger.Info("Get resived")
		return c.SendString("GET request")

	})

	app.Get("/:param", func(c *fiber.Ctx) error {
		return c.SendString("param: " + c.Params("param"))
	})

	app.Post("/", func(c *fiber.Ctx) error {
		return c.SendString("POST request")
	})

	log.Fatal(app.Listen(":3000"))
}
