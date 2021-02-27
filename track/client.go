package track

import (
	"fmt"
	"strings"
)

// Client contains information about a client.
type Client struct {
	Nickname  string
	ClientID  int
	ProjectID int
}

// String returns a string representation of the Client object.
func (client *Client) String() string {
	return fmt.Sprintf("Nickname: %s\nClient ID: %d\nProject ID: %d\n", client.Nickname, client.ClientID, client.ProjectID)
}

// AddClient adds a new client to a list of clients.
func AddClient(clients []Client, newClient Client) ([]Client, error) {
	for _, client := range clients {
		if client.Nickname == newClient.Nickname {
			return []Client{}, fmt.Errorf("client with nickname %s already exists", client.Nickname)
		}
	}
	return append(clients, newClient), nil
}

func matches(client Client, params Client) bool {
	if params.Nickname != "" && strings.ToLower(client.Nickname) == strings.ToLower(params.Nickname) {
		return true
	}
	if params.ClientID > 0 && client.ClientID == params.ClientID {
		return true
	}

	if params.ProjectID > 0 && client.ProjectID == params.ProjectID {
		return true
	}
	return false
}

// FilterClients filters clients by the client nickname, id or project id.
func FilterClients(clients []Client, params Client) []Client {
	if (params == Client{}) {
		return clients
	}
	res := []Client{}
	for _, client := range clients {
		if matches(client, params) {
			res = append(res, client)
		}
	}
	return res
}
