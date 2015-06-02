package bolt

import (
    "log"
    "fmt"
    "github.com/boltdb/bolt"
)
func main(){
db, err :=bolt.Open("dfm.db",0600,nil)
	if err!=nil{
	log.Fatal(err)
}
defer db.Close()

db.Update(func(tx *bolt.Tx) error {
    b, err := tx.CreateBucketIfNotExists([]byte("AdServer"))
    if err != nil {
        return err
    }
    b.Put([]byte("abc"), []byte("67"))
    b.Put([]byte("abc1"), []byte("97.8"))
    b.Put([]byte("abc2"), []byte("85.89"))
 return nil
})

db.Update(func(tx *bolt.Tx) error {
    b, err := tx.CreateBucketIfNotExists([]byte("AdServer"))
    if err != nil {
        return err
    }
  return b.Put([]byte("abc3"), []byte("68"))

})
err = db.View(func(tx *bolt.Tx) error {
      j:=tx.Bucket([]byte("AdServer"))
      c := j.Cursor()
    	for k, n := c.First(); k != nil; k, n = c.Next() {
        	  fmt.Printf("subline_id  %s bids %s.\n", k, n)
    		}
return nil
})
}