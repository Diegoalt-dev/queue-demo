package internal

import (
	"context"
	"encoding/json"
	"net/http"
)

type TaskRequest struct {
	Payload string `json:"payload"`
}

type SQSSender interface {
	Send(ctx context.Context, payload string) error
}

func CreateTask(sender SQSSender) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var taskRequest TaskRequest
		if err := json.NewDecoder(r.Body).Decode(&taskRequest); err != nil {
			http.Error(w, "invalid JSON", http.StatusBadRequest)
			return
		}

		if err := sender.Send(r.Context(), taskRequest.Payload); err != nil {
			http.Error(w, "failed to send task", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusAccepted)
	}
}
