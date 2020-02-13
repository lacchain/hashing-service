/*
	Hashing Service
	version 0.9
	author: Adrian Pareja Abarca
	email: adriancc5.5@gmail.com
*/

package model

type Metadata struct{
	Title string `json:"title"`
	Name  string `json:"name"`
	Organization string `json:"organization"`
	Author string `json:"author"`
	Document string `json:"hash"`
	ExpirationDate string `json:"expirationDate"` 
}