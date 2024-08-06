package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
	"github.com/go-chi/chi/v5"
	//"soc/api"
	"github.com/google/uuid"
	//"container/list"
	//"github.com/go-chi/chi/v5/middleware"
	//"github.com/PuerkitoBio/goquery"
	"github.com/fsnotify/fsnotify"
)

func main() {
	fmt.Println("Starting service on port 3000")
	createFolders() //if not exists
	//startNotifyWorker() // start notifing file create
	////////////////////
	watcher, err := fsnotify.NewWatcher()
    if err != nil {
        log.Fatal(err)
    }
    defer watcher.Close()

	go func() {
        for {
            select {
            case event, ok := <-watcher.Events:
                if !ok {
                    return
                }
                //log.Println("event:", event)
                if event.Has(fsnotify.Create) {
					fmt.Printf("New notification: %s\n", event.Name)
					//fmt.Printf("Task %s in progress\n", item.TaskId)
					//////////////////////
					// start workers?
					//////////////////////
                }
            case err, ok := <-watcher.Errors:
                if !ok {
                    return
                }
                log.Println("error:", err)
            }
        }
    }()

	err = watcher.Add("./tasks/worker")
	if err != nil {
        log.Fatal(err)
    }
	err = watcher.Add("./tasks/done")
    if err != nil {
        log.Fatal(err)
    }

	// Block main goroutine forever.
    // <-make(chan struct{})
	//////////////////////
	startWS() //start webserver

}

type hubRequest struct { //struct for post data
	UUID       uuid.UUID
	TaskId     string   `json:"task"` // example: ff81ac90-51b6-42e5-b42e-59b72e4b45d2
	SourceList []string `json:"sources"`
	Period     string   `json:"period"`
	ItemId     string   `json:"itemId"`
}

func startWS() { //start webserver and wait for post data in "/" route
	router := chi.NewRouter()

	router.Post("/", func(w http.ResponseWriter, r *http.Request) {
		var item hubRequest
		err := json.NewDecoder(r.Body).Decode(&item)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}

		item.UUID = uuid.New()
		createTask(item) // make go routine from this func?
		
		jsonItem, err := json.Marshal(item.TaskId)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("illegal json"))
			return
		}
		w.Write(jsonItem)
		fmt.Printf("Task %s in progress\n", item.TaskId)
	})

	err := http.ListenAndServe(":3000", router)
	if err != nil {
		log.Println(err)
	}
}

func createFolders() {
	_, err := os.Stat("tasks")
	if os.IsNotExist(err) {
		errDir := os.MkdirAll("tasks/api", 0755)
		if errDir != nil {
			log.Fatal(err)
		}
		errDir = os.MkdirAll("tasks/worker", 0755)
		if errDir != nil {
			log.Fatal(err)
		}
		errDir = os.MkdirAll("tasks/exporter", 0755)
		if errDir != nil {
			log.Fatal(err)
		}
		fmt.Println("Folders created")
	}
}

func createTask(item hubRequest) {
	f, err := os.Create("tasks/worker/"+item.TaskId + "." + item.UUID.String() + "." + strconv.Itoa(int(time.Now().Unix())) + ".json")
	
	if err != nil {
		log.Fatal(err)
	}
	
	b, err := json.Marshal(item)
	if err != nil {
		fmt.Println(err)
		return
	}

	f.Write(b)
	f.Close()
}

//func startNotifyWorker() {
	
//}

/*

func startNotifyWorker2(){
	watcher, err := fsnotify.NewWatcher() // Initialize an empty watcher
    if err != nil {
        log.Fatal(err)
	}
    defer watcher.Close() // Close watcher at the end of the program

    done := make(chan bool)
    go func() { // Start a coroutine to handle events sent by watcher separately
        for {
            select {
            case event, ok := <-watcher.Events: // Normal event processing logic
                if ! ok {
					return
                }
                log.Println("event:", event)
                if event.Op&fsnotify.Create == fsnotify.Create {
                    log.Println("created file:", event.Name)
					//////////////////////
									// start workers?
					fmt.Println("file created")
					/////////////////
                }
            case err, ok := <-watcher.Errors: // Processing logic when an error occurs
                if ! ok {
					return
                }
                log.Println("error:", err)
            }
        }
    }()
	
	fmt.Println("before watcher.Add")
    err = watcher.Add("./tasks/worker") // Enable watcher to monitor/TMP /foo
    if err != nil {
        log.Fatal(err)
    }
	fmt.Println("after watcher.Add")
    <-done // Make the main coroutine not exit
}

*/

/*2
	//url := "https://api.serpdog.io/search?api_key=66a3cfe3532d58cbb444cfb1&q=go+lang+tutorial&gl=us"

	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	req.Header.Set("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(string(body))
2*/

/*1
url := "https://www.google.com/search?q=go+tutorials&gl=us&hl=en"

req, err := http.NewRequest("GET", url, nil)
if err != nil {
	log.Fatal(err)
}

req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/101.0.4951.54 Safari/537.36")

client := &http.Client{}
res, err := client.Do(req)
if err != nil {
	log.Fatal(err)
}
defer res.Body.Close()

doc, err := goquery.NewDocumentFromReader(res.Body)
if err != nil {
	log.Fatal(err)
}

c := 0
doc.Find("div.g").Each(func(i int, result *goquery.Selection) {
	title := result.Find("h3").First().Text()
	link, _ := result.Find("a").First().Attr("href")
	snippet := result.Find(".VwiC3b").First().Text()

	fmt.Printf("Title: %s\n", title)
	fmt.Printf("Link: %s\n", link)
	fmt.Printf("Snippet: %s\n", snippet)
	fmt.Printf("Position: %d\n", c+1)
	fmt.Println()

	c++
})

1*/

//}
