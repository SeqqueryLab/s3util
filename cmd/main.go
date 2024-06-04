package main

import (
	"fmt"
	"log"
	"os"

	s3util "github.com/utubun/s3-util"
)

func main() {
	s := s3util.New()
	f, _ := os.Open("simple.csv")
	stat, _ := f.Stat()
	buff := make([]byte, stat.Size())
	f.Read(buff)
	s.UploadObject("ubot", "simple.csv", buff)

	/* 	err := s.UploadObjectMultipart("ubot", "sample.fastq.gz")
	   	if err != nil {
	   		fmt.Println(err)
	   	} */

	res, _ := s.ListObjectsBucket("ubot")
	for _, val := range res {
		fmt.Println(*val.Key)
	}

	_, err := s.SelectObjectContent("ubot", "simple.csv", "")
	log.Fatal(err)
}
