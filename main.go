package main

import (
	"errors"
	"flag"
	"github.com/PuerkitoBio/goquery"
	"github.com/bwmarrin/discordgo"
	"github.com/jonas747/dutil/commandsystem"
	"log"
	"strings"
)

var (
	flagToken string
	dgo       *discordgo.Session
)

func init() {
	flag.StringVar(&flagToken, "t", "", "Token to use")

	if !flag.Parsed() {
		flag.Parse()
	}
}

func main() {
	session, err := discordgo.New(flagToken)
	if err != nil {
		panic(err)
	}
	dgo = session

	cmd := &commandsystem.CommandDef{
		Arguments: []*commandsystem.ArgumentDef{
			&commandsystem.ArgumentDef{Name: "User", Type: commandsystem.ArgumentTypeUser},
		},
		RunInDm:                 true,
		IgnoreUserNotFoundError: true,
		RunFunc:                 HandleCommand,
	}

	cmdInvite := &commandsystem.CommandDef{
		Name:    "invite",
		RunInDm: true,
		RunFunc: func(parsed *commandsystem.ParsedCommand, m *discordgo.MessageCreate) {
			dgo.ChannelMessageSend(m.ChannelID, "Why would you even want me on your server? What the fuck is wrong with you? https://discordapp.com/oauth2/authorize?client_id=199556484550492161&scope=bot&permissions=101376")
		},
	}

	system := &commandsystem.CommandSystem{
		DefaultMentionHandler: cmd,
		Session:               session,
	}

	system.RegisterCommands(cmdInvite)

	dgo.AddHandler(system.HandleMessageCreate)
	dgo.AddHandler(HandleReady)
	dgo.AddHandler(HandleServerJoin)

	err = dgo.Open()
	if err != nil {
		panic(err)
	}
	log.Println("Started roody >:(")
	select {}
}

func HandleCommand(parsed *commandsystem.ParsedCommand, m *discordgo.MessageCreate) {
	mentionUser := m.Author.ID

	if parsed.Args[0] != nil {
		mentionUser = parsed.Args[0].DiscordUser().ID
	}

	insult, err := GetInsult()
	if err != nil {
		log.Println("Failed getting insult:", err)
		return
	}

	_, err = dgo.ChannelMessageSend(m.ChannelID, "<@"+mentionUser+"> "+insult)
	if err != nil {
		log.Println("Error sending insult:", err)
	}
}

func HandleReady(s *discordgo.Session, r *discordgo.Ready) {
	log.Println("Ready received! Connected to", len(s.State.Guilds), "Guilds")
}

func HandleServerJoin(s *discordgo.Session, g *discordgo.GuildCreate) {
	log.Println("Joined guild", g.Name, " Connected to", len(s.State.Guilds), "Guilds")
}

var ErrFailedFindingInsult = errors.New("Failed to find insult")

func GetInsult() (insult string, err error) {
	var doc *goquery.Document
	doc, err = goquery.NewDocument("http://www.insultgenerator.org/")
	if err != nil {
		return
	}

	doc.Find(".wrap").Each(func(i int, s *goquery.Selection) {
		// For each item found, get the band and title
		insult = strings.TrimSpace(s.Text())
	})

	if insult == "" {
		err = ErrFailedFindingInsult
	}
	return
}
