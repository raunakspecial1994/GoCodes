//maps 
package main
import "fmt"
func main(){
l:=make(map[string]string)
l["bandName"]="Pink floyd"
web:=l["bandName"]+""+"Music"
fmt.Println(web)
delete(l,"bandName")
fmt.Println(web)
fmt.Println(l["bandName"])

j := map[string]int{
      "raunak":88,
       "surekha":77,
        "rahul":78,
       "Sunil":86,
}
//random print
for key, value := range j {
fmt.Println("Key",key,"Value",value)
}
var n int
n=len(j)
fmt.Println(n)
}









 
