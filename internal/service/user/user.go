package skeleton

import (
	"context"
	"fmt"

	userEntity "go-tutorial-2020/internal/entity/user"
	"go-tutorial-2020/pkg/errors"
	"go-tutorial-2020/pkg/kafka"
	"net/http"
)

// UserData ...
type UserData interface {
	GetAllUsers(ctx context.Context) ([]userEntity.User, error)
	InsertUsers(ctx context.Context, user userEntity.User) error
	GetUserByNIP(ctx context.Context, NIP string) (userEntity.User, error)
	UpdateUserByNIP(ctx context.Context, NIP string, user userEntity.User) (userEntity.User, error)
	DeleteUserByNIP(ctx context.Context, NIP string) error
	InsertToFirebase(ctx context.Context, user userEntity.User, nipMax int) error
	GetUserFromFirebase(ctx context.Context, page int, size int, nip string) ([]userEntity.User, error)
	GetUserFromFirebaseByNIP(ctx context.Context, NIP string) (userEntity.User, error)
	UpdateuserFirebase(ctx context.Context, NIP string, user userEntity.User) error
	UpdateuserFire(ctx context.Context, user userEntity.User) (userEntity.User, error)
	UpdateuserFireRespon(ctx context.Context, user userEntity.User, respon userEntity.Respons) (userEntity.Respons, error)
	DeleteUserFromFirebase(ctx context.Context, NIP string) error
	NewNIPFirebase(ctx context.Context) (int, error)
	GetUserAPI(ctx context.Context, header http.Header) ([]userEntity.User, error)
	InsertMany(ctx context.Context, userList []userEntity.User) error
	IncrementID(ctx context.Context) (int, error) //buat sql increment NIP bukan ID
}

// Service ...
type Service struct {
	userData UserData
	kafka    *kafka.Kafka
}

// New ...
func New(userData UserData, kafka *kafka.Kafka) Service {
	return Service{
		userData: userData,
		kafka:    kafka,
	}
}

// GetAllUsers ...
func (s Service) GetAllUsers(ctx context.Context) ([]userEntity.User, error) {
	// Panggil method GetAllUsers di data layer user
	users, err := s.userData.GetAllUsers(ctx)
	// Error handling
	if err != nil {
		return users, errors.Wrap(err, "[SERVICE][GetAllUsers] fail query getalluser")
	}
	// Return users array
	return users, err
}

// InsertUsers ...
func (s Service) InsertUsers(ctx context.Context, user userEntity.User) error {
	maxNip, err := s.userData.IncrementID(ctx)

	user.NIP = "P" + fmt.Sprintf("%06d", maxNip)
	err = s.userData.InsertUsers(ctx, user)

	if err == nil {
		fmt.Println(user)
	} else {
		return errors.Wrap(err, "[SERVICE][InsertUsers]")
	}
	return err
}

//GetUserByNIP ...
func (s Service) GetUserByNIP(ctx context.Context, NIP string) (userEntity.User, error) {
	// user, err := s.userData.GetUserByNIP(ctx, NIP)

	// if err != nil {
	// 	return user, errors.Wrap(err, "[SERVICE][GetUserByNIP]")
	// }

	// return user, err
	users, err := s.userData.GetUserByNIP(ctx, NIP)
	if err != nil {
		return users, errors.Wrap(err, "SALAH")
	}
	return users, err

}

//UpdateUserByNIP ...
func (s Service) UpdateUserByNIP(ctx context.Context, NIP string, user userEntity.User) (userEntity.User, error) {
	user, err := s.userData.UpdateUserByNIP(ctx, NIP, user)

	if err != nil {
		return user, errors.Wrap(err, "salah")
	}
	return user, err

}

//DeleteUserByNIP ...
func (s Service) DeleteUserByNIP(ctx context.Context, NIP string) error {
	err := s.userData.DeleteUserByNIP(ctx, NIP)
	if err != nil {
		return errors.Wrap(err, "salah")
	}

	return err
}

//InsertToFirebase ...
func (s Service) InsertToFirebase(ctx context.Context, user userEntity.User) error {
	nipMax, err := s.userData.NewNIPFirebase(ctx)
	nipMax = nipMax + 1
	user.NIP = "P" + fmt.Sprintf("%06d", nipMax)
	err = s.userData.InsertToFirebase(ctx, user, nipMax)
	if err != nil {
		return errors.Wrap(err, "salah")
	}
	return err
}

//GetUserFromFirebase ...
func (s Service) GetUserFromFirebase(ctx context.Context, page int, size int, nip string) ([]userEntity.User, error) {
	//var users []userEntity.User

	users, err := s.userData.GetUserFromFirebase(ctx, page, size, nip)

	return users, err
}

//GetUserFromFirebaseByNIP ...
func (s Service) GetUserFromFirebaseByNIP(ctx context.Context, NIP string) (userEntity.User, error) {

	user, err := s.userData.GetUserFromFirebaseByNIP(ctx, NIP)

	return user, err
}

//UpdateuserFirebase ...
func (s Service) UpdateuserFirebase(ctx context.Context, NIP string, user userEntity.User) error {
	err := s.userData.UpdateuserFirebase(ctx, NIP, user)
	return err
}

func (s Service) UpdateuserFire(ctx context.Context, user userEntity.User) (userEntity.User, error) {
	users, err := s.userData.UpdateuserFire(ctx, user)

	return users, err
}

func (s Service) UpdateuserFireRespon(ctx context.Context, user userEntity.User, respon userEntity.Respons) (userEntity.Respons, error) {

	respon, err := s.userData.UpdateuserFireRespon(ctx,user,respon)

	return respon, err
}

//DeleteUserFromFirebase ...
func (s Service) DeleteUserFromFirebase(ctx context.Context, NIP string) error {
	err := s.userData.DeleteUserFromFirebase(ctx, NIP)
	return err
}

//GetUserAPI ...
func (s Service) GetUserAPI(ctx context.Context, header http.Header) ([]userEntity.User, error) {
	userList, err := s.userData.GetUserAPI(ctx, header)
	return userList, err
}

//PublishUser ...: yang generate date untuk topik
func (s Service) PublishUser(user userEntity.User) error {
	err := s.kafka.SendMessageJSON("New_User", user)

	if err != nil {
		return errors.Wrap(err, "[SERVICE][PublishUser]")
	}

	return err

}

//InsertMany ...
func (s Service) InsertMany(ctx context.Context, userList []userEntity.User) error {
	err := s.userData.InsertMany(ctx, userList)
	return err
}
