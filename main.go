package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

type DiscordGuildMember struct {
	Avatar                     interface{} `json:"avatar"`
	Banner                     interface{} `json:"banner"`
	CommunicationDisabledUntil interface{} `json:"communication_disabled_until"`
	Flags                      int         `json:"flags"`
	JoinedAt                   time.Time   `json:"joined_at"`
	Nick                       interface{} `json:"nick"`
	Pending                    bool        `json:"pending"`
	PremiumSince               interface{} `json:"premium_since"`
	Roles                      []string    `json:"roles"`
	UnusualDmActivityUntil     interface{} `json:"unusual_dm_activity_until"`
	User                       struct {
		Id                   string      `json:"id"`
		Username             string      `json:"username"`
		Avatar               *string     `json:"avatar"`
		Bot                  *bool       `json:"bot"`
		Discriminator        string      `json:"discriminator"`
		PublicFlags          int         `json:"public_flags"`
		Flags                int         `json:"flags"`
		Banner               interface{} `json:"banner"`
		AccentColor          interface{} `json:"accent_color"`
		GlobalName           *string     `json:"global_name"`
		AvatarDecorationData *struct {
			Asset     string      `json:"asset"`
			SkuId     string      `json:"sku_id"`
			ExpiresAt interface{} `json:"expires_at"`
		} `json:"avatar_decoration_data"`
		BannerColor interface{} `json:"banner_color"`
		Clan        interface{} `json:"clan"`
	} `json:"user"`
	Mute bool `json:"mute"`
	Deaf bool `json:"deaf"`
}

type GhosttyDiscordMember struct {
	Username string
	JoinedAt time.Time
}

func getGuildMembers(limit int, after string) []DiscordGuildMember {
	fmt.Println("Get guild members, after", after)
	req, err := http.NewRequest("GET", "https://discord.com/api/v10/guilds/"+os.Getenv("GUILD_ID")+"/members", nil)
	if err != nil {
		panic(err)
	}
	q := url.Values{}
	q.Add("limit", strconv.Itoa(limit))
	if after != "" {
		q.Add("after", after)
	}
	req.URL.RawQuery = q.Encode()
	req.Header.Add("Authorization", "Bot "+os.Getenv("BOT_TOKEN"))
	httpClient := &http.Client{}
	resp, err := httpClient.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	var guildMembers []DiscordGuildMember
	err = json.Unmarshal(body, &guildMembers)
	if err != nil {
		panic(err)
	}
	return guildMembers
}

func writeFile(filename string, content string) {
	err := os.WriteFile(filename, []byte(content), 0644)
	if err != nil {
		panic(err)
	}
}

func main() {
	limit := 1000
	limitVal := os.Getenv("LIMIT")
	if limitVal != "" {
		limit, _ = strconv.Atoi(limitVal)
	}

	var after string
	var guildMembers []DiscordGuildMember
	for {
		gm := getGuildMembers(limit, after)
		guildMembers = append(guildMembers, gm...)
		if len(gm) < limit {
			break
		}
		after = gm[len(gm)-1].User.Id
	}

	ghosttyWaitingMembers := make([]GhosttyDiscordMember, 0, len(guildMembers))
	testerCount := 0
	for _, member := range guildMembers {
		if len(member.Roles) > 0 {
			testerCount++
			continue
		}
		if member.User.Bot != nil && *member.User.Bot {
			continue
		}
		ghosttyWaitingMembers = append(ghosttyWaitingMembers, GhosttyDiscordMember{
			Username: member.User.Username,
			JoinedAt: member.JoinedAt,
		})
	}

	sort.Slice(ghosttyWaitingMembers, func(i, j int) bool {
		return ghosttyWaitingMembers[i].JoinedAt.Before(ghosttyWaitingMembers[j].JoinedAt)
	})

	var markdown strings.Builder
	markdown.WriteString(fmt.Sprintf("# Ghostty Waiting List, %s \n", time.Now().Format(time.DateTime)))
	markdown.WriteString(fmt.Sprintf("This is the _inofficial_ waiting list.\n"))
	markdown.WriteString(fmt.Sprintf("There are %d testers and %d in the queue.\n",
		testerCount, len(ghosttyWaitingMembers)))
	markdown.WriteString(fmt.Sprintf("If you have any questions, check the official Ghostty discord.\n"))
	markdown.WriteString(fmt.Sprintf("https://discord.com/channels/1005603569187160125/1140732798773248121\n\n"))

	markdown.WriteString("```mermaid\n")
	markdown.WriteString("pie title Beta participants\n")
	markdown.WriteString(fmt.Sprintf("    \"Tester\" : %d\n", testerCount))
	markdown.WriteString(fmt.Sprintf("    \"Aspirant\" : %d\n", len(ghosttyWaitingMembers)))
	markdown.WriteString("```\n\n")

	markdown.WriteString("|#|Username|Joined At|Days|\n")
	markdown.WriteString("|---|---|---|---|\n")
	for i, member := range ghosttyWaitingMembers {
		markdown.WriteString(
			fmt.Sprintf("|%d|%s|%s|%d|\n", i+1, member.Username,
				member.JoinedAt.Format(time.DateTime), int(time.Since(member.JoinedAt).Hours()/24)))
	}

	writeFile("list.md", markdown.String())
	writeFile(fmt.Sprintf("archive/%s.md", time.Now().Format(time.DateOnly)), markdown.String())
}
