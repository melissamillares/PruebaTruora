package main

import (
	"fmt"
	"time"	
	"log"	
	"net"		
	"net/url"
	"net/http"
	"io/ioutil"
	"strings"
	"github.com/likexian/whois-go"		
)

// verify if the string is a URL
func isURL(urlString string) (bool, error) {
	_, err := url.ParseRequestURI(urlString)

	if err != nil {			
		//panic(err)							
		return false, err
	} else {	
		return true, nil
	}		
}

//
func hostName(urlString string) string {
	u, err := url.Parse(urlString) 
	var hoststring string

	if err != nil {
		panic(err)
	} else {
		hoststring = u.Hostname()
	}	
	return hoststring
}

//
func getIP(urlString string) []net.IP {
	ips, err := net.LookupIP(urlString)

	if err != nil {			
		//panic(err)					
		log.Fatal(err)
	} 					
	
	return ips		
}

// get the information s from the query whois
// (e.g. s="Country: " returns the country associated with the IP)
func getInfoWhoIs(s string, ips []net.IP) string {
	for _, ip := range ips {
		result, err := whois.Whois(ip.String())		
		if err != nil {						
			log.Fatal(err)
		}
		// split the result from whois by \n
		splitResult := strings.Split(result, "\n")				
		// search in splitresult the string s
		for _, val := range splitResult {
			if strings.Contains(val, s) {				
				info := strings.Trim(val, s)				
				return info
			}
		}		
	}
	return ""	
}

// returns an array with the SSL grade of the host servers
// length: the length from the IPs array (associated with the host servers)
func getSSLGrade(host string, length int) []string {
	u := fmt.Sprintf("https://api.ssllabs.com/api/v3/analyze?host=%s", host)
	resp, err := http.Get(u)
	sslgrades := make([]string, length) // array with the length from the IPs array

	if err == nil {				
		defer resp.Body.Close()

		body, e := ioutil.ReadAll(resp.Body)

		if e == nil {						
			splitResult := strings.Split(string(body), ",")		
			for _, val := range splitResult {		
				if strings.Contains(val, "grade") {
					ssl := strings.Trim(val, "\"grade\":")
					for i := 0; i < length; i++ {
						sslgrades[i] = ssl
					}							
					return sslgrades
				}
			}
		}
	}

	return sslgrades
}

//
func isServerDown(urlString string) bool {
	_, err := http.Get(urlString)

	if err != nil {
		return true
	}
	
	return false
}

func compareOneHourAgo(serv1, serv2 Server) bool {
	resp := false
			
	y1, m1, d1 := serv1.Created.Date()
	serv1Date := string(y1) + string(m1) + string(d1)
	y2, m2, d2 := serv2.Created.Date()
	serv2Date := string(y2) + string(m2) + string(d2)

	if serv2.Created.Hour() - serv1.Created.Hour() >= 1 || serv1Date != serv2Date {
		// if the difference in update is more than 1 hour			
		resp = true								
	} else if (serv1.Updated != time.Time{}) && (serv2.Updated != time.Time{}) { // if the updated time is different from null time
		if serv2.Updated.Hour() - serv1.Updated.Hour() >= 1 {
			resp = true
		}
	}
	return resp
}

// 
func equalServers(s1, s2 Server) bool {
	res := false
	//if s1.ID == s2.ID {
		if s1.Address == s2.Address && s1.SSL_grade == s2.SSL_grade && s1.Country == s2.Country && s1.Owner == s2.Owner {
			res = true
		}
	//}
	return res
}
