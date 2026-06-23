package main

import (
	"fmt"
	"io"
	"net/http"
	"time"
)

func main() {
	consumerKey := "8Iec2PAnf2pSKDm0P12gsJbvqU4Fg7E1iKfslWVkoBNNqSh3"
	consumerSecret := "i80F5MlglXT20JSwBgJKcChXv0XaEm1yCCRSjCCATHTvru5nIu8nAHCDHh3hd7cm"

	fmt.Println("Testing M-Pesa API Connection (Production)...")
	fmt.Println("Consumer Key:", consumerKey[:10]+"...")
	fmt.Println("Consumer Secret:", consumerSecret[:10]+"...")

	// Production URL
	url := "https://api.safaricom.co.ke/oauth/v1/generate?grant_type=client_credentials"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}

	req.SetBasicAuth(consumerKey, consumerSecret)

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error making request:", err)
		return
	}
	defer resp.Body.Close()

	fmt.Println("Response Status Code:", resp.StatusCode)
	fmt.Println("Response Status:", resp.Status)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading body:", err)
		return
	}

	fmt.Println("Raw Response Body:", string(body))
}
