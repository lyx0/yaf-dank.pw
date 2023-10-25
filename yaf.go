package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/lyx0/yaf/exifscrubber"
	"github.com/lyx0/yaf/fileexpiration"
)

const (
	allowedChars = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	maxAge       = time.Hour * 24 * 7 // 7 days
)

type parameters struct {
	configFile string
}

func parseParams() *parameters {
	configFile := flag.String("configFile", "yaf.conf", "path to config file")
	flag.Parse()

	retval := &parameters{}
	retval.configFile = *configFile
	return retval
}

func main() {
	params := parseParams()

	// Read config
	config, err := ConfigFromFile(params.configFile)
	if err != nil {
		log.Fatalf("could not parse config file: %s\n", err.Error())
	}
	fd := config.FileDir
	fmt.Println(fd)

	handler := uploadHandler{
		config: config,
	}

	if config.ScrubExif {
		scrubber := exifscrubber.NewExifScrubber(config.ExifAllowedIds, config.ExifAllowedPaths)
		handler.exifScrubber = &scrubber
	}

	if config.FileExpiration {
		fmt.Println("FILE EXPIRATION ENABLED")
		go func() {
			for {
				<-time.After(time.Hour * 2)
				fileexpiration.DeleteExpired(fd, maxAge)
			}
		}()
	}

	router := httprouter.New()
	log.Printf("starting yaf on port: \t%d\n", config.Port)
	log.Printf("Maximum File Size:  \t%dMB\n", config.MaxFileSizeMB)
	router.HandlerFunc(http.MethodPost, "/upload", handler.PostUpload)
	router.HandlerFunc(http.MethodPost, "/uploadweb", handler.PostUploadRedirect)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", config.Port), router))
}
