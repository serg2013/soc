package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	//"net/url"
	"os"
	"strings"
	"time"
	"errors"
    "regexp"

	"github.com/go-chi/chi/v5"
	//"golang.org/x/text/date"

	//"soc/api"
	//"io"
	"github.com/google/uuid"
	//"container/list"
	//"github.com/go-chi/chi/v5/middleware"
	//"github.com/PuerkitoBio/goquery"
	"github.com/fsnotify/fsnotify"
	//api2captcha "github.com/2captcha/2captcha-go"
)

func main() {
	fmt.Println("Starting service on port 3000")
	/*
	clientCT := api2captcha.NewClient("004f1af9da5d225d5ca67f84f4f71ce7")
	//urlS = "https://www.google.com/sorry/index?continue=" + urlS
	captchaT := api2captcha.ReCaptcha {
		SiteKey: "6LfwuyUTAAAAAOAmoS0fdqijC2PbbdH4kjq62Y1b",
		Url: "https://www.google.com/sorry/index?continue=https://www.google.com/search%3Fq%3D21154521%26sca_esv%3Ddb73fa7fbeeadfbb%26sca_upv%3D1%26sxsrf%3DADLYWIKMUEeRZFt8J46mcy0t3UHd7S_5tA%253A1725166856290%26ei%3DCPXTZpK4EfiD1fIP8MyBcA%26ved%3D0ahUKEwiS3Jat-6CIAxX4QVUIHXBmAA4Q4dUDCA8%26uact%3D5%26oq%3D21154521%26gs_lp%3DEgxnd3Mtd2l6LXNlcnAiCDIxMTU0NTIxMgQQIxgnMgQQIxgnMgQQIxgnMggQABiABBiiBDIIEAAYgAQYogQyCBAAGIAEGKIEMggQABiABBiiBDIIEAAYgAQYogRI1RNQwwpYkw5wAXgAkAEAmAFFoAGFAaoBATK4AQPIAQD4AQGYAgKgAosBmAMAiAYBkgcBMqAH_A8%26sclient%3Dgws-wiz-serp&q=EgRT2cg9GP-d0bYGIjDl3vvqSmZDYSU6Q9l3FUf0MPU_B7MJJ8_y1uSSr56E35AVg1ZzvtsEVH70rYNklpYyAXJaAUM",
		DataS: "dpuvbjFlVr5Eub1add_JRcUkIjliBcvzNQQacszYXjUBLdrnIWxgI-fHNoTfDBRAmxqzFdqCxMkk30_o45ns4EyHDaMoTMG4CoDWUwmRDsi8ft4h5kRwbsDS0zewsigKxiTZTQeuBszZBZ807EepnkeMHziigCePVF0RR2NVh2khch7QRpEaTHRNsG_mj8mpE80w4khLHXs_LPVHexxJ-zB2mkui78BkbFqabLbm4cadWUh-R1cWQtn2V2h7vyIec9VQFMSdg40k-tQ0aYVfEL3cpItsZDs",
	}
	
	reqC := captchaT.ToRequest()

	//reqC.SetProxy("HTTP", "N5LCWU:vZNbygbj4z@188.130.142.27:1050")

	//Url: "https://www.google.com/recaptcha/api2/demo",
	//Url: "https://www.google.com/search?q=21154521+site:www.wildberries.ru+after:2019-01-01+before:2024-08-11&num=100",
	//code, err := clientCT.Solve(req)
	//https://www.google.com/sorry/index?continue=https://www.google.com/search%3Fq%3D21154521&q=EgRT2cg9GJKCzrYGIjCQRdzg27kXj5QOhTutRM-vqRWlMv0CxFm7TWFh0al-nwk6EmxxeDILJ21ufu4FpoIyAXJaAUM

	codeT, smthT, err := clientCT.Solve(reqC)
	if err != nil {
		log.Fatal(err);
	}
	fmt.Println("code "+codeT)
	fmt.Println("smth "+smthT)*/
	////////////////////////////////////////////////////////////////////////////////////
	//makeCaptchaAPICall()

	createFolders() //if not exists

	////////////////////
	watcher, err := fsnotify.NewWatcher()
    if err != nil {
        log.Fatal(err)
    }
    defer watcher.Close()

	// err = watcher.Add("./tasks/workerYandex")
	// if err != nil {
    //     log.Fatal(err)
    // }

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
						/*
						fmt.Printf("New notification: %s\n", event.Name)
							if strings.Contains(event.Name, "worker") {
								var hRworker hubRequest
				
								dSearch, err := os.Open(event.Name)
								if err != nil {
									log.Fatal(err)
								}
								
								err = json.NewDecoder(dSearch).Decode(&hRworker)
								if err != nil {
									log.Fatal(err)
								}
								
								dSearch.Close()
								
								searchGoogle(hRworker)
								//go searchGoogle(hRworker)

								eRem := os.Remove(event.Name) 
								if eRem != nil { 
									log.Fatal(eRem)
								}
								hRworker = hubRequest{}
							}
						*/
							if strings.Contains(event.Name, "done") {
								dat, err := os.ReadFile(event.Name)
								if err != nil {
									log.Fatal(err)
								}
								go sendResultsToHub(dat)			
							}
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
	SourceList 	string 		`json:"source"`
	PeriodStart	CustomTime	`json:"periodStart"`
	PeriodEnd	CustomTime	`json:"periodEnd"`
	ItemId     	string		`json:"itemId"`
}

type taskIdResponse struct {
	TaskId	string	`json:"task"`
	Status 	string	`json:"status"`
}

type CustomTime struct {
	time.Time
}

func (t *CustomTime) UnmarshalJSON(b []byte) (err error) {
	//s := strings.Trim(string(b), `"`)
	//fmt.Println("UNMARSHALED")
    date, err := time.Parse(`"2006-01-02"`, string(b))
	//date, err := time.Parse(`"2006-01-02"`, string(b))
    if err != nil {
        return err
    }
    t.Time = date
    return
}

func (t *CustomTime) MarshalJSON() ([]byte, error) {
	//fmt.Println("MARSHALED")
	if t.Time.IsZero() {
		return nil, nil
	}
	return []byte(fmt.Sprintf(`"%s"`, t.Time.Format("2006-01-02"))), nil
}

func (t *CustomTime) GoogleParameterDate() string {
    return t.Format("2006-01-02")
}
/*
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
*/
func IsExist(str, filepath string) bool {
	b, err := os.ReadFile(filepath)
	if err != nil {
			panic(err)
	}

	isExist, err := regexp.Match(str, b)
	if err != nil {
			panic(err)
	}
	return isExist
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
		if item.TaskId == "" {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("EMPTY TASK ID"))
			return
		}

		if IsExist(item.TaskId, "tasks/api/bd") {
			w.Write([]byte("TASK ALREADY EXISTS, CHECK RESULSTS BY GET REQUEST"))
			return
		}

		if item.SourceList == "" {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("EMPTY SOURCE"))
			return
		}
		if item.ItemId == "" {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("EMPTY ARTIClE"))
			return
		}
		item.UUID = uuid.New()
		
		resp.TaskId = item.TaskId
		resp.Status = "in progress"

		fBD, err := os.OpenFile("tasks/api/bd", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0755)
		if err != nil {
			log.Fatal(err)
		}

		defer fBD.Close()

		if _, err = fBD.WriteString(item.TaskId+"\n"); err != nil {
			log.Fatal(err)
		}
	
		jsonResp, err := json.Marshal(resp)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("ILLEGAL JSON"))
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

	/*router.Post("/test/", func(w http.ResponseWriter, r *http.Request) {
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
	})*/

	router.Get("/{task}", func(w http.ResponseWriter, r *http.Request) {
		
		var resP taskIdResponse

		task_id := chi.URLParam(r, "task") // 👈 getting path param
		pathFile := "tasks/done/"+task_id+".json"
		
		if !IsExist(task_id, "tasks/api/bd") {
			w.Write([]byte("THIS IS NEW TASK, PLEASE SEND via POST REQUEST"))
			return
		}

		if _, err := os.Stat(pathFile); errors.Is(err, os.ErrNotExist) {
			//fmt.Println("path/to/whatever does not exist")
			//w.Write([]byte("No such task ID! Try again, please."))
			resP.TaskId = task_id
			resP.Status = "in progress"
			jsonResp, err := json.Marshal(resP)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte("illegal json"))
				return
			}
			w.Write(jsonResp) 
		} else {
			dat, err := os.ReadFile(pathFile)
			if err != nil {
				log.Fatal(err)
			}
			w.Write(dat)
		}

		//dat, err := os.ReadFile(pathFile)


		/*
		jsonResp, err := json.Marshal(resOutPut)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("illegal json"))
			return
		}
		w.Write(jsonResp)
		*/
		//w.Write(dat)
		//opRes.Close()
	})

	err := http.ListenAndServe(":3000", router)
	if err != nil {
		log.Println(err)
	}
}

func createFolders() {
	if _, err := os.Stat("tasks"); os.IsNotExist(err) {
		errDir := os.MkdirAll("tasks/api", 0755)
		if errDir != nil {
			log.Fatal(errDir)
		}
		_, errF := os.Create("tasks/api/bd")
		if errF != nil {
			log.Fatal(errF)
		}
		errDir = os.MkdirAll("tasks/workerYandex", 0755)
		if errDir != nil {
			log.Fatal(errDir)
		}
		errDir = os.MkdirAll("tasks/workerGoogle", 0755)
		if errDir != nil {
			log.Fatal(errDir)
		}
		errDir = os.MkdirAll("tasks/done", 0755)
		if errDir != nil {
			log.Fatal(errDir)
		}
	 } 
}

func createTask(item hubRequest) {
	//var tPaths []string
	tPath := ""
	switch item.SourceList {
		case "www.wildberries.ru":
			tPath = "tasks/workerGoogle/"
			//tPaths = append(tPaths, "tasks/workerGoogle/")
		case "www.ozon.ru":
			tPath = "tasks/workerYandex/"
			//tPaths = append(tPaths, "tasks/workerGoogle/")
		case "vk.com":
			tPath = "tasks/workerYandex/"
		case "dzen.ru":
			tPath = "tasks/workerYandex/"
		case "ok.ru":
			tPath = "tasks/workerYandex/"
		case "t.me":
			tPath = "tasks/workerYandex/"
			//tPaths = append(tPaths, "tasks/workerYandex/")
		case "www.facebook.com":
			tPath = "tasks/workerGoogle/"
		case "www.instagram.com":
			tPath = "tasks/workerGoogle/"
		case "www.youtube.com":
			tPath = "tasks/workerGoogle/"
		case "www.tiktok.com":
			tPath = "tasks/workerGoogle/"
	}
	f, err := os.Create(tPath + item.TaskId + "." + item.UUID.String() + "." + time.Now().Format("2006-01-02T15:04:05") + ".json")

	if err != nil {
		log.Fatal(err)
	}
	//fmt.Println(item)
	err = json.NewEncoder(f).Encode(&item)
	if err != nil {
		fmt.Println(err)
		return
	}

	f.Close()
	//f.Write(b)	
}
/*
func writeTaskResults (it outPutResult) {
	f, err := os.Create("tasks/done/"+it.TaskId + ".json")
	//itUUID := uuid.New()
	//f, err := os.Create("tasks/done/" + it.TaskId + "." + itUUID.String() + "." + time.Now().Format("2006-01-02T15:04:05") + ".json")
	
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
*/
//https://www.google.com/search?q=21154521+site:www.wildberries.ru+before:2024-08-11&client=safari&sca_esv=f7c375419f60a470&sca_upv=1&sxsrf=ADLYWILBSZsVc_VsQTeQvIw9iIqzCBf37A%3A1723389303832&ei=d9W4ZoG7MqLUwPAPg7SFkA4&ved=0ahUKEwjB0o66ne2HAxUiKhAIHQNaAeIQ4dUDCA8&uact=5&oq=21154521+site%3Awww.wildberries.ru+after:2021-01-01+before:2024-08-11&num=100
/*
/

AEC=AVYB7coCTnjhExyZXrOAJKB1tEmb1xGkdHhr5Tm18pPJGuQ998VjQyaIWpw;NID=516=q2GOm6QjMcr4FHz9TZndMcepIBpO_ras7aN_yji4lz3YCKywfKY143PjNcxcnzuU939kIkYnsZTIrDCWFQsFBjepT0PYv282-4JOxacsonT5ji79rVs1v3GPugC0cQxFvxs0QFAfxiRsLwBDioBOVNbyfUfeHXaXnEnAW4USN-YR2H4bS32LjG651pA0BcllRqVgzOx62_DgV_hF31WNkaY

[
    {
        "domain": ".google.com",
        "expirationDate": 1739907915.434233,
        "hostOnly": false,
        "httpOnly": true,
        "name": "AEC",
        "path": "/",
        "sameSite": "lax",
        "secure": true,
        "session": false,
        "storeId": null,
        "value": "AVYB7coCTnjhExyZXrOAJKB1tEmb1xGkdHhr5Tm18pPJGuQ998VjQyaIWpw"
    },
    {
        "domain": ".google.com",
        "expirationDate": 1740167135.467891,
        "hostOnly": false,
        "httpOnly": true,
        "name": "NID",
        "path": "/",
        "sameSite": "no_restriction",
        "secure": true,
        "session": false,
        "storeId": null,
        "value": "516=q2GOm6QjMcr4FHz9TZndMcepIBpO_ras7aN_yji4lz3YCKywfKY143PjNcxcnzuU939kIkYnsZTIrDCWFQsFBjepT0PYv282-4JOxacsonT5ji79rVs1v3GPugC0cQxFvxs0QFAfxiRsLwBDioBOVNbyfUfeHXaXnEnAW4USN-YR2H4bS32LjG651pA0BcllRqVgzOx62_DgV_hF31WNkaY"
    }
]

/
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

/*
func captchaApiRequest (b []byte) {
	// var otRS outPutResult
	// fmt.Print(string(b))
	// err := json.Unmarshal(b, otRS)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	
	proxyPool := []string {
		"N5LCWU:vZNbygbj4z@https://188.130.142.27:1050", 
		"N5LCWU:vZNbygbj4z@https://109.248.12.98:1050",
		"N5LCWU:vZNbygbj4z@https://45.140.52.41:1050",
		"N5LCWU:vZNbygbj4z@https://45.134.180.201:1050",
		"N5LCWU:vZNbygbj4z@https://45.134.180.168:1050",
		"N5LCWU:vZNbygbj4z@https://188.130.188.148:1050",
		"N5LCWU:vZNbygbj4z@https://45.142.253.68:1050",
		"N5LCWU:vZNbygbj4z@https://45.134.252.187:1050",
		"N5LCWU:vZNbygbj4z@https://45.139.125.136:1050",
		"N5LCWU:vZNbygbj4z@https://188.130.221.225:1050",
	}

	proxyUrl, err := url.Parse(proxyPool[0])
	if err != nil {
		fmt.Println("Bad proxy URL", err)
		return
	}

	//client := &http.Client{}

	client := &http.Client{Transport: &http.Transport{Proxy: http.ProxyURL(proxyUrl)}}

	captchaUrl := "https://2captcha.com/in.php"
	////send post
	//jsonSend, err := json.Marshal(b)
	req, err := http.NewRequest("POST", captchaUrl, bytes.NewBuffer([]byte(b)))
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
/*
func searchGoogle (it hubRequest) {
	time.Sleep(3 * time.Second)
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
	//for iv, value := range it.SourceList {
		
		var outSrcs outPutSources
		var outLks outPutLinks

		outSrcs.Source = it.SourceList
		//println(outSrcs.Source)
		//Do something with index and value
		url := "https://www.google.com/search?q="+it.ItemId+"+site:"+it.SourceList+"+after:"+it.PeriodStart.GoogleParameterDate()+"+before:"+it.PeriodEnd.GoogleParameterDate()+"&num=100"
		//urlN := it.ItemId+"+site:"+value+"+after:"+it.PeriodStart.GoogleParameterDate()+"+before:"+it.PeriodEnd.GoogleParameterDate()+"&num=100"
		//fmt.Println(url)
		req, err := http.NewRequest("GET", url, nil)
	    if err != nil {
	    	log.Fatal(err)
	    }

		//req.Header = http.Header {
		//	"User-Agent": {"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/101.0.4951.54 Safari/537.36"},
		//	"Accept": {"application/json"},
		//	"Accept-Language": {"ru-RU,ru;q=0.8,en-US;q=0.5,en;q=0.3"},
		//}
		
		
		// cookie := http.Cookie {
		// 	"domain": ".google.com",
		// 	"expirationDate": 1740167135.467891,
		// 	"hostOnly": false,
		// 	"httpOnly": true,
		// 	"name": "NID",
		// 	"path": "/",
		// 	"sameSite": "no_restriction",
		// 	"secure": true,
		// 	"session": false,
		// 	"storeId": null,
		// 	"value": "516=q2GOm6QjMcr4FHz9TZndMcepIBpO_ras7aN_yji4lz3YCKywfKY143PjNcxcnzuU939kIkYnsZTIrDCWFQsFBjepT0PYv282-4JOxacsonT5ji79rVs1v3GPugC0cQxFvxs0QFAfxiRsLwBDioBOVNbyfUfeHXaXnEnAW4USN-YR2H4bS32LjG651pA0BcllRqVgzOx62_DgV_hF31WNkaY"
		// }
		// req.AddCookie(&cookie)
		
		
		//time.Sleep(1 * time.Second)
 		res, err := client.Do(req)
		if err != nil {
			fmt.Println("CAPTCHA2")
			log.Fatal(err)
		}
 		//defer res.Body.Close()
		
 		doc, err := goquery.NewDocumentFromReader(res.Body)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(doc)
		fmt.Println(doc.Find("div.g-recaptcha").Text())
		fmt.Println(doc.Find("div.g-recaptcha").Length())
		//////////////////////////
		if doc.Find("div.g-recaptcha").Length() > 0 {
			fmt.Println("CAPTCHA!")

			//htmlP, _ := doc.Find("div.g-recaptcha").Html()
			dataSiteKey, _ := doc.Find("captcha-widget").First().Attr("data-sitekey")
			dataS, _ := doc.Find("captcha-widget").First().Attr("data-s")
			urlS := url
			divInfo, _ := doc.Find("div#infoDiv").First().Attr("style")
			fmt.Println("dataSiteKey -- " + dataSiteKey)
			fmt.Println("dataS -- " + dataS)
			fmt.Println("divInfo -- " + divInfo)
			fmt.Println("urlS -- " + urlS)
			
			//testS := "feg"
			//or doc.Find("div#recaptcha") {
			//var docC *goquery.Document
			docC := makeCaptchaAPICall(dataSiteKey, dataS, urlS)
			/////////captcha//////////////
			
			//makeSerpDogCall(urlN)
		
			///////////////////////

			outSrcs.Total = 0
			//c := 0
			docC.Find("div.g").Each(func(i int, result *goquery.Selection) {
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

		//////////////////////////
		} else {
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
		}

		///////////////////////////
		
		res.Body.Close()
	//}
	
	//writeTaskResults(outResult)
}
*/

/*
func makeCaptchaAPICall (dataSiteKey string, dataS string, urlS string) (*goquery.Document) {
	//url *string, dataSite *string, dataS *string
	//data-s="NSVlTJncIGGdc6h70uaIiIS81ENv6mBwFAtqqos1cdUfUaJhArkqdkeG9QBs6sj8D24OlGdDaehOxhLunQQ5qCjU4C-IXTzUoij1pz_1VSuYXF3AIILFfwwwzyuTc8vHPiChjqIdZOvE1r-_ft1I7R_vZHKWr0g6hFRMannHa5Nv2oM8FmUhxNY7APPXFFUZ0bwCJrzWf3IXbEMrngO2MtpbO-tK8EK4uoNabOzVph8D1WYYx1V3804XO9yuSjv7v8l19gBNGiUykMFig1BoAztdkH-LqNA"
	//"___grecaptcha_cfg.clients['0']['L']['L']['callback']"
	clientC := api2captcha.NewClient("004f1af9da5d225d5ca67f84f4f71ce7")
	//clientC.Callback = "___grecaptcha_cfg.clients['0']['L']['L']['callback']"
	// clientC.DefaultTimeout = 120
	// clientC.RecaptchaTimeout = 600
	// clientC.PollingInterval = 100
	urlS = "https://www.google.com/sorry/index?continue=" + urlS
	captcha := api2captcha.ReCaptcha {
		SiteKey: dataSiteKey,
		Url: urlS,
		DataS: dataS,
		Invisible: true,
		Action: "verify",
	}

	code, smth, err := clientC.Solve(captcha.ToRequest())
	if err != nil {
		log.Fatal(err);
	}
	fmt.Println("code "+code)
	fmt.Println("smth "+smth)
	
	clientS := &http.Client{}

	urlC := urlS + "&g-recaptcha-response=" + code
		reqC, err := http.NewRequest("GET", urlC, nil)
	    if err != nil {
	    	log.Fatal(err)
	    }
		resC, err := clientS.Do(reqC)
		if err != nil {
			fmt.Println("CAPTCHA2")
			log.Fatal(err)
		}
		docC, err := goquery.NewDocumentFromReader(resC.Body)
		if err != nil {
			log.Fatal(err)
		}
		reqC.Body.Close()
		return docC
}
*/

/*
func makeSerpDogCall (url string) {
	urlN := "https://api.serpdog.io/search?api_key=66a3cfe3532d58cbb444cfb1&q=" + url

	serpClient := &http.Client{}
	req, err := http.NewRequest("GET", urlN, nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	req.Header.Set("Content-Type", "application/json")

	res, err := serpClient.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	//defer res.Body.Close()
	fmt.Println(res.Body)
	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("PARSED")
	fmt.Println(string(body))
	
	res.Body.Close()
}
*/

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
