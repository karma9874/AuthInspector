package http

import (
   "net/http"
   "net/url"
   "crypto/tls"
   "log"
   "sync"
   "bytes"
   "io/ioutil"
   "strconv"
   "github.com/google/uuid"
   "github.com/cheggaaa/pb/v3"
   "github.com/karma9874/AuthInspector/burp"
   "github.com/karma9874/AuthInspector/cmdOptions"
)

type HTTPClient struct {   
   UserAgent         string
   DefaultUserAgent  string
   Headers           []HTTPHeader
   Method            string
   Body              string
   URL               string
}

var defaultUserAgent = "Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/109.0.0.0 Safari/537.36" 

func MakeRequest(opt HTTPClient,c chan map[string]map[int]HttpResult,wg *sync.WaitGroup,uid string, burpItm burp.Item, idx int, cmdOptions *cmdOptions.CmdOpt,pbar *pb.ProgressBar)  {

   defer wg.Done()

   retValue := make(map[string]map[int]HttpResult)
   
   // log.Printf("Testing for %s %s\n", opt.Method, opt.URL)

   var errorStr = "" 
   //fmt.Println(opt.URL,opt.Body)
   req, _ := http.NewRequest(opt.Method,opt.URL,bytes.NewBuffer([]byte(opt.Body)))

   req.Header["User-Agent"] = []string{opt.DefaultUserAgent}

   if len(opt.Headers) > 0 {
      for _, h := range opt.Headers {
         req.Header[h.Name] = []string{h.Value}
      }
   }

   var proxyFunc = http.ProxyFromEnvironment
   if len(cmdOptions.IsProxy) > 0 {
      proxy, err := url.Parse(cmdOptions.IsProxy)
      if err != nil{
         log.Println("Error parsing proxy")
      }
      proxyFunc = http.ProxyURL(proxy)
   }

   tlsConfig := tls.Config{InsecureSkipVerify: true}

   client := &http.Client{Transport: &http.Transport{Proxy: proxyFunc,TLSClientConfig: &tlsConfig},Timeout: cmdOptions.TimeOut}
   resp, req_err := client.Do(req)
   
   if !cmdOptions.IsVerbose{
      pbar.Increment()   
   }

   if req_err != nil {
      errorStr = "timedOut"
      if cmdOptions.IsVerbose{log.Printf("[Failed] %s %s - Request Failed. Status Code: %d, Size: %d bytes\n", opt.Method, opt.URL, resp.StatusCode,0)}
      //log.Printf("Goroutine for %s %s finished\n", opt.Method, opt.URL)
      retValue[uid] = map[int]HttpResult{
            idx: {
                URL:         opt.URL,
                Method:      opt.Method,
                //Header:      nil,
                StatusCode:  errorStr,
                Size:        "\"\"",
                Body:        "",
                Error:       errorStr,
                RequestBody: opt.Body,
                Id:          uid,
                Idx:         idx,
                BurpItem:    burpItm,
            },
        }
     
      c <- retValue
      return
   }

   body,body_err := ioutil.ReadAll(resp.Body)
   //fmt.Println(string(body))
   if body_err != nil {
      body = []byte("Error on response")
   }
   if cmdOptions.IsVerbose{log.Printf("[Success] %s %s - Request completed. Status Code: %d, Size: %d bytes\n", opt.Method, opt.URL, resp.StatusCode, len(body))}
   //log.Printf("Goroutine for %s %s finished\n", opt.Method, opt.URL)
   retValue[uid] = map[int]HttpResult{
        idx: {
            URL:         opt.URL,
            Method:      opt.Method,
            //Header:      resp.Header,
            StatusCode:  strconv.Itoa(resp.StatusCode),
            Size:        strconv.Itoa(len(body)),
            Body:        string(body),
            Error:       errorStr,
            RequestBody: opt.Body,
            Id:          uid,
            Idx:         idx,
            BurpItem:    burpItm,
        },
    }

   c <- retValue

   return
}


func MakeRequestMultiAuth(postData string,k burp.Item,yamlAuthheaders []map[string]string,globalHeaders []map[string]string,HTTPHeader HTTPHeader,wg *sync.WaitGroup, resQueue chan map[string]map[int]HttpResult,cmdOpts *cmdOptions.CmdOpt,pbar *pb.ProgressBar) {

   id := uuid.New()
   //log.Printf("Testing for %s %s\n", k.Method, k.Url)
   for i := 0; i<len(yamlAuthheaders); i++ {
      header := initHeader(globalHeaders,yamlAuthheaders[i])
      if HTTPHeader.Name != "" && HTTPHeader.Value != ""{
         header = append(header,HTTPHeader)   
      }
      
      wg.Add(1)
      go func(data HTTPClient,c chan map[string]map[int]HttpResult,wg *sync.WaitGroup,uid string, burpItem burp.Item, idx int,cmdOpts *cmdOptions.CmdOpt,progressBar *pb.ProgressBar){
         
         MakeRequest(HTTPClient{URL: k.Url, Method: k.Method, DefaultUserAgent: defaultUserAgent, Body: postData, Headers: header},resQueue, wg,id.String(),k,idx,cmdOpts,pbar)
      }(HTTPClient{URL: k.Url, Method: k.Method, DefaultUserAgent: defaultUserAgent, Body: postData, Headers: header}, resQueue,wg,id.String(),k,i,cmdOpts,pbar)

   }
}
