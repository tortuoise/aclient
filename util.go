package aclient

import (
        "os"
        "net"
        "fmt"
)

func FatalCheck(err error) {
        if err != nil {
                fmt.Println("%v", err)
                os.Exit(1)
        }
}

func FilterIPV4(ips []net.IP) []string {
        ret := make([]string, 0)
        for _, ip := range ips {
                if ip.To4() != nil {
                        ret = append(ret, ip.String())
                }
        }
        return ret
}
