package control

import (
	"encoding/json"
	"net/http"

	"studentgit.kata.academy/Zhodaran/go-kata/internal/service"
)

type Controller struct {
	geoService service.GeoServicer
}

func NewController(geoService service.GeoServicer) *Controller {
	return &Controller{geoService: geoService}
}

type ErrorResponse struct {
	Message string `json:"message"` // Сообщение об ошибке
	Code    int    `json:"code"`    // Код ошибки
}

func (c *Controller) GetGeoCoordinatesAddress(w http.ResponseWriter, r *http.Request) {
	var req service.RequestAddressSearch
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	geo, err := c.geoService.GetGeoCoordinatesAddress(req.Query)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	jsonData, err := json.Marshal(geo)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)
}

func (c *Controller) GetGeoCoordinatesGeocode(w http.ResponseWriter, r *http.Request) {
	var req service.GeocodeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	geo, err := c.geoService.GetGeoCoordinatesGeocode(req.Lat, req.Lng)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	jsonData, err := json.Marshal(geo)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)
}
