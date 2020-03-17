package rooms

// ClientInfo è una struttura usata per aggiungere un client ad una room
type ClientInfo struct {
	Room string `json:"room"`
	Name string `json:"name"`
}
