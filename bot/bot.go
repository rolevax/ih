package main

import (
	"bufio"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"github.com/howeyc/gopass"
	"log"
	"math/rand"
	"net"
	"sync"
	"time"
)

var startGap = 500 * time.Millisecond
var lookAroundGap = 5 * time.Second
var stayLimit = 5
var maxBotPerTable = 2

var chLookAroundTicket = make(chan struct{})

type observe struct {
	wait    int
	waitBot int
	stay    int
	mutex   sync.Mutex
}

var prevs = [2]observe{}

func thinkGap(pass bool) time.Duration {
	//return time.Duration(10) * time.Millisecond
	if pass {
		return time.Duration(500+rand.Intn(500)) * time.Millisecond
	} else {
		r1 := rand.Intn(300)
		r2 := rand.Intn(300)
		r3 := rand.Intn(300)
		return time.Duration(1000+r1+r2+r3) * time.Millisecond
	}
}

type reqLogin struct {
	Type     string
	Username string
	Password []byte
	Version  string
}

type reqTypeOnly struct {
	Type string
}

type reqBook struct {
	Type     string
	BookType int
}

type reqAction struct {
	Type   string
	ActStr string
	ActArg string
	Nonce  int
}

func main() {
	rand.Seed(time.Now().UnixNano())
	fmt.Print("bot password: ")
	password, err := gopass.GetPasswd()
	if err != nil {
		log.Fatalln(err)
	}

	bots := []string{
		"手持两把锟斤拷", "鱼", "大章鱼", "京狗",
		"aa7", "ZzZzZ", "0--0--0", "X.X",
		"HasName", "喵打", "term", "职业菜鸡"}

	perm := rand.Perm(len(bots))

	for _, p := range perm {
		go loopBot(bots[p], string(password))
		time.Sleep(startGap)
	}

	for {
		chLookAroundTicket <- struct{}{}
		time.Sleep(lookAroundGap)
	}
}

type bot struct {
	conn     net.Conn
	chWrite  chan interface{}
	username string
}

func newBot(username, password string) *bot {
	bot := new(bot)
	bot.conn = login(username, password)
	bot.chWrite = make(chan interface{})
	bot.username = username
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
		log.Fatalln(err)
	}

	shaPw := sha256.Sum256([]byte(password))
	reqLogin := reqLogin{"login", username, shaPw[:], "0.7.6"}
	jsonb, _ := json.Marshal(reqLogin)
	conn.Write(append(jsonb, '\n'))

	_, err = bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		log.Fatalln("srv ----", err.Error())
	}
	log.Println("srv ++++ ", username)

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
	switch msg["Type"] {
	case "look-around":
		bot.handleLookAround(msg)
	case "start":
		msg := reqTypeOnly{"choose"} // girl index parsed 0
		bot.write(msg)
	case "chosen":
		msg := reqTypeOnly{"ready"}
		bot.write(msg)
	case "table":
		if msg["Event"] == "activated" {
			nonce := int(msg["Nonce"].(float64))
			args := msg["Args"].(map[string]interface{})
			action := args["action"].(map[string]interface{})
			_, pass := action["PASS"]
			msg := reqAction{"t-action", "BOT", "-1", nonce}
			time.Sleep(thinkGap(pass))
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
	books := msg["Books"].([]interface{})
	ds71 := books[0].(map[string]interface{})
	bot.tryBook(0, ds71, &prevs[0])
	cs71 := books[1].(map[string]interface{})
	bot.tryBook(1, cs71, &prevs[1])
}

func (bot *bot) tryBook(x int, xs71 map[string]interface{}, ob *observe) {
	bookable := xs71["Bookable"].(bool)

	if bookable {
		ob.mutex.Lock()
		defer ob.mutex.Unlock()

		wait := int(xs71["Book"].(float64))
		if wait == 0 { // not a strict cond, but fine
			ob.waitBot = 0
		}

		if wait == ob.wait {
			ob.stay++
			if ob.stay >= stayLimit {
				if ob.waitBot < maxBotPerTable {
					ob.waitBot++
					ob.stay = 0
					req := reqBook{"book", x}
					bot.write(req)
				}
			}
		} else {
			ob.wait = wait
			ob.stay = 0
		}
	}
}

func (bot *bot) lookAroundLoop() {
	for {
		<-chLookAroundTicket
		req := reqTypeOnly{"look-around"}
		bot.write(req)
	}
}
