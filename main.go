package main

import (
	"flag"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/crowdmob/goamz/aws"
	"github.com/crowdmob/goamz/s3"
	"github.com/gorilla/pat"
)

var s3client *s3.S3
var s3bucket *s3.Bucket
var s3prefix string

func main() {
	bindAddr := ":5253"

	auth, err := aws.EnvAuth()
	if err != nil {
		log.Fatalf("Error getting AWS env auth: [%s]", err)
	}

	s3client = s3.New(auth, aws.EUWest)

	bucket := ""
	prefix := ""

	flag.StringVar(&bindAddr, "bind", bindAddr, "bind address, e.g. :5253")
	flag.StringVar(&bucket, "bucket", bucket, "S3 bucket name")
	flag.StringVar(&prefix, "prefix", prefix, "S3 bucket prefix")
	flag.Parse()

	if len(bucket) == 0 {
		log.Fatal("You must specify a bucket using -bucket")
	}

	s3bucket = s3client.Bucket(bucket)
	s3prefix = prefix

	log.Printf("Listening on %s", bindAddr)

	p := pat.New()
	p.Path("/healthcheck").Methods("GET").HandlerFunc(healthcheck)
	p.Path("/images/{url:.+}").Methods("POST").HandlerFunc(upload)
	p.Path("/images/{url:.+}").Methods("head").HandlerFunc(head)
	p.Path("/images/{url:.+}").Methods("GET").HandlerFunc(download)
	p.Path("/images/{url:.+}").Methods("DELETE").HandlerFunc(remove)

	err = http.ListenAndServe(bindAddr, p)
	if err != nil {
		log.Printf("Error binding to %s: %v", bindAddr, err)
	}
}

func healthcheck(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(200)
}

func upload(w http.ResponseWriter, req *http.Request) {
	path := req.URL.Query().Get(":url")

	b, err := ioutil.ReadAll(req.Body)
	if err != nil {
		w.WriteHeader(400)
		return
	}
	req.Body.Close()

	if len(s3prefix) > 0 {
		path = s3prefix + "/" + path
	}

	res, err := s3bucket.Head(path, nil)
	if err != nil && err.Error() != "404 Not Found" {
		w.WriteHeader(500)
		w.Write([]byte("Error: calling Head: " + err.Error()))
		return
	}

	if res != nil && res.StatusCode < 400 {
		w.WriteHeader(400)
		w.Write([]byte("Error: file already exists"))
		return
	}

	err = s3bucket.Put(path, b, req.Header.Get("Content-Type"), "", s3.Options{})
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte("Error: calling Put: " + err.Error()))
		return
	}

	w.WriteHeader(201)
}

func download(w http.ResponseWriter, req *http.Request) {
	path := req.URL.Query().Get(":url")

	if len(s3prefix) > 0 {
		path = s3prefix + "/" + path
	}

	res, err := s3bucket.Head(path, nil)
	if err != nil && err.Error() != "404 Not Found" {
		w.WriteHeader(500)
		w.Write([]byte("Error: calling Head: " + err.Error()))
		return
	}
	if err != nil {
		w.WriteHeader(404)
		return
	}

	if res != nil && res.StatusCode > 399 {
		w.WriteHeader(404)
		return
	}

	b, err := s3bucket.Get(path)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte("Error: calling Get: " + err.Error()))
		return
	}

	w.WriteHeader(200)
	w.Write(b)
}

func head(w http.ResponseWriter, req *http.Request) {
	path := req.URL.Query().Get(":url")

	if len(s3prefix) > 0 {
		path = s3prefix + "/" + path
	}

	res, err := s3bucket.Head(path, nil)
	if err != nil && err.Error() != "404 Not Found" {
		w.WriteHeader(500)
		w.Write([]byte("Error: calling Head: " + err.Error()))
		return
	}
	if err != nil {
		w.WriteHeader(404)
		return
	}
	w.WriteHeader(res.StatusCode)
}

func remove(w http.ResponseWriter, req *http.Request) {
	path := req.URL.Query().Get(":url")

	if len(s3prefix) > 0 {
		path = s3prefix + "/" + path
	}

	res, err := s3bucket.Head(path, nil)
	if err != nil && err.Error() != "404 Not Found" {
		w.WriteHeader(500)
		w.Write([]byte("Error: calling Head: " + err.Error()))
		return
	}
	if err != nil {
		w.WriteHeader(404)
		return
	}

	if res != nil && res.StatusCode > 399 {
		w.WriteHeader(404)
		return
	}

	err = s3bucket.Del(path)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte("Error: calling Put: " + err.Error()))
		return
	}

	w.WriteHeader(200)
}
