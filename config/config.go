package config

import (
	"io/ioutil"
	"gopkg.in/yaml.v3"
	"fmt"
	"log"
	"os"
)

var sourceYAMLFileName = "init.yaml"

const defaultYAMLContent = `# Burp XML file name to be used in the authentication testing process.
source: burp_file.xml

# Headers with authentication information.
auth:
  - header_key: header_value
  - header_key: header_value  # Do not remove this header (use to check unauthenticated requests)

# Mime types(case sensitive, for more details list mime type check -listmime mode). The tool will focus on checking authentication issues only on specified mime types.
filterMimeTypes:
  - JSON
  - XML

# Global headers to be included in all requests.
headers:
  - User-Agent: "Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/109.0.0.0 Safari/537.36"
  - API-KEY: some_key
`


type yamlConfig struct {   
   SourceFileName    string                 `yaml:"source"`
   AuthHeaders       []map[string]string    `yaml:"auth"`
   FilterMimeTypes   []string               `yaml:"filterMimeTypes"`
   GlobalHeaders     []map[string]string    `yaml:"headers"`
}

func createExampleYAML() error {
   _, err := os.Stat("init.yaml")
   if os.IsNotExist(err) {
      err := ioutil.WriteFile("init.yaml", []byte(defaultYAMLContent), 0644)
      if err != nil {
         return err
      }
      fmt.Println("Example init.yaml file created successfully.")
   } else if err != nil {
      return err
   }

   return nil
}

func ReadConfigYaml() (yamlConfig){
   yamlFile, err := ioutil.ReadFile(sourceYAMLFileName)
   if err != nil {
      fmt.Println("File Not Found",err)
      fmt.Println("Creating example init.yaml file")
      errCreating := createExampleYAML()
      if errCreating != nil {
         fmt.Println("Error creating example init.yaml:",errCreating)
         log.Fatal("Unable to create file",errCreating)
      }
      log.Fatal("Exiting") 
   }

   config := yamlConfig{}
   err = yaml.Unmarshal(yamlFile,&config)

   if err != nil {
      log.Fatal("Not able to parse yaml file",err)
   }
   return config
}

