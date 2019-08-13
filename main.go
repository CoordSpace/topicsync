package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/lrstanley/girc"
	"github.com/spf13/viper"
)

type context struct {
	irc          *girc.Client
	discord      *discordgo.Session
	channelPairs map[string]string
	topic        string
	waitGroup    *sync.WaitGroup
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

func startIRC(c *context) {
	c.irc.Handlers.Add(girc.CONNECTED, c.joinHandlerIRC)
	c.irc.Handlers.Add(girc.TOPIC, c.topicHandlerIRC)
	c.irc.Handlers.AddTmp(girc.JOIN, time.Minute, c.inChannelHandlerIRC)
	for {
		if err := c.irc.Connect(); err != nil {
			log.Printf("error: %s", err)

			log.Printf("reconnecting in %d seconds...", viper.GetInt("timeout"))
			time.Sleep((time.Duration)(viper.GetInt("timeout")) * time.Second)
		} else {
			return
		}
	}
}

func (context *context) joinHandlerIRC(c *girc.Client, e girc.Event) {
	log.Print("Connected to IRC server")
	// auth with IRC services, if config supports
	if viper.IsSet("irc.auth") {
		log.Print("Authing with IRC bot ", viper.GetString("irc.auth.bot"))
		c.Cmd.Message(viper.GetString("irc.auth.bot"), viper.GetString("irc.auth.cmd"))
	}
	if viper.IsSet("irc.usermode") {
		log.Print("Setting usermode(s)")
		c.Cmd.Mode(c.Config.Nick, viper.GetString("irc.usermode"))
	}
	log.Print("Joining IRC channel")
	c.Cmd.Join(viper.GetString("irc.channel"))
}

func (context *context) inChannelHandlerIRC(c *girc.Client, e girc.Event) bool {
	if e.Source.Name == c.GetNick() {
		log.Print("IRC is connected and ready to go!")
		context.waitGroup.Done()
	}
	return true
}

func (context *context) topicHandlerIRC(c *girc.Client, e girc.Event) {
	// e.Params[channelName, ..., channeltopic]
	log.Print("Topic changed in IRC channel ", e.Params[0])
	topicClean := strings.TrimSpace(e.Last())
	if topicClean != context.topic {
		log.Print("Cleaned topic for discord: ", topicClean)
		update := discordgo.ChannelEdit{
			Topic: topicClean,
		}
		// Update the discord channel topic
		_, err := context.discord.ChannelEditComplex(context.channelPairs[e.Params[0]], &update)
		if err != nil {
			log.Panicf("Error updating discord channel!: %s \n", err)
		}
		// Update the context topic to eliminate any update loops
		context.topic = topicClean
	} else {
		// Post the topic change alert to the channel when the syncing dust settles
		if viper.IsSet("updateFormat") {
			_, err := context.postTopicAlertDiscord(context.channelPairs[e.Params[0]], topicClean)
			if err != nil {
				log.Panicf("Error posting discord topic alert: %s \n", err)
			}
		}
		log.Print("Duplicate topic, ignoring")
	}
}

func createDiscord() *discordgo.Session {
	// Create a new Discord session using the provided bot token.
	dg, err := discordgo.New("Bot " + viper.GetString("discord.token"))

	if err != nil {
		log.Fatalf("Error creating Discord session: %s \n", err)
	}
	return dg
}

func startDiscord(c *context) {
	// Register the callback for ChannelUpdate events.
	c.discord.AddHandler(c.channelUpdateHandlerDiscord)

	// Callback to handle the client Ready event.
	c.discord.AddHandlerOnce(c.clientReadyHandlerDiscord)

	log.Print("Connecting to Discord server...")
	// Open a websocket connection to Discord and begin listening.
	for {
		if err := c.discord.Open(); err != nil {
			log.Printf("error: %s", err)

			log.Printf("reconnecting in %d seconds...", viper.GetInt("timeout"))
			time.Sleep((time.Duration)(viper.GetInt("timeout")) * time.Second)
		} else {
			return
		}
	}
}

func (context *context) channelUpdateHandlerDiscord(s *discordgo.Session, u *discordgo.ChannelUpdate) {
	log.Print("Discord ChannelUpdate in ", u.Name, ": ", u.Topic)
	log.Print("Channel ID: ", u.ID, " Paired with: ", context.channelPairs[u.ID])

	// Clean the topic string since discord automatically strips whitespace.
	topicClean := strings.TrimSpace(u.Topic)
	if topicClean != context.topic {
		log.Print("New topic for ", u.Name, ": ", topicClean)
		channel, ok := context.channelPairs[u.ID]
		// only post updates to channels we're tracking, ignore otherwise
		if ok {
			context.irc.Cmd.Topic(channel, topicClean)
			// Update the context topic to eliminate any update loops
			context.topic = topicClean
		} else {
			log.Print("Channel not in list of tracked channels, ignoring topic event")
		}
	} else {
		// Post the topic change alert to the channel when the syncing dust settles
		if viper.IsSet("updateFormat") {
			_, err := context.postTopicAlertDiscord(u.ID, topicClean)
			if err != nil {
				log.Panicf("Error posting discord topic alert: %s \n", err)
			}
		}
		log.Print("Duplicate topic, ignoring")
	}
}

func (context *context) clientReadyHandlerDiscord(s *discordgo.Session, u *discordgo.Ready) {
	// Now that the bot is online with Discord, continue program execution
	log.Print("Discord is connected and ready to go!")
	context.waitGroup.Done()
}

func (context *context) postTopicAlertDiscord(channelID string, topic string) (*discordgo.Message, error) {
	log.Print("Posting alert message in Discord")
	// Build formatted topic change alert message.
	alertMsg := fmt.Sprintf(viper.GetString("updateFormat"), randomEmoji(viper.GetStringSlice("emojis")), topic)
	// Post the topic change alert to the channel
	return context.discord.ChannelMessageSend(channelID, alertMsg)
}

func setupConfig() {
	log.Print("Setting up viper config")
	// Defaults
	viper.SetDefault("nick", "TopicBot")
	viper.SetDefault("user", "TopicBot")
	viper.SetDefault("name", "A Topic Sync Bot")
	viper.SetDefault("emojis", []string{"🔔"})

	viper.SetConfigName("config")
	viper.AddConfigPath("/etc/topicsync/")
	viper.AddConfigPath("$HOME/.topicsync")
	viper.AddConfigPath(".")
	log.Print("Reading configuration file")
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("Fatal error config file: %s \n", err)
	}
}

func setupchannelPairs() map[string]string {
	// populate the channel connection mapping
	channelPairs := make(map[string]string)
	// Bi-directional mapping between channels that are paired
	channelPairs[viper.GetString("irc.channel")] = viper.GetString("discord.channelID")
	channelPairs[viper.GetString("discord.channelID")] = viper.GetString("irc.channel")
	return channelPairs
}

func randomEmoji(array []string) string {
	if len(array) < 1 {
		log.Fatal("No items in list to choose from, please check your config")
	} else if len(array) == 1 {
		return array[0]
	}
	s := rand.NewSource(time.Now().Unix())
	r := rand.New(s)
	return array[r.Intn(len(array))]
}

func main() {

	setupConfig()

	var waitGroup sync.WaitGroup
	context := context{
		irc:          createIRC(),
		discord:      createDiscord(),
		channelPairs: setupchannelPairs(),
		waitGroup:    &waitGroup,
	}

	context.waitGroup.Add(1)
	log.Print("Starting IRC...")
	go startIRC(&context)
	log.Print("Waiting on IRC...")
	context.waitGroup.Wait()

	context.waitGroup.Add(1)
	log.Print("Starting Discord...")
	go startDiscord(&context)
	log.Print("Waiting on Discord...")
	context.waitGroup.Wait()

	// Halt main() till a kill signal is received
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
