package cmdOptions

import (
   "os"
   "fmt"
   "strings"
   "log"
   "bufio"
)

func Banner(){
   var banner = `
    ___         __  __    ____                           __            
   /   | __  __/ /_/ /_  /  _/___  _________  ___  _____/ /_____  _____
  / /| |/ / / / __/ __ \ / // __ \/ ___/ __ \/ _ \/ ___/ __/ __ \/ ___/
 / ___ / /_/ / /_/ / / // // / / (__  ) /_/ /  __/ /__/ /_/ /_/ / /    
/_/  |_\__,_/\__/_/ /_/___/_/ /_/____/ .___/\___/\___/\__/\____/_/     
                                    /_/                                
                                                  - By karma9874`
   fmt.Println(banner)

}

func YesOrNo(s string) bool {
   reader := bufio.NewReader(os.Stdin)

   for {
      fmt.Printf("%s [y/n]: ", s)

      response, err := reader.ReadString('\n')
      if err != nil {
         log.Fatal(err)
      }

      response = strings.ToLower(strings.TrimSpace(response))

      if response == "y" || response == "yes" {
         return true
      } else if response == "n" || response == "no" {
         return false
      }
   }
   return true
}
