package routers

import (
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"

	"github.com/product-mgmt/common-service/middleware"
	"github.com/product-mgmt/common-service/storage"
	"github.com/product-mgmt/common-service/utils/commfunc"
)

type Storage struct {
	logger   *logrus.Logger
	router   *mux.Router
	sqlStore storage.MySQLStorage
}

func New(logger *logrus.Logger, router *mux.Router, sqlStore storage.MySQLStorage) *Storage {
	return &Storage{
		logger:   logger,
		router:   router,
		sqlStore: sqlStore,
	}
}

func (s *Storage) RegisterDefaultMiddleware() {
	midd := middleware.New(s.logger, s.sqlStore)

	s.router.MethodNotAllowedHandler = commfunc.MakeHTTPHandleFunc(midd.MethodNotAllowedHandler)
	s.router.NotFoundHandler = commfunc.MakeHTTPHandleFunc(midd.PathNotFoundHanler)
}
