package service

import (
	"HighArch/api"
	"HighArch/entity"
	"HighArch/storage"
	"encoding/json"
	"github.com/gorilla/websocket"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
	"sync"
	"time"
)

type FeedWsController interface {
	AddConnection(userId string, conn *websocket.Conn)
	HandleNewPostCreated(userId string, post entity.Post) error
}

type feedWsRabbitController struct {
	sync.RWMutex
	connections      map[string][]*websocket.Conn
	rabbitChannel    *amqp.Channel
	friendLinksStore storage.FriendLinksStore
}

func NewFeedWsController(rabbitChan *amqp.Channel, friendLinksStore storage.FriendLinksStore) FeedWsController {
	initRabbitExchange(rabbitChan)
	return &feedWsRabbitController{
		connections:      make(map[string][]*websocket.Conn),
		rabbitChannel:    rabbitChan,
		friendLinksStore: friendLinksStore,
	}
}

func (p *feedWsRabbitController) AddConnection(userId string, conn *websocket.Conn) {
	conn.SetCloseHandler(func(code int, text string) error {
		p.removeConnections(userId, []*websocket.Conn{conn})
		p.Lock()
		defer p.Unlock()
		conns, has := p.connections[userId]
		if !has || len(conns) == 0 {
			// no more connections, close rabbit consuming
			cancelRabbitQueue(p.rabbitChannel, userId)
		}
		return nil
	})
	// read and discard messages from the peer to handle socket disconnection
	go func() {
		for {
			if _, _, err := conn.NextReader(); err != nil {
				conn.Close()
				break
			}
		}
	}()
	p.Lock()
	defer p.Unlock()
	connections, has := p.connections[userId]
	if !has || len(connections) == 0 {
		// start Rabbit consuming only for first connection for the user
		listenRabbitQueue(p.rabbitChannel, userId, func(delivery amqp.Delivery) {
			p.sendMessageToConnections(userId, delivery.Body)
		})
	}
	p.connections[userId] = append(connections, conn)
}

func (p *feedWsRabbitController) HandleNewPostCreated(authorUserId string, post entity.Post) error {
	// TODO the same is doing at the same moment in FeedCacheController - could be optimized
	friendsIds, err := p.friendLinksStore.GetFriendsIds(authorUserId)
	if err != nil {
		log.Println(err)
		return err
	}
	var postModel = api.PostApiModel{
		Id:         post.Id,
		Text:       post.Text,
		AuthorId:   post.AuthorId,
		CreateTime: formatUnixTimestampToString(post.CreateTime, time.DateTime),
	}
	jsonPostModel, err := json.Marshal(postModel)
	if err != nil {
		return err
	}
	for _, friendId := range friendsIds {
		sendToRabbit(p.rabbitChannel, friendId, jsonPostModel)
	}
	return nil
}

func (p *feedWsRabbitController) sendMessageToConnections(userId string, message []byte) error {
	p.RLock()
	conns, has := p.connections[userId]
	p.RUnlock()
	if !has {
		return nil
	}
	for _, conn := range conns {
		log.Printf("WebSocket sending for %s", userId)
		err := conn.WriteMessage(websocket.TextMessage, message)
		if err != nil {
			log.Println(err)
		}
	}
	return nil
}

func (p *feedWsRabbitController) removeConnections(userId string, connections []*websocket.Conn) {
	p.Lock()
	defer p.Unlock()
	newConnections := make([]*websocket.Conn, 0)
	oldConnections := p.connections[userId]
	for _, conn := range oldConnections {
		var found = false
		for _, connToRemove := range connections {
			if conn == connToRemove {
				found = true
				break
			}
		}
		if !found {
			newConnections = append(newConnections, conn)
		}
	}
	p.connections[userId] = newConnections
}

// RabbitMQ interaction code

const (
	rabbitExchangeName string = "ws_feed_direct"
	rabbitExchangeType string = "direct"
)

func initRabbitExchange(rabbitChan *amqp.Channel) error {
	err := rabbitChan.ExchangeDeclare(
		rabbitExchangeName, // name
		rabbitExchangeType, // type
		true,               // durable
		false,              // auto-deleted
		false,              // internal
		false,              // no-wait
		nil,                // arguments
	)
	if err != nil {
		log.Println(err)
	}
	return err
}

func sendToRabbit(rabbitChan *amqp.Channel, userId string, message []byte) error {
	log.Printf("Send to Rabbit for %s", userId)
	err := rabbitChan.Publish(
		rabbitExchangeName, // exchange
		userId,             // routing key
		true,               // mandatory
		false,              // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        message,
		})
	if err != nil {
		log.Println(err)
	}
	return err
}

func listenRabbitQueue(rabbitChan *amqp.Channel, userId string, callback func(delivery amqp.Delivery)) {
	queue, err := rabbitChan.QueueDeclare(
		"",    // name
		false, // durable
		false, // delete when unused
		true,  // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		log.Println(err)
	}
	err = rabbitChan.QueueBind(
		queue.Name,         // queue name
		userId,             // routing key
		rabbitExchangeName, // exchange
		false,
		nil,
	)
	msgs, err := rabbitChan.Consume(
		queue.Name, // queue
		userId,     // consumer name
		true,       // auto ack
		false,      // exclusive
		false,      // no local
		false,      // no wait
		nil,        // args
	)
	go func() {
		for d := range msgs {
			log.Printf("Read from Rabbit for %s", userId)
			callback(d)
		}
	}()
}

func cancelRabbitQueue(rabbitChan *amqp.Channel, userId string) {
	log.Printf("Cancel Rabbit consuming for %s", userId)
	rabbitChan.Cancel(userId, true)
}
