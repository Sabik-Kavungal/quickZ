package handlers

import (
	"encoding/json"
	"fmt"

	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
)

const huggingFaceAPI = "https://api-inference.huggingface.co/models/facebook/blenderbot-400M-distill"
const HuggingFaceAPIKey = "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"

func AiChatbot(c *gin.Context) {
	var request struct {
		Message string `json:"message" binding:"required"`
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	client := resty.New()
	response, err := client.R().
		SetHeader("Authorization", "Bearer "+HuggingFaceAPIKey).
		SetHeader("Content-Type", "application/json").
		SetBody(map[string]interface{}{
			"inputs": request.Message,
			"parameters": map[string]interface{}{
				"max_length":  50,
				"temperature": 0.7,
			},
		}).
		Post(huggingFaceAPI)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to connect to AI service"})
		return
	}

	fmt.Println("Response Body:", string(response.Body())) // Debugging

	// Handle AI response
	var aiResponse []map[string]interface{}
	if err := json.Unmarshal(response.Body(), &aiResponse); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse AI response"})
		return
	}

	if len(aiResponse) > 0 {
		if generatedText, ok := aiResponse[0]["generated_text"].(string); ok {
			c.JSON(http.StatusOK, gin.H{
				"success":  true,
				"response": generatedText,
			})
			return
		}
	}

	c.JSON(http.StatusInternalServerError, gin.H{"error": "Unexpected AI response format"})
}
