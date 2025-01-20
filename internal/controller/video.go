package controller

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/gorilla/schema"
	"github.com/northmule/pinetree/internal/service"
	"golang.org/x/net/context"
)

const (
	VKGetVideoURL  = "https://api.vk.com/method/video.get"
	VKEditVideoURL = "https://api.vk.com/method/video.edit"
)

// Video контроллер
type Video struct {
	logger service.LogPusher
	client RequestSender
}

type RequestSender interface {
	Do(req *http.Request) (*http.Response, error)
}

// NewVideo конструктор
func NewVideo(client RequestSender, logger service.LogPusher) *Video {
	return &Video{
		client: client,
		logger: logger,
	}
}

type RequestGetAll struct {
	AccessToken string `schema:"access_token"`
	OwnerId     string `schema:"owner_id"`
	AlbumId     string `schema:"album_id"`
	Offset      int    `schema:"offset"`
	Count       int    `schema:"count"`
	Version     string `schema:"v"`
}

type ResponseGetAll struct {
	Response *Response      `json:"response"`
	Error    *ResponseError `json:"error"`
}

type ResponseError struct {
	ErrorCode     int    `json:"error_code"`
	ErrorMsg      string `json:"error_msg"`
	RequestParams []struct {
		Key   string `json:"key"`
		Value string `json:"value"`
	} `json:"request_params"`
}

// Response представляет структуру ответа.
type Response struct {
	// Count всего видео, не зависимо от смещения
	Count int    `json:"count"`
	Items []Item `json:"items"`
}

// Image представляет структуру изображения.
type Image struct {
	URL         string `json:"url"`
	Width       int    `json:"width"`
	Height      int    `json:"height"`
	WithPadding int    `json:"with_padding,omitempty"` // Поле может быть опущено в JSON.
}

// FirstFrame представляет структуру первого кадра видео.
type FirstFrame []Image

// Item представляет структуру одного элемента в списке.
type Item struct {
	AddingDate    int64      `json:"adding_date"`
	CanComment    int        `json:"can_comment"`
	CanLike       int        `json:"can_like"`
	CanRepost     int        `json:"can_repost"`
	CanSubscribe  int        `json:"can_subscribe"`
	CanAddToFaves int        `json:"can_add_to_faves"`
	CanAdd        int        `json:"can_add"`
	Comments      int        `json:"comments"`
	Date          int64      `json:"date"`
	Description   string     `json:"description"`
	Duration      int        `json:"duration"`
	Image         []Image    `json:"image"`
	FirstFrame    FirstFrame `json:"first_frame"`
	Width         int        `json:"width"`
	Height        int        `json:"height"`
	// ID уникальный ИД видео
	ID         int64  `json:"id"`
	OwnerID    int64  `json:"owner_id"`
	OV_ID      string `json:"ov_id"`
	Title      string `json:"title"`
	IsFavorite bool   `json:"is_favorite"`
	Player     string `json:"player"`
	Added      int    `json:"added"`
	Repeat     int    `json:"repeat"`
	Type       string `json:"type"`
	Views      int64  `json:"views"`
	Likes      struct {
		Count     int `json:"count"`
		UserLikes int `json:"user_likes"`
	} `json:"likes"`
	Reposts struct {
		Count        int `json:"count"`
		UserReposted int `json:"user_reposted"`
	} `json:"reposts"`
}

// GetAll все записи видео
func (c *Video) GetAll(requestData RequestGetAll) (*ResponseGetAll, error) {
	ctx := context.Background()
	var encoder = schema.NewEncoder()
	form := url.Values{}
	err := encoder.Encode(requestData, form)
	if err != nil {
		return nil, err
	}

	requestPrepare, err := http.NewRequestWithContext(ctx, http.MethodGet, VKGetVideoURL+"?"+form.Encode(), nil)
	if err != nil {
		c.logger.Error(err)
		return nil, err
	}

	response, err := c.client.Do(requestPrepare)
	if err != nil {
		c.logger.Error(err)
		return nil, err
	}

	bodyRaw, err := io.ReadAll(response.Body)
	if err != nil {
		c.logger.Error(err)
		return nil, err
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		c.logger.Warn(response.Status)
		return nil, fmt.Errorf("%s: %s", response.Status, bodyRaw)
	}
	responseData := new(ResponseGetAll)
	err = json.Unmarshal(bodyRaw, responseData)
	if err != nil {
		c.logger.Error(err)
		return nil, err
	}

	// items := responseData.Response.Items
	items := make([]Item, 0, responseData.Response.Count)
	items = append(items, responseData.Response.Items...)
	if len(items) < responseData.Response.Count {
		requestData.Offset = requestData.Count + requestData.Offset
		items, err = c.fillAllItems(requestData, items)
		if err != nil {
			return nil, err
		}
	}
	responseData.Response.Items = items
	return responseData, nil
}

func (c *Video) fillAllItems(requestData RequestGetAll, items []Item) ([]Item, error) {
	ctx := context.Background()
	var encoder = schema.NewEncoder()
	form := url.Values{}
	err := encoder.Encode(requestData, form)
	if err != nil {
		return nil, err
	}

	requestPrepare, err := http.NewRequestWithContext(ctx, http.MethodGet, VKGetVideoURL+"?"+form.Encode(), nil)
	if err != nil {
		c.logger.Error(err)
		return nil, err
	}

	response, err := c.client.Do(requestPrepare)
	if err != nil {
		c.logger.Error(err)
		return nil, err
	}

	bodyRaw, err := io.ReadAll(response.Body)
	if err != nil {
		c.logger.Error(err)
		return nil, err
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		c.logger.Warn(response.Status)
		return nil, fmt.Errorf("%s: %s", response.Status, bodyRaw)
	}
	responseData := new(ResponseGetAll)
	err = json.Unmarshal(bodyRaw, responseData)
	if err != nil {
		c.logger.Error(err)
		return nil, err
	}

	items = append(items, responseData.Response.Items...)
	if len(items) < responseData.Response.Count {
		requestData.Offset = requestData.Offset + requestData.Count
		return c.fillAllItems(requestData, items)
	}

	return items, nil
}

type RequestUpdateDescription struct {
	AccessToken string `schema:"access_token"`
	OwnerId     string `schema:"owner_id"`
	VideoId     int64  `schema:"video_id"`
	Desc        string `schema:"desc"`
	Version     string `schema:"v"`
}

type ResponseUpdateDescription struct {
	Response struct {
		Success   int    `json:"success"`
		AccessKey string `json:"access_key"`
	} `json:"response"`
}

// UpdateDescription обновляет описание для видео
func (c *Video) UpdateDescription(requestData RequestUpdateDescription) (*ResponseUpdateDescription, error) {
	ctx := context.Background()
	var encoder = schema.NewEncoder()
	form := url.Values{}
	err := encoder.Encode(requestData, form)
	if err != nil {
		return nil, err
	}

	requestPrepare, err := http.NewRequestWithContext(ctx, http.MethodGet, VKEditVideoURL+"?"+form.Encode(), nil)
	if err != nil {
		c.logger.Error(err)
		return nil, err
	}

	response, err := c.client.Do(requestPrepare)
	if err != nil {
		c.logger.Error(err)
		return nil, err
	}

	bodyRaw, err := io.ReadAll(response.Body)
	if err != nil {
		c.logger.Error(err)
		return nil, err
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		c.logger.Warn(response.Status)
		return nil, fmt.Errorf("%s: %s", response.Status, bodyRaw)
	}
	responseData := new(ResponseUpdateDescription)
	err = json.Unmarshal(bodyRaw, responseData)
	if err != nil {
		c.logger.Error(err)
		return nil, err
	}

	return responseData, nil
}
