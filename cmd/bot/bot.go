package main

import (
	"flag"
	"hash/fnv"
	"regexp"

	"github.com/Sirupsen/logrus"
	"github.com/thetooth/curveface"
	irc "github.com/thoj/go-ircevent"
)

var (
	server, channel, nick, identPassword, curvePubKey string

	help = ""
)

// Session contains nick auth data
type Session struct {
	Source string
	Pub    string
	Priv   string
}

// Users map
type Users map[string]Session

func init() {
	flag.StringVar(&server, "server", "10.0.1.1:6667", "")
	flag.StringVar(&channel, "channel", "#test", "")
	flag.StringVar(&nick, "nick", "facebot", "")
	flag.StringVar(&identPassword, "identify", "password111", "")
	flag.StringVar(&curvePubKey, "pub_key", "P9/NIhsJvQAvnxoc2177O/3aIzUHYOjcack2wkkpTn2Q=", "")
	flag.Parse()
}

func main() {
	users := make(Users, 0)

	irccon := irc.IRC(nick, "curvefacebot")
	irccon.VerboseCallbackHandler = false
	irccon.Debug = true
	irccon.Version = "curvefacebot"
	//irccon.UseTLS = true
	//irccon.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	irccon.AddCallback("001", func(e *irc.Event) {
		irccon.Privmsgf("nicksrv", "identify %s", identPassword)
		irccon.Join(channel)
	})

	irccon.AddCallback("CTCP_FINGER", func(e *irc.Event) {
		irccon.SendRawf("NOTICE %s :\x01FINGER %s\x01", e.Nick, "o-onii-chan no")
	})

	irccon.AddCallback("PRIVMSG", func(e *irc.Event) {
		// Reply to
		var replyChan string
		if e.Arguments[0] == nick {
			replyChan = e.Nick
		} else {
			replyChan = channel
		}

		// Channel commands
		switch e.Message() {
		case "!tlsface":
			if user, ok := users[e.Nick]; ok && user.Source == e.Source {
				if res, err := curveface.GetGface(
					"127.0.0.1:9001",
					user.Priv,
					user.Pub,
					curvePubKey,
				); err != nil {
					irccon.Privmsg(e.Nick, "⚠⚠⚠ https://www.youtube.com/watch?v=RfiQYRn7fBg ⚠⚠⚠")
				} else {
					h := fnv.New32a()
					h.Write([]byte(user.Priv))
					i := h.Sum32()
					irccon.Privmsgf(replyChan, "%s [%X]", res[:], i)
					return
				}
			} else {
				irccon.Privmsg(e.Nick, "You need to be authenticated to use this service, try \"register\"")
				return
			}
		}

		// Commands with params
		if replyChan == e.Nick {
			if ok, _ := regexp.MatchString(`^register`, e.Message()); ok {
				r := regexp.MustCompile(`register\s(.{45})\s(.{45})`)
				if input := r.FindAllStringSubmatch(e.Message(), 1); input != nil && len(input[0]) == 3 {
					user := Session{Source: e.Source, Priv: input[0][1], Pub: input[0][2]}
					if _, err := curveface.GetGface(
						"127.0.0.1:9001",
						user.Priv,
						user.Pub,
						curvePubKey,
					); err != nil {
						irccon.Privmsg(e.User, "You are not authenticated! 😢🖕")
					} else {
						users[e.Nick] = user
						irccon.Privmsg(e.User, "You are now authenticated! 😂👌")
					}
					return
				}
				irccon.Privmsg(e.User, "Usage: register {privkey} {pubkey}\n\n"+help)
				return
			}
		}
	})

	irccon.AddCallback("PART", func(e *irc.Event) {
		delete(users, e.Nick)
	})
	irccon.AddCallback("QUIT", func(e *irc.Event) {
		delete(users, e.Nick)
	})
	irccon.AddCallback("NICK", func(e *irc.Event) {
		delete(users, e.Nick)
	})

	err := irccon.Connect(server)
	if err != nil {
		logrus.Fatal(err)
	}

	irccon.Loop()
}
