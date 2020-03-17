package core

import (
	"log"
	"sync"
)

// Room rappresenta la room dei player, con la chat ed i dati live legati agli score dell'avversario
type Room struct {
	name    string
	Msgch   chan string // Channel usato per fare da broadcast ai canali dei client
	clients map[string]chan<- string
	Quit    chan struct{}
	*sync.RWMutex
}

// Metodo usato per mandare un messaggio ad ogni client nella room
func (r *Room) broadcastMsg(msg string) {
	r.RLock()
	defer r.RUnlock()
	for _, cc := range r.clients {
		go func(ch chan<- string) {
			ch <- msg
		}(cc)
	}
}

// closeChatRoomSync chiude il canale di broadcast ed elimina tutti i client da una chat room. (blocking call)
func (r *Room) closeChatRoomSync() {
	r.Lock()
	defer r.Unlock()
	close(r.Msgch)
	for name := range r.clients {
		delete(r.clients, name)
	}
}

// Run avvia la chat room
func (r *Room) Run() {
	log.Println("Starting chat room", r.name)
	// Ogni messaggio inviato nel canale principale è inviato a tutti i client tramite il metodo di broadcast
	go func() {
		for msg := range r.Msgch {
			r.broadcastMsg(msg)
		}
	}()

	//handle when the quit channel is triggered
	go func() {
		<-r.Quit
		r.closeChatRoomSync()
	}()
}

// CreateRoom crea un nuovo match
func CreateRoom(name string) *Room {
	room := &Room{
		name:    name,
		Msgch:   make(chan string),
		RWMutex: new(sync.RWMutex),
		clients: make(map[string]chan<- string),
		Quit:    make(chan struct{}),
	}
	room.Run()
	return room
}

// RemoveClient rimuove un client da una chat room. (blocking call)
func (r *Room) RemoveClient(name string) {
	r.Lock()
	defer r.Unlock()
	log.Printf("Removing client [%s]... ", name)
	delete(r.clients, name)
	log.Printf("Done\n")
}

// AddClient aggiunge un nuovo client alla chat room
func (r *Room) AddClient(clientname string) {
	r.Lock()
	defer r.Unlock()
	if _, ok := r.clients[clientname]; ok {
		log.Printf("Client [%s] already exist in chat room [%s]", clientname, r.name)
		return
	}
	log.Printf("Adding client [%s]... ", clientname)
	wc, done := make(chan string), make(chan struct{})
	r.clients[clientname] = wc
	log.Printf("Done\n")

	// Se viene inviato un messaggio nel channel 'done' rimuovo il client dalla room
	go func() {
		<-done
		r.RemoveClient(clientname)
		close(done)
	}()
}
