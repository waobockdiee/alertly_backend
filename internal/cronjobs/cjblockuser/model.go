package cjblockuser

// UserToBlock representa un usuario que ha sido reportado y debe ser bloqueado.
type UserToBlock struct {
	AccountID int64
	ReportCount int
}
