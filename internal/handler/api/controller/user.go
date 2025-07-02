package controller

import (
	"base/internal/domain/constants"
	"base/internal/handler/api/request"
	"encoding/json"
	"net/http"

	"go.opentelemetry.io/otel"
)

// RunPayroll godoc
// @Summary      Run payroll
// @Description  Admin runs payroll for a given attendance period
// @Tags         Payroll
// @Accept       json
// @Produce      json
// @Param        data  body      request.PayrollRequest  true  "RunPayroll payload"
// @Success      200   {object}  entity.PayrollResponse
// @Failure      400   {string}  string  "Bad request"
// @Failure      500   {string}  string  "Internal server error"
// @Router       /v1/payroll/run [post]
func (h *Handler) GetUser(w http.ResponseWriter, r *http.Request) {
	ctx, span := otel.Tracer(constants.Controller).Start(r.Context(), "Controller.CreateAttendancePeriod")
	defer span.End()
	var req request.GetByIdReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err := req.Validate()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	resp, err := h.userInteractor.GetByID(ctx, &req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(resp)
}
