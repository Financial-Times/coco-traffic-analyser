package resources

import (
	"encoding/json"
	"net/http"

	"github.com/Financial-Times/coco-traffic-analyser/analyser"
	log "github.com/Sirupsen/logrus"
)

type AnalyserHandler struct {
	analyser *analyser.StandardAnalyser
}

func NewAnalyserHandler(a *analyser.StandardAnalyser) *AnalyserHandler {
	return &AnalyserHandler{a}
}

func (h *AnalyserHandler) ServeTrafficGraph(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	graph, err := json.Marshal(h.analyser.TrafficGraph().Matrix())
	if err != nil {
		log.WithError(err).Warn("Traffic graph endpoint")
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	_, err = w.Write(graph)
	if err != nil {
		log.WithError(err).Warn("Traffic graph endpoint")
		http.Error(w, "", http.StatusInternalServerError)
	}
}
