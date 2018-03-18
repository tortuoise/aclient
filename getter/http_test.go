package getter

import (
    "bytes"
    "compress/gzip"
    "encoding/json"
    "errors"
    "fmt"
    "net/http"
	_ "reflect"
	"testing"
	_ "time"
    _"io/ioutil"
    "strconv"
    "sort"
    "sync"
    "time"

    "github.com/tortuoise/aclient/nse"
)

var (
    err error
	nsef = "https://nseindia.com/live_market/dynaContent/live_watch/get_quote/ajaxFOGetQuoteJSON.jsp?underlying=NIFTY&instrument=FUTIDX&type=-&strike=-&expiry="
	nsef1 = "https://nseindia.com/live_market/dynaContent/live_watch/get_quote/ajaxFOGetQuoteJSON.jsp?underlying=NIFTY&instrument=FUTIDX&expiry="
	nses = "https://nseindia.com/live_market/dynaContent/live_watch/get_quote/ajaxFOGetQuoteJSON.jsp?underlying="
	nses1 = "&instrument=FUTSTK&expiry="
	nses2 = "&type=SELECT&strike=SELECT"
	getter *HttpMultiGetter
)

func TestHttpGet(t *testing.T) {

        nseLive1 := []byte(nsef1)
        nses2b := []byte(nses2)
        var xprs [][]byte
        //_,x1 := nse.X1()
        x1 := "28MAR2018"
        _,x2 := nse.X2()
        _,x3 := nse.X3()
        xprs = append(xprs,[]byte(x1))
        xprs = append(xprs,[]byte(x2))
        xprs = append(xprs,[]byte(x3))
        xprs = xprs[:len(xprs)]
        url := string( append(append(nseLive1, xprs[0]...), nses2b...))
        t.Errorf("%v", url)
        getter,err := NewHttpGetter(url)
        if err != nil {
                t.Errorf("Error: %v", err)
        }
        err = getter.Get()
        if err != nil {
                t.Errorf("Error: %v", err)
        }
        err = getter.Unmarshal(&nse.OptionData{})
        if err != nil {
                t.Errorf("Error: %v", err)
                t.Errorf("Error: %v", string(getter.Ubs))
        } else {
        if getter.Ubs != nil {
                t.Errorf("Bytes: %v", getter.Display())
        }
        }

}

func TestHttpGoGet(t *testing.T) {

        nseLive1 := []byte(nsef)
        var xprs [][]byte
        _,x1 := nse.X1()
        _,x2 := nse.X2()
        _,x3 := nse.X3()
        xprs = append(xprs,[]byte(x1))
        xprs = append(xprs,[]byte(x2))
        xprs = append(xprs,[]byte(x3))
        xprs = xprs[:len(xprs)]
        url := string( append(nseLive1, xprs[0]...))
        getter,err := NewHttpGetter(url)
        if err != nil {
                t.Errorf("Error: %v", err)
        }
        doneChan := make(chan bool,1)
        errChan := make(chan error,1)
        getter.MultiGet(doneChan, errChan)
        <-doneChan
        //t.Errorf("%v", <-errChan)
        err = getter.Unmarshal(&nse.OptionData{})
        if err != nil {
                t.Errorf("Error: %v ", err)
        }
        if getter.Ubs != nil {
                t.Errorf("Bytes: %v", getter.Display())
        }

}

//ExampleHttMultiGet demonstrates how to make multiple http requests using goroutines using a single client/transport. go test -run HttpMultiGet github.com/tortuoise/aclient
func ExampleHttpMultiGet() {

        nseLive := []byte(nses)
        //raw, err := ioutil.ReadFile("nse_prtfl")
        raw, err := nse.Asset("static/nse_prtfl")
        if err != nil {
                fmt.Println(err)
        }
        sngls := bytes.Split(raw, []byte("\n"))
        sngls = sngls[:len(sngls)-1]
        //_, x1 := nse.X1()
        x1 := "28MAR2018"
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
                        nse.SetHeaders(req)
                        resp, err := Client.Do(req)
                        if err != nil {
                                errChan <- errors.New("GET"+err.Error())
                                return
                        }else {
                                if resp != nil {
                                        defer resp.Body.Close()
                                }
                                cl := resp.Header.Get(ContentLengthHeader)
                                icl, err := strconv.Atoi(cl)
                                if err != nil {
                                        //errChan <- err
                                        errChan <- errors.New("Strconv"+err.Error())
                                        return
                                }
                                ubs := make([]byte, icl*3)
                                ct := resp.Header.Get(ContentTypeHeader)
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
        strngs := make(nse.Datas, 0)
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
                                        od := &nse.OptionData{}
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
        close(doneChan)
        close(errChan)
        close(respChan)
        wg.Wait()
        sort.Sort(strngs)
        for _,strng := range strngs {
            fmt.Println(strng.String())
        }

        // Output: Varies
        // Varies
        // And varies some more

}

/*
func TestPersistence(t *testing.T) {
        t.Errorf("Still there")

}*/
	   /*nseLive1 := []byte(nsef)
	   var xprs [][]byte
	   _,x1 := x1()
	   _,x2 := x2()
	   _,x3 := x3()
	   xprs = append(xprs,[]byte(x1))
	   xprs = append(xprs,[]byte(x2))
	   xprs = append(xprs,[]byte(x3))
	   xprs = xprs[:len(xprs)]
	   urls := make([]string, len(xprs))
	   for i, xpr := range xprs {
	           url := string( append(nseLive1, xpr...))
	           urls[i] = url
	   }
	   getter,err := NewHttpMultiGetter(urls)
	   if err != nil {
	           t.Errorf("Error: %v", err)
	   }
	   doneChan := make(chan bool,1)
	   errChan := make(chan error,1)
	   getter.MultiGet(doneChan, errChan)
	   <-doneChan
	   err = getter.MultiUnmarshal(&OptionData{})
	   //err = getter.MultiUnmarshal(&Top10{})
	   if err != nil {
	           t.Errorf("Error: %v ", err)
	   } else if getter.Ubs != nil {
	           t.Errorf("%v", getter.Display())
	   }*/
