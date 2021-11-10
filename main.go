package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"
)

type Response struct {
	Photos []struct {
		ImageLink string `json:"img_src"`
	}
}

func main() {
	/*
		we will search the current day - 1,
		because on the current day there are no images
	*/
	earthDate := time.Now().AddDate(0, 0, -1)
	earthDateFormat := earthDate.Format("2006-01-02")

	apiURL := &url.URL{
		Scheme:   "https",
		Host:     "api.nasa.gov",
		Path:     "mars-photos/api/v1/rovers/curiosity/photos",
		RawQuery: "api_key=DEMO_KEY&earth_date=" + earthDateFormat,
	}

	url, err := url.Parse(apiURL.String())

	if err != nil {
		log.Fatal(err)
	}
	apiCall(url.String(), earthDateFormat)
}

func apiCall(urlApi string, earthDate string) {

	// Get request
	resp, err := http.Get(urlApi)
	if err != nil {
		fmt.Println("No response from request")
	}

	body, _ := ioutil.ReadAll(resp.Body)

	var result Response
	if err := json.Unmarshal(body, &result); err != nil {
		fmt.Println("Can not unmarshal JSON")
	}

	var photos [10]string

	if len(result.Photos) == 0 {
		fmt.Println("No response from request")
	} else {
		// In case we don't receive at least 10 links, we list all the links we received from the API
		fmt.Printf("Listing the first %v links to photos taken by NASA's Curiosity Mars Rover in %v\n", len(result.Photos), earthDate)

		var text string
		counter := 1

		text += "### Listing the first " + strconv.Itoa(len(result.Photos)) + " links to photos taken by NASA's Curiosity Mars Rover in " + earthDate + "\n\n"

		for i := 0; i < 10; i++ {
			/*
			 * Return when i < total of photos received from the API,
			 * preventing runtime error index out of range
			 */
			if i == len(result.Photos) {
				return
			}

			photos[i] = result.Photos[i].ImageLink
			text += "\t-> Link " + strconv.Itoa(counter) + ": " + photos[i] + "\n"

			fmt.Printf("\t -> Link #%v: %v\n", counter, photos[i])

			logIntoFile(text)
			counter++
		}
	}
}

func logIntoFile(text string) {
	f, err := os.Create("logLinks.txt")
	if err != nil {
		fmt.Println(err)
		return
	}

	l, err := f.WriteString(text)
	if err != nil {
		fmt.Println(err)
		f.Close()
		return
	}

	fmt.Println(l, "bytes written successfully!")
	err = f.Close()
	if err != nil {
		fmt.Println(err)
		return
	}
}
