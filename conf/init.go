package conf

import (
	"ELKExample/repository"
	"ELKExample/usecase"
	"github.com/elastic/go-elasticsearch/v7"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"log"
	"os"
)

func Init() {
	var err error
	logger := zerolog.New(os.Stderr).With().Timestamp().Logger()

	db, err := repository.InitDB("MYSQL_USER","MYSQL_PASSWORD","MYSQL_HOST","MYSQL_PORT","MYSQL_DB",logger)
	if err != nil {
		logger.Err(err).Msg("Connection failed")
		os.Exit(1)
	}
	logger.Info().Msg("Database connection established")

	esClient, err := elasticsearch.NewDefaultClient()
	if err != nil {
		log.Println(err)
		logger.Err(err).Msg("Connection failed")
		os.Exit(1)
	}

	repository.InitMigration(db.Conn)

	h :=  usecase.New(db, esClient, logger)
	router := gin.Default()
	rg := router.Group("/v1")
	h.Register(rg)
	router.Run(":8084")
}
