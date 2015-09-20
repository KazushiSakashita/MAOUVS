package main

//import "github.com/ChimeraCoder/anaconda"
import "fmt"
//import "net/url"
//import "strings"
import "time"
import "math/rand"

func testbool() bool{
  ans := true
  return ans
}


func main(){
  rand.Seed(time.Now().UnixNano())
  fmt.Println("random number is", rand.Intn(14))
}
