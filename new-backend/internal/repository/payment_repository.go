package repository

import (
	"context"
	"fmt"

	"event-management/backend/internal/domain"
	"event-management/backend/pkg/graphql"
)

type PaymentRepository interface {
	Create(ctx context.Context, payment *domain.Payment) (string, error)
	UpdateStatus(ctx context.Context, transactionRef, status string) error
	GetByTransactionRef(ctx context.Context, transactionRef string) (*domain.Payment, error)
	GetTicketByPaymentID(ctx context.Context, paymentID string) (*domain.Ticket, error)
	CheckDuplicateTicket(ctx context.Context, eventID, userID string) (bool, error)
	CreateTicket(ctx context.Context, ticket *domain.Ticket) (string, error)
	UpdateTicketStatus(ctx context.Context, paymentID, status string) error
	GetEventByID(ctx context.Context, eventID string) (*domain.Event, error)
}

type paymentRepository struct{}

func NewPaymentRepository() PaymentRepository {
	return &paymentRepository{}
}

func (r *paymentRepository) Create(ctx context.Context, payment *domain.Payment) (string, error) {
	mutation := `mutation InsertPayment($user_id: uuid!, $amount: numeric!, $currency: String!, $status: String!, $transaction_ref: String!, $payment_method: String!) {
        insert_payments_one(object: {
            user_id: $user_id
            amount: $amount
            currency: $currency
            status: $status
            transaction_ref: $transaction_ref
            payment_method: $payment_method
        }) {
            id
        }
    }`

	vars := map[string]interface{}{
		"user_id":         payment.UserID,
		"amount":          payment.Amount,
		"currency":        payment.Currency,
		"status":          payment.Status,
		"transaction_ref": payment.TransactionRef,
		"payment_method":  payment.PaymentMethod,
	}

	var resp struct {
		InsertPaymentsOne struct {
			ID string `json:"id"`
		} `json:"insert_payments_one"`
	}

	err := graphql.MutateRaw(ctx, mutation, vars, &resp)
	if err != nil {
		return "", fmt.Errorf("failed to create payment: %w", err)
	}

	return resp.InsertPaymentsOne.ID, nil
}

func (r *paymentRepository) UpdateStatus(ctx context.Context, transactionRef, status string) error {
	mutation := `mutation UpdatePayment($tx_ref: String!, $status: String!) {
        update_payments(where: { transaction_ref: { _eq: $tx_ref } }, _set: { status: $status }) {
            affected_rows
        }
    }`

	vars := map[string]interface{}{
		"tx_ref": transactionRef,
		"status": status,
	}

	var resp struct {
		UpdatePayments struct {
			AffectedRows int `json:"affected_rows"`
		} `json:"update_payments"`
	}

	err := graphql.MutateRaw(ctx, mutation, vars, &resp)
	if err != nil {
		return fmt.Errorf("failed to update payment status: %w", err)
	}

	if resp.UpdatePayments.AffectedRows == 0 {
		return fmt.Errorf("payment not found")
	}

	return nil
}

func (r *paymentRepository) GetByTransactionRef(ctx context.Context, transactionRef string) (*domain.Payment, error) {
	query := `query GetPaymentByRef($tx_ref: String!) {
        payments(where: { transaction_ref: { _eq: $tx_ref } }) {
            id
            user_id
            amount
            currency
            status
            transaction_ref
            payment_method
            created_at
            updated_at
        }
    }`

	vars := map[string]interface{}{
		"tx_ref": transactionRef,
	}

	var resp struct {
		Payments []domain.Payment `json:"payments"`
	}

	err := graphql.MutateRaw(ctx, query, vars, &resp)
	if err != nil {
		return nil, err
	}

	if len(resp.Payments) == 0 {
		return nil, nil
	}

	return &resp.Payments[0], nil
}

func (r *paymentRepository) GetTicketByPaymentID(ctx context.Context, paymentID string) (*domain.Ticket, error) {
	query := `query GetTicketByPayment($payment_id: uuid!) {
        tickets(where: { payment_id: { _eq: $payment_id } }) {
            id
            event_id
            user_id
            payment_id
            quantity
            total_price
            status
            created_at
            updated_at
        }
    }`

	vars := map[string]interface{}{
		"payment_id": paymentID,
	}

	var resp struct {
		Tickets []domain.Ticket `json:"tickets"`
	}

	err := graphql.MutateRaw(ctx, query, vars, &resp)
	if err != nil {
		return nil, err
	}

	if len(resp.Tickets) == 0 {
		return nil, nil
	}

	return &resp.Tickets[0], nil
}

func (r *paymentRepository) CheckDuplicateTicket(ctx context.Context, eventID, userID string) (bool, error) {
	query := `query CheckDuplicate($event_id: uuid!, $user_id: uuid!) {
        tickets(where: { 
            event_id: { _eq: $event_id }, 
            user_id: { _eq: $user_id },
            status: { _in: ["pending", "confirmed"] }
        }) {
            id
        }
    }`

	vars := map[string]interface{}{
		"event_id": eventID,
		"user_id":  userID,
	}

	var resp struct {
		Tickets []struct {
			ID string `json:"id"`
		} `json:"tickets"`
	}

	err := graphql.MutateRaw(ctx, query, vars, &resp)
	if err != nil {
		return false, err
	}

	return len(resp.Tickets) > 0, nil
}

func (r *paymentRepository) CreateTicket(ctx context.Context, ticket *domain.Ticket) (string, error) {
	mutation := `mutation InsertTicket($event_id: uuid!, $user_id: uuid!, $payment_id: uuid!, $quantity: Int!, $total_price: numeric!, $status: String!) {
        insert_tickets_one(object: {
            event_id: $event_id
            user_id: $user_id
            payment_id: $payment_id
            quantity: $quantity
            total_price: $total_price
            status: $status
        }) {
            id
        }
    }`

	vars := map[string]interface{}{
		"event_id":    ticket.EventID,
		"user_id":     ticket.UserID,
		"payment_id":  ticket.PaymentID,
		"quantity":    ticket.Quantity,
		"total_price": ticket.TotalPrice,
		"status":      ticket.Status,
	}

	var resp struct {
		InsertTicketsOne struct {
			ID string `json:"id"`
		} `json:"insert_tickets_one"`
	}

	err := graphql.MutateRaw(ctx, mutation, vars, &resp)
	if err != nil {
		return "", fmt.Errorf("failed to create ticket: %w", err)
	}

	return resp.InsertTicketsOne.ID, nil
}

func (r *paymentRepository) UpdateTicketStatus(ctx context.Context, paymentID, status string) error {
	mutation := `mutation UpdateTicket($payment_id: uuid!, $status: String!) {
        update_tickets(where: { payment_id: { _eq: $payment_id } }, _set: { status: $status }) {
            affected_rows
        }
    }`

	vars := map[string]interface{}{
		"payment_id": paymentID,
		"status":     status,
	}

	var resp struct {
		UpdateTickets struct {
			AffectedRows int `json:"affected_rows"`
		} `json:"update_tickets"`
	}

	err := graphql.MutateRaw(ctx, mutation, vars, &resp)
	if err != nil {
		return fmt.Errorf("failed to update ticket status: %w", err)
	}

	if resp.UpdateTickets.AffectedRows == 0 {
		return fmt.Errorf("ticket not found")
	}

	return nil
}

func (r *paymentRepository) GetEventByID(ctx context.Context, eventID string) (*domain.Event, error) {
	query := `query GetEvent($id: uuid!) {
        events(where: { id: { _eq: $id } }, limit: 1) {
            id
            title
            price
            status
        }
    }`

	vars := map[string]interface{}{
		"id": eventID,
	}

	var resp struct {
		Events []domain.Event `json:"events"`
	}

	err := graphql.MutateRaw(ctx, query, vars, &resp)
	if err != nil {
		return nil, err
	}

	if len(resp.Events) == 0 {
		return nil, nil
	}

	return &resp.Events[0], nil
}
