package main

import (
"runtime"

"github.com/pkg/profile"
"github.com/miekg/dns"
"config"
"http_dns"
"utils"	
"fmt"
"net"
)

// param: lower case + Upper Case ,No _ spliter
// Struct unit: Upper Case
// Func: golang style

func main() {
	config.InitConfig()

	if config.EnableProfile {
		defer profile.Start(profile.CPUProfile).Stop()
	}

	runtime.GOMAXPROCS(runtime.NumCPU() * 3)
	utils.InitLogger()	

	query_domain := "www.baidu.com"
	//query_domain := "api.weibo.cn"
	//query_domain := "www.taobao.com"
	//query_domain := "ww1.sinaimg.cn"
	//srcIP := "45.78.57.71"
	//srcIP := "202.118.10.23"
	srcIP := "61.135.152.203"
	//srcIP := "10.209.77.225"
	//srcIP := "219.142.118.233"

	if _, ok := dns.IsDomainName(query_domain); !ok {
		fmt.Printf("error domain name : %s\n ", query_domain)
		return
	}

	if srcIP == "" {
		fmt.Printf("error client ip : %s\n ", srcIP)
		return
	}
	fmt.Printf("src ip: %s query_domain: %s \n", string(srcIP), query_domain)
	if x := net.ParseIP(srcIP); x == nil {
		fmt.Printf("src ip : %s is not correct\n", srcIP)
		return
	}

	ok, re, e := http_dns.GetARecord(query_domain, srcIP)
	if ok {
		for _, ree := range re {
			//fmt.Printf("query result: %s\n", ree)
			if a, ok := ree.(*dns.A); ok {

				fmt.Printf("query result: %s\n ", a.A.String())
			} else {
				fmt.Printf("query result: %s\n ", ree.String())
			}
		}
	} else if e != nil {
		fmt.Printf("query domain: %s src_ip: %s  %s\n", query_domain, srcIP, e.Error())
	} else {
		fmt.Printf("query domain: %s src_ip: %s fail unkown error!\n", query_domain, srcIP)
	}

}
