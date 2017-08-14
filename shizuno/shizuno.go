package main

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"strconv"
	"time"

	"github.com/howeyc/gopass"
	"github.com/rolevax/hisa/model"
	"github.com/rolevax/hisa/netio"
)

func main() {
	fmt.Print("bot password: ")
	password, err := gopass.GetPasswd()
	if err != nil {
		log.Fatalln(err)
	}

	for i := 0; i < 40; i++ {
		username := "bot" + strconv.Itoa(i)
		go loopBot(username, string(password))
	}

	ch := make(chan struct{})
	<-ch
}

func signUp(username, password string) {
	conn, err := net.Dial("tcp", "127.0.0.1:6171")
	if err != nil {
		log.Fatalln(err)
	}

	shaPw := sha256.Sum256([]byte(password))
	reqLogin := &model.CsAuth{
		Type:     "sign-up",
		Version:  "0.8.2",
		Username: username,
		Password: base64.StdEncoding.EncodeToString(shaPw[:]),
	}
	jsonb, _ := json.Marshal(reqLogin)
	netio.Write(conn, jsonb)
	time.Sleep(1 * time.Second)
	conn.Close()
}

type bot struct {
	conn     net.Conn
	chWrite  chan interface{}
	username string
}

func newBot(username, password string) *bot {
	return &bot{
		conn:     login(username, password),
		chWrite:  make(chan interface{}),
		username: username,
	}
}

func (bot *bot) close() {
	bot.conn.Close()
}

func loopBot(username, password string) {
	bot := newBot(username, password)
	defer bot.close()

	go bot.lookAroundLoop()
	bot.readLoop()
}

func login(username, password string) net.Conn {
	conn, err := net.Dial("tcp", "127.0.0.1:6171")
	if err != nil {
		log.Fatalln(err)
	}

	shaPw := sha256.Sum256([]byte(password))
	reqLogin := &model.CsAuth{
		Type:     "login",
		Version:  "0.8.2",
		Username: username,
		Password: base64.StdEncoding.EncodeToString(shaPw[:]),
	}
	jsonb, _ := json.Marshal(reqLogin)
	netio.Write(conn, jsonb)

	_, err = netio.Read(conn)
	if err != nil {
		log.Fatalln(err.Error())
	}

	return conn
}

func (bot *bot) write(msg interface{}) {
	jsonb, err := json.Marshal(msg)
	if err != nil {
		log.Fatalln("write marshal:", err)
	}
	err = netio.Write(bot.conn, jsonb)
	if err != nil {
		log.Fatalln(err)
	}
	time.Sleep(10 * time.Millisecond)
	//log.Println(bot.username, "--->", string(jsonb))
}

func (bot *bot) readLoop() {
	for {
		reply, err := netio.Read(bot.conn)
		if err != nil {
			log.Fatalln("srv ---- ", err)
		}
		//log.Println(bot.username, "<---", string(reply))
		var msg map[string]interface{}
		if err := json.Unmarshal([]byte(reply), &msg); err != nil {
			log.Fatalln("unmarshal srv reply:", err)
		}
		bot.readSwitch(msg)
	}
}

func (bot *bot) readSwitch(msg map[string]interface{}) {
	switch msg["Type"] {
	case "look-around":
		bot.handleLookAround(msg)
	case "start":
		msg := &model.CsTypeOnly{"choose"} // girl index parsed 0
		bot.write(msg)
	case "chosen":
		msg := &model.CsTypeOnly{"ready"}
		bot.write(msg)
	case "table":
		if msg["Event"] == "activated" {
			nonce := int(msg["Nonce"].(float64))
			msg := &struct {
				Type, ActStr, ActArg string
				Nonce                int
			}{
				Type:   "t-action",
				ActStr: "BOT",
				ActArg: "-1",
				Nonce:  nonce,
			}
			bot.write(msg)
		}
	case "update-user":
	case "resume":
		// do nothing
	default:
		log.Fatalln("unknown reply type:", msg["Type"])
	}
}

func (bot *bot) handleLookAround(msg map[string]interface{}) {
	bot.tryBook(0)
	bot.tryBook(1)
	bot.tryBook(2)
}

func (bot *bot) tryBook(x model.BookType) {
	req := &struct {
		Type     string
		BookType model.BookType
	}{
		Type:     "book",
		BookType: x,
	}
	bot.write(req)
}

func (bot *bot) lookAroundLoop() {
	for {
		req := model.CsTypeOnly{"look-around"}
		bot.write(req)
		time.Sleep(1 * time.Second)
	}
}
