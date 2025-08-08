package cjbadgeearn

import (
	"fmt"
	"log"
)

// Service orquesta la lógica de otorgamiento de insignias.
type Service struct {
	repo *Repository
}

// NewService crea una nueva instancia de Service.
func NewService(r *Repository) *Service {
	return &Service{repo: r}
}

// Run ejecuta la lógica principal del cronjob.
func (s *Service) Run() {
	log.Println("cjbadgeearn: Running badge earning cronjob...")

	// 1. Cargar las definiciones de insignias
	badges := getBadges() // Función auxiliar para obtener las insignias

	// 2. Obtener la actividad de todos los usuarios activos
	usersActivity, err := s.repo.FetchUsersActivity()
	if err != nil {
		log.Printf("cjbadgeearn: Error fetching users activity: %v", err)
		return
	}

	if len(usersActivity) == 0 {
		log.Println("cjbadgeearn: No active users found to process.")
		return
	}

	// 3. Procesar cada usuario
	for _, user := range usersActivity {
		// Obtener las insignias que el usuario ya ha ganado
		earnedBadges, err := s.repo.GetEarnedBadgesForUser(user.AccountID)
		if err != nil {
			log.Printf("cjbadgeearn: Error fetching earned badges for account %d: %v", user.AccountID, err)
			continue // Continuar con el siguiente usuario
		}

		// Crear un mapa para una búsqueda rápida de insignias ya ganadas
	earnedMap := make(map[string]bool)
		for _, eb := range earnedBadges {
			// Usamos una combinación de Type (código de categoría) y Name (título de la insignia)
			// para identificar de forma única una insignia ganada.
			earnedMap[fmt.Sprintf("%s-%s-%d", eb.Type, eb.Name, eb.BadgeThreshold)] = true
		}

		// Iterar sobre todas las insignias posibles
		for _, badge := range badges {
			// Obtener el valor del contador del usuario para esta categoría de insignia
			userCounterValue := getCounterValue(user, badge.Code)

			// Verificar si el usuario ha alcanzado el umbral y no ha ganado esta insignia antes
			if userCounterValue >= badge.Number && !earnedMap[fmt.Sprintf("%s-%s-%d", badge.Code, badge.Title, badge.Number)] {
				// El usuario ha ganado una nueva insignia
				log.Printf("cjbadgeearn: Account %d earned new badge: %s (Category: %s, Threshold: %d)",
					user.AccountID, badge.Title, badge.Category, badge.Number)

				// Insertar la insignia en account_achievements
				err := s.repo.InsertEarnedBadge(user.AccountID, badge)
				if err != nil {
					log.Printf("cjbadgeearn: Error inserting earned badge for account %d, badge %s: %v", user.AccountID, badge.Title, err)
					continue // Continuar, pero registrar el error
				}

				// Insertar notificación push
				notificationTitle := "¡Felicidades! Nueva Insignia Ganada"
				notificationMessage := fmt.Sprintf("¡Has ganado la insignia \"%s\" por tu contribución en \"%s\"!", badge.Title, badge.Category)
				err = s.repo.InsertNotification(user.AccountID, notificationTitle, notificationMessage)
				if err != nil {
					log.Printf("cjbadgeearn: Error inserting notification for account %d, badge %s: %v", user.AccountID, badge.Title, err)
					// No es crítico, podemos continuar
				}
			}
		}
	}
	log.Println("cjbadgeearn: Badge earning cronjob finished.")
}

// getCounterValue es una función auxiliar para obtener el valor del contador
// de UserActivity basado en el código de la insignia.
func getCounterValue(user UserActivity, code string) int {
	switch code {
	case "total_reports":
		return user.CounterTotalIncidentsCreated
	case "incident_as_update":
		return user.IncidentAsUpdate
	case "crime":
		return user.Crime
	case "traffic_accident":
		return user.TrafficAccident
	case "medical_emergency":
		return user.MedicalEmergency
	case "fire_incident":
		return user.FireIncident
	case "vandalism":
		return user.Vandalism
	case "suspicious_activity":
		return user.SuspiciousActivity
	case "infrastructure_issues":
		return user.InfrastructureIssues
	case "extreme_weather":
		return user.ExtremeWeather
	case "community_events":
		return user.CommunityEvents
	case "dangerous_wildlife_sighting":
		return user.DangerousWildlifeSighting
	case "positive_actions":
		return user.PositiveActions
	case "lost_pet":
		return user.LostPet
	default:
		return 0
	}
}

// getBadges carga las definiciones de insignias.
// En un entorno de producción, esto podría cargarse desde un archivo de configuración
// o una base de datos para mayor flexibilidad.
func getBadges() []Badge {
	return []Badge{
		{Title: "Contributor", Category: "Total Reports", Image: "badges/1A.png", Description: "Awarded for your first 10 reports in any category.", Number: 10, Code: "total_reports"},
		{Title: "Enthusiast", Category: "Total Reports", Image: "badges/1B.png", Description: "Awarded for 25 total reports, showing consistent engagement.", Number: 50, Code: "total_reports"},
		{Title: "Power User", Category: "Total Reports", Image: "badges/1C.png", Description: "Awarded for 75 total reports, a top-tier contributor across the board.", Number: 100, Code: "total_reports"},
		{Title: "Evidence Collector", Category: "Visual Evidence", Image: "badges/2A.png", Description: "Awarded for your first photo or video upload.", Number: 10, Code: "incident_as_update"},
		{Title: "Proof Provider", Category: "Visual Evidence", Image: "badges/2B.png", Description: "Awarded for 10 visual-evidence uploads that helped confirm incidents.", Number: 25, Code: "incident_as_update"},
		{Title: "Master Investigator", Category: "Visual Evidence", Image: "badges/2C.png", Description: "Awarded for 50 visual-evidence uploads, your proof skills are unmatched.", Number: 50, Code: "incident_as_update"},
		{Title: "Crime Sentinel", Category: "Crime", Image: "badges/3A.png", Description: "Awarded for your very first crime report.", Number: 1, Code: "crime"},
		{Title: "Crime Guardian", Category: "Crime", Image: "badges/3B.png", Description: "Awarded for 10 crime reports, helping keep your streets safe.", Number: 10, Code: "crime"},
		{Title: "Crime Champion", Category: "Crime", Image: "badges/3C.png", Description: "Awarded for 50 crime reports, making you a true crime-fighting leader.", Number: 50, Code: "crime"},
		{Title: "Road Watcher", Category: "Traffic Accident", Image: "badges/4A.png", Description: "Awarded for your first traffic-accident report.", Number: 1, Code: "traffic_accident"},
		{Title: "Traffic Protector", Category: "Traffic Accident", Image: "badges/4B.png", Description: "Awarded for 10 traffic-accident reports, helping drivers stay alert.", Number: 10, Code: "traffic_accident"},
		{Title: "Road Guardian", Category: "Traffic Accident", Image: "badges/4C.png", Description: "Awarded for 50 traffic-accident reports, saving lives on the road.", Number: 50, Code: "traffic_accident"},
		{Title: "First Responder", Category: "Medical Emergency", Image: "badges/5A.png", Description: "Awarded for your first medical-emergency report.", Number: 1, Code: "medical_emergency"},
		{Title: "Health Advocate", Category: "Medical Emergency", Image: "badges/5B.png", Description: "Awarded for 10 medical-emergency reports, supporting timely care.", Number: 10, Code: "medical_emergency"},
		{Title: "Care Champion", Category: "Medical Emergency", Image: "badges/5C.png", Description: "Awarded for 50 medical-emergency reports, leading your community’s health efforts.", Number: 50, Code: "medical_emergency"},
		{Title: "Flame Watcher", Category: "Fire Incident", Image: "badges/6A.png", Description: "Awarded for your first fire-incident report.", Number: 1, Code: "fire_incident"},
		{Title: "Fire Guardian", Category: "Fire Incident", Image: "badges/6B.png", Description: "Awarded for 10 fire-incident reports, helping contain danger.", Number: 10, Code: "fire_incident"},
		{Title: "Fire Commander", Category: "Fire Incident", Image: "badges/6C.png", Description: "Awarded for 50 fire-incident reports, protecting lives and property.", Number: 50, Code: "fire_incident"},
		{Title: "Urban Defender", Category: "Vandalism", Image: "badges/7A.png", Description: "Awarded for your first vandalism report.", Number: 1, Code: "vandalism"},
		{Title: "City Protector", Category: "Vandalism", Image: "badges/7B.png", Description: "Awarded for 10 vandalism reports, preserving public spaces.", Number: 10, Code: "vandalism"},
		{Title: "Civic Champion", Category: "Vandalism", Image: "badges/7C.png", Description: "Awarded for 50 vandalism reports, defending your neighborhood’s integrity.", Number: 50, Code: "vandalism"},
		{Title: "Alert Watcher", Category: "Suspicious Activity", Image: "badges/8A.png", Description: "Awarded for your first suspicious-activity report.", Number: 1, Code: "suspicious_activity"},
		{Title: "Vigilance Guardian", Category: "Suspicious Activity", Image: "badges/8B.png", Description: "Awarded for 10 suspicious-activity reports, raising community alertness.", Number: 10, Code: "suspicious_activity"},
		{Title: "Vigilance Champion", Category: "Suspicious Activity", Image: "badges/8C.png", Description: "Awarded for 50 suspicious-activity reports, your vigilance keeps everyone safer.", Number: 50, Code: "suspicious_activity"},
		{Title: "Infrastructure Scout", Category: "Infrastructure Issue", Image: "badges/9A.png", Description: "Awarded for your first infrastructure-issue report.", Number: 1, Code: "infrastructure_issues"},
		{Title: "Urban Keeper", Category: "Infrastructure Issue", Image: "badges/9B.png", Description: "Awarded for 10 infrastructure-issue reports, maintaining community services.", Number: 10, Code: "infrastructure_issues"},
		{Title: "Civic Overseer", Category: "Infrastructure Issue", Image: "badges/9C.png", Description: "Awarded for 50 infrastructure-issue reports, safeguarding vital public assets.", Number: 50, Code: "infrastructure_issues"},
		{Title: "Weather Watcher", Category: "Extreme Weather", Image: "badges/10A.png", Description: "Awarded for your first extreme-weather report.", Number: 1, Code: "extreme_weather"},
		{Title: "Storm Guardian", Category: "Extreme Weather", Image: "badges/10B.png", Description: "Awarded for 10 extreme-weather reports, helping others prepare.", Number: 10, Code: "extreme_weather"},
		{Title: "Elemental Master", Category: "Extreme Weather", Image: "badges/10C.png", Description: "Awarded for 50 extreme-weather reports, leading resilience against nature’s forces.", Number: 50, Code: "extreme_weather"},
		{Title: "Community Connector", Category: "Community Events", Image: "badges/11A.png", Description: "Awarded for your first community-event report.", Number: 1, Code: "community_events"},
		{Title: "Neighborhood Champion", Category: "Community Events", Image: "badges/11B.png", Description: "Awarded for 10 community-event reports, bringing people together.", Number: 10, Code: "community_events"},
		{Title: "Civic Ambassador", Category: "Community Events", Image: "badges/11C.png", Description: "Awarded for 50 community-event reports, strengthening local bonds.", Number: 50, Code: "community_events"},
		{Title: "Wildlife Watcher", Category: "Dangerous Wildlife Sighting", Image: "badges/12A.png", Description: "Awarded for your first wildlife-sighting report.", Number: 1, Code: "dangerous_wildlife_sighting"},
		{Title: "Wildlife Guardian", Category: "Dangerous Wildlife Sighting", Image: "badges/12B.png", Description: "Awarded for 10 wildlife-sighting reports, alerting others to hazards.", Number: 10, Code: "dangerous_wildlife_sighting"},
		{Title: "Wildlife Defender", Category: "Dangerous Wildlife Sighting", Image: "badges/12C.png", Description: "Awarded for 50 wildlife-sighting reports, keeping people and animals safe.", Number: 50, Code: "dangerous_wildlife_sighting"},
		{Title: "Kindness Starter", Category: "Positive Actions", Image: "badges/13A.png", Description: "Awarded for your first positive-action report.", Number: 1, Code: "positive_actions"},
		{Title: "Altruism Advocate", Category: "Positive Actions", Image: "badges/13B.png", Description: "Awarded for 10 positive-action reports, spreading goodwill.", Number: 10, Code: "positive_actions"},
		{Title: "Community Hero", Category: "Positive Actions", Image: "badges/13C.png", Description: "Awarded for 50 positive-action reports, inspiring everyone around you.", Number: 50, Code: "positive_actions"},
		{Title: "Pet Savior", Category: "Lost Pet", Image: "badges/14A.png", Description: "Awarded for your first lost-pet report.", Number: 1, Code: "lost_pet"},
		{Title: "Animal Rescuer", Category: "Lost Pet", Image: "badges/14B.png", Description: "Awarded for 10 lost-pet reports, reuniting families.", Number: 10, Code: "lost_pet"},
		{Title: "Pet Champion", Category: "Lost Pet", Image: "badges/14C.png", Description: "Awarded for 50 lost-pet reports, your efforts make tails wag.", Number: 50, Code: "lost_pet"},
	}
}
