package services

import (
	"database/sql"
	"errors"
	"fmt"
	"jesterx-core/config"
	"jesterx-core/helpers"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/stripe/stripe-go/v82"
	stripecheckout "github.com/stripe/stripe-go/v82/checkout/session"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type CheckoutRequest struct {
	Plan string `json:"plan" binding:"required"`
}

type apiResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

func CreateCheckoutService(c *gin.Context) {
	userAny, ok := c.Get("user")
	if !ok {
		c.JSON(http.StatusUnauthorized, apiResponse{Success: false, Error: "unauthorized"})
		return
	}
	user, ok := userAny.(helpers.UserData)
	if !ok || user.Id == "" {
		c.JSON(http.StatusUnauthorized, apiResponse{Success: false, Error: "unauthorized"})
		return
	}

	var req CheckoutRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, apiResponse{Success: false, Error: err.Error()})
		return
	}

	planConfig, err := GetPlanConfig(c.Request.Context(), req.Plan)
	if err != nil {
		c.JSON(http.StatusBadRequest, apiResponse{Success: false, Error: "invalid plan"})
		return
	}
	if planConfig.PriceCents <= 0 {
		c.JSON(http.StatusBadRequest, apiResponse{Success: false, Error: "plan is not billable"})
		return
	}

	title := cases.Title(language.Portuguese).String(planConfig.Name)
	productName := fmt.Sprintf("Plano %s", title)

	params := &stripe.CheckoutSessionParams{
		Mode: stripe.String(string(stripe.CheckoutSessionModePayment)),
		SuccessURL: stripe.String(
			fmt.Sprintf("%s?session_id={CHECKOUT_SESSION_ID}", config.StripeSuccessURL),
		),
		CancelURL:         stripe.String(config.StripeCancelURL),
		ClientReferenceID: stripe.String(user.Id),
		LineItems: []*stripe.CheckoutSessionLineItemParams{
			{
				PriceData: &stripe.CheckoutSessionLineItemPriceDataParams{
					Currency:   stripe.String("brl"),
					UnitAmount: stripe.Int64(planConfig.PriceCents),
					ProductData: &stripe.CheckoutSessionLineItemPriceDataProductDataParams{
						Name: stripe.String(productName),
					},
				},
				Quantity: stripe.Int64(1),
			},
		},
	}

	params.AddMetadata("user_id", user.Id)
	params.AddMetadata("plan", req.Plan)

	sess, err := stripecheckout.New(params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, apiResponse{Success: false, Error: "stripe checkout session creation failed"})
		return
	}

	_, err = config.DB.Exec(`
		INSERT INTO payments (user_id, plan, provider, provider_payment_id, status, amount_cents, currency)
		VALUES ($1, $2, 'stripe', $3, 'pending', $4, 'BRL')
	`, user.Id, req.Plan, sess.ID, planConfig.PriceCents)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			c.JSON(http.StatusInternalServerError, apiResponse{Success: false, Error: "payment persistence failed"})
			return
		}
		c.JSON(http.StatusInternalServerError, apiResponse{Success: false, Error: "payment persistence failed"})
		return
	}

	c.JSON(http.StatusOK, apiResponse{
		Success: true,
		Data: gin.H{
			"provider":     "stripe",
			"session_id":   sess.ID,
			"checkout_url": sess.URL,
		},
	})
}
