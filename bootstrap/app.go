package bootstrap

import (
	"context"
	"log"

	"github.com/IBM/sarama"
	"github.com/afrizalsebastian/ai-cv-evaluator-with-go/config"
	chromaclient "github.com/afrizalsebastian/ai-cv-evaluator-with-go/modules/chroma-client"
	geminiclient "github.com/afrizalsebastian/ai-cv-evaluator-with-go/modules/gemini-client"
	gomysql "github.com/afrizalsebastian/ai-cv-evaluator-with-go/modules/go-mysql"
	ingestdocument "github.com/afrizalsebastian/ai-cv-evaluator-with-go/modules/ingest-document"
	"github.com/afrizalsebastian/ai-cv-evaluator-with-go/modules/kafka"
	"gorm.io/gorm"
)

type Application struct {
	ENV           *config.Config
	GeminiClient  geminiclient.IGeminiClient
	ChromaClient  chromaclient.IChromaClient
	Ingest        ingestdocument.IIngestFile
	DB            *gorm.DB
	KafkaProducer *kafka.Producer
}

func NewApp() *Application {
	ctx := context.Background()
	app := &Application{}

	if err := config.Init(); err != nil {
		log.Fatal("failed to initialize configuration")
	}

	app.ENV = config.Get()

	// Init DB
	dbConfig := &gomysql.MysqlConfig{
		DBUser:     app.ENV.DBUser,
		DBPassword: app.ENV.DBPassword,
		DBHost:     app.ENV.DBHost,
		DBPort:     app.ENV.DBPort,
		DBName:     app.ENV.DBName,
	}
	db, err := gomysql.NewDatabaseConnection(dbConfig)
	if err != nil {
		log.Fatal("failed to create db connection")
	}
	app.DB = db

	// Init Gemini Client
	geminiCient, err := geminiclient.NewGeminiAiCLient(ctx, app.ENV.GeminiApiKey, app.ENV.GeminiModel)
	if err != nil {
		log.Fatal("failed to init gemini client")
	}
	app.GeminiClient = geminiCient

	// Init chroma
	chromaClient, err := chromaclient.NewChromaClient(ctx, app.ENV.ChromaUrl)
	if err != nil {
		log.Fatalf("failed to init chroma client, %s", err.Error())
	}
	app.ChromaClient = chromaClient

	// Init ingestDocument
	ingesDocument := ingestdocument.NewIngestFile(chromaClient)
	app.Ingest = ingesDocument

	// Init Kafka Producer client
	kafkaProducer, err := kafka.NewProducer(
		app.ENV.KafkaBroker,
		func(c *sarama.Config) {
			c.Net.SASL.Enable = app.ENV.KafkaSASLEnable
			c.Net.SASL.Handshake = app.ENV.KafkaSASLHandshake
			c.Net.TLS.Enable = app.ENV.KafkaTLS
		},
	)
	if err != nil {
		log.Println("Kafka producer failed to initialize")
	}
	app.KafkaProducer = kafkaProducer

	return app
}
