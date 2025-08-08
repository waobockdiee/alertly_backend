package cjuserank

import (
	"fmt"
	"log"
)

// Service orquesta la lógica de otorgamiento de rangos de usuario.
type Service struct {
	repo *Repository
}

// NewService crea una nueva instancia de Service.
func NewService(r *Repository) *Service {
	return &Service{repo: r}
}

// Run ejecuta la lógica principal del cronjob.
func (s *Service) Run() {
	log.Println("cjuserank: Running user rank cronjob...")

	// 1. Cargar las definiciones de rangos
	ranks := getRanks() // Función auxiliar para obtener los rangos

	// 2. Obtener el score de todos los usuarios activos
	usersScore, err := s.repo.FetchUsersScore()
	if err != nil {
		log.Printf("cjuserank: Error fetching users score: %v", err)
		return
	}

	if len(usersScore) == 0 {
		log.Println("cjuserank: No active users found to process.")
		return
	}

	// 3. Procesar cada usuario
	for _, user := range usersScore {
		// Obtener los rangos que el usuario ya ha ganado
		earnedRanks, err := s.repo.GetEarnedRanksForUser(user.AccountID)
		if err != nil {
			log.Printf("cjuserank: Error fetching earned ranks for account %d: %v", user.AccountID, err)
			continue // Continuar con el siguiente usuario
		}

		// Crear un mapa para una búsqueda rápida de rangos ya ganados
	earnedMap := make(map[string]bool)
		for _, er := range earnedRanks {
			earnedMap[fmt.Sprintf("%s-%s", er.Type, er.Name)] = true
		}

		// Iterar sobre todos los rangos posibles
		for _, rank := range ranks {
			// Verificar si el usuario ha alcanzado el rango y no ha sido notificado antes
			if user.Score >= rank.ScoreMin && user.Score <= rank.ScoreMax && !earnedMap[fmt.Sprintf("%s-%s", "user_rank", rank.Title)] {
				// El usuario ha alcanzado un nuevo rango
				log.Printf("cjuserank: Account %d reached new rank: %s (Score: %d)",
					user.AccountID, rank.Title, user.Score)

				// Insertar el rango en account_achievements
				err := s.repo.InsertEarnedRank(user.AccountID, rank)
				if err != nil {
					log.Printf("cjuserank: Error inserting earned rank for account %d, rank %s: %v", user.AccountID, rank.Title, err)
					continue // Continuar, pero registrar el error
				}

				// Insertar notificación push
				notificationTitle := "¡Felicidades! Nuevo Rango Alcanzado"
				notificationMessage := fmt.Sprintf("¡Has alcanzado el rango de \"%s\" con un score de %d!", rank.Title, user.Score)
				err = s.repo.InsertNotification(user.AccountID, notificationTitle, notificationMessage)
				if err != nil {
					log.Printf("cjuserank: Error inserting notification for account %d, rank %s: %v", user.AccountID, rank.Title, err)
					// No es crítico, podemos continuar
				}
			}
		}
	}
	log.Println("cjuserank: User rank cronjob finished.")
}

// getRanks carga las definiciones de rangos.
func getRanks() []RankItem {
	return []RankItem{
		{Title: "New Neighbor", Code: "new_neighbor", BackgroundColor: "#CFD8DC", TextColor: "#546E7A", Description: "Earn up to 500 Citizen Score by reporting your first incidents and joining the community.", Image: "badges/1A.png", ScoreMin: 0, ScoreMax: 500},
		{Title: "Community Champion", Code: "community_champion", BackgroundColor: "#AED581", TextColor: "#33691E", Description: "Reach 501–1500 Citizen Score through consistent incident reporting and community support.", Image: "badges/1B.png", ScoreMin: 501, ScoreMax: 1500},
		{Title: "Neighborhood Legend", Code: "neighborhood_legend", BackgroundColor: "#FFF176", TextColor: "#F9A825", Description: "Earn 1500–3000 Citizen Score by making significant contributions in your neighborhood.", Image: "badges/1C.png", ScoreMin: 1500, ScoreMax: 3000},
		{Title: "Urban Guardian", Code: "urban_guardian", BackgroundColor: "#4FC3F7", TextColor: "#01579B", Description: "Collect 3000–6000 Citizen Score by actively reporting and verifying incidents across the city.", Image: "badges/1C.png", ScoreMin: 3000, ScoreMax: 6000},
		{Title: "Civic Hero", Code: "civic_hero", BackgroundColor: "#FF8A65", TextColor: "#BF360C", Description: "Achieve 6000–10000 Citizen Score by leading positive civic actions and reliable reporting.", Image: "badges/1C.png", ScoreMin: 6000, ScoreMax: 10000},
		{Title: "Maple Leaf Icon", Code: "maple_leaf_icon", BackgroundColor: "#E57373", TextColor: "#B71C1C", Description: "Gain 10000–20000 Citizen Score to be recognized as a national leader in civic engagement.", Image: "userrank/user_rank1.png", ScoreMin: 10000, ScoreMax: 20000},
	}
}
