package handler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/Kenya-i/twitter-clone/internal/domain"
	"github.com/gin-gonic/gin"
)

type TweetHandler struct {
	tweetUsecase domain.TweetUsecase
}

func NewTweetHandler(tweetUsecase domain.TweetUsecase) *TweetHandler {
	return &TweetHandler{tweetUsecase: tweetUsecase}
}

func (h *TweetHandler) Post(c *gin.Context) {
	var req struct {
		Content string `json:"content" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "contentは必須です"})
		return
	}

	userID := c.GetString("user_id")

	tweet, err := h.tweetUsecase.Post(userID, req.Content)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, tweet)
}

func (h *TweetHandler) GetTweet(c *gin.Context) {
	id := c.Param("id")
	userID := c.GetString("user_id")

	tweet, err := h.tweetUsecase.GetTweet(id, userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "ツイートが見つかりません"})
		return
	}

	c.JSON(http.StatusOK, tweet)
}

func (h *TweetHandler) GetTimeline(c *gin.Context) {
	userID := c.GetString("user_id")

	limit := 20
	if l := c.Query("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil && parsed > 0 {
			limit = parsed
		}
	}

	var cursor *time.Time
	if cs := c.Query("cursor"); cs != "" {
		parsed, err := time.Parse(time.RFC3339, cs)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "cursorの形式が不正です"})
			return
		}
		cursor = &parsed
	}

	tweets, err := h.tweetUsecase.GetTimeline(userID, cursor, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var nextCursor *time.Time
	if len(tweets) == limit {
		nextCursor = &tweets[len(tweets)-1].CreatedAt
	}

	c.JSON(http.StatusOK, gin.H{
		"tweets":      tweets,
		"next_cursor": nextCursor,
	})
}

func (h *TweetHandler) Search(c *gin.Context) {
	userID := c.GetString("user_id")
	query := c.Query("q")
	if query == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "qは必須です"})
		return
	}

	limit := 20
	if l := c.Query("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil && parsed > 0 {
			limit = parsed
		}
	}

	var cursor *time.Time
	if cs := c.Query("cursor"); cs != "" {
		parsed, err := time.Parse(time.RFC3339, cs)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "cursorの形式が不正です"})
			return
		}
		cursor = &parsed
	}

	tweets, err := h.tweetUsecase.Search(query, userID, cursor, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var nextCursor *time.Time
	if len(tweets) == limit {
		nextCursor = &tweets[len(tweets)-1].CreatedAt
	}

	c.JSON(http.StatusOK, gin.H{
		"tweets":      tweets,
		"next_cursor": nextCursor,
	})
}

func (h *TweetHandler) Update(c *gin.Context) {
	var req struct {
		Content string `json:"content" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "contentは必須です"})
		return
	}

	tweetID := c.Param("id")
	userID := c.GetString("user_id")

	tweet, err := h.tweetUsecase.Update(userID, tweetID, req.Content)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, tweet)
}

func (h *TweetHandler) Delete(c *gin.Context) {
	tweetID := c.Param("id")
	userID := c.GetString("user_id")

	if err := h.tweetUsecase.Delete(userID, tweetID); err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "削除しました"})
}

func (h *TweetHandler) Like(c *gin.Context) {
	tweetID := c.Param("id")
	userID := c.GetString("user_id")

	if err := h.tweetUsecase.Like(userID, tweetID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "いいねしました"})
}

func (h *TweetHandler) Unlike(c *gin.Context) {
	tweetID := c.Param("id")
	userID := c.GetString("user_id")

	if err := h.tweetUsecase.Unlike(userID, tweetID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "いいねを取り消しました"})
}
