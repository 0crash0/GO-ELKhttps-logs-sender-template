package main

import (
	"bufio"
	"fmt"
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

func main() {
	cfg := elasticsearch.Config{
		Addresses: []string{
			"https://192.168.1.34:9200",
		},
		APIKey:                 "ZkRZRDdvOEJCekRicVRYdUpKVHY6ZHdmYVNKdHVUVW1WdWs3Z3FJeGN0UQ==",
		CertificateFingerprint: "4A3618A97FF43A3D92A9E7691D56B49B6746EC29B4889CA6CBFEEB36FAA57D54",
		//Username:               "elastic",
		//Password:               "wXk51bclsiDOPsjG8h_x",
	}

	es, err := elasticsearch.NewClient(cfg)
	if err != nil {
		log.Fatal(err)
	}
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
				"my_index",                        // Index name
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

	////Offliene logger
	app := fiber.New()

	encoderConfig := ecszap.NewDefaultEncoderConfig()
	core := ecszap.NewCore(encoderConfig, os.Stdout, zap.DebugLevel)
	zapLogger := zap.New(core, zap.AddCaller())

	/*
		cfgZ := zap.Config{
			Encoding:         "json",
			Level:            zap.NewAtomicLevelAt(zap.InfoLevel),
			OutputPaths:      []string{"stdout"},
			ErrorOutputPaths: []string{"stderr"},
			EncoderConfig: zapcore.EncoderConfig{
				MessageKey:  "meassage",
				LevelKey:    "level",
				TimeKey:     "@timestamp",
				EncodeTime:  zapcore.ISO8601TimeEncoder,
				EncodeLevel: zapcore.CapitalColorLevelEncoder,
			},
		}

		enc := &prependEncoder{
			Encoder: zapcore.NewConsoleEncoder(cfgZ.EncoderConfig),
			pool:    buffer.NewPool(),
		}

		zapLogger := zap.New(
			zapcore.NewCore(
				enc,
				os.Stdout,
				zapcore.DebugLevel,
			),
			// this mimics the behavior of NewProductionConfig.Build
			zap.ErrorOutput(os.Stderr),
		)
	*/
	zapLogger.Info("this is info")
	zapLogger.Debug("this is debug")
	zapLogger.Warn("this is warn")

	/*cfgZ := zap.Config{
		Encoding:         "json",
		Level:            zap.NewAtomicLevelAt(zap.InfoLevel),
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
		EncoderConfig: zapcore.EncoderConfig{
			MessageKey:  "meassage",
			LevelKey:    "level",
			TimeKey:     "@timestamp",
			EncodeTime:  zapcore.ISO8601TimeEncoder,
			EncodeLevel: zapcore.CapitalColorLevelEncoder,
		},
	}

	zapLogger, err := cfgZ.Build()
	if err != nil {
		log.Panic(err)
	}*/

	zapLogger.Info("Server started",
		zap.String("logger", "ZAP"),
		zap.String("host", "localhost"),
		zap.String("port", "3000"),
	)

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

/*


type prependEncoder struct {
	// embed a zapcore encoder
	// this makes prependEncoder implement the interface without extra work
	zapcore.Encoder

	// zap buffer pool
	pool buffer.Pool
}

// implementing only EncodeEntry
func (e *prependEncoder) EncodeEntry(entry zapcore.Entry, fields []zapcore.Field) (*buffer.Buffer, error) {
	// new log buffer
	buf := e.pool.Get()

	// prepend the JournalD prefix based on the entry level
	buf.AppendString(e.toJournaldPrefix(entry.Level))
	buf.AppendString(" ")

	// calling the embedded encoder's EncodeEntry to keep the original encoding format
	consolebuf, err := e.Encoder.EncodeEntry(entry, fields)
	if err != nil {
		return nil, err
	}

	// just write the output into your own buffer
	_, err = buf.Write(consolebuf.Bytes())
	if err != nil {
		return nil, err
	}

	fmt.Println(buf)
	return buf, nil
}

// some mapper function
func (e *prependEncoder) toJournaldPrefix(lvl zapcore.Level) string {
	switch lvl {
	case zapcore.DebugLevel:
		return "<7>"
	case zapcore.InfoLevel:
		return "<6>"
	case zapcore.WarnLevel:
		return "<4>"
	}
	return ""
}
*/
