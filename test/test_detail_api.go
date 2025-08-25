// +build ignore

package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/pyhub-apps/sejong-cli/internal/api"
)

func main() {
	apiKey := os.Getenv("SEJONG_API_KEY")
	if apiKey == "" {
		fmt.Println("SEJONG_API_KEY environment variable not set")
		fmt.Println("Please set it with: export SEJONG_API_KEY=your_api_key")
		os.Exit(1)
	}

	// Create NLIC client
	client := api.NewNLICClient(apiKey)

	// Test with a known law ID
	lawID := "011357" // 개인정보 보호법
	fmt.Printf("Testing GetDetail with law ID: %s\n\n", lawID)

	ctx := context.Background()
	detail, err := client.GetDetail(ctx, lawID)
	if err != nil {
		log.Fatalf("Error getting law detail: %v", err)
	}

	// Display basic info
	fmt.Println("=== Basic Law Information ===")
	fmt.Printf("Law ID: %s\n", detail.LawInfo.ID)
	fmt.Printf("Law Name: %s\n", detail.LawInfo.Name)
	fmt.Printf("Department: %s\n", detail.LawInfo.Department)
	fmt.Printf("Law Type: %s\n", detail.LawInfo.LawType)
	fmt.Printf("Promulgation Date: %s\n", detail.LawInfo.PromulDate)
	fmt.Printf("Promulgation Number: %s\n", detail.LawInfo.PromulNo)
	fmt.Printf("Effective Date: %s\n", detail.LawInfo.EffectDate)
	fmt.Printf("Revision Type: %s\n", detail.LawInfo.Category)
	fmt.Printf("Serial Number: %s\n", detail.LawInfo.SerialNo)

	// Display articles summary
	fmt.Printf("\n=== Articles Summary ===\n")
	fmt.Printf("Total Articles: %d\n", len(detail.Articles))
	if len(detail.Articles) > 0 {
		fmt.Println("\nFirst 3 articles:")
		for i, article := range detail.Articles {
			if i >= 3 {
				break
			}
			fmt.Printf("\n%s %s\n", article.Number, article.Title)
			// Truncate content for display
			content := article.Content
			if len(content) > 100 {
				content = content[:100] + "..."
			}
			fmt.Printf("  Content: %s\n", content)
			if article.EffectDate != "" {
				fmt.Printf("  Effective Date: %s\n", article.EffectDate)
			}
		}
	}

	// Optional: Output full JSON for debugging
	if os.Getenv("DEBUG") == "1" {
		fmt.Println("\n=== Full JSON Output ===")
		jsonData, _ := json.MarshalIndent(detail, "", "  ")
		fmt.Println(string(jsonData))
	}
}