# AuthInspector

[![Apache-2.0 License](https://img.shields.io/badge/license-Apache2.0-blue.svg)](http://www.apache.org/licenses/) 
[![Twitter Follow](https://img.shields.io/twitter/follow/karma9874?label=Follow&style=social)](https://twitter.com/karma9874)
[![GitHub followers](https://img.shields.io/github/followers/karma9874?label=Follow&style=social)](https://github.com/karma9874)

AuthInspector is an advanced authorization and authentication testing tool designed to automate the assessment of authorization checks using multiple authentication headers. Seamlessly integrated with Burp Suite-generated requests file.

# Work flow
<h3 align="center">
  <img src="static/flow.png" width="700px"></a>
</h3>

# Installation

### Easy Installations
You can download the prebuilt binary from the [releases](https://github.com/karma9874/AuthInspector/releases) page.

### go install 
`go install github.com/karma9874/AuthInspector@latest`

### go build
`go get && go build`

# Usage
AuthInspector provides the following commands for customization:
```
-proxy		Set up a proxy for testing.
-respBody	Include response body in the output
-reqBody	Include request body in the output.
-timeout	Set the timeout for requests.
-threads	Specify the number of concurrent threads.
-listmime 	Lists the available mimetypes from the burp exported file
-verbose	Verbose output
```

## Running AuthInspector
`AuthInspector -proxy http://proxy.example.com -respBody -reqBody -time 5s -threads 20`

## Config Template
init.yaml
```
# Burp XML file name to be used in the authentication testing process.
source: example.xml

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
```



