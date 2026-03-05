package service

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/ViitoJooj/Jesterx/internal/config"
	"github.com/ViitoJooj/Jesterx/internal/domain"
	"github.com/ViitoJooj/Jesterx/internal/repository"
)

type PaymentService struct {
	paymentRepo repository.PaymentRepository
	userRepo    repository.UserRepository
}

type CheckoutSessionResult struct {
	SessionID   string
	CheckoutURL string
}

type stripeCheckoutSession struct {
	ID            string            `json:"id"`
	URL           string            `json:"url"`
	Mode          string            `json:"mode"`
	PaymentStatus string            `json:"payment_status"`
	Status        string            `json:"status"`
	AmountTotal   int64             `json:"amount_total"`
	Currency      string            `json:"currency"`
	Subscription  string            `json:"subscription"`
	Metadata      map[string]string `json:"metadata"`
}

type stripeCustomer struct {
	ID string `json:"id"`
}

type stripeListCustomersResponse struct {
	Data []stripeCustomer `json:"data"`
}

type stripeEvent struct {
	Type string `json:"type"`
	Data struct {
		Object stripeCheckoutSession `json:"object"`
	} `json:"data"`
}

func NewPaymentService(paymentRepo repository.PaymentRepository, userRepo repository.UserRepository) *PaymentService {
	return &PaymentService{
		paymentRepo: paymentRepo,
		userRepo:    userRepo,
	}
}

func (s *PaymentService) ListActivePlans() ([]domain.Plan, error) {
	return s.paymentRepo.ListActivePlans()
}

func (s *PaymentService) CreateCheckoutSession(userID string, planID int64, quantity int) (*CheckoutSessionResult, error) {
	if userID == "" {
		return nil, errors.New("invalid user")
	}
	if planID <= 0 {
		return nil, errors.New("invalid plan")
	}
	if quantity <= 0 {
		quantity = 1
	}
	if quantity > 100 {
		return nil, errors.New("quantity is too high")
	}

	user, err := s.userRepo.FindUserByID(userID)
	if err != nil {
		return nil, errors.New("internal error")
	}
	if user == nil {
		return nil, errors.New("user not found")
	}

	plan, err := s.paymentRepo.FindPlanByID(planID)
	if err != nil {
		return nil, errors.New("internal error")
	}
	if plan == nil || !plan.Active {
		return nil, errors.New("plan not available")
	}

	customerID, err := ensureStripeCustomer(*user)
	if err != nil {
		return nil, err
	}

	session, err := createStripeCheckoutSession(*user, *plan, customerID, quantity)
	if err != nil {
		return nil, err
	}

	paymentType := "plan"
	if isRecurringCycle(plan.BillingCycle) {
		paymentType = "subscription"
	}

	payment := domain.Payment{
		UserID:      user.Id,
		WebsiteID:   user.WebsiteId,
		ReferenceID: session.ID,
		Type:        paymentType,
		Quantity:    quantity,
		Amount:      plan.Price * float64(quantity),
		Currency:    "BRL",
		Status:      "pending",
	}

	if _, err := s.paymentRepo.CreatePayment(payment); err != nil {
		return nil, errors.New("internal error")
	}

	return &CheckoutSessionResult{
		SessionID:   session.ID,
		CheckoutURL: session.URL,
	}, nil
}

func (s *PaymentService) ProcessStripeWebhook(rawBody []byte, signature string) error {
	if config.StripeWebhookSecret != "" {
		if !validateStripeSignature(rawBody, signature, config.StripeWebhookSecret) {
			return errors.New("invalid stripe signature")
		}
	}

	var event stripeEvent
	if err := json.Unmarshal(rawBody, &event); err != nil {
		return errors.New("invalid event payload")
	}

	ref := event.Data.Object.ID
	if ref == "" {
		return errors.New("missing reference id")
	}

	switch event.Type {
	case "checkout.session.completed", "checkout.session.async_payment_succeeded":
		return s.paymentRepo.UpdatePaymentStatusByReference(ref, "completed")
	case "checkout.session.expired", "checkout.session.async_payment_failed":
		return s.paymentRepo.UpdatePaymentStatusByReference(ref, "canceled")
	default:
		return nil
	}
}

func (s *PaymentService) CancelSubscription(userID string) error {
	payment, err := s.paymentRepo.FindLatestCompletedPaymentByUserID(userID)
	if err != nil {
		return errors.New("internal error")
	}
	if payment == nil {
		return errors.New("no active subscription found")
	}

	session, err := retrieveStripeCheckoutSession(payment.ReferenceID)
	if err != nil {
		return err
	}
	if session.Subscription == "" {
		return errors.New("no stripe subscription associated")
	}

	if err := cancelStripeSubscription(session.Subscription); err != nil {
		return err
	}

	return s.paymentRepo.UpdatePaymentStatusByReference(payment.ReferenceID, "canceled")
}

func cancelStripeSubscription(subscriptionID string) error {
	endpoint := "https://api.stripe.com/v1/subscriptions/" + url.QueryEscape(subscriptionID)
	req, err := http.NewRequest(http.MethodDelete, endpoint, nil)
	if err != nil {
		return errors.New("internal error")
	}
	req.Header.Set("Authorization", "Bearer "+config.StripeSecret)

	client := &http.Client{Timeout: 15 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return errors.New("failed to reach stripe")
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("stripe error: %s", string(body))
	}
	return nil
}

func (s *PaymentService) ConfirmCheckoutSession(userID, sessionID string) error {
	if userID == "" || sessionID == "" {
		return errors.New("invalid request")
	}

	session, err := retrieveStripeCheckoutSession(sessionID)
	if err != nil {
		return err
	}

	metaUserID := session.Metadata["user_id"]
	if metaUserID == "" || metaUserID != userID {
		return errors.New("forbidden")
	}

	switch {
	case session.Status == "complete" && (session.PaymentStatus == "paid" || session.Mode == "subscription"):
		return s.paymentRepo.UpdatePaymentStatusByReference(session.ID, "completed")
	case session.Status == "expired" || session.PaymentStatus == "unpaid" || session.PaymentStatus == "no_payment_required":
		return s.paymentRepo.UpdatePaymentStatusByReference(session.ID, "canceled")
	default:
		return s.paymentRepo.UpdatePaymentStatusByReference(session.ID, "pending")
	}
}

func createStripeCheckoutSession(user domain.User, plan domain.Plan, customerID string, quantity int) (*stripeCheckoutSession, error) {
	cents := int64(plan.Price * 100)
	if cents <= 0 {
		return nil, errors.New("invalid plan price")
	}

	successURL := strings.TrimRight(config.FrontendURL, "/") + "/payment-success?session_id={CHECKOUT_SESSION_ID}&plan_id=" + strconv.FormatInt(plan.ID, 10)
	cancelURL := strings.TrimRight(config.FrontendURL, "/") + "/payment-cancel?plan_id=" + strconv.FormatInt(plan.ID, 10)

	form := url.Values{}
	form.Set("success_url", successURL)
	form.Set("cancel_url", cancelURL)
	form.Set("customer", customerID)
	form.Set("line_items[0][quantity]", strconv.Itoa(quantity))
	form.Set("line_items[0][price_data][currency]", "brl")
	form.Set("line_items[0][price_data][unit_amount]", strconv.FormatInt(cents, 10))
	form.Set("line_items[0][price_data][product_data][name]", plan.Name)

	if interval := stripeIntervalByBillingCycle(plan.BillingCycle); interval != "" {
		form.Set("mode", "subscription")
		form.Set("line_items[0][price_data][recurring][interval]", interval)
		form.Set("subscription_data[metadata][plan_id]", strconv.FormatInt(plan.ID, 10))
		form.Set("subscription_data[metadata][website_id]", user.WebsiteId)
		form.Set("subscription_data[metadata][user_id]", user.Id)
	} else {
		form.Set("mode", "payment")
	}

	if plan.Description != "" {
		form.Set("line_items[0][price_data][product_data][description]", plan.Description)
	}
	form.Set("metadata[user_id]", user.Id)
	form.Set("metadata[website_id]", user.WebsiteId)
	form.Set("metadata[plan_id]", strconv.FormatInt(plan.ID, 10))

	req, err := http.NewRequest(http.MethodPost, "https://api.stripe.com/v1/checkout/sessions", bytes.NewBufferString(form.Encode()))
	if err != nil {
		return nil, errors.New("internal error")
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Authorization", "Bearer "+config.StripeSecret)

	client := &http.Client{Timeout: 15 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, errors.New("failed to create checkout")
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.New("failed to create checkout")
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("stripe error: %s", string(body))
	}

	var session stripeCheckoutSession
	if err := json.Unmarshal(body, &session); err != nil {
		return nil, errors.New("failed to parse checkout session")
	}
	if session.ID == "" || session.URL == "" {
		return nil, errors.New("invalid checkout response")
	}
	return &session, nil
}

func retrieveStripeCheckoutSession(sessionID string) (*stripeCheckoutSession, error) {
	endpoint := "https://api.stripe.com/v1/checkout/sessions/" + url.QueryEscape(sessionID)
	req, err := http.NewRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, errors.New("internal error")
	}
	req.Header.Set("Authorization", "Bearer "+config.StripeSecret)

	client := &http.Client{Timeout: 15 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, errors.New("failed to reach stripe")
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.New("failed to parse stripe response")
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("stripe error: %s", string(body))
	}

	var session stripeCheckoutSession
	if err := json.Unmarshal(body, &session); err != nil {
		return nil, errors.New("failed to parse checkout session")
	}
	if session.ID == "" {
		return nil, errors.New("invalid checkout response")
	}
	return &session, nil
}

func isRecurringCycle(cycle string) bool {
	return stripeIntervalByBillingCycle(cycle) != ""
}

func stripeIntervalByBillingCycle(cycle string) string {
	switch strings.ToLower(strings.TrimSpace(cycle)) {
	case "monthly", "month":
		return "month"
	case "yearly", "annual", "year":
		return "year"
	default:
		return ""
	}
}

func ensureStripeCustomer(user domain.User) (string, error) {
	existingID, err := findStripeCustomerByEmail(user.Email)
	if err != nil {
		return "", err
	}
	if existingID != "" {
		return existingID, nil
	}
	return createStripeCustomer(user)
}

func findStripeCustomerByEmail(email string) (string, error) {
	endpoint := "https://api.stripe.com/v1/customers?email=" + url.QueryEscape(email) + "&limit=1"
	req, err := http.NewRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		return "", errors.New("internal error")
	}
	req.Header.Set("Authorization", "Bearer "+config.StripeSecret)

	client := &http.Client{Timeout: 15 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", errors.New("failed to reach stripe")
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", errors.New("failed to parse stripe response")
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return "", fmt.Errorf("stripe error: %s", string(body))
	}

	var payload stripeListCustomersResponse
	if err := json.Unmarshal(body, &payload); err != nil {
		return "", errors.New("failed to parse stripe response")
	}
	if len(payload.Data) == 0 {
		return "", nil
	}
	return payload.Data[0].ID, nil
}

func createStripeCustomer(user domain.User) (string, error) {
	form := url.Values{}
	form.Set("email", user.Email)
	form.Set("name", strings.TrimSpace(user.First_name+" "+user.Last_name))
	form.Set("metadata[user_id]", user.Id)
	form.Set("metadata[website_id]", user.WebsiteId)

	req, err := http.NewRequest(http.MethodPost, "https://api.stripe.com/v1/customers", bytes.NewBufferString(form.Encode()))
	if err != nil {
		return "", errors.New("internal error")
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Authorization", "Bearer "+config.StripeSecret)

	client := &http.Client{Timeout: 15 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", errors.New("failed to create customer")
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", errors.New("failed to create customer")
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return "", fmt.Errorf("stripe error: %s", string(body))
	}

	var customer stripeCustomer
	if err := json.Unmarshal(body, &customer); err != nil {
		return "", errors.New("failed to parse customer")
	}
	if customer.ID == "" {
		return "", errors.New("invalid customer response")
	}
	return customer.ID, nil
}

func validateStripeSignature(payload []byte, header string, secret string) bool {
	if header == "" || secret == "" {
		return false
	}

	var timestamp string
	signatures := make([]string, 0, 2)

	parts := strings.Split(header, ",")
	for _, part := range parts {
		kv := strings.SplitN(strings.TrimSpace(part), "=", 2)
		if len(kv) != 2 {
			continue
		}
		switch kv[0] {
		case "t":
			timestamp = kv[1]
		case "v1":
			signatures = append(signatures, kv[1])
		}
	}

	if timestamp == "" || len(signatures) == 0 {
		return false
	}

	ts, err := strconv.ParseInt(timestamp, 10, 64)
	if err != nil {
		return false
	}

	if time.Since(time.Unix(ts, 0)) > 5*time.Minute {
		return false
	}

	signedPayload := timestamp + "." + string(payload)
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(signedPayload))
	expected := hex.EncodeToString(mac.Sum(nil))

	for _, sig := range signatures {
		if subtle.ConstantTimeCompare([]byte(sig), []byte(expected)) == 1 {
			return true
		}
	}

	return false
}
