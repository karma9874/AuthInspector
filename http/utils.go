package http

import (
   "fmt"
   "slices"
   "log"
   "strings"
   "encoding/base64"
   "github.com/karma9874/AuthInspector/cmdOptions"
   "github.com/karma9874/AuthInspector/burp"
)

func GetPostDataAndContentType(request string) (string,HTTPHeader){
   var HTTPHeaderStruct HTTPHeader
   rawDecodedText,_ := base64.StdEncoding.DecodeString(request)
   token := strings.Split(string(rawDecodedText),"\r\n\r\n")

   header_token := strings.Split(token[0],"\r\n")
   for _,i := range(header_token){
      token_seg := strings.Split(i,":")
      if len(token_seg) > 1{
         if strings.TrimSpace(token_seg[0]) == "Content-Type"{
            HTTPHeaderStruct.Name =  strings.TrimSpace(token_seg[0])
            HTTPHeaderStruct.Value = strings.TrimSpace(token_seg[1])
         } 
      }
   }
   // if token[len(token)-2] == "\n" && token[len(token)-1] != "\n" {
   //    fmt.Println(token[len(token)-1])
   //    return token[1],HTTPHeaderStruct
   // }else{
      return token[1],HTTPHeaderStruct
   // }
}

func GetBurpResponseBody(Response string) string{

   rawDecodedText,err := base64.StdEncoding.DecodeString(Response)
   decoded := strings.Split(string(rawDecodedText),"\r\n")
   BurpRespBody_val := ""
   BurpRespBody_val = decoded[len(decoded)-1]   
   
   if err != nil{
      BurpRespBody_val = "Error decoding the data"
   }

   return BurpRespBody_val
}


func GetAllMimeTypes(burpData []burp.Item){
   Mimetype_list := []string{}
   for _,j := range burpData{
      if !slices.Contains(Mimetype_list,j.Mimetype) && len(j.Mimetype) > 2{
         Mimetype_list = append(Mimetype_list,j.Mimetype)
      }
   }
   fmt.Println("Available Mime:",strings.Join(Mimetype_list,","))
   res := cmdOptions.YesOrNo("Do you want to continue?")
   if !res{
      log.Fatal("Exiting")
   }
}

func GetAllReqCount(burpItem []burp.Item,FilterMimeTypes []string)int{
   total_req := 0
   for _,mt := range(FilterMimeTypes){
      for _,burpItemData := range(burpItem){
         if burpItemData.Mimetype == mt{
            total_req++
         }
      }
   }
   return total_req
}



func initHeader(globalHeaders []map[string]string, authheaders map[string]string) []HTTPHeader{
   header := []HTTPHeader{}

   for key,value := range(authheaders){
      header = append(header,HTTPHeader{Name: key, Value: value})
   }
   for _,j := range(globalHeaders){
      for key,value := range(j){
         header = append(header,HTTPHeader{Name: key, Value: value})   
      }
   }
   return header
}