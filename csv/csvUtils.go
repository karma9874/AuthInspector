package csv

import (
   "github.com/karma9874/AuthInspector/cmdOptions"
   "github.com/karma9874/AuthInspector/http"
   "strings"
   "strconv"
   "fmt"
   "os"
)


var CSVHeader = []string{"URL","HTTP Method","Burp Response Status Code", "Burp Response Content Length"}

func PopulateCSVHeaders(authHeadersLen int,cmdOpt *cmdOptions.CmdOpt) []string{
   if cmdOpt.IsRequestBody{
      CSVHeader = append(CSVHeader,"Burp Request Body")
   }
   if cmdOpt.IsResponseBody{
      CSVHeader = append(CSVHeader,"Burp Response Body")
   }
   for i:=0 ; i < authHeadersLen; i++{

      if i == authHeadersLen -1{
         CSVHeader = append(CSVHeader,"UnAuth Status Code")
         CSVHeader = append(CSVHeader,"UnAuth Content Length")
         if cmdOpt.IsResponseBody{	
            CSVHeader = append(CSVHeader,"UnAuth Response Body")  
         }         
      }else{
         CSVHeader = append(CSVHeader,"Auth"+strconv.Itoa(i+1)+" Status Code")
         CSVHeader = append(CSVHeader,"Auth"+strconv.Itoa(i+1)+" Content Length")
         if cmdOpt.IsResponseBody{
            CSVHeader = append(CSVHeader,"Auth"+strconv.Itoa(i+1)+" Response Body")  
         }
      }
   }

   CSVHeader = append(CSVHeader,"Result")
   return CSVHeader
}


func escapeCSVString(s string) string {
   s = strings.ReplaceAll(s,"\"", "'")
   s = strings.ReplaceAll(s,"\n", "\\n")
   s = strings.ReplaceAll(s,"\r", "\\r")
   s = strings.ReplaceAll(s,"\t", "\\t")
   return s
}


func printFormattedString(csvfileWriter *os.File,str string) {
   formatted_string := fmt.Sprintf(`"%s"`,escapeCSVString(strings.TrimSpace(str)))
   if len(formatted_string) > 30000 {
      fmt.Fprint(csvfileWriter,",","..............text is too long above 30k")
   } else {
      fmt.Fprint(csvfileWriter,",", formatted_string)
   }
}

func DoCSV(csvfileWriter *os.File,res map[string]map[int]http.HttpResult,authHeadersLen int, cmdOptions *cmdOptions.CmdOpt){
   for _, value := range res { 

      isBypassed := ""
      burpResponse := http.GetBurpResponseBody(value[0].BurpItem.Response)
       
      fmt.Fprint(csvfileWriter,strings.Join([]string{value[0].URL,value[0].Method,strconv.Itoa(value[0].BurpItem.Status),strconv.Itoa(len(burpResponse))},","))

      if cmdOptions.IsRequestBody{
         printFormattedString(csvfileWriter,value[0].RequestBody)
      }
      if cmdOptions.IsResponseBody{
         printFormattedString(csvfileWriter,burpResponse)
      }

      for i := 0 ;i < authHeadersLen; i++{
         
         fmt.Fprint(csvfileWriter,",",strings.Join([]string{value[i].StatusCode,value[i].Size},","))
         if cmdOptions.IsResponseBody{
            printFormattedString(csvfileWriter,value[i].Body)
         }
         if strconv.Itoa(len(burpResponse)) == value[i].Size{
            isBypassed = "Bypassed"
         }
         if i != 0{
            if (value[i].Size == value[i-1].Size) && (value[i].Size != strconv.Itoa(len(burpResponse))){
               isBypassed = "Sus"
            }   
         }
      }

   if len(isBypassed) == 0{
      isBypassed = "Enforced"
   }

   fmt.Fprint(csvfileWriter,",",isBypassed)  
   fmt.Fprintln(csvfileWriter)
   }
}