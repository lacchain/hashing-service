/*
	Hashing Service
	version 0.9
	author: Adrian Pareja Abarca
	email: adriancc5.5@gmail.com
*/
package lib

import(
	"crypto/sha256"
	"io"
	"log"
	"encoding/hex"
)

func Hash(file io.Reader)(string){
	hasher := sha256.New()
	if _, err := io.Copy(hasher, file); err != nil {
		log.Fatal(err)
	}
	//fmt.Printf("%x", hasher.Sum(nil))
	return hex.EncodeToString(hasher.Sum(nil))
}