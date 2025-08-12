package middleware

import (
	"alertly/internal/response"
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"regexp"
	"strings"

	"github.com/gin-gonic/gin"
)

// Profanity levels for English (Canada)
var profanityLists = map[string][]string{
	// Very offensive words - BLOCK completely
	"blocked": {
		"fuck", "shit", "bitch", "cunt", "pussy", "dick", "cock", "asshole", "motherfucker",
		"fucker", "whore", "slut", "nigger", "faggot", "cocksucker", "piss", "tits",
	},

	// Moderate words - CENSOR with asterisks
	"censored": {
		"damn", "hell", "ass", "bastard", "bitch", "bullshit", "crap", "dumb", "idiot",
		"stupid", "retard", "fat", "ugly", "loser", "jerk", "douche", "suck", "sucks",
	},

	// Context-sensitive words - ALLOW but flag for review
	"context": {
		"kill", "death", "die", "dead", "blood", "gun", "shoot", "bomb", "explode",
		"hate", "hateful", "racist", "sexist", "homophobic", "terrorist", "suicide",
	},
}

type ProfanityResult struct {
	HasProfanity bool   `json:"has_profanity"`
	Level        string `json:"level"`
	FilteredText string `json:"filtered_text"`
	Message      string `json:"message"`
}

// FilterProfanity validates and filters text for profanity
func FilterProfanity(text string) ProfanityResult {
	if text == "" {
		return ProfanityResult{
			HasProfanity: false,
			Level:        "",
			FilteredText: text,
			Message:      "",
		}
	}

	lowerText := strings.ToLower(text)
	hasBlocked := false
	hasCensored := false
	hasContext := false
	filteredText := text

	// Check for blocked words
	for _, word := range profanityLists["blocked"] {
		regex := regexp.MustCompile(`\b` + regexp.QuoteMeta(word) + `\b`)
		if regex.MatchString(lowerText) {
			hasBlocked = true
		}
	}

	// Check for censored words
	for _, word := range profanityLists["censored"] {
		regex := regexp.MustCompile(`\b` + regexp.QuoteMeta(word) + `\b`)
		if regex.MatchString(lowerText) {
			hasCensored = true
			// Replace with asterisks while preserving case
			replacement := strings.Repeat("*", len(word))
			filteredText = regex.ReplaceAllStringFunc(filteredText, func(match string) string {
				if strings.ToLower(match) == word {
					return replacement
				}
				return match
			})
		}
	}

	// Check for context words
	for _, word := range profanityLists["context"] {
		regex := regexp.MustCompile(`\b` + regexp.QuoteMeta(word) + `\b`)
		if regex.MatchString(lowerText) {
			hasContext = true
		}
	}

	// Determine profanity level and message
	var level string
	var message string

	if hasBlocked {
		level = "blocked"
		message = "This comment contains inappropriate language and cannot be posted."
		filteredText = "" // Return empty for blocked content
	} else if hasCensored {
		level = "censored"
		message = "Some words have been censored for community guidelines."
	} else if hasContext {
		level = "context"
		message = "This comment may be reviewed for context."
	}

	return ProfanityResult{
		HasProfanity: hasBlocked || hasCensored || hasContext,
		Level:        level,
		FilteredText: filteredText,
		Message:      message,
	}
}

// ProfanityFilterMiddleware validates comments before they reach the handler
func ProfanityFilterMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Only apply to comment endpoints
		if c.Request.URL.Path == "/api/cluster/send_comment" && c.Request.Method == "POST" {
			// Read the request body without consuming it
			bodyBytes, err := c.GetRawData()
			if err != nil {
				response.Send(c, http.StatusBadRequest, true, "Invalid request format", err.Error())
				c.Abort()
				return
			}

			// Parse the JSON manually
			var requestBody map[string]interface{}
			if err := json.Unmarshal(bodyBytes, &requestBody); err != nil {
				response.Send(c, http.StatusBadRequest, true, "Invalid JSON format", err.Error())
				c.Abort()
				return
			}

			// Extract comment text
			commentText, ok := requestBody["comment"].(string)
			if !ok {
				response.Send(c, http.StatusBadRequest, true, "Comment field is required", nil)
				c.Abort()
				return
			}

			// Filter profanity
			result := FilterProfanity(commentText)

			// If blocked, reject the request
			if result.Level == "blocked" {
				response.Send(c, http.StatusBadRequest, true, result.Message, nil)
				c.Abort()
				return
			}

			// If censored or context, update the request body with filtered text
			if result.HasProfanity {
				requestBody["comment"] = result.FilteredText
				requestBody["profanity_level"] = result.Level
				requestBody["profanity_message"] = result.Message

				// Update the request body in context
				c.Set("filtered_request", requestBody)
			}

			// Restore the body for the main handler
			c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
		}

		c.Next()
	}
}
