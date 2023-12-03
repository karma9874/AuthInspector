package burp

type Item struct{
   Url               string         `xml:"url"`
   Method            string         `xml:"method"` 
   Extension         string         `xml:"extension"`
   Request           string         `xml:"request"`
   Mimetype          string         `xml:"mimetype"`
   Response          string         `xml:"response"`
   Status            int            `xml:"status"`
   Responselength    string         `xml:"responselength"`
}


type Items struct {
   Item              []Item         `xml:"item"`     
}