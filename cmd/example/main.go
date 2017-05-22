package main

import (
	"encoding/json"
	"os"

	"github.com/wacul/go-loggly/loggly/retrieve"
)

func main() {
	client := retrieve.New(os.Getenv("LOGGLY_ACCOUNT"), os.Getenv("LOGGLY_USER_NAME"), os.Getenv("LOGGLY_PASSWORD"))
	searched, err := client.Search().Do(10, "tag:alfa")
	if err != nil {
		panic(err)
	}
	event, err := client.Events().Columns([]string{"syslog.host,syslog.timestamp"}).Do(searched.RSID.ID)
	if err != nil {
		panic(err)
	}
	buf, err := json.MarshalIndent(event, "", "  ")
	if err != nil {
		panic(err)
	}
	println(string(buf))
}
