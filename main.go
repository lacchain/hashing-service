/*
	Hashing Service
	version 0.9
	author: Adrian Pareja Abarca
	email: adriancc5.5@gmail.com
*/

package main

import (
    "fmt"
    "log"
    "bytes"
    "time"
    "os"
    "io/ioutil"
    "net/http"
    "github.com/gorilla/handlers"
    "github.com/gorilla/mux"
    "encoding/json"
    "github.com/go-openapi/strfmt"
    lib "github.com/lacchain/hashing-service/lib"
    model "github.com/lacchain/hashing-service/model"
)

func main(){
    router := mux.NewRouter()
    router.HandleFunc("/upload", uploadFile).Methods("POST","OPTIONS")
    router.HandleFunc("/validate", validateHash).Methods("POST","OPTIONS")
    
    headersOk := handlers.AllowedHeaders([]string{"Accept","Accept-Encoding","Accept-Language","Authorization","X-Requested-With","Content-Type","Origin","Access-Control-Request-Headers","Access-Control-Request-Method"})
    originsOk := handlers.AllowedOrigins([]string{"*"})
    methodsOk := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS"})

    err := http.ListenAndServeTLS(":9000", "/home/adrianpareja/certificate/server.crt", "/home/adrianpareja/certificate/server.key", handlers.CORS(originsOk, headersOk, methodsOk)(router))

    if err != nil {
        log.Fatal("ListenAndServe: ", err)
    }
}

func validateHash(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")

    // Parse our multipart form, 10 << 20 specifies a maximum
    // upload of 10 MB files.
    parseErr := r.ParseMultipartForm(10 << 20)
    if parseErr != nil{
        fmt.Println("error:",parseErr)
        http.Error(w, "failed to parse multipart message", http.StatusBadRequest)
        return
    }
    
    file, handler, err := r.FormFile("media")
    if err != nil {
        fmt.Println("Error Retrieving the File")
        fmt.Println(err)
        return
    }
    defer file.Close()
    fmt.Printf("Uploaded File: %+v\n", handler.Filename)
    fmt.Printf("File Size: %+v\n", handler.Size)
    fmt.Printf("MIME Header: %+v\n", handler.Header)

    hash := lib.Hash(file)

    response:="{\"hash\":\""+hash+"\"}"

    _, err = file.Seek(0, os.SEEK_SET)
    if err != nil {
        fmt.Println(err)
    }

    w.Write([]byte(response)) 
}

func uploadFile(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")

    // Parse our multipart form, 10 << 20 specifies a maximum
    // upload of 10 MB files.
    parseErr := r.ParseMultipartForm(10 << 20)
    if parseErr != nil{
        fmt.Println("error:",parseErr)
        http.Error(w, "failed to parse multipart message", http.StatusBadRequest)
        return
    }
    
    file, handler, err := r.FormFile("media")
    if err != nil {
        fmt.Println("Error Retrieving the File")
        fmt.Println(err)
        return
    }
    defer file.Close()
    fmt.Printf("Uploaded File: %+v\n", handler.Filename)
    fmt.Printf("File Size: %+v\n", handler.Size)
    fmt.Printf("MIME Header: %+v\n", handler.Header)

    metadata, errMeta := getMetadata(r)
	if errMeta != nil {
		http.Error(w, "failed to get metadata", http.StatusBadRequest)
		return
	}
	log.Println("Metadata:",string(metadata))

    contactInformation, errContact := getContactInformation(r)
	if errContact != nil {
		http.Error(w, "failed to get contact information", http.StatusBadRequest)
		return
	}
	log.Println("ContactInformation:",string(contactInformation))

    res := model.Metadata{}
    json.Unmarshal(metadata, &res)

    contact := model.Contact{}
    json.Unmarshal(contactInformation, &contact)
    
    res.Document = lib.Hash(file)

    _, err = file.Seek(0, os.SEEK_SET)
    if err != nil {
        fmt.Println(err)
    }

    fmt.Println("metadata json:",res)

    responseCredential := createCredential(&res,&contact) 

    fmt.Println("responseCredential:",responseCredential)

    w.Write([]byte(responseCredential)) 
}

func getMetadata(r *http.Request) ([]byte, error) {
	f, _, err := r.FormFile("metadata")
	if err != nil {
		return nil, fmt.Errorf("failed to get metadata form file: %v", err)
	}

	metadata, errRead := ioutil.ReadAll(f)
	if errRead != nil {
		return nil, fmt.Errorf("failed to read metadata: %v", errRead)
	}

	return metadata, nil
}

func getContactInformation(r *http.Request) ([]byte, error) {
	f, _, err := r.FormFile("contactInformation")
	if err != nil {
		return nil, fmt.Errorf("failed to get contact information form file: %v", err)
	}

	contactInformation, errRead := ioutil.ReadAll(f)
	if errRead != nil {
		return nil, fmt.Errorf("failed to read contactInformation: %v", errRead)
	}

	return contactInformation, nil
}

func createCredential(metadata *model.Metadata, contact *model.Contact)(string){
    credentials := make([]*model.CredentialSubject, 0, 50)
    credentialSubject := model.CredentialSubject{}
    credentialSubject.Type = "DocumentHashCredential"
    credentialSubject.Email = contact.Email
    
    fmt.Println("ffff:",time.Now().UTC().Format("2006-01-02T15:04:05Z"))
    issuanceDate, err := strfmt.ParseDateTime(time.Now().UTC().Format("2006-01-02T15:04:05Z"))
    if err != nil{
        fmt.Println("Error:",err)
    }

    expirationTime, err := time.Parse("2006-01-02T15:04:05Z", metadata.ExpirationDate)
    if err != nil {
        fmt.Println(err)
    }

    fmt.Println("expirationTime:",expirationTime.UTC().Format("2006-01-02T15:04:05Z"))

    expirationDate, err := strfmt.ParseDateTime(expirationTime.UTC().Format("2006-01-02T15:04:05Z"))
    if err != nil{
        fmt.Println("Error:",err)
    }

    credentialSubject.IssuanceDate = issuanceDate
    credentialSubject.ExpirationDate = expirationDate

    fmt.Println("credentialSubject:",credentialSubject.IssuanceDate)    

    credentialSubject.Content = &metadata
    credentials = append(credentials, &credentialSubject)
    jsonValue, _ := json.Marshal(credentials)
    fmt.Println("#####REQUEST####", string(jsonValue))
    
    timeout := time.Duration(15 * time.Second)
    client := http.Client{
        Timeout: timeout,
    }

    req, err := http.NewRequest("POST", "http://localhost:8000/v1/credential",  bytes.NewBuffer(jsonValue))
    req.Header.Set("Content-type","application/json")
    req.Header.Set("accept","application/json")

    response, err := client.Do(req)

    if err != nil {
        fmt.Printf("The HTTP request failed with error %s\n", err)
    } else {
        data, _ := ioutil.ReadAll(response.Body)
        fmt.Println(string(data))
        return string(data)
    }
    return "{}"
}