package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path"
	"strconv"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type Problem struct {
	FrontendQuestionID string `json:"frontendQuestionId"`
	TitleSlug          string `json:"titleSlug"`
}

type FsStruct struct {
	Problems   []Problem `json:"problems"`
	UpdateDate string    `json:"updateDate"`
}

var (
	BUCKET          string
	ACCOUNTID       string
	ACCESSKEYID     string
	ACCESSKEYSECRET string

	LOGPREFIX string

	Problems []Problem
)

var logger = log.New(os.Stdout, LOGPREFIX, log.LstdFlags)

func loadFromR2Once() {
	r2Client := getR2Client()
	if r2Client == nil {
		return
	}
	fname := "leetcode-problems.json"
	out, err := r2Client.GetObject(context.TODO(), &s3.GetObjectInput{
		Bucket: &BUCKET,
		Key:    &fname,
	})
	if err != nil {
		logger.Printf("loadFromR2Once: %v", err)
		return
	}

	var fsStruct FsStruct

	err = json.NewDecoder(out.Body).Decode(&fsStruct)

	if err != nil {
		logger.Printf("loadFromR2Once: %v", err)
		return
	}

	Problems = fsStruct.Problems

	logger.Printf("loadFromR2Once: loaded %d problems from R2", len(Problems))
}

func getR2Client() *s3.Client {
	r2Resolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
		return aws.Endpoint{
			URL: fmt.Sprintf("https://%s.r2.cloudflarestorage.com", ACCOUNTID),
		}, nil
	})

	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithEndpointResolverWithOptions(r2Resolver),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(ACCESSKEYID, ACCESSKEYSECRET, "")),
		config.WithRegion("auto"),
	)

	if err != nil {
		log.Printf("getR2Client: %v", err)
		return nil
	}

	return s3.NewFromConfig(cfg)
}

func getProblemSlug(questionID int) string {
	if quesSlugCache[questionID] != "" {
		return quesSlugCache[questionID]
	}

	if len(Problems) == 0 {
		loadFromR2Once()
	}

	for _, v := range Problems {
		if v.FrontendQuestionID == strconv.Itoa(questionID) {
			quesSlugCache[questionID] = v.TitleSlug
			return v.TitleSlug
		}
	}

	return ""
}

var quesSlugCache = map[int]string{}

func Handler(w http.ResponseWriter, r *http.Request) {
	quesStr := path.Base(r.URL.Path)
	logger.Printf("request quesStr: %s", quesStr)
	quesId, err := strconv.Atoi(quesStr)
	if err != nil {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}

	slug := getProblemSlug(quesId)

	if slug == "" {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("https://leetcode.com/problems/%s/", slug), http.StatusFound)
}
