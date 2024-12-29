package internalApi

//
//import (
//	"fmt"
//	"github.com/disgoorg/disgo/discord"
//	"github.com/disgoorg/snowflake/v2"
//	"net/http"
//	"strconv"
//)
//
//func (api *InternalApi) GetAdminPanel(res http.ResponseWriter, req *http.Request) {
//	client := *api.BotClient
//	//user, err := client.Rest().GetUser(client.ID())
//	//if err != nil {
//	//	panic(err)
//	//}
//	guildList := ""
//	i := 0
//	client.Caches().GuildCache().ForEach(func(g discord.Guild) {
//		guildList += fmt.Sprintf("<html><a href=\"/guild/%d\"><img src=\"%s\">%s — (%d)</a><br></html>", g.ID, *g.IconURL(), g.Name, g.ID)
//		i++
//	})
//	res.Write([]byte(guildList))
//	return
//}
//
//func (api *InternalApi) GetGuild(res http.ResponseWriter, req *http.Request) {
//	client := *api.BotClient
//	guildId, err := strconv.Atoi(req.RequestURI[len("/a/guild/"):])
//	if err != nil {
//		return
//	}
//	guild, found := client.Caches().GuildCache().Get(snowflake.ID(guildId))
//	if !found {
//		http.Redirect(res, req, "/a/", http.StatusNotFound)
//		return
//	}
//	res.Write([]byte(fmt.Sprintf("<html><img src=\"%s\">%s — %d<br>%d members<br><a href=\"/\">back</a></html><a href=\"/api/guild-leave/%d\">leave</a></html>", *guild.IconURL(), guild.Name, guild.ID, guild.MemberCount, guild.ID)))
//	return
//}
//
//func (api *InternalApi) GuildLeave(res http.ResponseWriter, req *http.Request) {
//	client := *api.BotClient
//	fmt.Println(req.RequestURI[len("/a/guild-leave/"):])
//	guildId, err := strconv.Atoi(req.RequestURI[len("/api/guild-leave/"):])
//	if err != nil {
//		return
//	}
//	fmt.Println("Leave", guildId)
//	_ = client.Rest().LeaveGuild(snowflake.ID(guildId))
//	http.Redirect(res, req, "/", http.StatusOK)
//}
