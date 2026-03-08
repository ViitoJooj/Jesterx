package jobs

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/ViitoJooj/Jesterx/internal/domain"
	"github.com/ViitoJooj/Jesterx/internal/repository"
	"github.com/ViitoJooj/Jesterx/internal/security"
	"github.com/ViitoJooj/Jesterx/internal/service"
)

func StartSalesDigestWorker(orderService *service.OrderService, userRepo repository.UserRepository, websiteRepo repository.WebsiteRepository) {
	go func() {
		for {
			now := time.Now()
			next := now.Truncate(2 * time.Hour).Add(2 * time.Hour)
			time.Sleep(time.Until(next))
			runSalesDigest(orderService, userRepo, websiteRepo)
		}
	}()
}

func runSalesDigest(orderService *service.OrderService, userRepo repository.UserRepository, websiteRepo repository.WebsiteRepository) {
	to := time.Now()
	from := to.Add(-2 * time.Hour)

	orders, err := orderService.GetSiteOrdersSince(from, to)
	if err != nil {
		log.Printf("[sales_digest] error fetching orders: %v", err)
		return
	}
	if len(orders) == 0 {
		return
	}

	byWebsite := make(map[string][]domain.Order)
	for _, o := range orders {
		byWebsite[o.WebsiteID] = append(byWebsite[o.WebsiteID], o)
	}

	for websiteID, siteOrders := range byWebsite {
		site, err := websiteRepo.FindWebSiteByID(websiteID)
		if err != nil || site == nil {
			continue
		}
		owner, err := userRepo.FindUserByID(site.Creator_id)
		if err != nil || owner == nil {
			continue
		}

		var total float64
		for _, o := range siteOrders {
			total += o.Total
		}

		subject := fmt.Sprintf("🛍 %d novo(s) pedido(s) em %s", len(siteOrders), site.Name)
		body := buildSalesDigestEmail(owner.First_name, site.Name, siteOrders, total, from, to)

		if err := security.SendSalesDigestEmail(owner.Email, subject, body); err != nil {
			log.Printf("[sales_digest] error sending email to %s: %v", owner.Email, err)
		} else {
			log.Printf("[sales_digest] sent digest to %s (%d orders, R$%.2f)", owner.Email, len(siteOrders), total)
		}
	}
}

func buildSalesDigestEmail(ownerName, siteName string, orders []domain.Order, total float64, from, to time.Time) string {
	var rows strings.Builder
	for _, o := range orders {
		itemNames := make([]string, 0, len(o.Items))
		for _, it := range o.Items {
			itemNames = append(itemNames, fmt.Sprintf("%dx %s", it.Qty, it.ProductName))
		}
		rows.WriteString(fmt.Sprintf(`
        <tr>
          <td style="padding:8px 12px;border-bottom:1px solid #eee;font-size:13px">%s</td>
          <td style="padding:8px 12px;border-bottom:1px solid #eee;font-size:13px">%s</td>
          <td style="padding:8px 12px;border-bottom:1px solid #eee;font-size:13px;font-weight:700;color:#ff5d1f">R$ %.2f</td>
        </tr>`, o.BuyerName, strings.Join(itemNames, ", "), o.Total))
	}

	return fmt.Sprintf(`<!DOCTYPE html>
<html><head><meta charset="utf-8"/></head>
<body style="margin:0;font-family:Inter,system-ui,sans-serif;background:#f5f7fa">
<div style="max-width:600px;margin:40px auto;background:#fff;border-radius:12px;overflow:hidden;box-shadow:0 2px 12px rgba(0,0,0,.08)">
  <div style="background:linear-gradient(135deg,#1a2740,#2d4070);padding:32px 40px;color:#fff">
    <h1 style="margin:0;font-size:22px">🛍 Resumo de Vendas</h1>
    <p style="margin:8px 0 0;opacity:.8;font-size:14px">%s → %s</p>
  </div>
  <div style="padding:32px 40px">
    <p style="font-size:16px;color:#1a2740">Olá, <strong>%s</strong>!</p>
    <p style="font-size:14px;color:#5a6379">Aqui estão os pedidos recebidos em <strong>%s</strong> nas últimas 2 horas:</p>
    <table style="width:100%%;border-collapse:collapse;margin:20px 0">
      <thead><tr style="background:#f5f7fa">
        <th style="padding:10px 12px;text-align:left;font-size:12px;color:#9aa5bc;font-weight:600">COMPRADOR</th>
        <th style="padding:10px 12px;text-align:left;font-size:12px;color:#9aa5bc;font-weight:600">ITENS</th>
        <th style="padding:10px 12px;text-align:left;font-size:12px;color:#9aa5bc;font-weight:600">TOTAL</th>
      </tr></thead>
      <tbody>%s</tbody>
    </table>
    <div style="background:#f5f7fa;border-radius:8px;padding:16px 20px;margin-top:16px">
      <span style="font-size:14px;color:#5a6379">Total do período: </span>
      <span style="font-size:20px;font-weight:700;color:#ff5d1f">R$ %.2f</span>
    </div>
  </div>
  <div style="padding:20px 40px;border-top:1px solid #f0f0f0;text-align:center;font-size:12px;color:#9aa5bc">
    Jesterx · Plataforma de Sites
  </div>
</div>
</body></html>`,
		from.Format("02/01 15:04"), to.Format("02/01 15:04"),
		ownerName, siteName,
		rows.String(),
		total,
	)
}
