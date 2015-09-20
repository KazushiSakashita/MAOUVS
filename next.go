package main
​
import "github.com/ChimeraCoder/anaconda"
import "fmt"
import "strconv"
import "github.com/toqueteos/webbrowser"
​
func main() {
	​anaconda.SetConsumerKey("j62kyFWBVqCiJElWSHvofSz59")
	anaconda.SetConsumerSecret("FmA9RF7MAEwLbvXfvYRgZMMstinToo6kbR7CcwfpdhWrsEPyrg")

	// webbrowser.Open()
	url, cred, _ := anaconda.AuthorizationURL("")
	webbrowser.Open(url)
​
	var pincode int64
	fmt.Scan(&pincode)
​
	newcred, _, _ := anaconda.GetCredentials(cred, strconv.FormatInt(pincode, 10))
​
	fmt.Printf("%v\n%v\n", newcred.Token, newcred.Secret)
}
