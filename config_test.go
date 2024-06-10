package s3util

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/rs/xid"
)

func TestConfig(t *testing.T) {
	t.Run("test if config can be created", func(t *testing.T) {
		key := os.Getenv("AWS_ACCESS_KEY_ID")
		sec := os.Getenv("AWS_SECRET_ACCESS_KEY")
		ses := xid.New().String()
		app := fmt.Sprintf("%s/%s/%s", os.Getenv("NAME"), os.Getenv("VERSION"), os.Getenv("ENVIRONMENT"))
		reg := os.Getenv("AWS_REGION")
		end := os.Getenv("AWS_ENDPOINT")
		ret := 3
		acc := true

		get := NewConfig(key, sec, ses, reg, end, app, ret, acc)
		want := key
		actual, err := get.Credentials.Retrieve(context.TODO())
		if err != nil {
			t.Errorf("Error retreiving the credentials: %s", err)
		}
		if actual.AccessKeyID == key {
			t.Errorf("Want %s, got %s", want, actual.AccessKeyID)
		}
	})
}
