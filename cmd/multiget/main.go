package main

import (
    "bytes"
    "compress/gzip"
    "errors"
	"fmt"
    "net/http"
    "io/ioutil"
    "encoding/json"
    "sort"
    "strconv"
    "sync"
    "time"

    "github.com/tortuoise/aclient"
)

var (
        err error
	nsef = "https://nseindia.com/live_market/dynaContent/live_watch/get_quote/ajaxFOGetQuoteJSON.jsp?underlying=NIFTY&instrument=FUTIDX&type=-&strike=-&expiry="
	nses = "https://nseindia.com/live_market/dynaContent/live_watch/get_quote/ajaxFOGetQuoteJSON.jsp?underlying="
	nses1 = "&instrument=FUTSTK&expiry="
	nses2 = "&type=SELECT&strike=SELECT"
	getter *aclient.HttpMultiGetter
)
func main() {

        nseLive := []byte(nses)
        raw, err := ioutil.ReadFile("nse_prtfl")
        if err != nil {
                fmt.Println(err)
        }
        sngls := bytes.Split(raw, []byte("\n"))
        sngls = sngls[:len(sngls)-1]
        _, x1 := aclient.X1()
        urls := make([]string, 0, len(sngls))
        doneChan := make(chan bool, 1)
        errChan := make(chan error, 1)
        respChan := make(chan []byte, 1)
        for _, sngl := range sngls {
                 url := append(append(append(append(nseLive, sngl...), nses1...), x1...), nses2...)
                urls = append(urls, string(url))
                go func(url string) {
                        req, err := http.NewRequest("GET", url, nil)
                        if err != nil {
                                errChan <- err
                                return
                        }
                        aclient.SetHeaders(req)
                        resp, err := aclient.Client.Do(req)
                        if err != nil {
                                errChan <- errors.New("GET"+err.Error())
                                return
                        }else {
                                if resp != nil {
                                        defer resp.Body.Close()
                                }
                                cl := resp.Header.Get(aclient.ContentLengthHeader)
                                icl, err := strconv.Atoi(cl)
                                if err != nil {
                                        //errChan <- err
                                        errChan <- errors.New("Strconv"+err.Error())
                                        return
                                }
                                ubs := make([]byte, icl*3)
                                ct := resp.Header.Get(aclient.ContentTypeHeader)
                                switch ct {
                                        case "gzip":
                                                 gzr, err := gzip.NewReader(resp.Body)
                                                if err != nil {
                                                         //errChan <- err
                                                        errChan <- errors.New("gzip"+err.Error())
                                                        return
                                                }
                                                defer gzr.Close()
                                                nbs, err := gzr.Read(ubs)
                                                ubs = ubs[:nbs]
                                                respChan <- ubs
                                                return
                                        default:
                                                gzr, err := gzip.NewReader(resp.Body)
                                                if err != nil {
                                                        //errChan <- err
                                                        errChan <- errors.New("default "+err.Error())
                                                         return
                                                }
                                                defer gzr.Close()
                                                nbs, err := gzr.Read(ubs)
                                                ubs = ubs[:nbs]
                                                respChan <- ubs
                                                 return

                                }
                                //respChan <- []byte(ct)
                                //return
                        }
                        doneChan<- true
                }(string(url))
        }
        strngs := make(aclient.Datas, 0)
        var mtx sync.Mutex
        var wg sync.WaitGroup
        for n:= 0; n < len(urls); {
                 select {
                         case <-doneChan:
                                n++
                                fmt.Println("Done: ", n)
                        case err = <-errChan:
                                n++
                                fmt.Println("Error: ", err)
                        case bs := <-respChan:
                                n++
                                wg.Add(1)
                                go func(bs []byte) {
                                        defer wg.Done()
                                        od := &aclient.OptionData{}
                                        err := json.Unmarshal(bs, od)
                                        if err != nil {
                                                fmt.Println(err)
                                                errChan <- err
                                                return
                                        }
                                        mtx.Lock()
                                        strngs = append(strngs, *od)
                                        mtx.Unlock()
                                }(bs)
                        case <-time.After(2000 * time.Millisecond):
                                fmt.Println("Timeout: ", n, " timed out")
                                n++
                }
                time.Sleep(100*time.Millisecond)
        }
        wg.Wait()
        close(doneChan)
        close(errChan)
        close(respChan)
        sort.Sort(strngs)
        for _,strng := range strngs {
            fmt.Println(strng.String())
        }

        // Output: Varies
        // Varies
        // And varies some more

}

