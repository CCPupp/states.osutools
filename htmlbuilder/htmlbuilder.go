package htmlbuilder

import (
	"fmt"
	"math"
	"net/http"
	"strconv"
	"strings"

	"states.osutools/discord"
	"states.osutools/player"
)

func BuildHTMLHeader(loop int, state string) string {
	backString := "../"
	finalBack := strings.Repeat(backString, loop)
	var finalHeader string = `<!DOCTYPE html>
	<html>
	<title>` + state + `</title>
	<meta charset="UTF-8" />
	<link rel="preconnect" href="https://fonts.gstatic.com">
	<link href="https://fonts.googleapis.com/css2?family=Roboto&display=swap" rel="stylesheet"> 
	<link rel="icon" href="` + finalBack + `web/media/favicon.png" type="image/x-icon"/>
	<meta name="viewport" content="width=device-width, initial-scale=1" />
	<link rel="stylesheet" type="text/css" href="` + finalBack + `web/css/main.css" />
	<link rel="stylesheet" type="text/css" href="` + finalBack + `web/css/flexbox.css" />
	<link rel="stylesheet" type="text/css" href="` + finalBack + `web/css/normalize.css" />
	<link rel="stylesheet" type="text/css" href="` + finalBack + `web/css/playercards.css" />
	<meta property="og:type" content="website">
	<meta property="og:title" content="State Leaderboard" />
	<meta property="og:description" content="A leaderboard for osu players split into each state" />
	<meta property="og:url" content="https://states.osutools.com" />
	<meta property="og:image" content="full thumbnail path" />
	<script src="https://ajax.googleapis.com/ajax/libs/jquery/1.12.4/jquery.min.js"></script>
	<script async src="https://pagead2.googlesyndication.com/pagead/js/adsbygoogle.js?client=ca-pub-0073930187680157" crossorigin="anonymous"></script>
	<script src="https://code.jquery.com/jquery-3.1.1.min.js"></script>
	<script src="` + finalBack + `web/scripts/main.js"></script>`
	return finalHeader
}

func BuildHTMLFooter() string {
	var finalFooter string = `
	
	</html>`
	return finalFooter
}

func BuildHTMLNavbar() string {
	finalString := `<body>
    <div class="navbar">
        <a href="/">Home</a>
        <a href="/all">All Users / Discords</a>
        <a href="https://twitter.com/ponparpanpor">Contact</a>
    </div>
	<br>
	<br>
	`
	return finalString
}

func CreateAllHTML(loop int) string {
	finalString := BuildHTMLHeader(loop, "All Players")
	finalString += BuildHTMLNavbar()
	finalString += `
	<br>
	<br>
	<div class='flex-container'>
	<ol>
	<b>Total Users: ` + strconv.Itoa(len(player.UserList)) + `</b><br><br>
	<b>Total Verified Users: ` + player.GetTotalVerified() + `</b><br><br>`

	users := player.UserList

	for i := len(users) - 1; i >= 0; i-- {
		finalString += ("<li><div style='height: 40px;' class='flex-center'><a href='https://osu.ppy.sh/users/" + strconv.Itoa(users[i].ID) + "' class='usercard'>" + users[i].Username + "</a>")
		finalString += ("<a href='/states/" + users[i].State + "'> State: " + users[i].State + getValidation(users[i]) + "</a></div></li>")
	}

	finalString += `</ol><ol>`

	for i := 0; i < len(discord.DiscordList); i++ {
		finalString += `<a href=` + discord.DiscordList[i].Link + `> ` + discord.DiscordList[i].State + `'s Discord Server </a><br><br>`
	}

	finalString += `</div></body>`
	finalString += BuildHTMLFooter()

	return finalString
}

func CreateStateHTML(w http.ResponseWriter, state, advstate, mode string, loop int) {
	discordString := ""
	youtubeString := ""
	for i := 0; i < len(discord.DiscordList); i++ {
		if discord.DiscordList[i].State == state {
			discordString += `<a href="` + discord.DiscordList[i].Link + `"> Discord Server </a>`
		}
	}
	for i := 0; i < len(discord.YoutubeList); i++ {
		if discord.YoutubeList[i].State == state {
			youtubeString += `<a href="` + discord.YoutubeList[i].Link + `"> Youtube Channel </a>`
		}
	}
	fmt.Fprint(w, BuildHTMLHeader(loop, state))

	users := player.SortUsersByRank()

	fmt.Fprint(w, `<body>
    <div class="navbar">
        <a href="/">Home</a>
		<a>`)
	fmt.Fprint(w, state+" / "+mode)
	fmt.Fprint(w, `</a>
	`+discordString+`
	`+youtubeString+`
	<a href="/states/`+state+`/osu">Standard</a>
	<a href="/states/`+state+`/mania">Mania</a>
	<a href="/states/`+state+`/catch">Catch</a>
	<a href="/states/`+state+`/taiko">Taiko</a>
	<a href="/login">Customize My Card</a>
    </div>
	<br>
	<p id="result"></p>
	<div class="playerlist">
	`)

	if discordString == "" {
		fmt.Fprint(w, `<b style="align-self: center;">There is no discord server for this state, try asking a player if it's invite only!</b>`)
	}

	rank := 0
	for i := 0; i < len(users); i++ {
		if mode == "all" {
			if (users[i].State == state) && (users[i].Statistics.Global_rank != 0) {
				if (advstate != "") && users[i].AdvState == advstate {
					rank++
					fmt.Fprint(w, CreateUser(users[i], 0))
				} else if advstate == "" {
					rank++
					fmt.Fprint(w, CreateUser(users[i], 0))
				}
			}
		} else {
			if (users[i].State == state) && (users[i].Statistics.Global_rank != 0) && (users[i].Playmode == mode) {
				rank++
				fmt.Fprint(w, CreateUser(users[i], rank))
			}
		}

	}
	fmt.Fprint(w, "</div>")

	fmt.Fprint(w, `</div></body>`)
	fmt.Fprint(w, BuildHTMLFooter())

}

func CreateStats(w http.ResponseWriter) {

	fmt.Fprint(w, `<body>
	<div class="navbar">
		<a href="/">Home</a>
	</div>
	<br>`)

	fmt.Fprint(w, `
		<h4>Total Users: `+strconv.Itoa(len(player.UserList))+`
	`)

	fmt.Fprint(w, `
		<h4>Total Verified: `+player.GetTotalVerified()+`
	`)
}

func CreateUser(user player.User, rank int) string {
	finalString := (`<div class="players-container" id="response">
						<div class="player">
							<div class="player-preview">
							<h4>#` + getModeRank(rank) + strconv.Itoa((player.GetUserStateRank(user.ID, user.State))) + `</h4>` + `
								<image loading="lazy" class="playerpfp" href="https://osu.ppy.sh/users/` + strconv.Itoa(user.ID) + `" src="http://s.ppy.sh/a/` + strconv.Itoa(user.ID) + `"></image>
								
							</div>
							<div loading="lazy" class="player-info" style="` + getBackground(user) + `">
								<div class="progress-container">
									<span class="progress-text hide-on-mobile">
										<h5>Mode: ` + user.Playmode + `</h5>
										<h5>Level ` + strconv.Itoa(user.Statistics.Level.Current) + `.` + strconv.Itoa(user.Statistics.Level.Progress) + `</h5>
										<h5>Discord: ` + user.Discord + getLink(user) + `</h5>
									</span>
								</div>
								<h6>` + user.State + getValidation(user) + `</h6>
								<a href="https://osu.ppy.sh/users/` + strconv.Itoa(user.ID) + `" target="_blank"><h2>` + user.Username + `</h2></a>
								<h4>Total PP: ` + floatToString(user.Statistics.Pp) + `</h4>
								<h4>Global Rank: ` + strconv.Itoa(user.Statistics.Global_rank) + `</h4>
								<h4>Accuracy: ` + floatToString(user.Statistics.Accuracy) + `</h4>
								<h4>Playcount: ` + strconv.Itoa(user.Statistics.Play_count) + `</h4>
							</div>
						</div>
						` + getBadges(user) + `
					</div>
				</div>`)
	return finalString
}

func CreateOptions(user player.User, token string) string {
	finalString := (`<br>
					<div class="settings-container players-container black-font">
						<div class="user-settings">
							<div class="settings-info">
								<p class="black-font">Hello ` + user.Username + `! Here you can change how your player card appears on the state leaderboard.</p>
								<input type="hidden" id="apitoken" value="` + token + `"/>
								<br>
								<select id="bg">
									<option value="` + user.Background + `" selected hidden>` + getBackgroundText(user) + `</option>
									<option value="true">On</option>
									<option value="false">Off</option>
								</select>
								<label>Background Image On/Off</label>
								<br>
								<br>
								<select id="state">
									<option value="` + user.State + `" selected hidden>` + user.State + `</option>
									<option value="Alabama">Alabama</option>
									<option value="Alaska">Alaska</option>
									<option value="Arizona">Arizona</option>
									<option value="Arkansas">Arkansas</option>
									<option value="California">California</option>
									<option value="Colorado">Colorado</option>
									<option value="Connecticut">Connecticut</option>
									<option value="Delaware">Delaware</option>
									<option value="Florida">Florida</option>
									<option value="Georgia">Georgia</option>
									<option value="Hawaii">Hawaii</option>
									<option value="Idaho">Idaho</option>
									<option value="Illinois">Illinois</option>
									<option value="Indiana">Indiana</option>
									<option value="Iowa">Iowa</option>
									<option value="Kansas">Kansas</option>
									<option value="Kentucky">Kentucky</option>
									<option value="Louisiana">Louisiana</option>
									<option value="Maine">Maine</option>
									<option value="Maryland">Maryland</option>
									<option value="Massachusetts">Massachusetts</option>
									<option value="Michigan">Michigan</option>
									<option value="Minnesota">Minnesota</option>
									<option value="Mississippi">Mississippi</option>
									<option value="Missouri">Missouri</option>
									<option value="Montana">Montana</option>
									<option value="Nebraska">Nebraska</option>
									<option value="Nevada">Nevada</option>
									<option value="New Hampshire">New Hampshire</option>
									<option value="New Jersey">New Jersey</option>
									<option value="New Mexico">New Mexico</option>
									<option value="New York">New York</option>
									<option value="North Carolina">North Carolina</option>
									<option value="North Dakota">North Dakota</option>
									<option value="Ohio">Ohio</option>
									<option value="Oklahoma">Oklahoma</option>
									<option value="Oregon">Oregon</option>
									<option value="Pennsylvania">Pennsylvania</option>
									<option value="Rhode Island">Rhode Island</option>
									<option value="South Carolina">South Carolina</option>
									<option value="South Dakota">South Dakota</option>
									<option value="Tennessee">Tennessee</option>
									<option value="Texas">Texas</option>
									<option value="Utah">Utah</option>
									<option value="Vermont">Vermont</option>
									<option value="Virginia">Virginia</option>
									<option value="Washington">Washington</option>
									<option value="West Virginia">West Virginia</option>
									<option value="Wisconsin">Wisconsin</option>
									<option value="Wyoming">Wyoming</option>
								</select>	
								<label>State Selection</label>
								<br>
								<br>
								<select id="adv">
									<option value="` + user.AdvState + `" selected hidden>` + user.AdvState + `</option>
									<option value="North">North</option>
									<option value="South">South</option>
								</select>
								<label>(California Only) North / South</label>
								<br>
								<br>
								<button id="update">Submit</button>
								<button id="delete">Delete Yourself</button>
							</div>
						</div>`)
	return finalString
}

func CreateAdminPanel(user player.User, token string) string {
	finalString := (`<div class="settings-container admin-container">
						<div class="user-settings black-font">
							<div class="settings-info">
								<p class="black-font">Hello ` + user.Username + `! Here is your admin panel to add other users with.</p>
								<input type="hidden" id="admintoken" value="` + token + `"/>
								<input id="playerid" value=""/>
								<label>User ID</label>
								<br>
								<br>
								<select id="adminstate">
									<option value="Alabama" selected>Alabama</option>
									<option value="Alaska">Alaska</option>
									<option value="Arizona">Arizona</option>
									<option value="Arkansas">Arkansas</option>
									<option value="California">California</option>
									<option value="Colorado">Colorado</option>
									<option value="Connecticut">Connecticut</option>
									<option value="Delaware">Delaware</option>
									<option value="Florida">Florida</option>
									<option value="Georgia">Georgia</option>
									<option value="Hawaii">Hawaii</option>
									<option value="Idaho">Idaho</option>
									<option value="Illinois">Illinois</option>
									<option value="Indiana">Indiana</option>
									<option value="Iowa">Iowa</option>
									<option value="Kansas">Kansas</option>
									<option value="Kentucky">Kentucky</option>
									<option value="Louisiana">Louisiana</option>
									<option value="Maine">Maine</option>
									<option value="Maryland">Maryland</option>
									<option value="Massachusetts">Massachusetts</option>
									<option value="Michigan">Michigan</option>
									<option value="Minnesota">Minnesota</option>
									<option value="Mississippi">Mississippi</option>
									<option value="Missouri">Missouri</option>
									<option value="Montana">Montana</option>
									<option value="Nebraska">Nebraska</option>
									<option value="Nevada">Nevada</option>
									<option value="New Hampshire">New Hampshire</option>
									<option value="New Jersey">New Jersey</option>
									<option value="New Mexico">New Mexico</option>
									<option value="New York">New York</option>
									<option value="North Carolina">North Carolina</option>
									<option value="North Dakota">North Dakota</option>
									<option value="Ohio">Ohio</option>
									<option value="Oklahoma">Oklahoma</option>
									<option value="Oregon">Oregon</option>
									<option value="Pennsylvania">Pennsylvania</option>
									<option value="Rhode Island">Rhode Island</option>
									<option value="South Carolina">South Carolina</option>
									<option value="South Dakota">South Dakota</option>
									<option value="Tennessee">Tennessee</option>
									<option value="Texas">Texas</option>
									<option value="Utah">Utah</option>
									<option value="Vermont">Vermont</option>
									<option value="Virginia">Virginia</option>
									<option value="Washington">Washington</option>
									<option value="West Virginia">West Virginia</option>
									<option value="Wisconsin">Wisconsin</option>
									<option value="Wyoming">Wyoming</option>
								</select>	
								<label>State Selection</label>
								<br>
								<br>
								<select id="adminadv">
									<option value="" selected hidden></option>
									<option value="North">North</option>
									<option value="South">South</option>
								</select>
								<label>(California Only) North / South</label>
								<br>
								<br>
								<button id="adminupdate">Submit</button>
							</div>
						</div>`)
	return finalString
}

func getBackground(user player.User) string {
	if user.Background == "true" || user.Background == "" {
		return `background-image: url('` + user.CoverURL + `'); `
	}
	return ""
}

func getValidation(user player.User) string {
	if user.Locks.State_Lock {
		return ` ✓`
	}
	return ""
}

func getLink(user player.User) string {
	if user.DiscordID != "" {
		return ` ✓`
	}
	return ""
}

func getModeRank(rank int) string {
	if rank != 0 {
		return strconv.Itoa(rank) + " | #"
	} else {
		return ""
	}
}

func getBadges(user player.User) string {
	finalString := ""
	if user.Badges != nil {
		finalString += `<div class="badges">`
		for i := 0; i < len(user.Badges); i++ {
			finalString += `<image class="badge" src="` + user.Badges[i].Image_URL + `">`
			if math.Mod(float64(i+1), 10) == 0 {
				finalString += `</div><div class="badges">`
			}
		}
		finalString += `</div>`
	}
	return finalString
}

func floatToString(input_num float64) string {
	// to convert a float number to a string
	return strconv.FormatFloat(input_num, 'f', 2, 64)
}

func getBackgroundText(bg player.User) string {
	if bg.Background == "true" {
		return "On"
	}
	return "Off"
}
