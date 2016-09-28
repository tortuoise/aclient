package aclient

import (
        "bytes"
        "compress/gzip"
        "encoding/json"
        "fmt"
        "net/http"
	_ "reflect"
	"testing"
	_ "time"
        "io/ioutil"
        "strconv"
        "time"
)

var (
        err error
	nsef   = "https://nseindia.com/live_market/dynaContent/live_watch/get_quote/ajaxFOGetQuoteJSON.jsp?underlying=NIFTY&instrument=FUTIDX&type=-&strike=-&expiry="
	nses = "https://nseindia.com/live_market/dynaContent/live_watch/get_quote/ajaxFOGetQuoteJSON.jsp?underlying="
	nses1 = "&instrument=FUTSTK&type=-&strike=-&expiry="
	getter *HttpMultiGetter
)

func TestHttpGet(t *testing.T) {

        nseLive1 := []byte(nsef)
        var xprs [][]byte
        _,x1 := x1()
        _,x2 := x2()
        _,x3 := x3()
        xprs = append(xprs,[]byte(x1))
        xprs = append(xprs,[]byte(x2))
        xprs = append(xprs,[]byte(x3))
        xprs = xprs[:len(xprs)]
        url := string( append(nseLive1, xprs[0]...))
        getter,err := NewHttpGetter(url)
        if err != nil {
                t.Errorf("Error: %v", err)
        }
        getter.Get()
        err = getter.Unmarshal(&OptionData{})
        if err != nil {
                t.Errorf("Error: %v", err)
        }
        if getter.Ubs != nil {
                t.Errorf("Bytes: %v", getter.Display())
        }

}

func TestHttpGoGet(t *testing.T) {

        nseLive1 := []byte(nsef)
        var xprs [][]byte
        _,x1 := x1()
        _,x2 := x2()
        _,x3 := x3()
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
        err = getter.Unmarshal(&OptionData{})
        if err != nil {
                t.Errorf("Error: %v ", err)
        }
        if getter.Ubs != nil {
                t.Errorf("Bytes: %v", getter.Display())
        }

}

//ExampleHttMultiGet demonstrates how to make multiple http requests using goroutines using a single client/transport.
func ExampleHttpMultiGet() {

        nseLive := []byte(nses)
        raw, err := ioutil.ReadFile("nse_prtfl")
        if err != nil {
                fmt.Println(err)
        }
        sngls := bytes.Split(raw, []byte("\n"))
        sngls = sngls[:len(sngls)-1]
        _, x1 := x1()
        urls := make([]string, 0, len(sngls))
        doneChan := make(chan bool, 1)
        errChan := make(chan error, 1)
        respChan := make(chan []byte, 1)
        for _, sngl := range sngls {
                url := append(append(append(nseLive, sngl...), nses1...), x1...)
                urls = append(urls, string(url))
                go func(url string) {
                        req, err := http.NewRequest("GET", url, nil)
                        if err != nil {
                                errChan <- err
                                return
                        }
                        setHeaders(req)
                        resp, err := client.Do(req)
                        if err != nil {
                                errChan <- err
                                return
                        }else {
                                if resp != nil {
                                        defer resp.Body.Close()
                                }
                                cl := resp.Header.Get(contentLengthHeader)
                                icl, err := strconv.Atoi(cl)
                                if err != nil {
                                        errChan <- err
                                        return
                                }
                                ubs := make([]byte, icl*3)
                                ct := resp.Header.Get(contentTypeHeader)
                                switch ct {
                                        case "gzip":
                                                gzr, err := gzip.NewReader(resp.Body)
                                                if err != nil {
                                                        errChan <- err
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
                                                        errChan <- err
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
        for n:= 0; n < len(urls); {
                select {
                        case <-doneChan:
                                n++
                                fmt.Println("Done: ", n)
                        case err = <-errChan:
                                n++
                                fmt.Println("Error: ", err)
                        case bs := <-respChan:
                                //fmt.Print("Success: ")
                                n++
                                go func(bs []byte) {
                                        od := &OptionData{}
                                        err := json.Unmarshal(bs, od)
                                        if err != nil {
                                                errChan <- err
                                                return
                                        }
                                        fmt.Println(od.String())
                                }(bs)
                        case <-time.After(500 * time.Millisecond):
                                fmt.Println("Timeout: ", n, " timed out")
                                n++
                }
        }
        close(doneChan)

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
