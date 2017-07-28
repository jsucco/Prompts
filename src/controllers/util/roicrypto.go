package util

import (
	"errors"
	"os"
	"github.com/gorilla/securecookie"
)

var (
	hashkey string
	blockkey string
	err error
)

func Decode(content string) (string, error) {

	hashkey, err = mustGetenv("ROI_HASHKEY")
	if err == nil {
		blockkey, err = mustGetenv("ROI_BLOCKKEY")
		if err == nil{
			var s= securecookie.New([]byte(hashkey), []byte(blockkey))

			value := make(map[string]string)
			if err := s.Decode("Key", content, &value); err == nil {
				return value["Key"], nil
			}
			return "", err
		} else {
			return "", errors.New("Failed to retrieve blockkey")
		}
	} else {
		return "", errors.New("Failed to retrieve hashkey")
	}
}

func Encode(content string) (string, error) {
	hashkey, err = mustGetenv("ROI_HASHKEY")
	if err == nil {
		blockkey, err = mustGetenv("ROI_BLOCKKEY")
		if err == nil{
			var s= securecookie.New([]byte(hashkey), []byte(blockkey))

			content_map := map[string]string{
				"Key": content,
			}
			if encoded, errs := s.Encode("Key", content_map); errs == nil {
				return encoded, nil
			}
			return "", errors.New("Encryption failed.")
		} else {
			return "", errors.New("Failed to retrieve blockkey")
		}
	} else {
		return "", errors.New("Failed to retrieve hashkey")
	}

}

func mustGetenv(k string) (string, error) {
	v := os.Getenv(k)
	if v == "" {
		return "", errors.New("environment variable not set.")
	}
	return v, nil
}