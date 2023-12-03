package cmdOptions

import (
   "time"
)

type CmdOpt struct{
   IsProxy           string
   TimeOut           time.Duration
   IsResponseBody    bool
   IsRequestBody     bool
   Threads           int
   MimeType          bool
   IsVerbose         bool
}
