package main
import (
	"fmt"
	"net/http"
	"io/ioutil"
	"time"
)

func main() {

	begin:= time.Now()
	resp, err := http.Get("http://www.01happy.com/demo/accept.php?id=1")
	if err != nil {
		fmt.Println( "http get failed")
	}
	fmt.Println( fmt.Sprintf("%s  %s","Response time: ",time.Since(begin).String()))

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		// handle error
	}

	fmt.Println(string(body))
}