package main

import (
	"net"
	"log"
	"time"
	"bufio"
	"encoding/json"
	"crypto/sha256"
)

var startGap = 5 * time.Second
var lookAroundGap = 5 * time.Second
var thinkGap = 10 * time.Millisecond

type reqLogin struct {
	Type		string
	Username	string
	Password	[]byte
	Version		string
}

type reqTypeOnly struct {
	Type		string
}

type reqBook struct {
	Type		string
	BookType	string
}

type reqAction struct {
	Type		string
	ActStr		string
	ActArg		string
	Nonce		int
}

func main() {
	bots := []string {
		"手持两把锟斤拷", "鱼", "大章鱼", "京狗",
		"aa7", "ZzZzZ", "0--0--0", "X.X",
		"HasName", "喵打", "term", "职业菜鸡"}

	for _, b := range bots {
		go loopBot(b, "iamarobot")
		time.Sleep(startGap)
	}

	forever := make(chan struct{})
	<-forever
}

type bot struct {
	conn		net.Conn
	chWrite		chan interface{}
}

func newBot(username, password string) *bot {
	bot := new(bot)
	bot.conn = login(username, password)
	bot.chWrite = make(chan interface{})
	return bot
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
		log.Fatalln(err);
	}

	shaPw := sha256.Sum256([]byte(password))
	reqLogin := reqLogin{"login", username, shaPw[:], "0.7.0"}
	jsonb, _ := json.Marshal(reqLogin)
	conn.Write(append(jsonb, '\n'))

	reply, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		log.Fatalln("srv ----", err.Error())
	}
	log.Print("srv ++++ " + string(reply))

	return conn
}

func (bot *bot) write(msg interface{}) {
	jsonb, err := json.Marshal(msg)
	if err != nil {
		log.Fatalln("write marshal:", err)
	}
	_, err = bot.conn.Write(append(jsonb, '\n'))
	if err != nil {
		log.Fatalln(err)
	}
}

func (bot *bot) readLoop() {
	reader := bufio.NewReader(bot.conn)
	for {
		reply, err := reader.ReadString('\n')
		if err != nil {
			log.Fatalln("srv ---- ", err)
		}
		var msg map[string]interface{}
		if err := json.Unmarshal([]byte(reply), &msg); err != nil {
			log.Fatalln("unmarshal srv reply:", err)
		}
		bot.readSwitch(msg)
	}
}

func (bot *bot) readSwitch(msg map[string]interface{}) {
	switch (msg["Type"]) {
	case "look-around":
		bot.handleLookAround(msg)
	case "start":
		msg := reqTypeOnly{"ready"}
		bot.write(msg)
	case "table":
		if msg["Event"] == "activated" {
			nonce := int(msg["Nonce"].(float64))
			msg := reqAction{"t-action","BOT","-1",nonce}
			time.Sleep(thinkGap)
			bot.write(msg)
		}
	case "update-user":
		// do nothing
	default:
		log.Fatalln("unknown reply type:", msg["Type"])
	}
}

func (bot *bot) handleLookAround(msg map[string]interface{}) {
	books := msg["Books"].(map[string]interface{})
	ds71 := books["DS71"].(map[string]interface{})
	bookable := ds71["Bookable"].(bool)
	if bookable {
		req := reqBook{"book", "DS71"}
		bot.write(req)
	}
}

func (bot *bot) lookAroundLoop() {
	for {
		req := reqTypeOnly{"look-around"}
		bot.write(req)
		time.Sleep(lookAroundGap)
	}
}


