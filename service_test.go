package s3util

import (
	"fmt"
	"os"
	"reflect"
	"testing"

	"github.com/rs/xid"
)

func TestService(t *testing.T) {
	t.Run("test if config can be created", func(t *testing.T) {
		key := os.Getenv("AWS_ACCESS_KEY_ID")
		sec := os.Getenv("AWS_SECRET_ACCESS_KEY")
		ses := xid.New().String()
		uid := fmt.Sprintf("%s/%s/%s", os.Getenv("NAME"), os.Getenv("VERSION"), os.Getenv("ENVIRONMENT"))
		reg := os.Getenv("AWS_REGION")
		end := os.Getenv("AWS_ENDPOINT")
		ret := 3
		acc := false

		srv := NewFromConfig(uid, key, sec, ses, reg, end, ret, acc)

		res, err := srv.ListBucket()
		if err != nil {
			t.Errorf("failed to initialize the client: %s", err)
		}

		want := "[]types.Bucket"
		got := reflect.TypeOf(res).String()

		if want != got {
			t.Errorf("failed to list buckets with the client: got %s want %s", want, got)
		}
	})
}
