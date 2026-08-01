package backend

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	geoloc "github.com/IPGeolocation/ip-geolocation-go-sdk/sdk"
	"github.com/joho/godotenv"
)

// https://app.ipgeolocation.io/api
func AnalyzeIpAddress() {
	godotenv.Load(".env")
	ctx := geoloc.WithAPIKey(context.Background(), os.Getenv("GEO_LOCATION_API_KEY"))
	config := geoloc.NewConfiguration()
	client := geoloc.NewAPIClient(config)
	a, b, err := client.IPGeolocationAPI.GetIpGeolocation(ctx).Ip("91.165.54.16").Execute()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	responsJson, err := json.MarshalIndent(a, "", " ")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Print(string(responsJson), b.Status)
}
