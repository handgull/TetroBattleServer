package core

import (
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	roomsModels "github.com/handgull/TetroBattleServer/models/rooms"
)

// Hub definisce la struttura dove raccolgo tutte le chat room che si creano nel gioco
type Hub struct {
	rooms map[string]*Room
	*sync.RWMutex
}

// New fa partire una nuova istanza dell'hub dei match
func New() *Hub {
	hub := &Hub{
		rooms:   make(map[string]*Room),
		RWMutex: new(sync.RWMutex),
	}
	go func() {
		ch := make(chan os.Signal)
		signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM) // Se mi arriva SIGINT o SIGTERM lo notifico tramite channel
		<-ch                                               // Richiesta bloccante di un segnale sul channel
		hub.RLock()
		defer hub.RUnlock()
		for name, room := range hub.rooms { // Chiudo il canale Quit di ogni match (e quindi li termino)
			log.Printf("Closing room [%s] \n", name)
			room.Quit <- struct{}{}
		}
		os.Exit(0) // Termino l'esecuzione
	}()

	return hub
}

// AddClient aggiunge un nuovo client all'interno di una delle room in base ai parametri della richiesta
func (h *Hub) AddClient(cinf roomsModels.ClientInfo) {
	h.Lock()
	defer h.Unlock()
	room, ok := h.rooms[cinf.Room]
	if !ok {
		room = CreateRoom(cinf.Room)
	}
	room.AddClient(cinf.Name)
	h.rooms[cinf.Room] = room
}
