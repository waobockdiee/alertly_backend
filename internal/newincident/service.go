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
	MINIMUM_DURATION_HOURS = 24
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
	// ✅ OPTIMIZACIÓN: Geocoding asíncrono compatible
	// 1) Guardar incidente inmediatamente con dirección temporal
	addr, city, prov, postal := "Processing...", "Processing...", "Processing...", "..."

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
		// ✅ FIX: Asegurar que el InclId se mantiene para la respuesta
	} else {
		// 3) Lógica habitual de NUEVO CLUSTER o VOTO bayesiano
		cluster, err := s.repo.CheckAndGetIfClusterExist(incident)
		if err != nil && err != sql.ErrNoRows {
			return IncidentReport{}, err
		}

		if err == sql.ErrNoRows {
			// crear cluster nuevo…
			// 1. Obtenemos la duración "ideal" de la subcategoría
			proposedDuration, err := s.repo.GetDurationForSubcategory(incident.SubcategoryCode)
			if err != nil {
				return IncidentReport{}, fmt.Errorf("error getting subcategory duration: %w", err)
			}

			// 2. Aplicamos la lógica de duración mínima garantizada
			finalDurationHours := proposedDuration
			if finalDurationHours < MINIMUM_DURATION_HOURS {
				finalDurationHours = MINIMUM_DURATION_HOURS
			}

			// 3. Calculamos la fecha de finalización
			now := time.Now().UTC()
			end := now.Add(time.Duration(finalDurationHours) * time.Hour)

			cluster = Cluster{
				AccountId:       incident.AccountId,
				CreatedAt:       &now,
				StartTime:       &now,
				EndTime:         &end,
				MediaUrl:        incident.Media.Uri,
				CenterLatitude:  incident.Latitude,
				CenterLongitude: incident.Longitude,
				InsuId:          incident.InsuId,
				IsActive:        true, // Cluster activo por defecto
				MediaType:       incident.MediaType,
				EventType:       incident.EventType,
				Description:     incident.Description,
				Address:         addr, // Dirección temporal
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

	// ✅ FIX: Asegurar que InclId esté asignado en TODOS los paths
	// En el path de "update de posición", incident.InclId ya viene del frontend
	// En el path de "nuevo cluster", se asigna desde cluster.InclId arriba

	// 4) Ahora grabamos siempre el incident_report con dirección temporal
	incident.MediaUrl = incident.Media.Uri
	incident.Address = addr
	incident.City = city
	incident.Province = prov
	incident.PostalCode = postal

	fmt.Printf("IsAnonymousssss: %v\n", incident.IsAnonymous)

	inreId, err := s.repo.Save(incident)
	if err != nil {
		return IncidentReport{}, err
	}
	incident.InreId = inreId

	// 5) Actualizamos contador de perfil
	profSvc := profile.NewService(profile.NewRepository(database.DB))
	_ = profSvc.UpdateTotalIncidents(incident.AccountId)

	// ✅ OPTIMIZACIÓN: Geocoding asíncrono en background
	go func() {
		// Esperar un poco para no sobrecargar Nominatim
		time.Sleep(2 * time.Second)

		// Realizar geocoding
		realAddr, realCity, realProv, realPostal, err := common.ReverseGeocode(incident.Latitude, incident.Longitude)
		if err != nil {
			fmt.Printf("⚠️ Geocoding failed for incident %d: %v\n", inreId, err)
			return
		}

		// Actualizar cluster con dirección real
		if incident.InclId != 0 {
			if err := s.repo.UpdateClusterAddress(incident.InclId, realAddr, realCity, realProv, realPostal); err != nil {
				fmt.Printf("⚠️ Failed to update cluster address for %d: %v\n", incident.InclId, err)
			} else {
				fmt.Printf("✅ Geocoding completed for cluster %d: %s, %s\n", incident.InclId, realAddr, realCity)
			}
		}

		// Actualizar incident report con dirección real
		if err := s.repo.UpdateIncidentAddress(inreId, realAddr, realCity, realProv, realPostal); err != nil {
			fmt.Printf("⚠️ Failed to update incident address for %d: %v\n", inreId, err)
		} else {
			fmt.Printf("✅ Geocoding completed for incident %d: %s, %s\n", inreId, realAddr, realCity)
		}
	}()

	return incident, nil
}
