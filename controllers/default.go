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

func (c *MainController) GetRandomCats() {
	count, _ := c.GetInt("count", 1)
	breedIds := c.GetString("breed_ids")
	apiKey := beego.AppConfig.DefaultString("cat_api_key", "")
	catsChan := make(chan []CatImage)

	go func() {
		url := fmt.Sprintf("https://api.thecatapi.com/v1/images/search?limit=%d", count)
		if breedIds != "" {
			url += fmt.Sprintf("&breed_ids=%s", breedIds)
		}
		req, _ := http.NewRequest("GET", url, nil)
		req.Header.Add("x-api-key", apiKey)

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			c.Data["json"] = map[string]string{"error": err.Error()}
			c.ServeJSON()
			return
		}
		defer resp.Body.Close()

		body, _ := ioutil.ReadAll(resp.Body)
		var cats []CatImage
		json.Unmarshal(body, &cats)

		catsChan <- cats
	}()

	cats := <-catsChan
	c.Data["json"] = cats
	c.ServeJSON()
}

func (c *MainController) GetCatBreeds() {
	apiKey := beego.AppConfig.DefaultString("cat_api_key", "")
	breedsChan := make(chan []map[string]interface{})

	go func() {
		url := "https://api.thecatapi.com/v1/breeds"
		req, _ := http.NewRequest("GET", url, nil)
		req.Header.Add("x-api-key", apiKey)

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			c.Data["json"] = map[string]string{"error": err.Error()}
			c.ServeJSON()
			return
		}
		defer resp.Body.Close()

		body, _ := ioutil.ReadAll(resp.Body)
		var breeds []map[string]interface{}
		json.Unmarshal(body, &breeds)

		breedsChan <- breeds
	}()

	breeds := <-breedsChan
	c.Data["json"] = breeds
	c.ServeJSON()
}

type VoteData struct {
	ImageID string `json:"image_id"`
	SubID   string `json:"sub_id"`
	Value   int    `json:"value"`
}

func (c *MainController) RecordVote() {
	var voteData VoteData
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &voteData); err != nil {
		c.Data["json"] = map[string]string{"error": "Invalid request body"}
		c.ServeJSON()
		return
	}

	apiKey, _ := beego.AppConfig.String("cat_api_key")
	url := "https://api.thecatapi.com/v1/votes"

	jsonData, _ := json.Marshal(voteData)
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", apiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		c.Data["json"] = map[string]string{"error": err.Error()}
		c.ServeJSON()
		return
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	c.Ctx.Output.Body(body)
}

func (c *MainController) GetVotes() {
	subID := c.GetString("sub_id")
	apiKey, _ := beego.AppConfig.String("cat_api_key")
	url := fmt.Sprintf("https://api.thecatapi.com/v1/votes?sub_id=%s", subID)

	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("x-api-key", apiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		c.Data["json"] = map[string]string{"error": err.Error()}
		c.ServeJSON()
		return
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	c.Ctx.Output.Body(body)
}
func (c *MainController) GetConfig() {
	apiKey, _ := beego.AppConfig.String("cat_api_key")

	config := map[string]string{
		"catapi_key": apiKey,
	}
	c.Data["json"] = config
	c.ServeJSON()

}

// New function to add a favorite
func (c *MainController) AddFavorite() {
	var favoriteData FavoriteData
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &favoriteData); err != nil {
		c.Data["json"] = map[string]string{"error": "Invalid request body"}
		c.ServeJSON()
		return
	}

	apiKey := getAPIKey()
	if apiKey == "" {
		c.Data["json"] = map[string]string{"error": "API key is not configured"}
		c.ServeJSON()
		return
	}

	url := "https://api.thecatapi.com/v1/favourites"
	jsonData, _ := json.Marshal(favoriteData)

	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", apiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		c.Data["json"] = map[string]string{"error": err.Error()}
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
	if apiKey == "" {
		c.Data["json"] = map[string]string{"error": "API key is not configured"}
		c.ServeJSON()
		return
	}
	fmt.Println(apiKey)
	url := fmt.Sprintf("https://api.thecatapi.com/v1/favourites?sub_id=%s", subID)

	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("x-api-key", apiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		c.Data["json"] = map[string]string{"error": err.Error()}
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

// New function to delete a favorite
func (c *MainController) DeleteFavorite() {
	favoriteID := c.Ctx.Input.Param(":favoriteId")
	//apiKey := beego.AppConfig.DefaultString("cat_api_key", "")
	url := fmt.Sprintf("https://api.thecatapi.com/v1/favourites/%s", favoriteID)

	req, _ := http.NewRequest("DELETE", url, nil)
	req.Header.Set("Content-Type", "application/json")
	//setAPIKeyHeader(req)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		c.Data["json"] = map[string]string{"error": err.Error()}
		c.ServeJSON()
		return
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		var errorResponse map[string]interface{}
		json.Unmarshal(body, &errorResponse)
		c.Data["json"] = map[string]interface{}{
			"error":  fmt.Sprintf("API Error: %v", errorResponse),
			"status": resp.StatusCode,
		}
		c.ServeJSON()
		return
	}

	c.Data["json"] = map[string]string{"message": "Favorite deleted successfully"}
	c.ServeJSON()
}

