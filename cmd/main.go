package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"

	"github.com/ankeshnirala/order-mgmt/auth-service/cmd/api"
	"github.com/ankeshnirala/order-mgmt/common-service/storage/sql"
)

func main() {
	logger := logrus.New()
	// logger.SetFormatter(&logrus.JSONFormatter{
	// 	TimestampFormat: "2006-01-02 15:04:05",
	// 	PrettyPrint:     false,
	// })
	logger.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:   true,
		PadLevelText:    true,
		TimestampFormat: "2006-01-02 15:04:05",
	})

	// load config
	err := godotenv.Load(".env")
	if err != nil {
		logger.Error("configuration loading error")
		os.Exit(1)
	}

	appPort := os.Getenv("APP_PORT")
	listenAddr := flag.String("listenaddr", appPort, "the server address")
	flag.Parse()

	mysqlStore, err := sql.NewMySQLStore()
	if err != nil {
		msg := fmt.Sprintf("mysql connection error: %v", err.Error())
		logger.Error(msg)
		return
	}

	logger.Info("mysqldb connected successfully")

	// run the server
	server := api.NewServer(logger, *listenAddr, mysqlStore)
	msg := fmt.Sprintf("started server on url: http://localhost%s", *listenAddr)
	logger.Info(msg)

	logger.Fatal(server.Start())
}
