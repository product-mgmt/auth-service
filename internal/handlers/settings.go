package handlers

import (
	"github.com/sirupsen/logrus"

	"github.com/ankeshnirala/order-mgmt/common-service/storage"
)

type Storage struct {
	logger   *logrus.Logger
	sqlStore storage.MySQLStorage
}

func New(logger *logrus.Logger, sqlStore storage.MySQLStorage) *Storage {
	return &Storage{
		logger:   logger,
		sqlStore: sqlStore,
	}
}
