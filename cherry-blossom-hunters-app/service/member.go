package service

import (
    // "encoding/json"
    "net/url"
    "fmt"
    "context"
    "github.com/bwmarrin/discordgo"
	"cherry-blossom-hunters-app/appConfig"
    "cherry-blossom-hunters-app/logger"
)

type ComplianceCheckResultResponse struct {
    Message string   
    Diff    []string 
}

type MemberService struct {
	discordToken    string
	memberChannelID string
	ruleChannelID   string
	masterUserID    string
    goodReactionUrlEncodeString string
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
        goodReactionUrlEncodeString: config.Discord.GoodReactionUrlEncodeString,
    }
}

func botTokenStringBuilder(s *MemberService) string {
    return fmt.Sprintf("Bot %s", s.discordToken)
}

func (s *MemberService) CheckMemberComplianceWithContext(ctx context.Context) (*ComplianceCheckResultResponse, error) {
    dg, err := discordgo.New(botTokenStringBuilder(s))
    if err != nil {
        return nil, fmt.Errorf("Discord session creation error: %v", err)
    }

    members, err := s.getChannelMembers(dg, s.memberChannelID)
    if err != nil {
        return nil, fmt.Errorf("Discord session creation error: %v", err)
    }
    ruleReaderMembers, err := s.getRuleReaders(dg, s.ruleChannelID)
    result, equals := s.diffSlices(members, ruleReaderMembers)

    msg := "差分があります"
    if equals {
        msg = "差分はありません"
    }

    res := &ComplianceCheckResultResponse{
        Message: msg,
        Diff:    result,
    }
    return res, nil
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
    return memberNames, nil
}

func (s *MemberService) isGoodReaction(emojiAPIName string, goodEmojiEncodeString string) bool {
    if emojiAPIName == goodEmojiEncodeString {
        return true
    }
    
    decoded, err := url.QueryUnescape(goodEmojiEncodeString)
    if err == nil && emojiAPIName == decoded {
        return true
    }
    
    return false
}

func (s *MemberService) getRuleReaders(dg *discordgo.Session, channelId string) ([]string, error) {
    messages, err := dg.ChannelMessages(channelId, 100, "", "", "")
    if err != nil {
        return nil, err
    }

    ruleReadersMap := make(map[string]bool)
    for _, message := range messages {
        for _, reaction := range message.Reactions {
            if s.isGoodReaction(reaction.Emoji.APIName(), s.goodReactionUrlEncodeString){
                reactionMembers, err := dg.MessageReactions(channelId, message.ID, reaction.Emoji.APIName(), 100, "", "")
                if err != nil {
                    continue
                }

                for _, reactionMember := range reactionMembers {
                    if !reactionMember.Bot {
                        ruleReadersMap[reactionMember.GlobalName] = true
                    }
                }
            }
        }
    }
    var ruleReaders []string
    for username := range ruleReadersMap {
        ruleReaders = append(ruleReaders, username)
    }
    logger.Logging("debug", "%v", ruleReaders)    
    return ruleReaders, nil
}

func (s *MemberService) diffSlices(a, b []string) ([]string, bool) {
    countA := make(map[string]int)
    countB := make(map[string]int)
    var onlyInA []string

    for _, item := range a {
        if item == "" {
            continue
        }
        countA[item]++
    }
    for _, item := range b {
        if item == "" {
            continue
        }
        countB[item]++
    }

    for item, ca := range countA {
        if diff := ca - countB[item]; diff > 0 {
            for i := 0; i < diff; i++ {
                onlyInA = append(onlyInA, item)
            }
        }
    }

    equal := len(onlyInA) == 0
    return onlyInA, equal
}
