package aclient

import (
        "fmt"
        "net/http"
	"log"
	"strconv"
	"strings"
	"time"
)

type OptionData struct {
	Valid          string `json:"valid,omitempty"` // field names must begin with upper case for (un)marshaling to work
	IsinCode       string `json:"isinCode,omitempty"`
	LastUpdateTime string `json:"lastUpdateTime,omitempty"`
	OcLink         string `json:"ocLink,omitempty"`
	TradedDate     string `json:"tradedDate,omitempty"`
	Book           []data `json:"data,omitempty"`
	CompanyName    string `json:"companyName,omitempty"`
	EqLink         string `json:"eqLink,omitempty"`
}

type data struct {
	Underlying           string `json:"underlying,omitempty"`
	UnderlyingValue      string `json:"underlyingValue,omitempty"`
	Ltp                  string `json:"ltp,omitempty"`
	Change               string `json:"change,omitempty"`
	AnnualisedVolatility string `json:"annualisedVolatility,omitempty"`
	ImpliedVolatility    string `json:"impliedVolatility,omitempty"`
	DailyVolatility      string `json:"dailyVolatility,omitempty"`
	OptionType           string `json:"optionType,omitempty"`
	PrevClose            string `json:"prevClose,omitempty"`
	PChange              string `json:"pChange,omitempty"`
	LastPrice            string `json:"lastPrice,omitempty"`
	HighPrice            string `json:"highPrice,omitempty"`
	LowPrice             string `json:"lowPrice,omitempty"`
	StrikePrice          string `json:"strikePrice,omitempty"`

	BestSell          string `json:"bestSell,omitempty"`
	BestBuy           string `json:"bestBuy,omitempty"`
	BuyPrice1         string `json:"buyPrice1,omitempty"`
	BuyPrice2         string `json:"buyPrice2,omitempty"`
	BuyPrice3         string `json:"buyPrice3,omitempty"`
	BuyPrice4         string `json:"buyPrice4,omitempty"`
	BuyPrice5         string `json:"buyPrice5,omitempty"`
	BuyQuantity1      string `json:"buyQuantity1,omitempty"`
	BuyQuantity2      string `json:"buyQuantity2,omitempty"`
	BuyQuantity3      string `json:"buyQuantity3,omitempty"`
	BuyQuantity4      string `json:"buyQuantity4,omitempty"`
	BuyQuantity5      string `json:"buyQuantity5,omitempty"`
	TotalBuyQuantity  string `json:"totalBuyQuantity,omitempty"`
	SellPrice1        string `json:"sellPrice1,omitempty"`
	SellPrice2        string `json:"sellPrice2,omitempty"`
	SellPrice3        string `json:"sellPrice3,omitempty"`
	SellPrice4        string `json:"sellPrice4,omitempty"`
	SellPrice5        string `json:"sellPrice5,omitempty"`
	SellQuantity1     string `json:"sellQuantity1,omitempty"`
	SellQuantity2     string `json:"sellQuantity2,omitempty"`
	SellQuantity3     string `json:"sellQuantity3,omitempty"`
	SellQuantity4     string `json:"sellQuantity4,omitempty"`
	SellQuantity5     string `json:"sellQuantity5,omitempty"`
	TotalSellQuantity string `json:"totalSellQuantity,omitempty"`
	Vwap              string `json:"vwap,omitempty"`

	MarketLot                string `json:"marketLot,omitempty"`
	NumberOfContractsTraded  string `json:"numberOfContractsTraded,omitempty"`
	TurnoverinRsLakhs        string `json:"turnoverinRsLakhs,omitempty"`
	MarketWidePositionLimits string `json:"marketWidePositionLimits,omitempty"`
	ClientWisePositionLimits string `json:"clientWisePositionLimits,omitempty"`
	OpenInterest             string `json:"openInterest,omitempty"`
	PchangeinOpenInterest    string `json:"pchangeinOpenInterest,omitempty"`
	SettlementPrice          string `json:"settlementPrice,omitempty"`
	InstrumentType           string `json:"instrumentType,omitempty"`
	ExpiryDate               string `json:"expiryDate,omitempty"`
	Strike                   string `json:"strike,omitempty"`
}

type Top10 struct {
	Time string  `json:"time"`
	Data []Stock `json:"data"`
}

type Stock struct {
	Symbol                   string `json:"symbol"`
	Series                   string `json:"series"`
	OpenPrice                string `json:"openPrice"`
	HighPrice                string `json:"highPrice"`
	LowPrice                 string `json:"lowPrice"`
	Ltp                      string `json:"ltp"`
	PreviousPrice            string `json:"previousPrice"`
	NetPrice                 string `json:"netPrice"`
	TradedQuantity           string `json:"tradedQuantity"`
	TurnoverInLakhs          string `json:"turnoverInLakhs"`
	LastCorpAnnouncement     string `json:"lastCorpAnnouncement"`
	LastCorpAnnouncementDate string `json:"lastCorpAnnouncementDate"`
}

func (od *OptionData) String() string{

	return fmt.Sprintf(" %10s %v %8s %10s %10s %8s %10s %10s ", od.Book[0].Underlying, od.Book[0].ExpiryDate, od.Book[0].BuyQuantity1, od.Book[0].BuyPrice1, od.Book[0].SellPrice1, od.Book[0].SellQuantity1, od.Book[0].LowPrice, od.Book[0].HighPrice)

}

//Datas implements sort.Interface
type Datas []OptionData

func (ods Datas) Len() int {
    return len(ods)
}

func (ods Datas) Swap(i,j int) {
        ods[i], ods[j] = ods[j], ods[i]
}

func (ods Datas) Less(i, j int) bool {
    return ods[i].Book[0].Underlying < ods[j].Book[0].Underlying
}

func lastThurs(month time.Month) int {
	this := time.Date(time.Now().Year(), month, dayspMonth(month), 12, 0, 0, 0, time.UTC)
	if isThurs(this) {
		return this.Day()
	}
	for day := dayspMonth(month); day > 0; day-- {
		if this.AddDate(0, 0, -1).Weekday() == time.Thursday {
			return day - 1
		}
		this = this.AddDate(0, 0, -1)
	}
	return 0
}

func isThurs(now time.Time) bool {
	if now.Weekday() == time.Thursday {
		return true
	}
	return false
}

func isLastThurs(now time.Time) bool {
	nextThurs := now.AddDate(0, 0, 7)
	if nextThurs.Month() != now.Month() {
		return true
	}
	return false
}

func x1() (time.Time, string) {
	now := time.Now()
	lt := lastThurs(now.Month())
	mth := make([]byte, 3)
	if lt >= now.Day() {
		_, err := strings.NewReader(now.Month().String()).Read(mth)
		if err != nil {
			log.Fatal(err)
		}
		return time.Date(now.Year(), now.Month(), lt, 12, 0, 0, 0, time.UTC), strconv.Itoa(lt) + strings.ToUpper(string(mth)) + strconv.Itoa(now.Year())
	}
	now = now.AddDate(0, 0, 7)
	lt = lastThurs(now.Month())
	_, err := strings.NewReader(now.Month().String()).Read(mth)
	if err != nil {
		log.Fatal(err)
	}
	return time.Date(now.Year(), now.Month(), lt, 12, 0, 0, 0, time.UTC), strconv.Itoa(lt) + strings.ToUpper(string(mth)) + strconv.Itoa(now.Year())
}

func x2() (time.Time, string) {
	x1, _ := x1()
	x2 := x1.AddDate(0, 0, 28)
	mth := make([]byte, 3)
	_, err := strings.NewReader(x2.Month().String()).Read(mth)
	if err != nil {
		log.Fatal(err)
	}
	if isLastThurs(x2) {
		return x2, strconv.Itoa(x2.Day()) + strings.ToUpper(string(mth)) + strconv.Itoa(x2.Year())
	}
	x2 = x2.AddDate(0, 0, 7)
	_, err = strings.NewReader(x2.Month().String()).Read(mth)
	if err != nil {
		log.Fatal(err)
	}
	return x2, strconv.Itoa(x2.Day()) + strings.ToUpper(string(mth)) + strconv.Itoa(x2.Year())
}

func x3() (time.Time, string) {
	x2, _ := x2()
	x3 := x2.AddDate(0, 0, 28)
	mth := make([]byte, 3)
	_, err := strings.NewReader(x3.Month().String()).Read(mth)
	if err != nil {
		log.Fatal(err)
	}
	if isLastThurs(x3) {
		return x3, strconv.Itoa(x3.Day()) + strings.ToUpper(string(mth)) + strconv.Itoa(x3.Year())
	}
	x3 = x3.AddDate(0, 0, 7)
	_, err = strings.NewReader(x3.Month().String()).Read(mth)
	if err != nil {
		log.Fatal(err)
	}
	return x3, strconv.Itoa(x3.Day()) + strings.ToUpper(string(mth)) + strconv.Itoa(x3.Year())
}

func dayspMonth(month time.Month) int {
	switch month {
	case time.January, time.March, time.May, time.July, time.August, time.October, time.December:
		return 31
	case time.April, time.June, time.September, time.November:
		return 30
	case time.February:
		if isLeap(time.Now().Year()) {
			return 29
		}
		return 28
	}
	return 30
}

func isLeap(year int) bool {
	t := time.Date(year, time.February, 29, 12, 0, 0, 0, time.UTC)
	if t.Month() == time.February {
		return true
	}
	return false
}

func setHeaders(req *http.Request) {
	//req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 6.1; rv:39.0) Gecko/20100101 Firefox/39.0")
	//req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64; rv:31.0) Gecko/20100101 Firefox/31.0 Iceweasel/31.8.0i")
	//req.Header.Set("Host", "nseindia.com")
	//req.Header.Set("DNT", "1")
	//req.Header.Set("Connection", "keep-alive")
	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64; rv:31.0) Gecko/20100101 Firefox/58.0")
	req.Header.Set("Host", "www.nseindia.com")
	req.Header.Set("Accept", "*/*")
    req.Header.Set("X-Requested-With", "XMLHttpRequest")
    req.Header.Set("Referer", req.URL.String())
	//req.Header.Set("Cache-Control", "max-age=0")
	req.Header.Set("Accept-Language", "en-US,en;q=0.5")
	req.Header.Set("Accept-Encoding", "gzip, deflate")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*,q=0.8")
}

