package main

import (
   "fmt"
   "strings"
   "sync"
   "time"
   "flag" 
   "log"
   "os"
   "github.com/cheggaaa/pb/v3"
   "github.com/karma9874/AuthInspector/http"
   "github.com/karma9874/AuthInspector/burp"
   "github.com/karma9874/AuthInspector/cmdOptions"
   "github.com/karma9874/AuthInspector/config"
   "github.com/karma9874/AuthInspector/csv"
)

func main() {

   wg := &sync.WaitGroup{}
   cmdOptions := new(cmdOptions.CmdOpt)
   var progressBar *pb.ProgressBar
   
   flag.StringVar(&cmdOptions.IsProxy, "proxy", "", "Use proxy")
   flag.BoolVar(&cmdOptions.IsResponseBody, "respBody",false, "Get response body")
   flag.BoolVar(&cmdOptions.IsRequestBody, "reqBody",false, "Get request body")
   flag.BoolVar(&cmdOptions.MimeType, "listmime",false, "List all mime type from burp file")
   flag.BoolVar(&cmdOptions.IsVerbose, "verbose",false, "Enable verbose output")
   flag.DurationVar(&cmdOptions.TimeOut, "timeout",5*time.Second, "Set timeout for request in sec")
   flag.IntVar(&cmdOptions.Threads, "threads", 10, "Number of concurrent threads")
   flag.Parse()
   
   
   
   successfulRequests := 0
   errorCount  := 0
   

   
   yamlObj := config.ReadConfigYaml()
   xmlData := burp.ReadBurpXML(yamlObj.SourceFileName)

   if cmdOptions.MimeType{
      http.GetAllMimeTypes(xmlData.Item)
   }


   startTime := time.Now()

   totalRequests := http.GetAllReqCount(xmlData.Item,yamlObj.FilterMimeTypes)*len(yamlObj.AuthHeaders)

   if !cmdOptions.IsVerbose{
      progressBar = pb.StartNew(totalRequests)   
   }else{
      progressBar = nil
   }
   

   if cmdOptions.IsVerbose{
      log.Println("Configurations:")
      log.Printf("Source YAML File: %s\n", yamlObj.SourceFileName)
      log.Printf("Threads: %d, Timeout: %s\n", cmdOptions.Threads, cmdOptions.TimeOut)
   }

   fmt.Println("\nInspecting authentication issues...")

   resQueue := make(chan map[string]map[int]http.HttpResult, cmdOptions.Threads)
   
   for _,mt := range(yamlObj.FilterMimeTypes){
      for _,burpItemData := range(xmlData.Item){
         if burpItemData.Mimetype == mt{
            postData,req_HTTPHeader := http.GetPostDataAndContentType(burpItemData.Request)      
            http.MakeRequestMultiAuth(postData,burpItemData,yamlObj.AuthHeaders,yamlObj.GlobalHeaders,req_HTTPHeader,wg,resQueue,cmdOptions,progressBar)
         }
      }
   }

   go func() {
      wg.Wait()
      close(resQueue)
   }()

   finalResult := make(map[string]map[int]http.HttpResult)
   for elem := range resQueue {

      for key, value := range elem {
            if _, ok := finalResult[key]; !ok {
                finalResult[key] = make(map[int]http.HttpResult)
            }
            for idx, res := range value {
               if res.StatusCode == "timedOut"{
                  errorCount++
               }else{
                  successfulRequests++
               }
               finalResult[key][idx] = res
            }
        }
   }
   if !cmdOptions.IsVerbose{progressBar.Finish()}

   csvfileWriter, csvwriter_err := os.Create("Output.csv")
   if csvwriter_err != nil{
      panic(csvwriter_err)
   }
   defer csvfileWriter.Close()

   fmt.Fprintln(csvfileWriter,strings.Join(csv.PopulateCSVHeaders(len(yamlObj.AuthHeaders),cmdOptions),","))
   csv.DoCSV(csvfileWriter,finalResult,len(yamlObj.AuthHeaders),cmdOptions)
   fmt.Println("\nAuthInspector completed successfully.")
   
   fmt.Println("Output stored in Output.csv")
   elapsedTime := time.Since(startTime)
   fmt.Printf("Summary: Total Requests: %d, Successful Requests: %d, Errors: %d", totalRequests, successfulRequests, errorCount)
   fmt.Println("\nRequest completed in", elapsedTime)
}


// func checkCORS(respObj *http.Response, returnedURL string) {
//    if len(respObj.Header["Access-Control-Allow-Origin"]) > 0{
//       if returnedURL != respObj.Header["Access-Control-Allow-Origin"][0] {
//          fmt.Println("CORS issue received:",respObj.Header["Access-Control-Allow-Origin"])
//       }else{
//          fmt.Println("Origin input returned in Access-Control-Allow-Origin",respObj.Header["Access-Control-Allow-Origin"])
//       }
//    }
// }
