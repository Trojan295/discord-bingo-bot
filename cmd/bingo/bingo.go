package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/Trojan295/discord-bingo-bot/pkg/bingo/discord"
	"github.com/Trojan295/discord-bingo-bot/pkg/bingo/game"
	"github.com/Trojan295/discord-bingo-bot/pkg/bingo/repository"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/bwmarrin/discordgo"
	"github.com/sirupsen/logrus"
)

var (
	token      string
	tableName  string
	controller *game.Controller
)

func init() {
	debug := flag.Bool("debug", false, "")
	flag.Parse()

	if *debug {
		logrus.SetLevel(logrus.DebugLevel)
	}

	var ok bool
	if token, ok = os.LookupEnv("DISCORD_BOT_TOKEN"); !ok {
		panic("Missing DISCORD_BOT_TOKEN")
	}

	if tableName, ok = os.LookupEnv("DYNAMODB_TABLE_NAME"); !ok {
		tableName = "bingo-games"
	}

	awsSession, err := session.NewSession()
	if err != nil {
		logrus.Panicf("cannot create AWS session: %s", err.Error())
		panic("Cannot create AWS session")
	}

	gameRepository, err := repository.NewDynamoDBGameRepository(awsSession, tableName)
	if err != nil {
		panic(err)
	}

	controller = game.NewController(gameRepository)

	rand.Seed(time.Now().UnixNano())
}

func main() {
	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}

	dg.AddHandler(messageCreate)

	err = dg.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	dg.Close()
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {

	if m.Author.ID == s.State.User.ID || !strings.HasPrefix(m.Content, ".bingo") {
		return
	}

	response, err := controller.ProcessMessage(m.ChannelID, m.Content)
	if err != nil {
		log.Printf("cannot process message: %v", err.Error())
		s.ChannelMessageSend(m.ChannelID, "Bingo bot currently not available...")
	}

	if response == nil {
		return
	}

	for _, msg := range discord.PrintMessage(response, m) {
		s.ChannelMessageSend(m.ChannelID, msg)
	}
}
