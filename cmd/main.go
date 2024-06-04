package main

import (
	"fmt"
	"log"

	s3util "github.com/utubun/s3-util"
)

func main() {
	s := s3util.New()

	/* 	err := s.UploadObjectMultipart("ubot", "sample.fastq.gz")
	   	if err != nil {
	   		fmt.Println(err)
	   	} */

	res, err := s.ListObjectsBucket("seqquery-earlyapp", "ubot")
	if err != nil {
		log.Fatal(err)
	}
	for _, val := range res {
		fmt.Println(*val.Key)
	}
}
