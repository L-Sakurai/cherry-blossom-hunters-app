package service

import (
    // "encoding/json"
    // "strings"
    "fmt"
    "context"
    "github.com/bwmarrin/discordgo"
	"cherry-blossom-hunters-app/appConfig"
    "cherry-blossom-hunters-app/logger"

)

type MemberService struct {
	discordToken    string
	memberChannelID string
	ruleChannelID   string
	masterUserID    string
}

type MemberComplianceData struct {
    RuleReaders    []string `json:"rule_readers"`    // ルールを読んだメンバー（グッドマーク付き）
    ChannelMembers []string `json:"channel_members"` // チャネルの全メンバー
    MissingMembers []string `json:"missing_members"` // ルール未読メンバー
    CheckedAt      string   `json:"checked_at"`
}

func NewMemberService(appConfig *appConfig.Config) *MemberService {
	config := appConfig
    return &MemberService {
        discordToken:   config.Discord.Token,
        memberChannelID:    config.Discord.MemberChannelID,
        ruleChannelID:  config.Discord.RuleChannelID,
        masterUserID: config.Discord.MasterUserID,
    }
}

func botTokenStringBuilder(s *MemberService) string {
    return fmt.Sprintf("Bot %s", s.discordToken)
}

func (s *MemberService) CheckMemberComplianceWithContext(ctx context.Context) error {
    dg, err := discordgo.New(botTokenStringBuilder(s))
    if err != nil {
        return fmt.Errorf("Discord session creation error: %v", err)
    }
    s.getChannelMembers(dg, s.memberChannelID)
    return nil
}

func (s *MemberService) getChannelMembers(dg *discordgo.Session, channelId string) ([]string, error){
    channel, err := dg.Channel(channelId)
    if err != nil {
        return nil, err
    }

    rawMembers, err := dg.GuildMembers(channel.GuildID, "", 1000)
    
    var memberNames []string
    for _, m := range rawMembers {
        memberNames = append(memberNames, m.User.GlobalName)
    }
     logger.Logging("debug", "Users: %v", memberNames)
    return memberNames, nil
}

