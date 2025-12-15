package cardingester
import (
	"net/http"
	"io"
	"fmt"
)

func getBlob(url string, cache *Cache, client *http.Client) ([]byte, error) {
	if blob, inDb := cache.Get(url); inDb {
		return blob, nil
	}

	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}
	request.Header.Set("User-Agent", "cardingest/1.0")
	request.Header.Set("Accept", "application/json")

	res, err := client.Do(request)
	if err != nil {
		return nil, fmt.Errorf("error making request: %v", err)
	}
	defer res.Body.Close()
	blob, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("io error: %v", err)
	}

	cache.Add(url, blob)
	return blob, nil
}
