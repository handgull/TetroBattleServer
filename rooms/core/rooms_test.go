package core

import "testing"

func TestRoom(t *testing.T) {
	room := CreateRoom("test")

	if name := room.name; name != "test" {
		t.Errorf("Expected room name equals to \"test\", but got %v", name)
	}

	room.AddClient("c1")
	room.AddClient("c2")
	room.AddClient("c3")
	_, ok1 := room.clients["c1"]
	_, ok2 := room.clients["c2"]
	_, ok3 := room.clients["c3"]
	if nclients := len(room.clients); !ok1 || !ok2 || !ok3 || nclients != 3 {
		t.Errorf("Expected room have 3 clients in it, but got %v", nclients)
	}

	room.RemoveClient("c1")
	if _, okr1 := room.clients["c1"]; okr1 {
		t.Errorf("\"c1\" still present in the room")
	}

	room.Quit <- struct{}{}
	_, okq1 := room.clients["c1"]
	_, okq2 := room.clients["c2"]
	_, okq3 := room.clients["c3"]
	if nclients := len(room.clients); okq1 || okq2 || okq3 || nclients != 0 {
		t.Errorf("Expected room have 0 clients in it, but got %v", nclients)
	}
}
