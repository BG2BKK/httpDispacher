package http_dns

import (
	"fmt"
	"MyError"
	"utils"
	"github.com/miekg/dns"
)

const CNAME_CHAIN_LENGTH = 10


func GetARecord(d string, srcIP string) (bool, []dns.RR, *MyError.MyError) {
	var bigloopflag bool = false // big loop flag
	var c = 0                    //big loop count
	//fmt.Println(utils.GetDebugLine(), "^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^")

	//Can't loop for CNAME chain than bigger than 10
	for dst := d; (bigloopflag == false) && (c < CNAME_CHAIN_LENGTH); c++ {
		fmt.Println(utils.GetDebugLine(), "GetARecord : ", dst, " srcIP: ", srcIP)

		soa, ns, e := QuerySOA(dst)

		fmt.Println(utils.GetDebugLine(), "GetARecord: GetSOARecord return : ", soa, " error: ", e)

		if e != nil {
			fmt.Println(utils.GetDebugLine(), "GetSOARecord error: %s", e.Error())
		}

		fmt.Println(utils.GetDebugLine(), "Got dst: ", dst, " srcIP: ", srcIP, " soa.NS: ", ns)
		var ns_a []string
		//todo: may the soa.NS be nil ?
		for _, x := range ns {
			fmt.Println(x)
			ns_a = append(ns_a, x.Ns)
		}
		ok, rr_i, rtype, ee := GetAFromDNSBackend(dst, srcIP, ns_a)
		if ok && rtype == dns.TypeA {
			return true, rr_i, nil
		} else if ok && rtype == dns.TypeCNAME {
			dst = rr_i[0].(*dns.CNAME).Target
			continue
		} else if !ok && rr_i == nil && ee != nil && ee.ErrorNo == MyError.ERROR_NORESULT {
			continue
		} else {
			return false, nil, MyError.NewError(MyError.ERROR_UNKNOWN, "Unknown error")
		}
	}
	return false, nil, MyError.NewError(MyError.ERROR_UNKNOWN, "Unknown error")
}

func GetAFromDNSBackend(	dst, srcIP string, ns_a []string ) (bool, []dns.RR, uint16, *MyError.MyError) {

	var reE *MyError.MyError = nil
	var rtype uint16
	rr, edns_h, edns, e := QueryA(dst, srcIP, ns_a, "53")
	if e == nil && rr != nil {
		var rr_i []dns.RR
		if a, ok := ParseA(rr, dst); ok {
			//rr is A record
			fmt.Println(utils.GetDebugLine(), "GetAFromDNSBackend : typeA record: ", a, " dns.TypeA: ", ok)
			//utils.ServerLogger.Debug("GetAFromDNSBackend : typeA record: ", a, " dns.TypeA: ", ok)
			for _, i := range a {
				rr_i = append(rr_i, dns.RR(i))
			}
			//if A ,need parse edns client subnet
			//			return true,rr_i,nil
			rtype = dns.TypeA
		} else if b, ok := ParseCNAME(rr, dst); ok {
			//rr is CNAME record
			fmt.Println(utils.GetDebugLine(), "GetAFromDNSBackend: typeCNAME record: ", b, " dns.TypeCNAME: ", ok)
			//utils.ServerLogger.Debug("GetAFromDNSBackend: typeCNAME record: ", b, " dns.TypeCNAME: ", ok)
			dst = b[0].Target
			for _, i := range b {
				rr_i = append(rr_i, dns.RR(i))
			}
			rtype = dns.TypeCNAME
			reE = MyError.NewError(MyError.ERROR_NOTVALID,
				"Got CNAME result for dst : "+dst+" with srcIP : "+srcIP)
			//if CNAME need parse edns client subnet
		} else {
			//error return and retry
			fmt.Println(utils.GetDebugLine(), "GetAFromDNSBackend: ", rr)
			//utils.ServerLogger.Debug("GetAFromDNSBackend: ", rr)
			return false, nil, uint16(0), MyError.NewError(MyError.ERROR_NORESULT,
				"Got error result, need retry for dst : "+dst+" with srcIP : "+srcIP)
		}

		// Parse edns client subnet
		fmt.Println(utils.GetDebugLine(), "GetAFromDNSBackend: ", " edns_h: ", edns_h, " edns: ", edns)
		//utils.ServerLogger.Debug("GetAFromDNSBackend: ", " edns_h: ", edns_h, " edns: ", edns)

		return true, rr_i, rtype, reE
	}
	return false, nil, rtype, MyError.NewError(MyError.ERROR_UNKNOWN, utils.GetDebugLine()+"Unknown error")
}
