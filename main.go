package main

import (
	"fmt"
	"sync"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
	"github.com/bwmarrin/discordgo"
	"github.com/lrstanley/girc"
	"github.com/spf13/viper"
)

type ircClient struct {
	client		*girc.Client
	isReady		bool
}

type discordClient struct {
	session		*discordgo.Session
	isReady		bool
}

type topicInfo struct {
	topic		string
	timeSet		time.Time
	mutex		sync.Mutex
}

var (
	i	*ircClient
	d	*discordClient
	t	*topicInfo
)

func StartIRC() {
	log.Print("Starting up IRC...")
	config := girc.Config{
		Server: viper.GetString("irc.server"),
		Port:   viper.GetInt("irc.port"),
		Nick:   viper.GetString("irc.nick"),
		User:   viper.GetString("irc.user"),
		Name:   viper.GetString("irc.name"),
	}

	// Setup our IRC client instance
	i := &ircClient {
		client: girc.New(config),
		isReady: false,
	}

	log.Print("Setting up IRC Handlers...")
	i.client.Handlers.Add(girc.CONNECTED, func(c *girc.Client, e girc.Event) {
		c.Cmd.Join(viper.GetString("irc.channel"))
		// Set IRC client as ready
		i.isReady = true
	})

	i.client.Handlers.Add(girc.TOPIC, func (c *girc.Client, e girc.Event) {
		log.Println("Topic changed! - Channel: ", e.Params[0], ". Topic: ", e.Last())
	})

	log.Print("Connecting to IRC server...")
	if err := i.client.Connect(); err != nil {
		log.Fatalf("an error occurred while attempting to connect to %s: %s", i.client.Server(), err)
	}
}

func StartDiscord() {
	log.Print("Starting up Discord...")

	// Create a new Discord session using the provided bot token.
	dg, err := discordgo.New("Bot " + viper.GetString("discord.token"))

	if err != nil {
		log.Fatalf("Error creating Discord session,", err)
		return
	}

	d := &discordClient {
		session: dg,
		isReady: false,
	}

	// Register the messageCreate func as a callback for ChannelUpdate events.
	d.session.AddHandler(func(s *discordgo.Session, u *discordgo.ChannelUpdate) {
		log.Print("Channel Update in ", u.Name, ": ", u.Topic)
		// check the event channel name against the one in the config
		// then start the mutex topic update
	})

	// Callback to handle the ready event when the bot is connected
	d.session.AddHandler(func(s *discordgo.Session, m *discordgo.Ready) {
		// set Discord client as ready
		d.isReady = true
	})

	log.Print("Connecting to Discord server...")
	// Open a websocket connection to Discord and begin listening.
	err = d.session.Open()
	if err != nil {
		log.Fatalf("Error opening connection to Discord,", err)
		return
	}
}

func SetupConfig() {
	log.Print("Setting up viper config...")
	// TODO: Add defaults for bot branding
	viper.SetConfigName("config")
	viper.AddConfigPath("/etc/topicsync/")
	viper.AddConfigPath("$HOME/.topicsync")
	viper.AddConfigPath(".")
	log.Print("Reading configuration file...")
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
}

func main() {
	log.Print("Starting the service...")

	// load config from file
	SetupConfig()

	// initialize the topic checker
	t := &topicInfo{
		topic: "",
		timeSet: time.Now(),
	}

	log.Print(t.timeSet)

	// start up irc goroutine
	go StartIRC()

	// start up discord goroutine
	go StartDiscord()

	// Halt main() till kill is received
	log.Print("Program will now run till sigterm is received...")
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	killSignal := <-interrupt
	switch killSignal {
	case os.Interrupt:
		log.Print("Got SIGINT...")
	case syscall.SIGTERM:
		log.Print("Got SIGTERM...")
	}

	log.Print("The service is shutting down...")
	// gracefully kill discord & irc goroutines
	log.Print("Done")
}