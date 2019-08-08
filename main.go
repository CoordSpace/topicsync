package main

import (
	"sync"
	"log"
	"os"
	"os/signal"
	"syscall"
	//"time"
	"github.com/bwmarrin/discordgo"
	"github.com/lrstanley/girc"
	"github.com/spf13/viper"
)

type Context struct {
	irc *girc.Client
	discord *discordgo.Session
	channelMap 	map[string]string 
	wg *sync.WaitGroup
}

func createIRC() *girc.Client {
	// create a new IRC client using the config information
	config := girc.Config{
		Server: viper.GetString("irc.server"),
		Port:   viper.GetInt("irc.port"),
		Nick:   viper.GetString("irc.nick"),
		User:   viper.GetString("irc.user"),
		Name:   viper.GetString("irc.name"),
	}
	// Setup our IRC client instance
	return girc.New(config)
}

func StartIRC(c *Context) {
	c.irc.Handlers.Add(girc.CONNECTED, c.joinHandlerIRC)

	c.irc.Handlers.Add(girc.TOPIC, c.topicHandlerIRC)

	c.irc.Handlers.Add(girc.JOIN, c.inChannelHandlerIRC)

	if err := c.irc.Connect(); err != nil {
		log.Fatalf("an error occurred while attempting to connect to %s: %s", c.irc.Server(), err)
	}
}

func (context *Context) joinHandlerIRC(c *girc.Client, e girc.Event) {
	log.Print("Connected to IRC server!")
	log.Print("Joining IRC channel...")
	c.Cmd.Join(viper.GetString("irc.channel"))
}

func (context *Context) topicHandlerIRC(c *girc.Client, e girc.Event) {
	log.Print("Topic changed in IRC!")
	context.discord.ChannelMessageSend(context.channelMap[e.Params[0]], "Topic changed in IRC!")
}

func (context *Context) inChannelHandlerIRC(c *girc.Client, e girc.Event) {
	if(e.Source.Name == c.GetNick()) {
		log.Print("Joined Channel!")
		context.wg.Done()
	}
}

func createDiscord() *discordgo.Session {
	// Create a new Discord session using the provided bot token.
	dg, err := discordgo.New("Bot " + viper.GetString("discord.token"))

	if err != nil {
		log.Fatalf("Error creating Discord session,", err)
	}
	return dg
}

func StartDiscord(c *Context) {
	// Register the messageCreate func as a callback for ChannelUpdate events.
	c.discord.AddHandler(c.channelUpdateHandlerDiscord)

	// Callback to handle the ready event when the bot is connected
	c.discord.AddHandler(c.clientReadyHandlerDiscord)

	log.Print("Connecting to Discord server...")
	// Open a websocket connection to Discord and begin listening.
	err := c.discord.Open()
	if err != nil {
		log.Fatalf("Error opening connection to Discord,", err)
	}
}

func (context *Context) channelUpdateHandlerDiscord(s *discordgo.Session, u *discordgo.ChannelUpdate) {
	log.Print("Channel Update in ", u.Name, ": ", u.Topic)
	// check the event channel name against the one in the config
	// then start the mutex topic update
	log.Print("Channel ID: ", u.ID, " Paired with: ", context.channelMap[u.ID])
	context.irc.Cmd.Message(context.channelMap[u.ID], "Topic changed in Discord!")
}

func (context *Context) clientReadyHandlerDiscord(s *discordgo.Session, u *discordgo.Ready) {
	// set Discord client as ready
	log.Print("Discord client is ready.")
	context.wg.Done()
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
		log.Fatalf("Fatal error config file: %s \n", err)
	}
}

func setupChannelMap() map[string]string {
	// populate the channel connection mapping
	channelMap := make(map[string]string)
	// Bi-directional mapping between channels that are paired
	channelMap[viper.GetString("irc.channel")] = viper.GetString("discord.channelID")
	channelMap[viper.GetString("discord.channelID")] = viper.GetString("irc.channel")
	log.Print(channelMap)
	return channelMap
}

func main() {

	SetupConfig()

	var wg sync.WaitGroup
	context := Context {
		irc: 		createIRC(),
		discord: 	createDiscord(),
		channelMap:	setupChannelMap(),
		wg: 		&wg,
	}

	context.wg.Add(1)
	log.Print("Starting IRC...")
	go StartIRC(&context)
	log.Print("Waiting on IRC...")
	context.wg.Wait()
	log.Print("IRC is connected and ready to go!")

	context.wg.Add(1)
	log.Print("Starting Discord...")
	go StartDiscord(&context)
	log.Print("Waiting on Discord...")
	context.wg.Wait()
	log.Print("Discord is connected and ready to go!")

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
}