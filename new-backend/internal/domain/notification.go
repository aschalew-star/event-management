package domain

// Event Trigger Types
type EventTriggerPayload struct {
	Event struct {
		Op  string `json:"op"` // INSERT, UPDATE, DELETE
		New struct {
			ID          string  `json:"id"`
			Title       string  `json:"title"`
			Description string  `json:"description"`
			UserID      string  `json:"user_id"`
			Venue       string  `json:"venue"`
			Address     string  `json:"address"`
			EventDate   string  `json:"event_date"`
			StartTime   *string `json:"start_time"`
			EndTime     *string `json:"end_time"`
			Status      string  `json:"status"`
			Price       float64 `json:"price"`
			IsFree      bool    `json:"is_free"`
		} `json:"new"`
		Old struct {
			ID     string `json:"id"`
			Title  string `json:"title"`
			UserID string `json:"user_id"`
			Status string `json:"status"`
		} `json:"old"`
	} `json:"event"`
	Table struct {
		Name   string `json:"name"`
		Schema string `json:"schema"`
	} `json:"table"`
	CreatedAt string `json:"created_at"`
}

type Users struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	AvatarURL string `json:"avatar_url"`
}

type Follow struct {
	Follower User `json:"follower"`
}

// Email Types
type EmailRequest struct {
	ToEmail     string
	ToName      string
	Subject     string
	HTMLContent string
	PlainText   string
}

type EventNotification struct {
	EventID     string
	EventTitle  string
	EventUserID string
	EventVenue  string
	EventDate   string
	EventStatus string
	Operation   string
}

// Response Types
type NotificationResponse struct {
	Status  string `json:"status"`
	EventID string `json:"event_id"`
	Message string `json:"message"`
}
