package services

import (
	"net/http"
	"strings"

	"jesterx-core/responses"

	"github.com/gin-gonic/gin"
)

type updatePlanBody struct {
	Name        string   `json:"name"`
	PriceCents  int64    `json:"price_cents"`
	Description string   `json:"description"`
	Features    []string `json:"features"`
	SiteLimit   int      `json:"site_limit"`
}

func AdminListPlansService(c *gin.Context) {
	plans, err := ListPlanConfigs(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Failed to list plans"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    toPlanResponse(plans),
	})
}

func AdminUpdatePlanService(c *gin.Context) {
	planID := strings.TrimSpace(c.Param("plan_id"))
	if planID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Plan id required"})
		return
	}

	var body updatePlanBody
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Invalid body"})
		return
	}

	updated, err := UpdatePlanConfig(c.Request.Context(), planID, PlanConfig{
		ID:          planID,
		Name:        strings.TrimSpace(body.Name),
		PriceCents:  body.PriceCents,
		Description: strings.TrimSpace(body.Description),
		Features:    body.Features,
		SiteLimit:   body.SiteLimit,
	})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    toPlanResponse([]PlanConfig{updated}),
	})
}

func ListPlansService(c *gin.Context) {
	plans, err := ListPlanConfigs(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Failed to list plans"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    toPlanResponse(plans),
	})
}

func toPlanResponse(plans []PlanConfig) []responses.PlanConfigResponse {
	var resp []responses.PlanConfigResponse
	for _, plan := range plans {
		resp = append(resp, responses.PlanConfigResponse{
			ID:          plan.ID,
			Name:        plan.Name,
			PriceCents:  plan.PriceCents,
			Description: plan.Description,
			Features:    plan.Features,
			SiteLimit:   plan.SiteLimit,
		})
	}
	return resp
}
