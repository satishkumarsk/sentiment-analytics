package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Azure/azure-sdk-for-go/services/cognitiveservices/v2.1/textanalytics"
	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/to"
	"github.com/montanaflynn/stats"

	"sort"
	"strconv"
)

func main() {
	client := textanalytics.New("https://sathish.cognitiveservices.azure.com/")
	authorizer := autorest.NewCognitiveServicesAuthorizer("80bf268424b146ec8d844344a8ab7300")
	client.Authorizer = authorizer
	ctx := context.Background()
	idToMessage := []string{"This is a very good team. They take a lot of time to deliver & sometimes the quality is also low. But the team is good.", " test", "happy", "really happy", "he is bad", "this product is not so great", "what's the problem?", "Thank you I already have a penguin"}
	var inputDocuments []textanalytics.MultiLanguageInput
	var scores []float64

	for index, val := range idToMessage {

		inputDocuments = append(inputDocuments, textanalytics.MultiLanguageInput{
			Language: to.StringPtr("en"),
			ID:       to.StringPtr(string(index)),
			Text:     to.StringPtr(val),
		},
		)

	}

	result, err := client.Sentiment(ctx, to.BoolPtr(false), &textanalytics.MultiLanguageBatchInput{Documents: &inputDocuments})
	if err != nil {
		fmt.Print(err)
	}
	batchResult := textanalytics.SentimentBatchResult{}
	jsonString, _ := json.Marshal(result.Value)
	json.Unmarshal(jsonString, &batchResult)

	// Printing sentiment results
	for _, document := range *batchResult.Documents {
		fmt.Printf("Document ID: %s ", *document.ID)
		id, _ := strconv.Atoi(*document.ID)
		fmt.Print(id)
		fmt.Printf("Document Message: %s ", idToMessage[id])

		fmt.Printf("Sentiment Score: %f\n", *document.Score)
		scores = append(scores, *document.Score)
	}
	sort.Float64s(scores)
	percent, _ := stats.Percentile(stats.Float64Data(scores), 50)
	if percent > 0.5 {
		fmt.Println("poisitve person")
	} else {
		fmt.Println("Negative person")
	}

	// Printing document errors
	fmt.Println("Document Errors")
	for _, error := range *batchResult.Errors {
		fmt.Printf("Document ID: %s Message : %s\n", *error.ID, *error.Message)
	}

}
