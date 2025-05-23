package newincident

import (
	"alertly/internal/common"
	"alertly/internal/database"
	"alertly/internal/profile"
	"database/sql"
	"fmt"
	"time"
)

const (
	CLUSTER_EXPIRES_IN = 48 // in hours
)

type Service interface {
	Save(incident IncidentReport) (IncidentReport, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

// service.go (paquete newincident)

// service.go

func (s *service) Save(incident IncidentReport) (IncidentReport, error) {
	// 1) Geocoding (lo necesitas siempre para el report)
	addr, city, prov, postal, err := common.ReverseGeocode(incident.Latitude, incident.Longitude)
	if err != nil {
		return IncidentReport{}, err
	}

	// 2) **Si viene incl_id Y NO viene vote, es solo un update de posición**
	if incident.InclId != 0 && incident.Vote == nil {
		// actualizamos únicamente la ubicación del cluster
		if _, err := s.repo.UpdateClusterLocation(
			incident.InclId,
			incident.Latitude,
			incident.Longitude,
		); err != nil {
			return IncidentReport{}, fmt.Errorf("updating cluster location: %w", err)
		}
		// después seguimos a grabar el report
	} else {
		// 3) Lógica habitual de NUEVO CLUSTER o VOTO bayesiano
		cluster, err := s.repo.CheckAndGetIfClusterExist(incident)
		if err != nil && err != sql.ErrNoRows {
			return IncidentReport{}, err
		}

		if err == sql.ErrNoRows {
			// crear cluster nuevo…
			now := time.Now().UTC()
			end := now.Add(CLUSTER_EXPIRES_IN * time.Hour)
			cluster = Cluster{
				CreatedAt:       &now,
				StartTime:       &now,
				EndTime:         &end,
				MediaUrl:        incident.Media.Uri,
				CenterLatitude:  incident.Latitude,
				CenterLongitude: incident.Longitude,
				InsuId:          incident.InsuId,
				MediaType:       incident.MediaType,
				EventType:       incident.EventType,
				Description:     incident.Description,
				Address:         addr,
				City:            city,
				Province:        prov,
				PostalCode:      postal,
				SubcategoryName: incident.SubCategoryName,
				SubcategoryCode: incident.SubcategoryCode,
				CategoryCode:    incident.CategoryCode,
			}
			cluster.InclId, err = s.repo.SaveCluster(cluster, incident.AccountId)
			if err != nil {
				return IncidentReport{}, err
			}
		} else {
			// existe → aplicamos voto si viene y no ha votado ya
			voted, _, err := s.repo.HasAccountVoted(cluster.InclId, incident.AccountId)
			if err != nil {
				return IncidentReport{}, fmt.Errorf("checking vote history: %w", err)
			}
			if !voted && incident.Vote != nil {
				if *incident.Vote {
					_, err = s.repo.UpdateClusterAsTrue(
						cluster.InclId,
						incident.AccountId,
						incident.Latitude,
						incident.Longitude,
					)
				} else {
					_, err = s.repo.UpdateClusterAsFalse(
						cluster.InclId,
						incident.AccountId,
						incident.Latitude,
						incident.Longitude,
					)
				}
				if err != nil {
					return IncidentReport{}, fmt.Errorf("update cluster vote: %w", err)
				}
			}
		}
		// nos aseguramos de fijar el clusterId para el report
		incident.InclId = cluster.InclId
	}

	// 4) Ahora grabamos siempre el incident_report
	incident.MediaUrl = incident.Media.Uri
	incident.Address = addr
	incident.City = city
	incident.Province = prov
	incident.PostalCode = postal

	inreId, err := s.repo.Save(incident)
	if err != nil {
		return IncidentReport{}, err
	}
	incident.InreId = inreId

	// 5) Actualizamos contador de perfil
	profSvc := profile.NewService(profile.NewRepository(database.DB))
	_ = profSvc.UpdateTotalIncidents(incident.AccountId)

	return incident, nil
}
