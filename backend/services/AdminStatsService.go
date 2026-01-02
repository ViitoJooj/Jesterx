package services

import (
	"context"
	"jesterx-core/config"
	"jesterx-core/responses"
	"net/http"

	"github.com/gin-gonic/gin"
)

func AdminOverviewService(c *gin.Context) {
	ctx := c.Request.Context()
	var overview responses.AdminOverviewResponse

	if err := config.DB.QueryRowContext(ctx, `SELECT COUNT(*) FROM users`).Scan(&overview.TotalUsers); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Failed to load stats"})
		return
	}

	if err := config.DB.QueryRowContext(ctx, `SELECT COUNT(*) FROM users WHERE banned = TRUE`).Scan(&overview.BannedUsers); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Failed to load stats"})
		return
	}

	overview.ActiveUsers = overview.TotalUsers - overview.BannedUsers

	if err := config.DB.QueryRowContext(ctx, `SELECT COALESCE(SUM(amount_cents), 0) FROM payments WHERE status = 'paid'`).Scan(&overview.PaidTotalCents); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Failed to load stats"})
		return
	}

	if err := config.DB.QueryRowContext(ctx, `SELECT COALESCE(SUM(amount_cents), 0) FROM payments WHERE status = 'paid' AND created_at >= NOW() - INTERVAL '30 days'`).Scan(&overview.PaidLast30DaysCents); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Failed to load stats"})
		return
	}

	if err := config.DB.QueryRowContext(ctx, `SELECT COALESCE(COUNT(*),0) FROM users WHERE created_at >= NOW() - INTERVAL '30 days'`).Scan(&overview.NewUsersLast30Days); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Failed to load stats"})
		return
	}

	if err := config.DB.QueryRowContext(ctx, `SELECT COALESCE(COUNT(*),0) FROM users WHERE created_at >= NOW() - INTERVAL '24 hours'`).Scan(&overview.CreatedLast24h); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Failed to load stats"})
		return
	}

	if err := config.DB.QueryRowContext(ctx, `SELECT COALESCE(SUM(amount_cents), 0) FROM payments WHERE status = 'paid' AND created_at >= NOW() - INTERVAL '24 hours'`).Scan(&overview.PaymentsLast24hCents); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Failed to load stats"})
		return
	}

	if err := config.DB.QueryRowContext(ctx, `SELECT COALESCE(AVG(amount_cents), 0) FROM payments WHERE status = 'paid'`).Scan(&overview.AverageTicketCents); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Failed to load stats"})
		return
	}

	if err := config.DB.QueryRowContext(ctx, `SELECT COALESCE(COUNT(DISTINCT user_id), 0) FROM payments WHERE status = 'paid'`).Scan(&overview.PayingUsers); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Failed to load stats"})
		return
	}

	newUsersSeries, err := loadDailySeries(ctx, `
		SELECT to_char(date_trunc('day', created_at), 'YYYY-MM-DD') AS label, COUNT(*) AS value
		FROM users
		WHERE created_at >= NOW() - INTERVAL '14 days'
		GROUP BY label
		ORDER BY label DESC
	`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Failed to load stats"})
		return
	}
	overview.NewUsersSeries = newUsersSeries

	paymentSeries, err := loadDailySeries(ctx, `
		SELECT to_char(date_trunc('day', created_at), 'YYYY-MM-DD') AS label, COALESCE(SUM(amount_cents),0) AS value
		FROM payments
		WHERE status = 'paid' AND created_at >= NOW() - INTERVAL '14 days'
		GROUP BY label
		ORDER BY label DESC
	`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Failed to load stats"})
		return
	}
	overview.PaymentsSeries = paymentSeries

	planUsage, err := loadDailySeries(ctx, `
		SELECT plan AS label, COUNT(*) AS value
		FROM users
		GROUP BY plan
		ORDER BY value DESC
	`)
	if err == nil {
		overview.PlansByUsage = planUsage
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    overview,
	})
}

func loadDailySeries(ctx context.Context, query string) ([]responses.AdminMetricPoint, error) {
	rows, err := config.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []responses.AdminMetricPoint
	for rows.Next() {
		var item responses.AdminMetricPoint
		if err := rows.Scan(&item.Label, &item.Value); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, rows.Err()
}
