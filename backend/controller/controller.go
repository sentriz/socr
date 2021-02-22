package controller

import (
	"go.senan.xyz/socr/imagery"

	"github.com/blevesearch/bleve"
	"github.com/gorilla/websocket"
)

type Controller struct {
	Index          bleve.Index
	Directories    map[string]string
	SocketUpgrader websocket.Upgrader
	// SocketClientsSettings   map[*websocket.Conn]struct{}
	// SocketClientsScreenshot map[string]map[*websocket.Conn]struct{}
	// SocketUpdatesSettings   chan struct{}
	// SocketUpdatesScreenshot chan *screenshot.Screenshot
	HMACSecret    string
	LoginUsername string
	LoginPassword string
	APIKey        string
	DefaultFormat imagery.Format
}

// func (c *Controller) EmitUpdatesSettings() error {
// 	for range c.SocketUpdatesSettings {
// 		for client := range c.SocketClientsSettings {
// 			if err := client.WriteMessage(websocket.TextMessage, []byte(nil)); err != nil {
// 				log.Printf("error writing to socket client: %v", err)
// 				client.Close()
// 				delete(c.SocketClientsSettings, client)
// 				continue
// 			}
// 		}
// 	}
// 	return nil
// }

// func (c *Controller) EmitUpdatesScreenshot() error {
// 	for screenshot := range c.SocketUpdatesScreenshot {
// 		for client := range c.SocketClientsScreenshot[screenshot.ID] {
// 			if err := client.WriteMessage(websocket.TextMessage, []byte(nil)); err != nil {
// 				log.Printf("error writing to socket client: %v", err)
// 				client.Close()
// 				delete(c.SocketClientsScreenshot[screenshot.ID], client)
// 				continue
// 			}
// 		}
// 	}
// 	return nil
// }
