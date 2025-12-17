package services

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"gen-you-ecommerce/config"
	"io"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/stripe/stripe-go/v82"
	"github.com/stripe/stripe-go/v82/webhook"
)

func PaymentWebhookService(c *gin.Context) {
	const maxBodyBytes = int64(1 << 20)

	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, maxBodyBytes)
	raw, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, apiResponse{Success: false, Error: "invalid body"})
		return
	}

	sig := c.GetHeader("Stripe-Signature")
	if strings.TrimSpace(sig) == "" {
		c.JSON(http.StatusBadRequest, apiResponse{Success: false, Error: "missing signature"})
		return
	}

	event, err := webhook.ConstructEvent(raw, sig, config.StripeWebhookSecret)
	if err != nil {
		c.JSON(http.StatusBadRequest, apiResponse{Success: false, Error: "signature verification failed"})
		return
	}

	switch event.Type {
	case "checkout.session.completed", "checkout.session.async_payment_succeeded", "checkout.session.async_payment_failed":
		var sess stripe.CheckoutSession
		if err := json.NewDecoder(bytes.NewReader(event.Data.Raw)).Decode(&sess); err != nil {
			c.JSON(http.StatusBadRequest, apiResponse{Success: false, Error: "invalid event payload"})
			return
		}
		if sess.ID == "" {
			c.JSON(http.StatusBadRequest, apiResponse{Success: false, Error: "missing session id"})
			return
		}

		newStatus := "pending"
		if event.Type == "checkout.session.async_payment_failed" {
			newStatus = "failed"
		} else if sess.PaymentStatus == stripe.CheckoutSessionPaymentStatusPaid {
			newStatus = "paid"
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
		var currentStatus string

		err = tx.QueryRow(`
			SELECT id, user_id, plan, status
			FROM payments
			WHERE provider = 'stripe'
			  AND provider_payment_id = $1
			FOR UPDATE
		`, sess.ID).Scan(&paymentID, &userID, &plan, &currentStatus)

		if err != nil {
			if err == sql.ErrNoRows {
				_ = tx.Commit()
				c.JSON(http.StatusOK, apiResponse{Success: true})
				return
			}
			c.JSON(http.StatusInternalServerError, apiResponse{Success: false, Error: "database error"})
			return
		}

		if currentStatus == "paid" {
			_ = tx.Commit()
			c.JSON(http.StatusOK, apiResponse{Success: true})
			return
		}

		_, err = tx.Exec(`
			UPDATE payments
			SET status = $1, updated_at = NOW()
			WHERE id = $2
		`, newStatus, paymentID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, apiResponse{Success: false, Error: "payment update failed"})
			return
		}

		if newStatus == "paid" {
			_, err = tx.Exec(`
				UPDATE users
				SET plan = $1, updated_at = NOW()
				WHERE id = $2
			`, plan, userID)
			if err != nil {
				c.JSON(http.StatusInternalServerError, apiResponse{Success: false, Error: "user update failed"})
				return
			}
		}

		if err := tx.Commit(); err != nil {
			c.JSON(http.StatusInternalServerError, apiResponse{Success: false, Error: "transaction commit failed"})
			return
		}

		c.JSON(http.StatusOK, apiResponse{Success: true})
		return

	default:
		c.JSON(http.StatusOK, apiResponse{Success: true})
		return
	}
}
