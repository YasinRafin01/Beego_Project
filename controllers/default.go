package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

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

func (c *MainController) Get() {
	c.TplName = "index.tpl"
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

	apiKey := beego.AppConfig.DefaultString("cat_api_key", "")
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
	apiKey := beego.AppConfig.DefaultString("cat_api_key", "")
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
