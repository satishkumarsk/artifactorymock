package main

import (
	"fmt"

	"github.com/jfrog/jfrog-client-go/artifactory"
	"github.com/jfrog/jfrog-client-go/artifactory/auth"
	"github.com/jfrog/jfrog-client-go/artifactory/services"
	"github.com/jfrog/jfrog-client-go/artifactory/services/utils"
	"github.com/jfrog/jfrog-client-go/config"
)

func searchArtifacts(aurl string, port string) error {
	rtDetails := auth.NewArtifactoryDetails()
	rtDetails.SetUrl(fmt.Sprintf("http://%s:%s/%s/", aurl, port, "artifactory"))
	rtDetails.SetUser("")
	rtDetails.SetPassword("")
	serviceConfig := config.NewConfigBuilder()
	serviceConfig.SetServiceDetails(rtDetails)
	sconf, err := serviceConfig.Build()
	if err != nil {
		return err
	}
	sManager, err := artifactory.New(sconf)
	if err != nil {
		return err
	}
	params := services.NewSearchParams()
	params.Pattern = "repo/*/*.zip"
	// Filter the files by properties.
	params.Props = "key1=val1;key2=val2"
	params.Recursive = true
	reader, err := sManager.SearchFiles(params)
	if err != nil {
		return err
	}
	// Iterate over the results.
	for currentResult := new(utils.ResultItem); reader.NextRecord(currentResult) == nil; currentResult = new(utils.ResultItem) {
		fmt.Printf("Found artifact: %s of type: %s\n", currentResult.Name, currentResult.Type)
	}
	if err := reader.GetError(); err != nil {
		return err
	}

	// Resets the reader pointer back to the beginning of the output. Make sure not to call this method after the reader had been closed using ```reader.Close()```
	reader.Reset()
	return nil

}

func main() {

	err := searchArtifacts("127.0.0.1", "8080")
	if err != nil {
		panic(err)
	}
}
