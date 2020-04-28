package user

import "time"

// JAVA EQUIVALENT -> MODEL

// User object model
type User struct {
	ID       int    `db:"id" json:"user_id"`
	NIP      string `db:"nip" json:"nip"`
	Nama     string `db:"nama_lengkap" json:"nama_lengkap"`
	TglLahir time.Time `db:"tanggal_lahir" json:"tgl_lahir"`
	Jabatan  string `db:"jabatan" json:"jabatan"`
	Email    string `db:"email" json:"email"`
}

//DataResp ...
type DataResp struct {
	Data     []User      `json:"data"`
	Metadata interface{} `json:"metadata"`
	Error    interface{} `json:"error"`
}

type Respons struct {
	ID  int    `db:"id" json:"user_id"`
	NIP string `db:"nip" json:"nip"`
}
