package get

import (
	"log/slog"
	"net/http"
	"strconv"

	"github.com/headblockhead/landmine/backend"
	"github.com/headblockhead/landmine/models"
	"github.com/headblockhead/landmine/respond"
)

func New(log *slog.Logger, backend backend.Backend) Handler {
	return Handler{
		log:     log,
		backend: backend,
	}
}

type Handler struct {
	log     *slog.Logger
	backend backend.Backend
}

func intOrDefault(s string, def int) int {
	if s == "" {
		return def
	}
	i, err := strconv.Atoi(s)
	if err != nil {
		return def
	}
	return i
}

func parseBool(s string) (v, ok bool) {
	if s != "true" && s != "false" {
		return false, false
	}
	return s == "true", true
}

func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var req models.ListGetRequest
	req.BaseID = r.PathValue("baseID")
	req.TableIDOrName = r.PathValue("tableIDOrName")

	q := r.URL.Query()
	validationFailures := make(map[string]string)

	var ok bool

	req.TimeZone = q.Get("timeZone")
	req.UserLocale = q.Get("userLocale")
	req.PageSize = intOrDefault(q.Get("pageSize"), 100)
	req.MaxRecords = intOrDefault(q.Get("maxRecords"), -1)
	req.Offset = q.Get("offset")
	req.View = q.Get("view")
	req.Sort = models.NewSortDirection(q.Get("sort"))
	req.FilterByFormula = q.Get("filterByFormula")
	req.CellFormat = models.NewCellFormat(q.Get("cellFormat"))
	req.Fields = q["fields"]
	if q.Get("returnFieldsByFieldId") != "" {
		if req.ReturnFieldsByFieldId, ok = parseBool(q.Get("returnFieldsByFieldId")); !ok {
			validationFailures["returnFieldsByFieldId"] = "invalid boolean"
		}
	}
	req.RecordMetadata = q["recordMetadata"]

	if len(validationFailures) > 0 {
		br := models.NewBadRequest("failed to evaluate query strings", validationFailures)
		respond.WithJSON(w, br, http.StatusBadRequest)
		return
	}

	resp, err := h.backend.List(r.Context(), req)
	if err != nil {
		h.log.Error("failed to list records", slog.Any("error", err))
		respond.WithJSON(w, models.NewError("failed to list records", http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	respond.WithJSON(w, resp, http.StatusOK)
}
