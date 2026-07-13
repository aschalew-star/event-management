package domain

// Request types
type UploadEventImagesRequest struct {
	EventID string   `json:"event_id"`
	Images  []string `json:"images"` // Base64 image payload strings
}

type UploadEventImagesPayload struct {
	Input            UploadEventImagesRequest `json:"input"`
	SessionVariables map[string]string        `json:"session_variables"`
}

type DeleteEventImageRequest struct {
	ImageID string `json:"image_id"`
}

type DeleteEventImagePayload struct {
	Input            DeleteEventImageRequest `json:"input"`
	SessionVariables map[string]string       `json:"session_variables"`
}

type SetFeaturedImageRequest struct {
	EventID string `json:"event_id"`
	ImageID string `json:"image_id"`
}

type SetFeaturedImagePayload struct {
	Input            SetFeaturedImageRequest `json:"input"`
	SessionVariables map[string]string       `json:"session_variables"`
}

// Response types
type UploadImageResponse struct {
	Success bool     `json:"success"`
	Message string   `json:"message"`
	URLs    []string `json:"urls"`
}

type ImageUploadResult struct {
	URL      string `json:"url"`
	PublicID string `json:"public_id"`
}

// Domain models
type EventImage struct {
	ID         string `json:"id"`
	EventID    string `json:"event_id"`
	ImageURL   string `json:"image_url"`
	IsFeatured bool   `json:"is_featured"`
}
