package proxy

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

var FHOST string
var ADDBYTEURL string

func InitUser() {
	url := os.Getenv("CONFIGURL")
	FHOST = os.Getenv("FAILTOHOST")
	ADDBYTEURL = os.Getenv("ADDBYTEURL")

	if url == "" {
		log.Println("CONFIGURL no set")
	} else {
		loadUserFromUrl(url)
		go TrafficNotify()
	}
}

func TrafficNotify() {
	var ut = make(map[string]uint64)
	for {
		for u, i := range Users {
			if inc := i.Traf - ut[u]; inc > 0 {
				AddTraf(u, i.CurrentRemote, inc)
			} else {
				continue
			}
			ut[u] = i.Traf
		}
		time.Sleep(time.Second * 60 * 2)
	}
}

func ParseHeader(from net.Conn) (user string, buf []byte) {
	buf = make([]byte, 100)
	n, err := from.Read(buf)
	if !(n > 0 && err == nil) {
		log.Println(err)
	}
	user = parseUser(string(buf))
	return user, buf[0:n]
}

func parseUser(dat string) string {
	lines := strings.Split(dat, "\n")
	if len(lines) > 2 {
		if strings.Contains(lines[1], "Host") {
			sl := strings.Split(lines[1], ":")
			return strings.Split(strings.TrimSpace(sl[1]), ".")[0]
		}
	}
	return ""
}

type Info struct {
	Traf          uint64
	PreferRemote  string
	CurrentRemote string
	Expired       time.Time
	TrafficLeft   uint64
}

var Users = make(map[string]*Info)

func loadUserFromUrl(url string) {
	rsp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
		return
	}

	if err := json.NewDecoder(rsp.Body).Decode(&Users); err != nil {
		log.Fatal(err)
		return
	}
	defer rsp.Body.Close()
}

func changeRemote(u, r string) (newRemote string) {
	newRemote = r
	if user, ok := Users[u]; ok {
		if user.PreferRemote != "" {
			newRemote = user.PreferRemote
			return
		}
	}
	return
}

func AddTraf(user, userAddr string, val uint64) {
	userIP := strings.Split(userAddr, ":")[0]
	req, err := http.NewRequest("GET", ADDBYTEURL, nil)
	if err != nil {
		log.Println(err)
	}
	q := req.URL.Query()
	q.Add("user", user)
	q.Add("userIP", userIP)
	q.Add("byte", fmt.Sprint(val))
	req.URL.RawQuery = q.Encode()
	client := http.Client{}
	if rsp, err := client.Do(req); err == nil {
		if dat, e := ioutil.ReadAll(rsp.Body); e == nil {
			if trf, ee := strconv.Atoi(string(dat)); ee == nil {
				Users[user].TrafficLeft = uint64(trf)
			}
		}
		rsp.Body.Close()
	}
}

func isValidUser(u string) bool {
	if u == "" {
		return false
	}
	if _, ok := Users[u]; !ok {
		return false
	}
	return Users[u].Expired.Unix() > time.Now().Unix()
}
