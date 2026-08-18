package api

import (
	"github.com/marc-pango/marc-chatbot/backend/internal/bot"
	"github.com/marc-pango/marc-chatbot/backend/internal/domain"
	"net/http"
	"strconv"
	"time"
)

func (h *Handler) listMessages(w http.ResponseWriter, r *http.Request) {
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	messages, err := h.service.ListMessages(r.Context(), r.URL.Query().Get("groupId"), domain.MessageDirection(r.URL.Query().Get("direction")), limit)
	if err != nil {
		handleError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, messages)
}
func (h *Handler) mockMessage(w http.ResponseWriter, r *http.Request) {
	var body struct {
		ProviderMessageID string `json:"providerMessageId"`
		ProviderGroupID   string `json:"providerGroupId"`
		GroupName         string `json:"groupName"`
		SenderID          string `json:"senderId"`
		SenderName        string `json:"senderName"`
		Text              string `json:"text"`
		Self              bool   `json:"self"`
	}
	if err := decodeJSON(r, &body); err != nil {
		writeError(w, http.StatusBadRequest, "invalid_json", err.Error())
		return
	}
	err := h.service.InjectMockMessage(r.Context(), bot.Message{ProviderMessageID: body.ProviderMessageID, ProviderGroupID: body.ProviderGroupID, GroupName: body.GroupName, SenderID: body.SenderID, SenderName: body.SenderName, Text: body.Text, MessageType: "text", OccurredAt: time.Now(), Self: body.Self})
	if err != nil {
		handleError(w, err)
		return
	}
	writeJSON(w, http.StatusAccepted, map[string]string{"status": "accepted"})
}
