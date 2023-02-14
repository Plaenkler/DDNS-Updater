package ddns

type UpdateRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Domain   string `json:"domain"`
	IP       string `json:"ip"`
}
