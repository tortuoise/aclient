package aclient

import (
        "testing"
        _"time"

)

var (

	nsef = "https://nseindia.com/live_market/dynaContent/live_watch/get_quote/ajaxFOGetQuoteJSON.jsp?underlying=NIFTY&instrument=FUTIDX&type=-&strike=-&expiry="
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
        err = getter.Unmarshal(&OptionData{})
        if err != nil {
                t.Errorf("Error: %v ", err)
        }
        if getter.Ubs != nil {
                t.Errorf("Bytes: %v", getter.Display())
        }

}

/*func TestHttpMultiGet(t *testing.T) {

        nseLive1 := []byte(nsef)
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
        //err = getter.MultiUnmarshal(&OptionData{})
        //if err != nil {
        //        t.Errorf("Error: %v ", err)
        //}
        if getter.Ubs != nil {
                t.Errorf("Bytes: %v", getter.Display())
        }

}*/
