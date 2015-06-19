package main

import (
	"encoding/binary"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"os"
	"time"

	"github.com/boltdb/bolt"
)

// Set this to see how the counts are actually updated.
const verbose = false

// Counter updates a counter in Bolt for every URL path requested.
type counter struct {
	db *bolt.DB
}

func (c counter) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	// Communicates the new count from a successful database
	// transaction.
	var result uint64

	increment := func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte("hits"))
		if err != nil {
			return err
		}
		key := []byte(req.URL.String())
		// Decode handles key not found for us.
		count := decode(b.Get(key)) + 1
		b.Put(key, encode(count))
		// All good, communicate new count.
		result = count
		return nil
	}
	if err := c.db.Batch(increment); err != nil {
		http.Error(rw, err.Error(), 500)
		return
	}

	if verbose {
		log.Printf("server: %s: %d", req.URL.String(), result)
	}

	rw.Header().Set("Content-Type", "application/octet-stream")
	fmt.Fprintf(rw, "%d\n", result)
}

func client(id int, base string, paths []string) error {
	// Process paths in random order.
	rng := rand.New(rand.NewSource(int64(id)))
	permutation := rng.Perm(len(paths))

	for i := range paths {
		path := paths[permutation[i]]
		resp, err := http.Get(base + path)
		if err != nil {
			return err
		}
		defer resp.Body.Close()
		buf, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		if verbose {
			log.Printf("client: %s: %s", path, buf)
		}
	}
	return nil
}

func main() {
	// Open the database.
	db, _ := bolt.Open("asd.db", 0666, nil)
	defer os.Remove(db.Path())
	defer db.Close()

	// Start our web server
	count := counter{db}
	srv := httptest.NewServer(count)
	defer srv.Close()

	// Get every path multiple times concurrently.
	const clients = 20
	paths := []string{
		"/foo",
		"/bar",
		"/baz",
		"/quux",
		"/thud",
		"/xyzzy",
	}
startTime := time.Now()
	errors := make(chan error, clients)
	for i := 0; i < clients; i++ {
		go func(id int) {
			errors <- client(id, srv.URL, paths)
		}(i)
	}
	// Check all responses to make sure there's no error.
	for i := 0; i < clients; i++ {
		if err := <-errors; err != nil {
			fmt.Printf("client error: %v", err)
			return
		}
		
	}
log.Printf("it only took %.3f seconds!", time.Since(startTime).Seconds())


	// Check the final result
	db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("hits"))
		c := b.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			fmt.Printf("bid for %s: %d\n", k, decode(v))
		}
		return nil
	})

	// Output:
	// bids for /bar: 20
	// bids for /baz: 20
	// bids for /foo: 20
	// bids for /quux: 20
	// bids for /thud: 20
	// bids for /xyzzy: 20
}

// encode marshals a counter.
func encode(n uint64) []byte {
	buf := make([]byte, 8)
	binary.BigEndian.PutUint64(buf, n)
	return buf
}

// decode unmarshals a counter. Nil buffers are decoded as 0.
func decode(buf []byte) uint64 {
	if buf == nil {
		return 0
	}
	return binary.BigEndian.Uint64(buf)
}
