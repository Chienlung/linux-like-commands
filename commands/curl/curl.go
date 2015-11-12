package main
import (
	"flag"
	"net/http"
	"log"
	"io/ioutil"
	"fmt"
	"strings"
	"bytes"
)

type headers []string
func (h *headers) String() string {
	return fmt.Sprint(*h)
}
func (h *headers) Set(value string) error {
//	if len(*h) > 0 {
//		return errors.New("headers flag has been set")
//	}

	for _, v := range strings.Split(value, ",") {
		*h = append(*h, v)
	}
	return nil
}
var (i bool
	x string
	h headers
	d string
	url string)

func init() {
	flag.BoolVar(&i, "i", false, "include the http-headers of the output?")
	flag.StringVar(&x, "X", "GET", "http method")
	flag.Var(&h, "H", "http header")
	flag.StringVar(&d, "d", "", "post body")
//	flag.StringVar(&url, "url", "", "url")
}
func main() {
	flag.Parse()
	url = flag.Arg(flag.NArg() - 1)
	if x == "GET" {
		get()
	} else if x == "POST" {
		post()
	}

}

func get() {
	resp, err := http.Get(url)
	if err != nil {
		log.Fatalln(err)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	if i {
		fmt.Println(resp.Proto, resp.Status)
		for k, v := range resp.Header {
			fmt.Print(k, ":")
			for _, value := range v {
				fmt.Print("\t", value)
			}
			fmt.Print("\n")
		}
		fmt.Print("\n")
	}
	fmt.Println(string(body))
}

func post() {
	// set the body
	var body *bytes.Reader
	if strings.HasPrefix(d, "@") {
		content, err := ioutil.ReadFile(strings.TrimLeft(d, "@"))
		if err != nil {
			log.Fatalln(err)
		}
		body = bytes.NewReader(content)
	} else {
		body = bytes.NewReader([]byte(d))
	}
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		log.Fatalln(err)
	}

	// set the headers
	for _, header := range h {
		kv := strings.Split(header, ":")
		if len(kv) != 2 {
			log.Fatalln("invalid header:", header)
		}
		req.Header.Add(strings.Trim(kv[0], " "), strings.Trim(kv[1], " "))
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalln(err)
	}
	defer resp.Body.Close()

	respContent, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	if i {
		fmt.Println(resp.Proto, resp.Status)
		for k, v := range resp.Header {
			fmt.Print(k, ":")
			for _, value := range v {
				fmt.Print("\t", value)
			}
			fmt.Print("\n")
		}
		fmt.Print("\n")
	}
	fmt.Println(string(respContent))
}