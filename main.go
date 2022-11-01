package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"sort"
	"strconv"
	"strings"
)

type AuthorResponse struct {
	Documents []Authors `json:"docs"`
}
type Authors struct {
	Authors []string `json:"author_name"`
}

type WorkResponse struct {
	Documents []Work `json:"docs"`
}
type Work struct {
	Key   string `json:"key"`
	Title string `json:"title"`
}

type RevisionResponse struct {
	Revision int `json:"revision"`
}

type WorkSort struct {
	work     Work
	revision int
}

// type skuska struct {
// 	Key string
// 	Value int
// }



func main() {
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("Enter books name:")
	nazovKnihy, err := reader.ReadString('\n')
	if err != nil {
		log.Fatal(err)
	}
	nazovKnihy = strings.TrimSuffix(nazovKnihy, "\n")

	fmt.Println("Enter works sorting  (asc, desc):")
	zoradenie, err := reader.ReadString('\n')
	if err != nil {
		log.Fatal(err)
	}
	zoradenie = strings.TrimSuffix(zoradenie, "\n")

	var baseUrl = "https://openlibrary.org/search.json?title=" + url.QueryEscape(nazovKnihy)

	fmt.Println("Searching: " + baseUrl)

	titleResponse, err := http.Get(baseUrl)
	if err != nil {
		log.Fatal(err)
	}

	defer titleResponse.Body.Close()
	titleResponseData, err := ioutil.ReadAll(titleResponse.Body)
	if err != nil {
		log.Fatal(err)
	}
	var authorResponse AuthorResponse
	json.Unmarshal(titleResponseData, &authorResponse)

	var authors []string
	for i := 0; i < len(authorResponse.Documents); i++ {
		for j := 0; j < len(authorResponse.Documents[i].Authors); j++ {
			var authorName = authorResponse.Documents[i].Authors[j]
			authors = append(authors, authorName)
		}
	}

	authors = removeDuplicateStringValues(authors)
	sort.Strings(authors)
	for i := 0; i < len(authors); i++ {
		fmt.Println(authors[i] + ":")

		var url2 = "https://openlibrary.org/search.json?author=" + url.QueryEscape(authors[i])

		authorResponse, err := http.Get(url2)
		if err != nil {
			log.Fatal(err)
		}
		defer authorResponse.Body.Close()
		authorResponseData, err := ioutil.ReadAll(authorResponse.Body)
		if err != nil {
			log.Fatal(err)
		}
		var workResponse WorkResponse
		json.Unmarshal(authorResponseData, &workResponse)

		var books []Work
		books = removeDuplicateValues(workResponse.Documents)

		var sortBooks []WorkSort
		for k := 0; k < len(books); k++ {
			var url3 = "https://openlibrary.org" + books[k].Key + ".json"

			workResponse, err := http.Get(url3)
			if err != nil {
				log.Fatal(err)
			}
			defer workResponse.Body.Close()
			workResponseData, err := ioutil.ReadAll(workResponse.Body)
			if err != nil {
				log.Fatal(err)
			}

			var revResponse RevisionResponse
			json.Unmarshal(workResponseData, &revResponse)

			var workSort WorkSort
			workSort.work = books[k]
			workSort.revision = revResponse.Revision

			sortBooks = append(sortBooks, workSort)
		}

		// here todo > sort sortBooks by value
		// reference https://www.geeksforgeeks.org/how-to-sort-golang-map-by-keys-or-values/
		
		if zoradenie == "acs" {
		sort.Slice(sortBooks, func(i, j int) bool {
			return sortBooks[i].revision < sortBooks[j].revision
		})
	}		else if zoradenie == "desc"{
		sort.Slice(sortBooks, func(i, j int) bool {
			return sortBooks[i].revision > sortBooks[j].revision
		})
	}
		for k := 0; k < len(sortBooks); k++ {
		fmt.Println("- " + sortBooks[k].work.Title + " (" + strconv.Itoa(sortBooks[k].revision) + ")")

		}		
		fmt.Println("")
		fmt.Println("")
		

		
		
	}
}

func removeDuplicateStringValues(stringSlice []string) []string {
	keys := make(map[string]bool)
	list := []string{}

	// If the key(values of the slice) is not equal
	// to the already present value in new slice (list)
	// then we append it. else we jump on another element.
	for _, entry := range stringSlice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}

func removeDuplicateValues(workSlice []Work) []Work {
	keys := make(map[string]bool)
	list := []Work{}

	// If the key(values of the slice) is not equal
	// to the already present value in new slice (list)
	// then we append it. else we jump on another element.
	for _, entry := range workSlice {
		if _, value := keys[entry.Title]; !value {
			keys[entry.Title] = true
			list = append(list, entry)
		}
	}
	return list

	
}




// func sortbyrevision(records map[string]int64)
	