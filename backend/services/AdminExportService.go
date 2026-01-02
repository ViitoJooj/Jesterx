package services

import (
	"bytes"
	"fmt"
	"jesterx-core/config"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/xuri/excelize/v2"
)

func AdminExportUsersService(c *gin.Context) {
	rows, err := config.DB.QueryContext(c.Request.Context(), `
		SELECT 
			u.id,
			u.email,
			COALESCE(p.first_name, ''),
			COALESCE(p.last_name, ''),
			COALESCE(u.plan, 'free'),
			COALESCE(u.role, 'platform_user'),
			COALESCE(u.banned, FALSE),
			u.created_at
		FROM users u
		LEFT JOIN user_profiles p ON p.user_id = u.id
		ORDER BY u.created_at DESC
	`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Failed to export users"})
		return
	}
	defer rows.Close()

	f := excelize.NewFile()
	sheet := "users"
	f.SetSheetName("Sheet1", sheet)

	headers := []string{"ID", "Email", "First Name", "Last Name", "Plan", "Role", "Banned", "Created At"}
	for idx, header := range headers {
		cell := fmt.Sprintf("%s1", string(rune('A'+idx)))
		f.SetCellValue(sheet, cell, header)
	}

	rowIndex := 2
	for rows.Next() {
		var id, email, firstName, lastName, plan, role string
		var banned bool
		var createdAt time.Time

		if err := rows.Scan(&id, &email, &firstName, &lastName, &plan, &role, &banned, &createdAt); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Failed to export users"})
			return
		}

		values := []interface{}{id, email, firstName, lastName, plan, role, banned, createdAt.Format(time.RFC3339)}
		for i, val := range values {
			cell := fmt.Sprintf("%s%d", string(rune('A'+i)), rowIndex)
			f.SetCellValue(sheet, cell, val)
		}
		rowIndex++
	}

	var buf bytes.Buffer
	if err := f.Write(&buf); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Failed to generate file"})
		return
	}

	c.Header("Content-Disposition", "attachment; filename=\"users.xlsx\"")
	c.Data(http.StatusOK, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", buf.Bytes())
}
