package play

import (
   "fmt"
   "os"
   "testing"
   "time"
)

func TestSync(t *testing.T) {
   home, err := os.UserHomeDir()
   if err != nil {
      t.Fatal(err)
   }
   home += "/google-play"
   for _, abi := range Abis {
      fmt.Println(abi)
      Device.Abi = abi
      var data []byte
      _, err := Device.Checkin(&data)
      if err != nil {
         t.Fatal(err)
      }
      err = os.WriteFile(home+"/"+abi+".txt", data, os.ModePerm)
      if err != nil {
         t.Fatal(err)
      }
      time.Sleep(time.Second)
      var checkin GoogleCheckin
      err = checkin.Unmarshal(data)
      if err != nil {
         t.Fatal(err)
      }
      err = Device.Sync(&checkin)
      if err != nil {
         t.Fatal(err)
      }
      time.Sleep(time.Second)
   }
}
