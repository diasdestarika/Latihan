package user

import (
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"

	userEntity "go-tutorial-2020/internal/entity/user"
	"go-tutorial-2020/pkg/response"
)

// IUserSvc is an interface to User Service
type IUserSvc interface {
	GetAllUsers(ctx context.Context) ([]userEntity.User, error)
	InsertUsers(ctx context.Context, user userEntity.User) error
	GetUserByNIP(ctx context.Context, NIP string) (userEntity.User, error)
	UpdateUserByNIP(ctx context.Context, NIP string, user userEntity.User) (userEntity.User, error)
	DeleteUserByNIP(ctx context.Context, NIP string) error
	InsertToFirebase(ctx context.Context, user userEntity.User) error
	GetUserFromFirebase(ctx context.Context, page int, size int, nip string) ([]userEntity.User, error)
	GetUserFromFirebaseByNIP(ctx context.Context, NIP string) (userEntity.User, error)
	UpdateuserFirebase(ctx context.Context, NIP string, user userEntity.User) error
	UpdateuserFire(ctx context.Context, user userEntity.User) (userEntity.User, error)
	UpdateuserFireRespon(ctx context.Context, user userEntity.User, respon userEntity.Respons) (userEntity.Respons, error)
	UpdateTglLahir(ctx context.Context, structUpdate []userEntity.UserUpdate, respon userEntity.Respons) (userEntity.Respons, error)
	DeleteUserFromFirebase(ctx context.Context, NIP string) error
	GetUserAPI(ctx context.Context, header http.Header) ([]userEntity.User, error)
	PublishUser(user userEntity.User) error
	InsertMany(ctx context.Context, userList []userEntity.User) error
}

type (
	// Handler ...
	Handler struct {
		userSvc IUserSvc
	}
)

// New for user domain handler initialization
func New(is IUserSvc) *Handler {
	return &Handler{
		userSvc: is,
	}
}

// UserHandler will return user data
func (h *Handler) UserHandler(w http.ResponseWriter, r *http.Request) {
	var (
		resp     *response.Response
		metadata interface{}
		result   interface{}
		err      error
		errRes   response.Error
		user     userEntity.User
		update   []userEntity.UserUpdate
		page     int
		size     int
		respon   userEntity.Respons
	)
	// Make new response object
	resp = &response.Response{}
	body, _ := ioutil.ReadAll(r.Body)
	// Defer will be run at the end after method finishes
	defer resp.RenderJSON(w, r)

	switch r.Method {
	// Check if request method is GET
	case http.MethodGet:

		var _type string

		if _, getOK := r.URL.Query()["Get"]; getOK {
			_type = r.FormValue("Get")
		}

		switch _type {
		case "sqlAll":
			result, err = h.userSvc.GetAllUsers(context.Background())
		case "sqlNIP":
			result, err = h.userSvc.GetUserByNIP(context.Background(), r.FormValue("NIP"))
		case "firebaseAll":
			page, err = strconv.Atoi(r.FormValue("page"))
			size, err = strconv.Atoi(r.FormValue("size"))
			result, err = h.userSvc.GetUserFromFirebase(context.Background(), page, size, r.FormValue("nipCari"))
		case "firebaseNIP":
			result, err = h.userSvc.GetUserFromFirebaseByNIP(context.Background(), r.FormValue("NIP"))
		case "API":
			json.Unmarshal(body, &user)

			result, err = h.userSvc.GetUserAPI(context.Background(), r.Header)

		}

	case http.MethodPost:
		var (
			_type    string
			userList []userEntity.User
		)

		if _, fireOK := r.URL.Query()["Insert"]; fireOK {
			_type = r.FormValue("Insert")
		}
		switch _type {

		case "sql":
			json.Unmarshal(body, &user)
			err = h.userSvc.InsertUsers(context.Background(), user)
		case "firebase":
			json.Unmarshal(body, &user)
			err = h.userSvc.InsertToFirebase(context.Background(), user)
		case "kafka":
			json.Unmarshal(body, &user)
			err = h.userSvc.PublishUser(user)
		case "many":
			json.Unmarshal(body, &userList)
			err = h.userSvc.InsertMany(context.Background(), userList)
		}

	case http.MethodPut:
		var (
			_type string
		)

		if _, updateOK := r.URL.Query()["Update"]; updateOK {
			_type = r.FormValue("Update")
		}

		switch _type {
		case "sql":
			json.Unmarshal(body, &user)
			result, err = h.userSvc.UpdateUserByNIP(context.Background(), r.FormValue("NIP"), user)

		case "firebase":
			json.Unmarshal(body, &user)
			err = h.userSvc.UpdateuserFirebase(context.Background(), r.FormValue("NIP"), user)
		case "fire":
			json.Unmarshal(body, &user)
			result, err = h.userSvc.UpdateuserFire(context.Background(), user)
		case "fireRespon":
			json.Unmarshal(body, &user)
			result, err = h.userSvc.UpdateuserFireRespon(context.Background(), user, respon)
		case "updateTgl":
			json.Unmarshal(body, &update)
			result, err = h.userSvc.UpdateTglLahir(context.Background(), update, respon )

		}
	case http.MethodDelete:
		var (
			_type string
		)
		if _, fireOK := r.URL.Query()["Delete"]; fireOK {
			_type = r.FormValue("Delete")
		}
		switch _type {
		case "sql":
			err = h.userSvc.DeleteUserByNIP(context.Background(), r.FormValue("NIP"))
		case "firebase":
			err = h.userSvc.DeleteUserFromFirebase(context.Background(), r.FormValue("NIP"))

		}
	default:
		err = errors.New("400")
	}

	// If anything from service or data return an error
	if err != nil {
		// Error response handling
		errRes = response.Error{
			Code:   101,
			Msg:    "Data Not Found",
			Status: true,
		}
		// If service returns an error
		if strings.Contains(err.Error(), "service") {
			// Replace error with server error
			errRes = response.Error{
				Code:   201,
				Msg:    "Failed to process request due to server error",
				Status: true,
			}
		}

		// Logging
		log.Printf("[ERROR] %s %s - %v\n", r.Method, r.URL, err)
		resp.Error = errRes
		return
	}

	// Inserting data to response
	resp.Data = result
	resp.Metadata = metadata
	// Logging
	log.Printf("[INFO] %s %s\n", r.Method, r.URL)
}
