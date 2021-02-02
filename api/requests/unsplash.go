package requests

import (
	"encoding/json"
	"net/http"
)

func ReadRespBody(url string, t interface{}) {
	res, _ := http.Get(url)
	defer res.Body.Close()
	json.NewDecoder(res.Body).Decode(&t)
}
