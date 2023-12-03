package http

import(
   "github.com/karma9874/AuthInspector/burp"
)

type HttpResult struct {
   Idx               int
   Id                string
   URL               string
   Method            string
   Header            HTTPHeader
   StatusCode        string
   Size              string
   Body              string
   Error             string
   RequestBody       string
   BurpItem          burp.Item
}

type HTTPHeader struct {
   Name              string
   Value             string
}