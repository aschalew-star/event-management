package domain

// Request types
type ProcessPaymentArgs struct {
	EventID   string `json:"eventId"`
	Quantity  int    `json:"quantity"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

type HasuraProcessPaymentPayload struct {
	Input struct {
		Input ProcessPaymentArgs `json:"input"`
	} `json:"input"`
	SessionVariables map[string]string `json:"session_variables"`
}

type VerifyPaymentArgs struct {
	TransactionRef string `json:"transactionRef"`
}

type HasuraVerifyPaymentPayload struct {
	Input struct {
		Input VerifyPaymentArgs `json:"input"`
	} `json:"input"`
	SessionVariables map[string]string `json:"session_variables"`
}

// Domain models
type Events struct {
	ID     string  `json:"id"`
	Title  string  `json:"title"`
	Price  float64 `json:"price"`
	Status string  `json:"status"`
}

type Payment struct {
	ID             string  `json:"id"`
	UserID         string  `json:"user_id"`
	Amount         float64 `json:"amount"`
	Currency       string  `json:"currency"`
	Status         string  `json:"status"`
	TransactionRef string  `json:"transaction_ref"`
	PaymentMethod  string  `json:"payment_method"`
	CreatedAt      string  `json:"created_at"`
	UpdatedAt      string  `json:"updated_at"`
}

type Ticket struct {
	ID         string  `json:"id"`
	EventID    string  `json:"event_id"`
	UserID     string  `json:"user_id"`
	PaymentID  string  `json:"payment_id"`
	Quantity   int     `json:"quantity"`
	TotalPrice float64 `json:"total_price"`
	Status     string  `json:"status"`
	CreatedAt  string  `json:"created_at"`
	UpdatedAt  string  `json:"updated_at"`
}

// Chapa API models
type ChapaRequest struct {
	Amount        float64                `json:"amount"`
	Currency      string                 `json:"currency"`
	Email         string                 `json:"email"`
	FirstName     string                 `json:"first_name"`
	LastName      string                 `json:"last_name"`
	PhoneNumber   string                 `json:"phone_number"`
	TxRef         string                 `json:"tx_ref"`
	CallbackURL   string                 `json:"callback_url"`
	ReturnURL     string                 `json:"return_url"`
	Customization map[string]interface{} `json:"customization"`
	Meta          map[string]interface{} `json:"meta"`
}

type ChapaResponse struct {
	Status  string      `json:"status"`
	Message interface{} `json:"message"`
	Data    struct {
		CheckoutURL string `json:"checkout_url"`
	} `json:"data"`
}

type ChapaVerifyResponse struct {
	Status string `json:"status"`
	Data   struct {
		Status    string  `json:"status"`
		Reference string  `json:"reference"`
		Amount    float64 `json:"amount"`
	} `json:"data"`
}

// Response types
type ProcessPaymentResponse struct {
	Success        bool   `json:"success"`
	Message        string `json:"message"`
	CheckoutURL    string `json:"checkout_url"`
	TransactionRef string `json:"transaction_ref"`
}

type PaymentVerificationResult struct {
	Success      bool    `json:"success"`
	Message      string  `json:"message"`
	Status       string  `json:"status"`
	PaymentID    string  `json:"payment_id"`
	TicketID     string  `json:"ticket_id"`
	EventID      string  `json:"event_id"`
	Quantity     int     `json:"quantity"`
	TotalPrice   float64 `json:"total_price"`
	TicketStatus string  `json:"ticket_status"`
}

type VerifyPaymentResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Status  string `json:"status"`
	Ticket  struct {
		ID         string  `json:"id"`
		EventID    string  `json:"event_id"`
		Quantity   int     `json:"quantity"`
		TotalPrice float64 `json:"total_price"`
		Status     string  `json:"status"`
	} `json:"ticket"`
}
