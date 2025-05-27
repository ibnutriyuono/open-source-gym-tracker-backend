package health

import (
	"caloria-backend/internal/helper/response"
	"encoding/json"
	"expvar"
	"net/http"
)

type HealthController struct{}

// HealthCheck caloria
//
// @Summary Get health
// @Description Get health status
// @Tags health
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /health [get]
func (uc *HealthController) HealthCheck(w http.ResponseWriter, r *http.Request) {

	selectedKeys := map[string]bool{
		"version":    true,
		"database":   true,
		"goroutines": true,
		"mem_stats":  true,
	}

	expvarMap := make(map[string]interface{})
	expvar.Do(func(kv expvar.KeyValue) {
		if selectedKeys[kv.Key] {
			var value interface{}
			_ = json.Unmarshal([]byte(kv.Value.String()), &value)
			expvarMap[kv.Key] = value
		}
	})

	response.SendJSON(w, http.StatusOK, expvarMap, "OK")
}
