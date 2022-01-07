package usecase

import (
	"ELKExample/models"
	"context"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

func (h Handler) CreatePost(c *gin.Context) {
	post := &models.Posts{}

	body, _ := ioutil.ReadAll(c.Request.Body)
	if err := json.Unmarshal(body, post); err != nil {
		h.Logger.Err(err).Msg("could not parse request body")
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("invalid request body: %s", err.Error())})
		return
	}

	err := h.DB.SavePost(post)

	if err != nil {
		h.Logger.Err(err).Msg("could not save post")
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("could not save post: %s", err.Error())})
	} else {
		c.JSON(http.StatusCreated, gin.H{"post": post})
	}
}

func (h *Handler) UpdatePosts(c *gin.Context) {
	var id int
	var post models.Posts
	var err error

	if id, err = strconv.Atoi(c.Param("id")); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid post id"})
		return
	}

	if err = c.ShouldBindJSON(&post); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("could not parse request: %s", err.Error())})
		return
	}

	post.ID = id

	err = h.DB.UpdatePost(&post)

	if err != nil {
		h.Logger.Err(err).Msg("could not update post")
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("could not update post: %s", err.Error())})
		return
	} else {
		c.JSON(http.StatusOK, gin.H{"post": post})
	}
}

func (h *Handler) DeletePost(c *gin.Context) {
	var id int
	var err error
	if id, err = strconv.Atoi(c.Param("id")); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid post id"})
		return
	}
	err = h.DB.DeletePost(id)

	if err != nil {
		h.Logger.Err(err).Msg("could not update post")
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("%s", err.Error())})
		return
	} else {
		c.JSON(http.StatusOK, gin.H{"status": "OK"})
	}
}

func (h *Handler) GetPosts(c *gin.Context) {
	posts, err := h.DB.GetPost()
	if err != nil {
		h.Logger.Err(err).Msg("Could not fetch posts")
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	} else {
		c.JSON(http.StatusOK, gin.H{"data": posts})
	}
}

func (h *Handler) SearchPosts(c *gin.Context) {
	var query string

	if query, _ = c.GetQuery("q"); query == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "no search query present"})
		return
	}

	body := fmt.Sprintf(
		`{"query": {"multi_match": {"query": "%s", "fields": ["title", "body"]}}}`,
		query)
	res, err := h.ESClient.Search(
		h.ESClient.Search.WithContext(context.Background()),
		h.ESClient.Search.WithIndex("posts"),
		h.ESClient.Search.WithBody(strings.NewReader(body)),
		h.ESClient.Search.WithPretty(),
	)
	if err != nil {
		h.Logger.Err(err).Msg("elasticsearch error")
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	defer res.Body.Close()
	if res.IsError() {
		var e map[string]interface{}
		if err := json.NewDecoder(res.Body).Decode(&e); err != nil {
			h.Logger.Err(err).Msg("error parsing the response body")
		} else {
			h.Logger.Err(fmt.Errorf("[%s] %s: %s",
				res.Status(),
				e["error"].(map[string]interface{})["type"],
				e["error"].(map[string]interface{})["reason"],
			)).Msg("failed to search query")
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": e["error"].(map[string]interface{})["reason"]})
		return
	}

	h.Logger.Info().Interface("res", res.Status())

	var r map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		h.Logger.Err(err).Msg("elasticsearch error")
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": r["hits"]})
}

func (h *Handler) Ping(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"Status": "OK"})
}



