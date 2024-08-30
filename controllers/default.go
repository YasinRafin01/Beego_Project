package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/beego/beego/v2/core/logs"
	//beego "github.com/beego/beego/v2/core/logs"
	beego "github.com/beego/beego/v2/server/web"
)

type MainController struct {
	beego.Controller
}

type CatImage struct {
	ID     string `json:"id"`
	URL    string `json:"url"`
	Width  int    `json:"width"`
	Height int    `json:"height"`
}
type Breed struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Origin      string `json:"origin"`
}
type FavoriteData struct {
	ImageID string `json:"image_id"`
	SubID   string `json:"sub_id"`
}

type FavoriteResponse struct {
	ID string `json:"id"`
}
type VoteData struct {
	ImageID string `json:"image_id"`
	SubID   string `json:"sub_id"`
	Value   int    `json:"value"`
}

func (c *MainController) Get() {
	c.TplName = "index.tpl"
}

func getAPIKey() string {
	apiKey := beego.AppConfig.DefaultString("cat_api_key", "")
	if apiKey == "" {
		logs.Error("API key is not set in the configuration")
	}
	return apiKey
}

func makeAPIRequest(method, url string, body []byte, apiKey string) chan *http.Response {
	responseChan := make(chan *http.Response)
	go func() {
		req, _ := http.NewRequest(method, url, bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("x-api-key", apiKey)

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			logs.Error("Error making API request:", err)
			responseChan <- nil
		} else {
			responseChan <- resp
		}
	}()
	return responseChan
}

func (c *MainController) GetRandomCats() {
	count, _ := c.GetInt("count", 1)
	breedIds := c.GetString("breed_ids")
	apiKey := getAPIKey()

	url := fmt.Sprintf("https://api.thecatapi.com/v1/images/search?limit=%d", count)
	if breedIds != "" {
		url += fmt.Sprintf("&breed_ids=%s", breedIds)
	}

	responseChan := makeAPIRequest("GET", url, nil, apiKey)
	resp := <-responseChan

	if resp == nil {
		c.Data["json"] = map[string]string{"error": "Failed to make API request"}
		c.ServeJSON()
		return
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	var cats []CatImage
	json.Unmarshal(body, &cats)

	c.Data["json"] = cats
	c.ServeJSON()
}

func (c *MainController) GetCatBreeds() {
	apiKey := getAPIKey()
	url := "https://api.thecatapi.com/v1/breeds"

	responseChan := makeAPIRequest("GET", url, nil, apiKey)
	resp := <-responseChan

	if resp == nil {
		c.Data["json"] = map[string]string{"error": "Failed to make API request"}
		c.ServeJSON()
		return
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	var breeds []map[string]interface{}
	json.Unmarshal(body, &breeds)

	c.Data["json"] = breeds
	c.ServeJSON()
}

func (c *MainController) RecordVote() {
	var voteData VoteData
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &voteData); err != nil {
		c.Data["json"] = map[string]string{"error": "Invalid request body"}
		c.ServeJSON()
		return
	}

	apiKey := getAPIKey()
	url := "https://api.thecatapi.com/v1/votes"
	jsonData, _ := json.Marshal(voteData)

	responseChan := makeAPIRequest("POST", url, jsonData, apiKey)
	resp := <-responseChan

	if resp == nil {
		c.Data["json"] = map[string]string{"error": "Failed to make API request"}
		c.ServeJSON()
		return
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	c.Ctx.Output.Body(body)
}

func (c *MainController) GetVotes() {
	subID := c.GetString("sub_id")
	apiKey := getAPIKey()
	url := fmt.Sprintf("https://api.thecatapi.com/v1/votes?sub_id=%s", subID)

	responseChan := makeAPIRequest("GET", url, nil, apiKey)
	resp := <-responseChan

	if resp == nil {
		c.Data["json"] = map[string]string{"error": "Failed to make API request"}
		c.ServeJSON()
		return
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	c.Ctx.Output.Body(body)
}

func (c *MainController) GetConfig() {
	apiKey := getAPIKey()
	config := map[string]string{
		"catapi_key": apiKey,
	}
	c.Data["json"] = config
	c.ServeJSON()
}

func (c *MainController) AddFavorite() {
	var favoriteData FavoriteData
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &favoriteData); err != nil {
		c.Data["json"] = map[string]string{"error": "Invalid request body"}
		c.ServeJSON()
		return
	}

	apiKey := getAPIKey()
	url := "https://api.thecatapi.com/v1/favourites"
	jsonData, _ := json.Marshal(favoriteData)

	responseChan := makeAPIRequest("POST", url, jsonData, apiKey)
	resp := <-responseChan

	if resp == nil {
		c.Data["json"] = map[string]string{"error": "Failed to make API request"}
		c.ServeJSON()
		return
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		c.Data["json"] = map[string]interface{}{
			"error":  fmt.Sprintf("API Error: %s", string(body)),
			"status": resp.StatusCode,
		}
		c.ServeJSON()
		return
	}

	var favoriteResponse FavoriteResponse
	if err := json.Unmarshal(body, &favoriteResponse); err != nil {
		c.Data["json"] = map[string]string{"error": "Failed to parse API response"}
		c.ServeJSON()
		return
	}

	c.Data["json"] = favoriteResponse
	c.ServeJSON()
}

func (c *MainController) GetFavorites() {
	subID := c.GetString("sub_id")
	apiKey := getAPIKey()
	url := fmt.Sprintf("https://api.thecatapi.com/v1/favourites?sub_id=%s", subID)

	responseChan := makeAPIRequest("GET", url, nil, apiKey)
	resp := <-responseChan

	if resp == nil {
		c.Data["json"] = map[string]string{"error": "Failed to make API request"}
		c.ServeJSON()
		return
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		c.Data["json"] = map[string]interface{}{
			"error":  fmt.Sprintf("API Error: %s", string(body)),
			"status": resp.StatusCode,
		}
		c.ServeJSON()
		return
	}

	var favorites []map[string]interface{}
	if err := json.Unmarshal(body, &favorites); err != nil {
		c.Data["json"] = []map[string]interface{}{} // Return empty array if parsing fails
	} else {
		c.Data["json"] = favorites
	}
	c.ServeJSON()
}

func (c *MainController) DeleteFavorite() {
	favoriteID := c.Ctx.Input.Param(":favoriteId")
	apiKey := getAPIKey()
	url := fmt.Sprintf("https://api.thecatapi.com/v1/favourites/%s", favoriteID)

	responseChan := makeAPIRequest("DELETE", url, nil, apiKey)
	resp := <-responseChan

	if resp == nil {
		c.Data["json"] = map[string]string{"error": "Failed to make API request"}
		c.ServeJSON()
		return
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	logs.Info("Response status:", resp.Status)
	logs.Info("Response body:", string(body))

	if resp.StatusCode != http.StatusOK {
		c.Data["json"] = map[string]interface{}{
			"error":  fmt.Sprintf("API Error: %s", string(body)),
			"status": resp.StatusCode,
		}
		c.ServeJSON()
		return
	}

	c.Data["json"] = map[string]string{"message": "Favorite deleted successfully"}
	c.ServeJSON()
}
