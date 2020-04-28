package user

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"go-tutorial-2020/pkg/errors"

	firebaseclient "go-tutorial-2020/pkg/firebaseClient"

	httpclient "go-tutorial-2020/pkg/httpClient"

	userEntity "go-tutorial-2020/internal/entity/user"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"

	"github.com/jmoiron/sqlx"
)

type (
	// Data ...
	Data struct {
		db     *sqlx.DB
		fb     *firestore.Client
		stmt   map[string]*sqlx.Stmt
		client *httpclient.Client
	}

	// statement ...
	statement struct {
		key   string
		query string
	}
)

const (
	// get ALl
	getAllUsers  = "GetAllUsers"
	qGetAllUsers = "SELECT * FROM user_test"

	// insert new user

	insertUsers  = "InsertUsers"
	qInsertUsers = "INSERT INTO user_test VALUES (null,?,?,?,?,?)"

	//getUserByNIP

	getUserByNIP  = "GetUserByNIP"
	qGetUserByNIP = "SELECT * FROM user_test WHERE nip = ?"

	//updateUserByNIP

	updateUserByNIP  = "UpdateUserByNIP"
	qUpdateUserByNIP = "UPDATE user_test SET nama_lengkap = ?, tanggal_lahir = ?, jabatan = ?, email = ? WHERE nip = ? "

	//deleteUserByNIP

	deleteUserByNIP  = "DeleteUserByNIP"
	qDeleteUserByNIP = "DELETE FROM user_test WHERE nip= ?"

	//auto incrementID sql
	incrementNIP  = "IncrementID"
	qIncrementNIP = "SELECT MAX(CAST(RIGHT(nip,6)AS INT)) + 1 FROM user_test"
)

var (
	readStmt = []statement{
		{getAllUsers, qGetAllUsers},
		{insertUsers, qInsertUsers},
		{getUserByNIP, qGetUserByNIP},
		{updateUserByNIP, qUpdateUserByNIP},
		{deleteUserByNIP, qDeleteUserByNIP},
		{incrementNIP, qIncrementNIP},
	}
)

// New ...
func New(db *sqlx.DB, fb *firebaseclient.Client, client *httpclient.Client) Data {
	d := Data{
		db:     db,
		fb:     fb.Client,
		client: client,
	}

	d.initStmt()
	return d
}

func (d *Data) initStmt() {
	var (
		err   error
		stmts = make(map[string]*sqlx.Stmt)
	)

	for _, v := range readStmt {
		stmts[v.key], err = d.db.PreparexContext(context.Background(), v.query)
		if err != nil {
			log.Fatalf("[DB] Failed to initialize statement key %v, err : %v", v.key, err)
		}
	}

	d.stmt = stmts
}

// GetAllUsers digunakan untuk mengambil semua data user
func (d Data) GetAllUsers(ctx context.Context) ([]userEntity.User, error) {
	var (
		user  userEntity.User
		users []userEntity.User
		err   error
	)

	// Query ke database
	rows, err := d.stmt[getAllUsers].QueryxContext(ctx)

	// Looping seluruh row data
	for rows.Next() {
		// Insert row data ke struct user
		if err = rows.StructScan(&user); err != nil {
			return users, errors.Wrap(err, "[DATA][GetAllUsers] ")
		}
		// Tambahkan struct user ke array user
		users = append(users, user)
	}
	// Return users array
	return users, err
}

//ExecContext : cuma ngejalannin querynya

//QueryRowxContext : ada balikin data dari database

//InsertUsers ...
func (d Data) InsertUsers(ctx context.Context, user userEntity.User) error {
	_, err := d.stmt[insertUsers].ExecContext(ctx,

		user.NIP,
		user.Nama,
		user.TglLahir,
		user.Jabatan,
		user.Email,
	)
	return err
}

//GetUserByNIP ...
func (d Data) GetUserByNIP(ctx context.Context, NIP string) (userEntity.User, error) {

	var (
		user userEntity.User
		err  error
	)
	if err = d.stmt[getUserByNIP].QueryRowxContext(ctx,

		NIP,
	).StructScan(&user); err != nil {

		return user, errors.Wrap(err, "SALAH")
	}
	return user, err

}

//UpdateUserByNIP ...
func (d Data) UpdateUserByNIP(ctx context.Context, NIP string, user userEntity.User) (userEntity.User, error) {
	_, err := d.stmt[updateUserByNIP].ExecContext(ctx,
		user.Nama,
		user.TglLahir,
		user.Jabatan,
		user.Email,
		NIP,
	)

	return user, err
}

//DeleteUserByNIP ...
func (d Data) DeleteUserByNIP(ctx context.Context, NIP string) error {
	_, err := d.stmt[deleteUserByNIP].QueryxContext(ctx,
		NIP,
	)

	return err

}

//IncrementID ...
func (d Data) IncrementID(ctx context.Context) (int, error) {
	var maxNip int

	err := d.stmt[incrementNIP].QueryRowxContext(ctx).Scan(&maxNip)

	return maxNip, err
}

//InsertToFirebase ...
func (d Data) InsertToFirebase(ctx context.Context, user userEntity.User, nipMax int) error {
	_, err := d.fb.Collection("maxNip").Doc("Nip").Update(ctx, []firestore.Update{{
		Path: "NIP", Value: nipMax,
	}})
	_, err = d.fb.Collection("user_test").Doc(user.NIP).Set(ctx, user)
	return err
}

//GetUserFromFirebase ...
func (d Data) GetUserFromFirebase(ctx context.Context, page int, size int, nip string) ([]userEntity.User, error) {
	var (
		users []userEntity.User
		err   error
	)
	if page == 1 {
		iter := d.fb.Collection("user_test").OrderBy("NIP", firestore.Asc).Limit(size).Documents(ctx)
		for {
			var user userEntity.User
			doc, err := iter.Next()
			if err == iterator.Done {
				break
			}
			err = doc.DataTo(&user)
			users = append(users, user)
		}
	} else {
		dsnap, _ := d.fb.Collection("user_test").Doc(nip).Get(ctx)
		log.Println(dsnap)
		iter := d.fb.Collection("user_test").OrderBy("NIP", firestore.Asc).Limit(size).StartAfter(dsnap.Data()["NIP"]).Documents(ctx)
		for {
			var user userEntity.User
			doc, err := iter.Next()
			if err == iterator.Done {
				break
			}
			err = doc.DataTo(&user)
			users = append(users, user)
		}
	}

	return users, err
}

//GetUserFromFirebaseByNIP ...
func (d Data) GetUserFromFirebaseByNIP(ctx context.Context, NIP string) (userEntity.User, error) {
	var (
		err  error
		user userEntity.User
	)

	iter := d.fb.Collection("user_test").Where("NIP", "==", NIP).Documents(ctx)

	for {
		doc, err := iter.Next()

		if err == iterator.Done {
			break
		}

		log.Println(doc)
		err = doc.DataTo(&user)
		log.Println(user)
	}

	return user, err
}

//UpdateuserFirebase ...
func (d Data) UpdateuserFirebase(ctx context.Context, NIP string, user userEntity.User) error {
	iter, err := d.fb.Collection("user_test").Doc(NIP).Get(ctx)

	userValidate := iter.Data()

	if userValidate == nil {
		return errors.Wrap(err, "ga ada")
	}
	user.NIP = NIP
	_, err = d.fb.Collection("user_test").Doc(NIP).Set(ctx, user)

	return err
}

func (d Data) UpdateuserFire(ctx context.Context, user userEntity.User) (userEntity.User, error) {
	var users userEntity.User

	NIP := user.NIP
	iter, err := d.fb.Collection("user_test").Doc(NIP).Get(ctx)

	userValidate := iter.Data()

	if userValidate == nil {
		return users, errors.Wrap(err, "ga ada")
	}

	date := "04/22/2020"
	parse, _ := time.Parse(time.RFC3339, date)

	//user.NIP = NIP
	//_, err = d.fb.Collection("user_test").Doc(NIP).Set(ctx, user)

	_, err = d.fb.Collection("user_test").Doc(NIP).Update(ctx, []firestore.Update{
		{Path: "TglLahir", Value: parse},
	})

	return users, err
}

func (d Data) UpdateuserFireRespon(ctx context.Context, user userEntity.User, respon userEntity.Respons) (userEntity.Respons, error) {
	//var users userEntity.User

	NIP := user.NIP
	iter, err := d.fb.Collection("user_test").Doc(NIP).Get(ctx)

	userValidate := iter.Data()

	if userValidate == nil {
		return respon, errors.Wrap(err, "ga ada")
	}

	tgl := user.TglLahir.Format(time.RFC3339)

	fmt.Println(tgl)
	
	layout := "01-02-2006"

	parse, _ := time.Parse(layout, tgl)

	// user.NIP = NIP
	// _, err = d.fb.Collection("user_test").Doc(NIP).Set(ctx, user)
	_, err = d.fb.Collection("user_test").Doc(NIP).Update(ctx, []firestore.Update{
		{Path: "TglLahir", Value: parse},
	})

	//fmt.Println(user.TglLahir)
	respon.ID = user.ID
	respon.NIP = user.NIP

	return respon, err
}

//DeleteUserFromFirebase ...
func (d Data) DeleteUserFromFirebase(ctx context.Context, NIP string) error {

	iter, err := d.fb.Collection("user_test").Doc(NIP).Get(ctx)

	userValidate := iter.Data()

	if userValidate == nil {
		return errors.Wrap(err, "ga ada")
	}
	_, err = d.fb.Collection("user_test").Doc(NIP).Delete(ctx)

	return err

}

//NewNIPFirebase ...
func (d Data) NewNIPFirebase(ctx context.Context) (int, error) {
	doc, err := d.fb.Collection("maxNip").Doc("Nip").Get(ctx)
	max, err := doc.DataAt("NIP")
	nipMax := int(max.(int64))
	return nipMax, err
}

//GetUserAPI ...
func (d Data) GetUserAPI(ctx context.Context, header http.Header) ([]userEntity.User, error) {
	var resp userEntity.DataResp
	var endpoint = "http://10.0.111.143:8888/users?GET=SQL"

	_, err := d.client.GetJSON(ctx, endpoint, header, &resp)

	if err != nil {
		return []userEntity.User{}, errors.Wrap(err, "[DATA][GetUserAPI]")
	}

	return resp.Data, err
}

//InsertMany ...
func (d Data) InsertMany(ctx context.Context, userList []userEntity.User) error {
	var (
		err error
	)
	for _, i := range userList {
		_, err = d.fb.Collection("user_test").Doc(i.NIP).Set(ctx, i)
	}
	return err
}
