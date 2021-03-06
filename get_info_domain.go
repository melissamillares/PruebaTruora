package main

import (		
	"net/http"
	"golang.org/x/net/html"
	"io"
	"io/ioutil"
	"strings"
	"sort"
)

// ssl: array with the SSL grade from the servers
func getLowerGrade(ssl []string) string {
	sort.Strings(ssl) // sorts the array in increasing order
	
	if len(ssl) == 1 { // if the array only has one element, returns that element
	    return ssl[0]
	}
	last_index := len(ssl) - 1
	
	for _, val := range ssl {
	    if val == "" { // if the array has one grade equal to "", returns ""
		return val
	    }
	}
	// compares if the last position is equal to the previous position concatenated 
	// with the symbol "+", returns the previous position
	if ssl[last_index] == ssl[last_index- 1]+"+" {
	    return ssl[last_index - 1]
	} 

	return ssl[last_index]
}

func getLogo(urlString string) string {	
	var resultLogo string	
	resp, err := http.Get(urlString)	

	if err != nil {		
		return ""
	} else {
		defer resp.Body.Close()

		body, _ := ioutil.ReadAll(resp.Body)	
		split := strings.Split(string(body), "<")

		for _, val := range split {		
			if strings.Contains(val, "rel=\"fluid-icon\"") {
				val = strings.Trim(val, ">")
				splitResult := strings.Split(val, " ")
				for _, v := range splitResult {					
					if strings.Contains(v, ".png") || strings.Contains(v, ".jpg") || strings.Contains(v, ".svg") || strings.Contains(v, ".ico") {
						logo := strings.Trim(v, "href=")											
						resultLogo = logo
					}
				}
			} else if strings.Contains(val, "rel=\"shortcut icon\"") {
				val = strings.Trim(val, ">")
				splitResult := strings.Split(val, " ")
				for _, v := range splitResult {					
					if strings.Contains(v, ".png") || strings.Contains(v, ".jpg") || strings.Contains(v, ".svg") || strings.Contains(v, ".ico") {						
						logo := strings.Trim(v, "href=")											
						resultLogo = logo				
					}				
				}			
			} else if strings.Contains(val, "rel=\"apple-touch-icon\"") {
				val = strings.Trim(val, ">")
				splitResult := strings.Split(val, " ")
				for _, v := range splitResult {					
					if strings.Contains(v, ".png") || strings.Contains(v, ".jpg") || strings.Contains(v, ".svg") || strings.Contains(v, ".ico") {						
						logo := strings.Trim(v, "href=")											
						resultLogo = logo				
					}				
				}			
			} else if strings.Contains(val, "rel=\"icon\"") {
				val = strings.Trim(val, ">")
				splitResult := strings.Split(val, " ")
				for _, v := range splitResult {					
					if strings.Contains(v, ".png") || strings.Contains(v, ".jpg") || strings.Contains(v, ".svg") || strings.Contains(v, ".ico") {						
						logo := strings.Trim(v, "href=")											
						resultLogo = logo				
					}				
				}			
			} else if strings.Contains(val, "rel=\"mask-icon\"") {
				val = strings.Trim(val, ">")
				splitResult := strings.Split(val, " ")
				for _, v := range splitResult {					
					if strings.Contains(v, ".png") || strings.Contains(v, ".jpg") || strings.Contains(v, ".svg") || strings.Contains(v, ".ico") {						
						logo := strings.Trim(v, "href=")											
						resultLogo = logo				
					}				
				}			
			}
		}
	}			
	return resultLogo
}

// Reads HTML file and return the title 
func getTitle(urlString string) (string, bool) {	
	isdown := false
	t := "title"
	var resultTitle string
	resp, err := http.Get(urlString)

	if err != nil {		
		isdown = true		
	} else {
		defer resp.Body.Close()
		//create a new tokenizer over the response body
		tokenizer := html.NewTokenizerFragment(resp.Body, "head")

		for {
			tokenType := tokenizer.Next() // get the token type
			
			if tokenType == html.ErrorToken {
				err := tokenizer.Err()			
				if err == io.EOF {				
					break //end of the file
				}			
			}				
			if tokenType == html.StartTagToken {			
				token := tokenizer.Token()	// get the token	
				// if the name of the element is "title"
				if t == token.Data {				
					tokenType = tokenizer.Next() 	// get the type of the next token			
					//get the page title
					resultTitle = tokenizer.Token().Data
					break																		
				}  						
			}	
		}
	}	
	return resultTitle, isdown
}
