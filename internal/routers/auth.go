package routers

import (
	"github.com/product-mgmt/auth-service/internal/handlers"
	"github.com/product-mgmt/common-service/constants/endpoints"
	"github.com/product-mgmt/common-service/middleware"
	"github.com/product-mgmt/common-service/utils/commfunc"
)

func (s *Storage) RegisterAuthRoutes() {
	ctrl := handlers.New(s.logger, s.sqlStore)
	midd := middleware.New(s.logger, s.sqlStore)

	publicRoute := s.router.PathPrefix(endpoints.AUTH_BASE_PATH).Subrouter()
	privateRoute := s.router.PathPrefix(endpoints.AUTH_BASE_PATH).Subrouter()

	privateRoute.Use(midd.Authenticate)

	// public routes
	publicRoute.HandleFunc(endpoints.SIGNUP_PATH, commfunc.MakeHTTPHandleFunc(ctrl.SignupHandler)).Methods("POST")
	publicRoute.HandleFunc(endpoints.SIGNIN_PATH, commfunc.MakeHTTPHandleFunc(ctrl.SigninHandler)).Methods("POST")

	// // private routes
	privateRoute.HandleFunc(endpoints.PROFILE_PATH, commfunc.MakeHTTPHandleFunc(ctrl.ProfileHandler)).Methods("GET")
}
