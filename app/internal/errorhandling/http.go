package errorhandling

type Error struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}
