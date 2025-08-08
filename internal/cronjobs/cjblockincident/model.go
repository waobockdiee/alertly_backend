package cjblockincident

// IncidentToReject representa un incidente que ha sido reportado y debe ser rechazado.
type IncidentToReject struct {
	IncidentID int64
	FlagCount  int
}
