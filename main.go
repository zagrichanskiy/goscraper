package main

const (
	// RootURL to look up builds.
	RootURL = "http://dehil-ae03.debads.europe.delphiauto.net:82/" +
		"platforms/csp/csp/latest-release/"
	// SdkURL url to the sdk folder.
	SdkURL = RootURL + "sdk"
	// SdkReg regexp for the sdk file.
	SdkReg = `\.sh$`
	// Blade1URL url to the blade 1
	Blade1URL = RootURL + "csp-image-blade-i-intel-corei7-64-dom0.wic"
	// Blade2URL url to the blade 3
	Blade2URL = RootURL + "csp-image-blade-ii-intel-corei7-64-dom0.wic"
	// Blade3URL url to the blade 3
	Blade3URL = RootURL + "csp-image-blade-iii-intel-corei7-64-dom0.wic"
)

// func main() {
// 	// Getting the web page.
// 	fmt.Println("Openning url: ", DefaultURL)
// 	resp, err := http.Get(DefaultURL)
// 	if err != nil {
// 		fmt.Println("Can't open the url: ", err)
// 		return
// 	}
// 	defer resp.Body.Close()

// 	// Getting the list of urls on the page.
// 	z := html.NewTokenizer(resp.Body)
// 	urls := make([]string, 0)
// 	for tt := z.Next(); tt != html.ErrorToken; tt = z.Next() {
// 		if tt == html.StartTagToken {
// 			if t := z.Token(); t.Data == "a" {
// 				fillURLArray(&t, &urls)
// 			}
// 		}
// 	}

// 	for _, url := range urls {
// 		fmt.Println(" - ", url)
// 	}

// 	// Finding the link to the balde3
// 	fmt.Println(DefaultURL + getURL(Blade3Reg, urls))
// }

// func fillURLArray(t *html.Token, urls *[]string) {
// 	for _, a := range t.Attr {
// 		if a.Key == "href" {
// 			*urls = append(*urls, a.Val)
// 			break
// 		}
// 	}
// }

// func getURL(reg string, urls []string) string {
// 	var ret string
// 	re := regexp.MustCompile(reg)
// 	for _, url := range urls {
// 		if re.MatchString(url) {
// 			ret = url
// 			break // We suppose there is only one match in array.
// 		}
// 	}
// 	return ret
// }
