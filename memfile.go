package sdk

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
	"time"
)

var MFileNotFound error = fmt.Errorf("file not found in memory")

var files map[string]*file
var mux sync.Mutex

type file struct {
	mtime   time.Time
	content []byte
}

func init() {
	files = make(map[string]*file)
}

func mlock() {
	mux.Lock()
}

func munlock() {
	mux.Unlock()
}

func mfileDownloadRequired(filename string) bool {
	required := false
	file, found := files[filename]
	if !found {
		required = true
	} else {
		mtime := file.mtime
		now := time.Now()
		diff := now.Sub(mtime)
		if int(diff.Seconds()) >= UPDATE_INTERVAL {
			required = true
		}
	}
	return required
}

func downloadMfile(filename string, url string) error {
	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// read the file from response
	responseContent, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	// Create the file
	files[filename] = &file{
		mtime:   time.Now(),
		content: responseContent,
	}

	return err
}

func parseMfile(filename string) (map[string]*Properties, error) {
	items := make(map[string]*Properties)

	file, found := files[filename]
	if !found {
		return nil, MFileNotFound
	}

	err := json.Unmarshal(file.content, &items)
	if err != nil {
		return nil, err
	}
	return items, nil
}
