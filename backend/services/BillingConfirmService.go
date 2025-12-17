package services

import (
	"database/sql"
	"gen-you-ecommerce/config"
	"gen-you-ecommerce/helpers"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/stripe/stripe-go/v82"
	stripecheckout "github.com/stripe/stripe-go/v82/checkout/session"
)

type ConfirmRequest struct {
	SessionID string `json:"session_id" binding:"required"`
}

func ConfirmCheckoutService(c *gin.Context) {
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

	var req ConfirmRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, apiResponse{Success: false, Error: err.Error()})
		return
	}

	tx, err := config.DB.Begin()
	if err != nil {
		c.JSON(http.StatusInternalServerError, apiResponse{Success: false, Error: "transaction start failed"})
		return
	}
	defer tx.Rollback()

	var paymentID string
	var userID string
	var plan string
	var status string

	err = tx.QueryRow(`
		SELECT id, user_id, plan, status
		FROM payments
		WHERE provider = 'stripe'
		  AND provider_payment_id = $1
		FOR UPDATE
	`, req.SessionID).Scan(&paymentID, &userID, &plan, &status)

	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, apiResponse{Success: false, Error: "payment not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, apiResponse{Success: false, Error: "database error"})
		return
	}

	if userID != user.Id {
		c.JSON(http.StatusForbidden, apiResponse{Success: false, Error: "forbidden"})
		return
	}

	if status == "paid" {
		_ = tx.Commit()
		c.JSON(http.StatusOK, apiResponse{
			Success: true,
			Data: gin.H{
				"status": "paid",
				"plan":   plan,
			},
		})
		return
	}

	sess, err := stripecheckout.Get(req.SessionID, nil)
	if err != nil {
		c.JSON(http.StatusBadRequest, apiResponse{Success: false, Error: "invalid session"})
		return
	}

	if uid, ok := sess.Metadata["user_id"]; ok && uid != user.Id {
		c.JSON(http.StatusForbidden, apiResponse{Success: false, Error: "forbidden"})
		return
	}

	if sess.PaymentStatus != stripe.CheckoutSessionPaymentStatusPaid {
		_ = tx.Commit()
		c.JSON(http.StatusOK, apiResponse{
			Success: true,
			Data: gin.H{
				"status": string(sess.PaymentStatus),
				"plan":   plan,
			},
		})
		return
	}

	_, err = tx.Exec(`
		UPDATE payments
		SET status = 'paid', updated_at = NOW()
		WHERE id = $1
	`, paymentID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, apiResponse{Success: false, Error: "payment update failed"})
		return
	}

	_, err = tx.Exec(`
		UPDATE users
		SET plan = $1, updated_at = NOW()
		WHERE id = $2
	`, plan, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, apiResponse{Success: false, Error: "user update failed"})
		return
	}

	if err := tx.Commit(); err != nil {
		c.JSON(http.StatusInternalServerError, apiResponse{Success: false, Error: "transaction commit failed"})
		return
	}

	c.JSON(http.StatusOK, apiResponse{
		Success: true,
		Data: gin.H{
			"status": "paid",
			"plan":   plan,
		},
	})
}
