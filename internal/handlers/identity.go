package handlers

import (
	"bitespeed-identity-reconciliation/internal/models"
    "bitespeed-identity-reconciliation/internal/services"
    "bitespeed-identity-reconciliation/pkg/utils"
    "net/http"
)

type IdentifyHandler struct {
    identityService *services.IdentityService
}

func NewIdentifyHandler() *IdentifyHandler {
    return &IdentifyHandler{
        identityService: services.NewIdentityService(),
    }
}

func (h *IdentifyHandler) Identify(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        utils.WriteError(w, http.StatusMethodNotAllowed, 
            nil, "Method not allowed. Use POST.")
        return
    }

	var req models.IdentifyRequest
    if err := utils.ParseJSON(r, &req); err != nil {
        utils.WriteError(w, http.StatusBadRequest, err, 
            "Invalid JSON in request body")
        return
    }

	response, err := h.identityService.IdentifyContact(&req)
    if err != nil {
        utils.WriteError(w, http.StatusBadRequest, err, 
            "Failed to process identity request")
        return
    }

	utils.WriteJSON(w, http.StatusOK, response)
}