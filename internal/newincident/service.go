package newincident

import (
	"alertly/internal/common"
	"alertly/internal/database"
	"alertly/internal/profile"
	"database/sql"
	"errors"
	"fmt"
	"time"
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

var cluster Cluster
var clusterId int64

func (s *service) Save(incident IncidentReport) (IncidentReport, error) {
	address, city, province, postalCode, errGeo := common.ReverseGeocode(incident.Latitude, incident.Longitude)

	if errGeo != nil {
		return IncidentReport{}, errGeo
	}

	var err error

	cluster, err = s.repo.CheckAndGetIfClusterExist(incident)
	if err != nil {

		// return incident, err
		if err == sql.ErrNoRows {
			// No existe un cluster: se debe crear uno nuevo
			// Ejemplo: cluster, err = s.repo.CreateCluster(incident)
			now := time.Now().UTC()
			endTime := time.Now().Add(24 * time.Hour)

			fmt.Println("SUBCATEGORY", incident.SubcategoryCode)
			fmt.Println("CATEGORY", incident.CategoryCode)
			cluster = Cluster{
				CreatedAt:       &now,
				StartTime:       &now,
				EndTime:         &endTime,
				MediaUrl:        incident.Media.Uri,
				CenterLatitude:  incident.Latitude,
				CenterLongitude: incident.Longitude,
				InsuId:          incident.InsuId,
				MediaType:       incident.MediaType,
				EventType:       incident.EventType,
				Description:     incident.Description,
				Address:         address,
				City:            city,
				Province:        province,
				PostalCode:      postalCode,
				SubcategoryName: incident.SubCategoryName,
				SubcategoryCode: incident.SubcategoryCode,
				CategoryCode:    incident.CategoryCode,
			}
			cluster.InclId, err = s.repo.SaveCluster(cluster, incident.AccountId)
			clusterId = cluster.InclId

			if err != nil {
				return IncidentReport{}, err
			}

		} else {
			// Ocurrió otro error en la consulta
			return IncidentReport{}, err
		}

	} else {
		// EXISTE y actualiza el cluster
		clusterId = cluster.InclId
		result, err := s.repo.UpdateCluster(cluster.InclId, incident)
		if err != nil {
			return IncidentReport{}, errors.New("error updating cluster. try later")
		}

		rowsAffected, err := result.RowsAffected()
		if err != nil || rowsAffected == 0 {
			return IncidentReport{}, fmt.Errorf("no rows affected or error getting rows affected: %w", err)
		}

	}

	// path := incident.Media.Uri
	// if _, err := os.Stat(path); err != nil {
	// 	fmt.Printf("Error: file does not exist or cannot be accessed: %v\n", err)
	// 	return IncidentReport{}, err
	// }

	incident.MediaUrl = incident.Media.Uri
	incident.InclId = clusterId
	incident.Address = address
	incident.City = city
	incident.Province = province
	incident.PostalCode = postalCode

	resultt, err := s.repo.Save(incident)
	if err != nil {
		return IncidentReport{}, err
	}

	repo := profile.NewRepository(database.DB)
	cs := profile.NewService(repo)
	err = cs.UpdateTotalIncidents(incident.AccountId)

	if err != nil {
		fmt.Println(err)
	}

	incident.InreId = resultt
	return incident, nil
}
