package api

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/nordsyd/cli/core/manifest"
)

// FileUpload represents a file to upload
type FileUpload struct {
	Key       string `json:"key"`
	UploadURL string `json:"upload_url"`
}

// DeployResponse represents the response from the Nordsyd API
type DeployResponse struct {
	DeployHash string       `json:"deploy_hash"`
	ToUpload   []FileUpload `json:"to_upload"`
}

// CreateDeploy registers a deploy for a given site
func CreateDeploy(siteSlug string, manifest manifest.Manifest) DeployResponse {
	jsonManifest, _ := json.Marshal(manifest)

	payload := map[string]string{
		"manifest": string(jsonManifest),
	}

	response, error := Post("/site/"+siteSlug+"/deploy", payload)

	if error != nil {
		panic(error)
	}

	var deployResponse DeployResponse

	json.Unmarshal([]byte(response), &deployResponse)

	return deployResponse
}

// FinaliseDeploy signals the Nordsyd API that the deploy is ready to be activated
func FinaliseDeploy(siteSlug string, deployHash string) {
	_, error := Post("/site/"+siteSlug+"/deploy/"+deployHash+"/finish", map[string]string{})

	if error != nil {
		panic(error)
	}
}

// UploadFile uploads a file given key and upload URL
func UploadFile(key string, uploadURL string) {
	//fmt.Println("Uploading file: ", key)
	data, err := os.Open(key)
	if err != nil {
		log.Fatal(err)
	}
	defer data.Close()

	fileStat, _ := os.Stat(key)

	contentLength := fileStat.Size()

	//fmt.Println(uploadURL)

	req, err := http.NewRequest("PUT", uploadURL, data)

	// Set Content-Length to avoid "Transfer-Encoding: chunked"
	req.ContentLength = contentLength

	//fmt.Println(req.Header)

	if err != nil {
		//fmt.Println("error creating request", uploadURL)

		//fmt.Println(err)

		return
	}

	response, error := http.DefaultClient.Do(req)
	if error != nil {
		//fmt.Println("failed making request")

		//fmt.Println(error)

		return
	}

	_, uploadError := ioutil.ReadAll(response.Body)

	if uploadError != nil {
		panic(uploadError)
	}

	//fmt.Println(string(actualResponse))

	//fmt.Println("File uploaded successfully: ", key)
}
