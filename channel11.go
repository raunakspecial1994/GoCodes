package main
import "fmt"
func add(a []int,c chan int){
sum:=0
for _, v:=range a{
sum+=v
}
c<-sum
}
func main(){
a :=[]int{3,5,-5,3,-9,2}
c:=make(chan int)
go add(a[:len(a)/2],c)
go add(a[len(a)/2:],c)
x , y :=<-c,<-c
fmt.Println(x,y,x+y)
}

