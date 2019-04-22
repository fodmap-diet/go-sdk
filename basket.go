package sdk

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	google "google.golang.org/appengine"
)

const (
	BASKET_REPO     = "https://raw.githubusercontent.com/fodmap-diet/basket/master/"
	UPDATE_INTERVAL = 120
)

var ItemNotFound error = fmt.Errorf("item not found")
var Failed error = fmt.Errorf("failed to get items")

func downloadFile(filename string, url string) error {
	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Create the file
	out, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	return err
}

func fileDownloadRequired(filename string) bool {
	required := false
	stat, err := os.Stat(filename)
	if os.IsNotExist(err) {
		required = true
	} else {
		mtime := stat.ModTime()
		now := time.Now()
		diff := now.Sub(mtime)
		if int(diff.Seconds()) >= UPDATE_INTERVAL {
			required = true
		}
	}
	return required
}

func parseFile(filename string) (map[string]*Properties, error) {
	items := make(map[string]*Properties)
	jsonfile, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer jsonfile.Close()
	byteValue, _ := ioutil.ReadAll(jsonfile)
	err = json.Unmarshal([]byte(byteValue), &items)
	if err != nil {
		return nil, err
	}
	return items, nil
}

func SearchItem(name string) (*Properties, error) {
	var err error

	// Download the file
	filename := string(name[0]) + ".json"

	var items map[string]*Properties

	// Check if running in app engine
	if google.IsAppEngine() || google.IsDevAppServer() {
		mlock()

		// download file if required
		if mfileDownloadRequired(filename) {
			fileurl := BASKET_REPO + filename
			err = downloadMfile(filename, fileurl)
			if err != nil {
				log.Println(err.Error())
				munlock()
				return nil, Failed
			}
		}

		// parse json file to items
		items, err = parseMfile(filename)
		if err != nil {
			log.Println(err.Error())
			munlock()
			return nil, Failed
		}

		munlock()
	} else {
		// download file if required
		if fileDownloadRequired(filename) {
			fileurl := BASKET_REPO + filename
			err = downloadFile(filename, fileurl)
			if err != nil {
				log.Println(err.Error())
				return nil, Failed
			}
		}

		// parse json file to items
		items, err = parseFile(filename)
		if err != nil {
			log.Println(err.Error())
			return nil, Failed
		}
	}

	// get the item
	item, found := items[name]
	if !found {
		return nil, ItemNotFound
	}

	return item, nil
}
