package burp

import (
	"log"
	"io/ioutil"
	"encoding/xml"
)

func ReadBurpXML(filename string) Items {
   xmlFile, err := ioutil.ReadFile(filename)
   if err != nil {
      log.Fatal("File Not Found",err)
   }
   var config Items
   err = xml.Unmarshal(xmlFile,&config)

   if err != nil {
      log.Fatal("Error parsing xml file",err)
   }
   return config
}

