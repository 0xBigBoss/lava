package proxy

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

var MockFolder string = "testutil/e2e/proxy/mockMaps/"

var responsesChanged bool = false
var total uint = 0
var realCount int = 0
var cacheCount int = 0
var fakeCount int = 0

var fakeResponse bool = false

var saveJsonEvery int = 10 // in seconds
var epochTime int = 3      // in seconds
var epochCount int = 0     // starting epoch count
var proxies []Proxy = []Proxy{}

type Proxy struct {
	id           string
	port         string
	host         string
	mockfile     string
	mock         *mockMap
	handler      func(http.ResponseWriter, *http.Request)
	malicious    bool
	cache        bool
	strict       bool
	noSave       bool
	latency      uint
	availability uint
}

type jsonStruct struct {
	Jsonrpc string        `json:"jsonrpc"`
	ID      int           `json:"id,omitempty"`
	Method  string        `json:"method"`
	Params  []interface{} `json:"params,omitempty"`
}

func getDomain(s string) (domain string) {
	parts := strings.Split(s, ".")
	if len(parts) >= 2 {
		return parts[len(parts)-2]
	}
	return s
}

var current *Proxy

func main() {

	// CLI ARGS
	host := flag.String("host", "", "HOST (required) - Which host do you wish to proxy\nUsage Example:\n	$ go run proxy.go http://google.com/")
	port := flag.String("p", "1111", "PORT")
	id := flag.String("id", "", "ID (optional) - will set the default host id instead of the full domain name"+
		"\nUsage Example:\n	$ go run proxy.go randomnumberapi.com -id random -cache")
	cache := flag.Bool("cache", false, "CACHE (optional) - This will make proxy return from cache if possible "+
		"(from default [host].json unless -alt was set)\nUsage Example:\n	$ go run proxy.go http://google.com/ -cache")
	alt := *flag.String("alt", "", "ALT (optional) [JSONFILE] - This will make proxy return from alternative cache file if possible"+
		"\nUsage Example:\n	$ go run proxy.go http://google.com/ -cache -alt ./mockMaps/google_alt.json		# respond from google_alt.json")
	latency := flag.Uint("l", 0, "min latency of the proxy answer")
	availability := flag.Uint("a", 100, "availability, every x messages will be intercepted")
	strict := flag.Bool("strict", false, "STRICT (optional) - This will make proxy return ONLY from cache, no external calls")
	help := flag.Bool("h", false, "Shows this help message")
	noSave := *flag.Bool("no-save", false, "NO-SAVE (optional) will not store any data from proxy")

	flag.Parse()
	if *help || (*host == "" && flag.NArg() == 0) {
		fmt.Println("\ngo run proxy.go [host] -p [port] OPTIONAL -cache -alt [JSONFILE] -strict")
		fmt.Println("	Usage Example:")
		fmt.Println("	$ go run proxy.go -host google.com/ -p 1111 -cache")
		flag.Usage()
	} else if *host == "" {
		if len(os.Args) > 0 {
			if os.Args[1] != "-host" {
				*host = os.Args[1]
				flag.CommandLine.Parse(append([]string{"-host"}, os.Args[1:]...))
			} else {
				*host = os.Args[1]
			}
		}
	}
	println()

	domain := getDomain(*host)
	if *id != "" {
		domain = *id
	} else {
		*id = domain
	}

	mockfile := GetMockFilePath(domain, alt)

	if *host == "" {
		println("\n [host] is required. Exiting")
		os.Exit(1)
	}
	malicious := false // default

	StartEpochUpdate(noSave)

	process := Proxy{
		id:           domain,
		port:         *port,
		host:         *host,
		mockfile:     mockfile,
		mock:         &mockMap{requests: map[string]string{}},
		malicious:    false,
		cache:        *cache,
		strict:       *strict,
		noSave:       noSave,
		latency:      *latency,
		availability: *availability,
	}
	// current = &process

	if !malicious {
		process.handler = process.LavaTestProxy
	} else {
		println()
		println("MMMMMMMMMMMMMMM MALICIOUS MMMMMMMMMMMMMMM PORT", port)
		println()
		// TODO: Make malicious proxy
		process.handler = process.LavaTestProxy
	}
	StartProxy(&process)
}

func NewProxy(id string,
	port int64,
	host string,
	mockfile string,
	malicious bool,
	cache bool,
	strict bool,
	noSave bool,
	latency uint,
	availability uint) *Proxy {
	proc := &Proxy{
		id:           id,
		port:         strconv.FormatInt(port, 10),
		host:         host,
		mockfile:     mockfile,
		mock:         &mockMap{requests: map[string]string{}},
		malicious:    false,
		cache:        cache,
		strict:       strict,
		noSave:       noSave,
		latency:      latency,
		availability: availability,
	}

	proc.handler = proc.LavaTestProxy
	return proc
}

func GetMockFilePath(domain string, alt string) string {
	mockfile := MockFolder + domain + ".json"
	if alt != "" {
		mockfile = MockFolder + alt
	}
	return mockfile
}

func StartProxy(process *Proxy) *http.Server {
	proxies = append(proxies, *process)
	process.mock.requests = jsonFileToMap(process.mockfile)
	if process.mock.requests == nil {
		process.mock.requests = map[string]string{}
	}
	if process.malicious {
		fakeResponse = true
	}
	fmt.Println(":::::::::::::::::::::::::::::::::::::::::::::::") // HOST ", process.host)
	fmt.Println("::::::::::::: Mock Proxy Started ::::::::::::::") // CACHE", fmt.Sprintf("%d", len(process.mock.requests)))
	fmt.Println(":::::::::::::::::::::::::::::::::::::::::::::::") // PORT ", process.port)
	println()
	fmt.Print(fmt.Sprintf(" ::: Proxy ID 		::: %s", process.id) + "\n")
	fmt.Print(fmt.Sprintf(" ::: Proxy Host 	::: %s", process.host) + "\n")
	fmt.Print(fmt.Sprintf(" ::: Return Cache 	::: %t", process.cache) + "\n")
	fmt.Print(fmt.Sprintf(" ::: Strict Mode 	::: %t", process.strict) + "\n")
	fmt.Print(fmt.Sprintf(" ::: Saving	 	::: %t", !process.noSave) + "\n")
	if !process.noSave || process.cache {
		fmt.Print(fmt.Sprintf(" ::: Cache File 	::: %s", process.mockfile) + "\n")
		fmt.Print(fmt.Sprintf(" ::: Loaded Responses 	::: %d", len(process.mock.requests)) + "\n")
	}
	println()
	fmt.Print(fmt.Sprintf(" ::: Proxy Started! 	::: ID: %s", process.id) + "\n")
	fmt.Print(fmt.Sprintf(" ::: Listening On 	::: %s", "http://0.0.0.0:"+process.port+"/") + "\n")

	mux := http.NewServeMux()
	mux.HandleFunc("/", process.handler)
	srv := &http.Server{Handler: mux}

	l, err := net.Listen("tcp", ":"+process.port)
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		if err := srv.Serve(l); err != nil {
			log.Fatal(err)
		}
	}()

	return srv
}

func getMockBlockNumber() (block string) {
	return "0xe" + fmt.Sprintf("%d", epochCount) // return "0xe39ab8"
}

func StartEpochUpdate(noSave bool) {
	go func() {
		count := 0
		for {
			wait := 1 * time.Second
			time.Sleep(wait)
			if count%epochTime == 0 {
				epochCount += 1
			}
			if !noSave && responsesChanged && count%saveJsonEvery == 0 {
				for _, process := range proxies {
					mapToJsonFile(*process.mock, process.mockfile)
				}
				responsesChanged = false
			}
			count += 1
		}
	}()
}

func fakeResult(val string, fake string) string {
	parts := strings.Split(val, ",")
	found := -1
	for i, part := range parts {
		if strings.Contains(part, "result") {
			found = i
		}
	}
	if found != -1 {
		parts[found] = fmt.Sprintf("\"result\":\"%s\"}", fake)
	}
	return strings.Join(parts, ",")
}

func trunc(str string) string {
	if len(str) > 200 {
		return strings.Clone(str[:200]) + "..." //too long is ugly
	} else {
		return str
	}
}

func (p Proxy) LavaTestProxy(rw http.ResponseWriter, req *http.Request) {

	host := p.host
	mock := p.mock

	// Get request body
	rawBody := getDataFromIORead(&req.Body, true)
	// TODO: set all ids to 1
	rawBodyS := string(rawBody)
	// sep := "id\":"
	// first := strings.Split(rawBodyS, sep)[0]
	// second := strings.Split(rawBodyS, sep)[1]
	// second = strings.Join(strings.Split(second, ",")[1:], ",")
	// rawBodyS = first + sep + "1," + second

	// println()
	// println(" ::: "+p.port+" ::: "+p.id+" ::: INCOMING PROXY MSG :::", rawBodyS)

	// TODO: make generic
	// Check if asking for blockNumber
	if fakeResponse && strings.Contains(rawBodyS, "blockNumber") {
		println("!!!!!!!!!!!!!! block number")
		p.returnResponse(rw, 200, []byte(fmt.Sprintf("{\"jsonrpc\":\"2.0\",\"id\":1,\"result\":\"%s\"}", getMockBlockNumber())))
	} else {
		// Return Cached data if found in history and fromCache is set on
		jStruct := &jsonStruct{}
		json.Unmarshal([]byte(rawBodyS), jStruct)
		jStruct.ID = 0
		rawBodySNoID, _ := json.Marshal(jStruct)
		if val, ok := mock.requests[string(rawBodySNoID)]; ok && p.cache {
			// println(" ::: "+p.port+" ::: "+p.id+" ::: Cached Response ::: ", trunc(val))
			cacheCount += 1

			// Change Response
			if fakeResponse {
				val = fakeResult(val, "0xe000000000000000000")
				// val = "{\"jsonrpc\":\"2.0\",\"id\":1,\"result\":\"0xe000000000000000000\"}"
				// println(p.port+" ::: Fake Response ::: ", val)
				fakeCount += 1
			}
			p.returnResponse(rw, 200, []byte(val))

		} else {
			// Recreating Request
			proxyRequest, err := createProxyRequest(req, host, rawBodyS)
			if err != nil {
				println(err.Error())
			} else {

				// Send Request to Host & Get Response
				proxyRes, err := sendRequest(proxyRequest)
				// respBody := []byte("error")
				var respBody []byte
				respBodyStr := "xxxxxx"
				status := 400
				if err != nil {
					println(err.Error())
					respBody = []byte(err.Error())
				} else {
					status = proxyRes.StatusCode
					respBody = getDataFromIORead(&proxyRes.Body, true)
					respBodyStr = string(respBody)
					mock.requests[rawBodyS] = respBodyStr
					realCount += 1
					println(" ::: "+p.port+" ::: "+p.id+" ::: Real Response ::: ", trunc(respBodyStr))
				}

				// Check if response is not good, if not - try again
				if false && (strings.Contains(string(respBody), "error") || strings.Contains(string(respBody), "Error")) {
					println("XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX Got error in response - retrying request")

					// Recreating Request
					proxyRequest, err = createProxyRequest(req, host, rawBodyS)
					if err != nil {
						println(err.Error())
						respBody = []byte(err.Error())
					} else {

						// Send Request to Host & Get Response
						proxyRes, err = sendRequest(proxyRequest)
						if err != nil {
							println(err.Error())
							respBody = []byte("error: " + err.Error())
						} else {
							respBody = getDataFromIORead(&proxyRes.Body, true)
							mock.requests[rawBodyS] = string(respBody)
							status = proxyRes.StatusCode
						}
						realCount += 1
						println(" ::: "+p.port+" ::: "+p.id+" ::: Real Response ::: ", trunc(string(respBody)))

						// TODO: Check if response is good, if not - try again
						if strings.Contains(string(respBody), "error") || strings.Contains(string(respBody), "Error") {
							println("XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX Got another error in response ")
							println()
						} else {
							println("YYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYY SUCCESS - no error in response ")
							println()
						}
					}

				}
				// Change Response
				if fakeResponse {
					// respBody = []byte("{\"jsonrpc\":\"2.0\",\"id\":1,\"result\":\"0xe000000000000000000\"}")
					respBody = []byte(fakeResult(respBodyStr, "0xe000000000000000000"))
					println(" ::: "+p.port+" ::: "+p.id+" ::: Fake Response ::: ", string(respBody))
					fakeCount += 1
				}
				responsesChanged = true

				//Return Response
				if respBody == nil {
					respBody = []byte("error")
				}

				p.returnResponse(rw, status, respBody)
			}
		}
	}
	if realCount > 0 || cacheCount > 0 {
		// id := ""
		// id = p.id + ":" + p.port
		// fmt.Println("_________________________________", realCount, "/", cacheCount, ": proxy sent (new/from cache)", id, "\n")
	}
}

func (p Proxy) returnResponse(rw http.ResponseWriter, status int, body []byte) {
	time.Sleep(time.Duration(p.latency) * time.Millisecond)

	total++
	if total%p.availability != 0 {
		fmt.Printf("Return Response sent %s %d\n", time.Now().String(), total)
		rw.WriteHeader(status)
		rw.Write(body)
	} else {
		fmt.Printf("Return Response intercepted %s %d\n", time.Now().String(), total)
	}

}
