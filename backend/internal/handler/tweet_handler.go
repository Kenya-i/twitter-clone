package handler

import (
	"net/http"

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

	tweet, err := h.tweetUsecase.GetTweet(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "ツイートが見つかりません"})
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
