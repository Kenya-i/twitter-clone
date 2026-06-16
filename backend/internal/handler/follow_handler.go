package handler

import (
	"net/http"

	"github.com/Kenya-i/twitter-clone/internal/domain"
	"github.com/gin-gonic/gin"
)

type FollowHandler struct {
	followUsecase domain.FollowUsecase
}

func NewFollowHandler(followUsecase domain.FollowUsecase) *FollowHandler {
	return &FollowHandler{followUsecase: followUsecase}
}

func (h *FollowHandler) Follow(c *gin.Context) {
	followerID := c.GetString("user_id")
	followingID := c.Param("id")

	if err := h.followUsecase.Follow(followerID, followingID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "フォローしました"})
}

func (h *FollowHandler) Unfollow(c *gin.Context) {
	followerID := c.GetString("user_id")
	followingID := c.Param("id")

	if err := h.followUsecase.Unfollow(followerID, followingID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "フォローを解除しました"})
}

func (h *FollowHandler) GetFollowInfo(c *gin.Context) {
	currentUserID := c.GetString("user_id")
	targetUserID := c.Param("id")

	followers, following, err := h.followUsecase.GetFollowCounts(targetUserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	isFollowing, err := h.followUsecase.IsFollowing(currentUserID, targetUserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"followers_count": followers,
		"following_count": following,
		"is_following":    isFollowing,
	})
}
