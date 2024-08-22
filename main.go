package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	//"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	//"golang.org/x/text/date"

	//"soc/api"
	"github.com/google/uuid"
	//"container/list"
	//"github.com/go-chi/chi/v5/middleware"
	"github.com/PuerkitoBio/goquery"
	"github.com/fsnotify/fsnotify"
)

func main() {
	fmt.Println("Starting service on port 3000")
	createFolders() //if not exists

	////////////////////
	watcher, err := fsnotify.NewWatcher()
    if err != nil {
        log.Fatal(err)
    }
    defer watcher.Close()

	err = watcher.Add("./tasks/worker")
	if err != nil {
        log.Fatal(err)
    }
	err = watcher.Add("./tasks/done")
    if err != nil {
        log.Fatal(err)
    }

	go func() {
		for {
            select {
            case event, ok := <-watcher.Events:
                if !ok {
                    return
                }
                //log.Println("event:", event)
                if event.Has(fsnotify.Create) {
					//fmt.Printf("New notification: %s\n", event.Name)
						if strings.Contains(event.Name, "worker") {
							var hRworker hubRequest
							dSearch, err := os.Open("./"+event.Name)
							if err != nil {
								log.Fatal(err)
							}
							//fmt.Print(string(dat))
							//err = json.Unmarshal(dat, &hRworker)
							err = json.NewDecoder(dSearch).Decode(&hRworker)
							if err != nil {
								log.Fatal(err)
							}
							//fmt.Print(hRworker)
							dSearch.Close()

							go searchGoogle(hRworker)

							eRem := os.Remove("./"+event.Name) 
							if eRem != nil { 
								log.Fatal(eRem) 
							}
							hRworker = hubRequest{}
						}

						if strings.Contains(event.Name, "done") {
							dat, err := os.ReadFile("./"+event.Name)
							if err != nil {
								log.Fatal(err)
							}
							/*
							/////////////////////////
							var otRS outPutResult
							fSendRes, err := os.Open("./"+event.Name)
							if err != nil {
								log.Fatal(err)
							}
							//fmt.Print(string(dat))
							//err = json.Unmarshal(dat, &hRworker)
							err = json.NewDecoder(fSendRes).Decode(&otRS)
							if err != nil {
								log.Fatal(err)
							}
							//fmt.Print(string(dat))
							fSendRes.Close()
							fmt.Print(string(dRes))
							*/
							
							go sendResultsToHub(dat)

						}
					//////////////////////
					// start workers in new pool routine?
					/*
					f, err := os.Create("tasks/worker/"+item.TaskId + "." + item.UUID.String() + "." + time.Now().Format("2006-01-02T15:04:05") + ".json")
	
					if err != nil {
						log.Fatal(err)
					}
					
					b, err := json.Marshal(item)
					if err != nil {
						fmt.Println(err)
						return
					}
					//test
					f.Write(b)
					f.Close()
					*/
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

	// Block main goroutine forever.
    // <-make(chan struct{})

	//////////////////////
	//startNotifyWorker()
	startWS() //start webserver
}

type hubRequest struct { //struct for post data
	UUID       	uuid.UUID
	TaskId     	string   	`json:"task"` // example: ff81ac90-51b6-42e5-b42e-59b72e4b45d2
	SourceList 	[]string 	`json:"sources"`
	PeriodStart	CustomTime	`json:"periodStart"`
	PeriodEnd	CustomTime	`json:"periodEnd"`
	ItemId     	string		`json:"itemId"`
}

type taskIdResponse struct {
	TaskId	string	`json:"task"`
}

type CustomTime struct {
	time.Time
}

func (t *CustomTime) UnmarshalJSON(b []byte) (err error) {
	//s := strings.Trim(string(b), `"`)
    date, err := time.Parse(`"2006-01-02"`, string(b))
	//date, err := time.Parse(`"2006-01-02"`, string(b))
    if err != nil {
        return err
    }
    t.Time = date
    return
}

func (t *CustomTime) MarshalJSON() ([]byte, error) {
	if t.Time.IsZero() {
		return nil, nil
	}
	return []byte(fmt.Sprintf(`"%s"`, t.Time.Format("2006-01-02"))), nil
}

func (t *CustomTime) GoogleParameterDate() string {
    return t.Format("2006-01-02")
}

type outPutLinks struct {
	Link				string	`json:"link"`
	IndexedDate			string	`json:"indexedDate"`
	IndexedDateFormat	string 	`json:"indexedDateFormat"`
}

type outPutSources struct { //struct for post data
	Source	string   		`json:"sourceName"` // example: ff81ac90-51b6-42e5-b42e-59b72e4b45d2
	Links	[]outPutLinks	`json:"links"`
	Total	int   			`json:"total"`
}

type outPutResult struct { //struct for post data
	TaskId			string   		`json:"task"` // example: ff81ac90-51b6-42e5-b42e-59b72e4b45d2
	Sources			[]outPutSources	`json:"sources"`
	TotalResults	int				`json:"totalResults"`
}

func startWS() { //start webserver and wait for post data in "/" route
	router := chi.NewRouter()
	
	router.Post("/", func(w http.ResponseWriter, r *http.Request) {
		var item hubRequest
		var resp taskIdResponse
		err := json.NewDecoder(r.Body).Decode(&item)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}

		item.UUID = uuid.New()
		
		resp.TaskId = item.TaskId
		jsonResp, err := json.Marshal(resp)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("illegal json"))
			return
		}
		w.Write(jsonResp)

		go createTask(item) // make go routine from this func?

		//var outRes outPutResult = searchGoogle(item)

		//go searchGoogle(item)/////////////////////////////////////////////////////////////////////////

		// jsonItem, err := json.Marshal(outRes)
		// if err != nil {
		// 	w.WriteHeader(http.StatusInternalServerError)
		// 	w.Write([]byte("illegal json"))
		// 	return
		// }
		// w.Write(jsonItem)

	})

	router.Post("/test/", func(w http.ResponseWriter, r *http.Request) {
		var chP outPutResult
		err := json.NewDecoder(r.Body).Decode(&chP)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}
		jsonResp, err := json.Marshal(chP)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("illegal json"))
			return
		}
		w.Write(jsonResp)
		//fmt.Println(chP)
		chP = outPutResult{}
	})

	router.Get("/{task}", func(w http.ResponseWriter, r *http.Request) {
		
		//var resOutPut outPutResult

		task_id := chi.URLParam(r, "task") // 👈 getting path param

		dat, err := os.ReadFile("./tasks/done/"+task_id+".json")
		if err != nil {
			log.Fatal(err)
		}

		/*
		jsonResp, err := json.Marshal(resOutPut)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("illegal json"))
			return
		}
		w.Write(jsonResp)
		*/
		w.Write(dat)
		//opRes.Close()
	})

	err := http.ListenAndServe(":3000", router)
	if err != nil {
		log.Println(err)
	}
}

func createFolders() {
	_, err := os.Stat("tasks/api")
	if os.IsNotExist(err) {
		errDir := os.MkdirAll("tasks/api", 0755)
		if errDir != nil {
			log.Fatal(err)
		}
		errDir = os.MkdirAll("tasks/worker", 0755)
		if errDir != nil {
			log.Fatal(err)
		}
		errDir = os.MkdirAll("tasks/done", 0755)
		if errDir != nil {
			log.Fatal(err)
		}
		//fmt.Println("Folders created")
	}
}

func createTask(item hubRequest) {
	f, err := os.Create("tasks/worker/"+item.TaskId + "." + item.UUID.String() + "." + time.Now().Format("2006-01-02T15:04:05") + ".json")
	
	if err != nil {
		log.Fatal(err)
	}
	//fmt.Println(item)
	err = json.NewEncoder(f).Encode(&item)
	if err != nil {
		fmt.Println(err)
		return
	}

	//f.Write(b)
	f.Close()
}

func writeTaskResults (it outPutResult) {
	f, err := os.Create("tasks/done/"+it.TaskId + ".json")
	
	if err != nil {
		log.Fatal(err)
	}
	//fmt.Println(it)
	err = json.NewEncoder(f).Encode(&it)
	if err != nil {
		fmt.Println(err)
		return
	}

	//f.Write(b)
	f.Close()
}

//https://www.google.com/search?q=21154521+site:www.wildberries.ru+before:2024-08-11&client=safari&sca_esv=f7c375419f60a470&sca_upv=1&sxsrf=ADLYWILBSZsVc_VsQTeQvIw9iIqzCBf37A%3A1723389303832&ei=d9W4ZoG7MqLUwPAPg7SFkA4&ved=0ahUKEwjB0o66ne2HAxUiKhAIHQNaAeIQ4dUDCA8&uact=5&oq=21154521+site%3Awww.wildberries.ru+after:2021-01-01+before:2024-08-11&num=100
/*

type outPutLinks struct {
	Link     	string
	indexedDate string	
}

type outPutSources struct { //struct for post data
	Source			string   	`json:"sourceName"` // example: ff81ac90-51b6-42e5-b42e-59b72e4b45d2
	Links		[]outPutLinks{}	`json:"links"`
	Total			int   		`json:"total"`
}

type outPutResult struct { //struct for post data
	TaskId			string   	`json:"task"` // example: ff81ac90-51b6-42e5-b42e-59b72e4b45d2
	Sources		[]outPutSources{}	`json:"sources"`
}

*/



/*
func rusToEngMonthFormat (sD string) string {
	
	/*monthsMap := map[string]string {
		"янв.":	"January",
		"фев.":	"February",
		"мар.":	"March",
		"апр.":	"April",
		"мая":	"May",
		"июн.":	"June",
		"июл.":	"July",
		"авг.":	"August",
		"сен.": "September",
		"окт.":	"October",
		"ноя.":	"November",
		"дек.":	"December",
	}
	fmt.Println("sD:  ")
	fmt.Println(sD)
	monthsNumMap := map[string]string {
		"янв.":	"01",
		"фев.":	"02",
		"мар.":	"03",
		"апр.":	"04",
		"мая":	"05",
		"июн.":	"06",
		"июл.":	"07",
		"авг.":	"08",
		"сен.": "09",
		"окт.":	"10",
		"ноя.":	"11",
		"дек.":	"12",
	}

	var res string

	rsD := strings.Split(sD, " ")
	if len(rsD[0]) < 2 {
		//res = "0" + rsD[0] + " " + monthsNumMap[rsD[1]] + " " + rsD[2]
		res = rsD[2] + "-" + monthsNumMap[rsD[1]] + "-" + "0" + rsD[0]
	}

	fmt.Println(res)

	return res
}
*/

func sendResultsToHub (b []byte) {
	// var otRS outPutResult
	// fmt.Print(string(b))
	// err := json.Unmarshal(b, otRS)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	
	client := &http.Client{}
	url := "http://localhost.ru/test/:3000"
	////send post
	//jsonSend, err := json.Marshal(b)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(b)))
	if err != nil {
		log.Fatal(err)
	}

	req.Header = http.Header {
		"Accept": {"application/json"},
		"Accept-Language": {"ru-RU,ru;q=0.8,en-US;q=0.5,en;q=0.3"},
	}

	res, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	res.Body.Close()

}


//func searchGoogle (it hubRequest) outPutResult {
func searchGoogle (it hubRequest) {
	time.Sleep(1 * time.Second)
	var outResult outPutResult
	outResult.TaskId = it.TaskId
	outResult.TotalResults = 0

	monthsNumMap := map[string]string {
		"янв.":	"01",
		"фев.":	"02",
		"мар.":	"03",
		"апр.":	"04",
		"мая":	"05",
		"июн.":	"06",
		"июл.":	"07",
		"авг.":	"08",
		"сен.": "09",
		"окт.":	"10",
		"ноя.":	"11",
		"дек.":	"12",
	}

	var rsD []string

	client := &http.Client{}
	for iv, value := range it.SourceList {
		
		var outSrcs outPutSources
		var outLks outPutLinks

		outSrcs.Source = it.SourceList[iv]
		//println(outSrcs.Source)
		//Do something with index and value
		url := "https://www.google.com/search?q="+it.ItemId+"+site:"+value+"+after:"+it.PeriodStart.GoogleParameterDate()+"+before:"+it.PeriodEnd.GoogleParameterDate()+"&num=100"
		//fmt.Println(url)
		req, err := http.NewRequest("GET", url, nil)
	    if err != nil {
	    	log.Fatal(err)
	    }
		
		req.Header = http.Header {
			"User-Agent": {"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/101.0.4951.54 Safari/537.36"},
			"Accept": {"application/json"},
			"Accept-Language": {"ru-RU,ru;q=0.8,en-US;q=0.5,en;q=0.3"},
		}

 		res, err := client.Do(req)
		if err != nil {
			log.Fatal(err)
		}
 		//defer res.Body.Close()

 		doc, err := goquery.NewDocumentFromReader(res.Body)
		if err != nil {
			log.Fatal(err)
		}
		//fmt.Println(doc)
		outSrcs.Total = 0
		//c := 0
		doc.Find("div.g").Each(func(i int, result *goquery.Selection) {
			outLks.Link, _ = result.Find("a").First().Attr("href")
			outLks.IndexedDate = result.Find("span.Sqrs4e").Contents().First().Text()
			
			rsD = strings.Split(outLks.IndexedDate, " ")
			rsYear := rsD[2]
			
			if len(rsD[0]) < 2 {
				outLks.IndexedDateFormat = rsYear[0:4] + "-" + monthsNumMap[rsD[1]] + "-" + "0" + rsD[0]
			}
			outLks.IndexedDateFormat = rsYear[0:4] + "-" + monthsNumMap[rsD[1]] + "-" + rsD[0]
			
			outSrcs.Links = append(outSrcs.Links, outLks)
			//c++
			//outSrcs.Total = c
			outSrcs.Total++
			outResult.TotalResults++
			//fmt.Println(outResult.TotalResults)
		})
		//fmt.Println(outSrcs)
		outResult.Sources = append(outResult.Sources, outSrcs)
		
		res.Body.Close()
	}
	//fmt.Println(outResult)
	writeTaskResults(outResult)
}

/*
func startNotifyWorker() {
	
}

UUID       	uuid.UUID
	TaskId     	string   	`json:"task"` // example: ff81ac90-51b6-42e5-b42e-59b72e4b45d2
	SourceList 	[]string 	`json:"sources"`
	PeriodStart	CustomTime	`json:"periodStart"`
	PeriodEnd	CustomTime	`json:"periodEnd"`
	ItemId     	string		`json:"itemId"`

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
