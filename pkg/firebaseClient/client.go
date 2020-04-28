package firebaseclient

import (
	"context"
	"encoding/json"

	"go-tutorial-2020/internal/config"
	"go-tutorial-2020/pkg/errors"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/option"
)

var (
	sharedClient = &firestore.Client{}
	credentials  = map[string]string{
		"type": "service_account",
		"project_id": "fbtest-4c1d3",
		"private_key_id": "befefdbaa53f7a18e0b01efef60527a65b17116f",
		"private_key": "-----BEGIN PRIVATE KEY-----\nMIIEvQIBADANBgkqhkiG9w0BAQEFAASCBKcwggSjAgEAAoIBAQDV8Z7ah8ZG3G3R\n8lwy1CQJygZKqOs30zyD8obc8Xb+yJsqGlHW2I2OeU1raz9KxCm9SmiG/A/LLhAe\n8wxgJNDZ5lHvfbFBmdxnlAzI7Xh2isCdPF+KIUqAmkVSgl8P7fnUriXrBv7vwCRg\n8Wd9VyNYDOpQkFqvNcqYmgAdR9v6PSG+cgFDTbTYzc1r54M0Dud6NQRa0z4esZsr\nZSCvQiVMJ3dnGfg6JNMdPyUHEpbGr/DqjLG0NUR108s0dv/Yhm1FiJdwx6wo90NK\nRAQl4XDTf6xgjPcSqFwHmMZDKEgcZ9krZwEuFR/q2/KDOIydvN6FQU/fXkYwPn7L\n1GUO4bz1AgMBAAECggEABNsKw1lyXSswEo7ZBliiZ8F4zUAntIORZko+bYHNH6u5\nunsl8TBxiAgF7eKdfrV/NVoHZMw65zOIr0x4swRAOY3SEr4sJCRukurn3mCQtm9g\nP0FT8XHfe8CPQw0370feqXPOuFaXOdBHKGxvhZymbfS/FayCY7h0EcIUWhXGhEJN\nI1TlMuQVmUaJ5gbIjznhnMLyLR9shqYIo2G67Z7+ypkhSyUYZB+lHwYWKMXVQ9Vu\nG3UdiGQMsrGLk0f/u2NjA2Ecpvi93Lvqm9VQKx2dD8Bp9oMd0+OXHAJsWaoNQwhg\nvBiLh8anFCKNizuEXFrSo/NjfnXzvcccFipWZxm+QQKBgQDvONm+xI44zwuQyB81\nVHpqPVjyeKfBo9VYkvdFlMesUCNCjmx5f15O5U2jk1XzlQIqxWLJThU/fjnNxUk3\nlMyLPOKUCZ1H13+85QmSF7Q8O0580TRGF7cHEYeVBW2PHVSqRApks5a3mDRDvQj0\n1BajOgr8/iCbsU6oL6QVdNbWNQKBgQDk8uhqO2dv/F4q7gukV8fdAcl6cW9CAds/\n9OOyvY1KbUOpgRthoklZm0ImpR6a2MAKztOJ6MHREA6FFlQPbEfcxG2c7O2UU69f\niqH76e1vy4c/Q26bCQMz9Mj5HjipsVQjOBc+Swu8OGKxaDeqgy545BWW2KxetpyP\naM/9Mz0jwQKBgQCA7KqtXVEo3KznAnOPUlAHIbjmNJB0k89PRSVuophaDXZzUD61\n3Cb/biVBmw4fkJbyZh6vTx20clrEwyaKhe8Wu2GBVw0kwsddDjLyQUQpkezi5/y8\nKdvCO3hOn/ZDwxL2EGVpkEASAj1opGBHUmZA4e86GduJDS3PBp3v0mBWYQKBgCGQ\nmxI380ovrX6Nt5c4Z0y3XlpdFvqOWx5dQKSLtZMbwbev/duqdyZz5JbVzk7VSBJN\nkCW/wepseDR6uYgpT7/F7Gv9MDd2rVdMc8MC4JRrOkDEGgsQny+Wy3/6NkRqgvNG\n3eF8DxRhD9cCeGa/JKkEh0W+LkcUbo93xkZQpL4BAoGAK6kmLedZB8cxbIXYw701\nWbxZJaxbOc8cdUu5rST0oy3Al49nivJE2hSZaDtmZLD2uoxatkMQjz5Si85h5NhP\n7Q1Pmmyjf/dN4/B36/5gjkSRUOJCra4aCaqINCq0DIWfZG8KZnB9WgywSzF3gEwR\nqJe1M0+dtxDNR9QZfj9s5do=\n-----END PRIVATE KEY-----\n",
		"client_email": "firebase-adminsdk-ow4la@fbtest-4c1d3.iam.gserviceaccount.com",
		"client_id": "113984826458601744666",
		"auth_uri": "https://accounts.google.com/o/oauth2/auth",
		"token_uri": "https://oauth2.googleapis.com/token",
		"auth_provider_x509_cert_url": "https://www.googleapis.com/oauth2/v1/certs",
		"client_x509_cert_url": "https://www.googleapis.com/robot/v1/metadata/x509/firebase-adminsdk-ow4la%40fbtest-4c1d3.iam.gserviceaccount.com",
		}
)

// Client ...
type Client struct {
	Client *firestore.Client
}

// NewClient ...
func NewClient(cfg *config.Config) (*Client, error) {
	var c Client
	cb, err := json.Marshal(credentials)
	if err != nil {
		return &c, errors.Wrap(err, "[FIREBASE] Failed to marshal credentials!")
	}
	option := option.WithCredentialsJSON(cb)
	c.Client, err = firestore.NewClient(context.Background(), cfg.Firebase.ProjectID, option)
	if err != nil {
		return &c, errors.Wrap(err, "[FIREBASE] Failed to initiate firebase client!")
	}
	return &c, err
}
