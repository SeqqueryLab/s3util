package main

import (
	"fmt"
	"log"
	"time"

	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/widget"
	s3util "github.com/utubun/s3-util"
)

func listFiles(storage *widget.Label, s3 *s3util.Service) {
	res, err := s3.ListObjectsBucket("ubot")
	if err != nil {
		storage.SetText(fmt.Sprint(err))
	}

	files := make(map[string]interface{})
	for _, val := range res {
		files[*val.Key] = struct {
			Size    int64
			Created string
		}{
			Size:    *val.Size,
			Created: val.LastModified.Format("Jan 01, 2001 03:00:03"),
		}

		log.Printf("Map: %+v", files)
	}

	storage.SetText(fmt.Sprintf("%+v", files))
}

func main() {
	myApp := app.New()
	w := myApp.NewWindow("Storage")

	s3 := s3util.New()
	fls, _ := s3.ListBucket()
	log.Printf("Buckets: %+v", fls)

	storage := widget.NewLabel("")
	listFiles(storage, s3)

	w.SetContent(storage)
	go func() {
		for range time.Tick(1 * time.Second) {
			listFiles(storage, s3)
		}
	}()

	w.ShowAndRun()
	tidyUp()
}

func tidyUp() {
	fmt.Println("Exited")
}
