package cjblockuser

import (
	"log"
)

// Service orquesta la lógica de bloqueo de usuarios.
type Service struct {
	repo *Repository
}

// NewService crea una nueva instancia de Service.
func NewService(r *Repository) *Service {
	return &Service{repo: r}
}

// Run ejecuta la lógica principal del cronjob.
func (s *Service) Run() {
	log.Println("cjblockuser: Running user blocking cronjob...")

	// 1. Obtener usuarios que necesitan ser bloqueados
	usersToBlock, err := s.repo.FetchUsersToBlock()
	if err != nil {
		log.Printf("cjblockuser: Error fetching users to block: %v", err)
		return
	}

	if len(usersToBlock) == 0 {
		log.Println("cjblockuser: No users found to block.")
		return
	}

	// 2. Bloquear a cada usuario identificado
	for _, user := range usersToBlock {
		log.Printf("cjblockuser: Blocking account %d (reported %d times)", user.AccountID, user.ReportCount)
		err := s.repo.BlockUser(user.AccountID)
		if err != nil {
			log.Printf("cjblockuser: Error blocking account %d: %v", user.AccountID, err)
			continue // Continuar con el siguiente usuario a pesar del error
		}
		// Opcional: Insertar una notificación para el usuario bloqueado o para administradores
		// s.repo.InsertNotification(user.AccountID, "Tu cuenta ha sido bloqueada", "Debido a múltiples reportes.")
	}

	log.Printf("cjblockuser: User blocking cronjob finished. %d users blocked.", len(usersToBlock))
}
