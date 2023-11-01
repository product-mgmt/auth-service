package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/ankeshnirala/sqlscan"

	"github.com/ankeshnirala/order-mgmt/common-service/constants/messages"
	"github.com/ankeshnirala/order-mgmt/common-service/constants/procedures"
	"github.com/ankeshnirala/order-mgmt/common-service/constants/tables"
	"github.com/ankeshnirala/order-mgmt/common-service/types"
	"github.com/ankeshnirala/order-mgmt/common-service/utils/commfunc"
	"github.com/ankeshnirala/order-mgmt/common-service/utils/jwtauth"
)

func (s *Storage) SignupHandler(w http.ResponseWriter, r *http.Request) error {

	// sync request body data with SignupRequest
	req := new(types.SignupRequest)
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		return err
	}

	// create a context to timeout db operation once work end
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// checking user is registered or not, if not then throw USERNOTREGISTERED error
	rows, err := s.sqlStore.GetRecordByArgs(ctx, procedures.SP_GETRECORD, tables.USERS, "email", req.Email)
	if err != nil {
		return err
	}

	if rows.Next() {
		return fmt.Errorf(messages.EMAILREGISTERED, req.Email)
	}

	// create a new user
	user, err := types.NewUser(req.Name, req.Email, req.Password)
	if err != nil {
		return err
	}

	// register new user in db
	createdUser, err := s.sqlStore.AddReord(ctx, procedures.SP_CREATE_USER, user.Name, user.Email, user.Password)
	if err != nil {
		return err
	}

	var output types.RegisterOutput
	if err := sqlscan.Row(&output, createdUser); err != nil {
		s.logger.Error(fmt.Errorf(messages.USERCREATING, err.Error()))
		return fmt.Errorf(messages.SOMETHINGWENTWRONG)
	}

	resp := types.SignupResponse{
		Message: output.Message,
	}

	return commfunc.WriteJSON(w, http.StatusOK, resp)
}

func (s *Storage) SigninHandler(w http.ResponseWriter, r *http.Request) error {
	// sync request body data with LoginRequest
	req := new(types.SigninRequest)
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		return err
	}

	// create a context to timeout db operation once work end
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// checking user is registered or not, if not then throw USERNOTREGISTERED error
	rows, err := s.sqlStore.GetRecordByArgs(ctx, procedures.SP_GETRECORD, tables.USERS, "email", req.Email)
	if err != nil {
		s.logger.Error(fmt.Errorf("GetRecordByArgs.Error: %s", messages.WRONGSIGNINDETAILS))
		return fmt.Errorf(messages.WRONGSIGNINDETAILS)
	}

	// maping db user user details to user struct
	var user types.User
	if err := sqlscan.Row(&user, rows); err != nil {
		s.logger.Error(fmt.Errorf("sqlscan.Error: %s", messages.WRONGSIGNINDETAILS))
		return fmt.Errorf(messages.WRONGSIGNINDETAILS)
	}

	// checking password
	if ok := user.ValidPassword(req.Passowrd); !ok {
		s.logger.Error(fmt.Errorf("ValidPassword.Error: %s", messages.WRONGSIGNINDETAILS))
		return fmt.Errorf(messages.WRONGSIGNINDETAILS)
	}

	// // generate jwt token and send it in response to login
	token, err := jwtauth.CreateJWT(user.ID)
	if err != nil {
		return err
	}

	// configure signin response
	resp := types.SigninResponse{
		Message: messages.USERSIGNEDIN,
		Token:   token,
	}

	return commfunc.WriteJSON(w, http.StatusOK, resp)
}

func (s *Storage) ProfileHandler(w http.ResponseWriter, r *http.Request) error {
	userId := r.Context().Value(types.CTXKey{Key: "userID"})

	// create a context to timeout db operation once work end
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// checking user is registered or not, if not then throw USERNOTREGISTERED error
	rows, err := s.sqlStore.GetRecordByArgs(ctx, procedures.SP_GETRECORD, tables.USERS, "id", userId)
	if err != nil {
		return fmt.Errorf(messages.USERNOTFOUND)
	}

	// maping db user user details to user struct
	var user types.User
	if err := sqlscan.Row(&user, rows); err != nil {
		return fmt.Errorf(messages.USERNOTFOUND)
	}

	resp := types.ProfileResponse{
		Message: messages.RECORDFETCHED,
		User:    user,
	}

	return commfunc.WriteJSON(w, http.StatusOK, resp)
}
