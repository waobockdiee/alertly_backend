package cjblockincident

import (
	"log"
)

// Service orquesta la lógica de rechazo de incidentes.
type Service struct {
	repo *Repository
}

// NewService crea una nueva instancia de Service.
func NewService(r *Repository) *Service {
	return &Service{repo: r}
}

// Run ejecuta la lógica principal del cronjob.
func (s *Service) Run() {
	log.Println("cjblockincident: Running incident blocking cronjob...")

	// 1. Obtener incidentes que necesitan ser rechazados
	incidentsToReject, err := s.repo.FetchIncidentsToReject()
	if err != nil {
		log.Printf("cjblockincident: Error fetching incidents to reject: %v", err)
		return
	}

	if len(incidentsToReject) == 0 {
		log.Println("cjblockincident: No incidents found to reject.")
		return
	}

	// 2. Rechazar cada incidente identificado
	for _, incident := range incidentsToReject {
		log.Printf("cjblockincident: Rejecting incident %d (flagged %d times)", incident.IncidentID, incident.FlagCount)
		err := s.repo.RejectIncident(incident.IncidentID)
		if err != nil {
			log.Printf("cjblockincident: Error rejecting incident %d: %v", incident.IncidentID, err)
			continue // Continuar con el siguiente incidente a pesar del error
		}
		// Opcional: Insertar una notificación para el creador del incidente o para administradores
	}

	log.Printf("cjblockincident: Incident blocking cronjob finished. %d incidents rejected.", len(incidentsToReject))
}
