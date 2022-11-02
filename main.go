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
	"regexp"
	"sort"
	"strconv"
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

func main() {
	reader := bufio.NewReader(os.Stdin)
	re := regexp.MustCompile(`\r?\n`)

	fmt.Println("Enter books name:")
	booksName, err := reader.ReadString('\n')
	if err != nil {
		log.Fatal(err)
	}
	booksName = re.ReplaceAllString(booksName, "")

	fmt.Println("Enter works sorting  (asc, desc):")
	order, err := reader.ReadString('\n')
	if err != nil {
		log.Fatal(err)
	}
	order = re.ReplaceAllString(order, "")

	var baseUrl = "https://openlibrary.org/search.json?title=" + url.QueryEscape(booksName)

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

		if order == "asc" {
			sort.Slice(sortBooks, func(i, j int) bool {
				return sortBooks[i].revision < sortBooks[j].revision
			})
		}
		if order == "desc" {
			sort.Slice(sortBooks, func(i, j int) bool {
				return sortBooks[i].revision > sortBooks[j].revision
			})
		}

		for k := 0; k < len(sortBooks); k++ {
			fmt.Println("   - " + sortBooks[k].work.Title + " (" + strconv.Itoa(sortBooks[k].revision) + ")")

		}
		fmt.Println("")
		fmt.Println("")
	}

}

func removeDuplicateStringValues(stringSlice []string) []string {
	keys := make(map[string]bool)
	list := []string{}

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

	for _, entry := range workSlice {
		if _, value := keys[entry.Title]; !value {
			keys[entry.Title] = true
			list = append(list, entry)
		}
	}
	return list

}
