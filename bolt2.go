package main
import (
"fmt"
"log"
"github.com.boltdb.bolt"
)
var world=[]byte["hello"]
func main(){
 db ,err:=bolt.Open(b1.db,0644,nil)
if err!=nil{
log.Fatal(err)}
defer db.Close()
key :=[]byte("hello")
value :=[]byte("hello")

err = db.Update(func(tx *bolt.Tx) error {
        bucket, err := tx.CreateBucketIfNotExists(world)
        if err != nil {
            return err
        }
err =Bucket.Put(key,value)


