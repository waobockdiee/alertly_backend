package getincidentsasreels

type Inputs struct {
	MinLatitude  float64 `uri:"min_latitude" binding:"required"`
	MaxLatitude  float64 `uri:"max_latitude" binding:"required"`
	MinLongitude float64 `uri:"min_longitude" binding:"required"`
	MaxLongitude float64 `uri:"max_longitude" binding:"required"`
}
