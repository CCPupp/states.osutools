// websockets.go
package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"secret"

	"states.osutools/api"
	"states.osutools/discord"
	"states.osutools/htmlbuilder"
	"states.osutools/player"
	"states.osutools/updater"
	"states.osutools/validations"
)

const OnlyBot = false

func main() {
	player.InitializeUserList()
	if !secret.IS_TESTING {
		go updater.StartUpdate()
	}
	// go stats.StartStats()
	go discord.StartBot()
	// Handler points to available directories
	http.Handle("/web/html", http.StripPrefix("/web/html", http.FileServer(http.Dir("web/html"))))
	http.Handle("/web/scripts/", http.StripPrefix("/web/scripts/", http.FileServer(http.Dir("web/scripts"))))
	http.Handle("/web/css/", http.StripPrefix("/web/css/", http.FileServer(http.Dir("web/css"))))
	http.Handle("/web/media/", http.StripPrefix("/web/media/", http.FileServer(http.Dir("web/media"))))
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path[1:] == "" {
			http.ServeFile(w, r, "web/html/index.html")
		} else if r.URL.Path[1:4] == "all" {
			fmt.Fprint(w, htmlbuilder.CreateAllHTML(1))
		} else if r.URL.Path[1:5] == "user" {
			user(w, r)
		} else if r.URL.Path[1:6] == "login" {
			http.Redirect(w, r, "https://osu.ppy.sh/oauth/authorize?response_type=code&client_id="+strconv.Itoa(secret.OSU_CLIENT_ID)+"&redirect_uri="+secret.REDIRECT_URL+"/user&scope=public", http.StatusSeeOther)
		} else if r.URL.Path[1:6] == "stats" {
			fmt.Fprint(w, (htmlbuilder.BuildHTMLHeader(0, "stats")))
			htmlbuilder.BuildHTMLNavbar()
			htmlbuilder.CreateStats(w)

		} else if r.URL.Path[1:7] == "states" {
			if r.URL.Path[8:] != "" {
				if r.URL.Path[len(r.URL.Path)-3:] == "osu" {
					htmlbuilder.CreateStateHTML(w, r.URL.Path[8:len(r.URL.Path)-4], "", "osu", 2)
				} else if r.URL.Path[len(r.URL.Path)-5:] == "mania" {
					htmlbuilder.CreateStateHTML(w, r.URL.Path[8:len(r.URL.Path)-6], "", "mania", 2)
				} else if r.URL.Path[len(r.URL.Path)-5:] == "catch" {
					htmlbuilder.CreateStateHTML(w, r.URL.Path[8:len(r.URL.Path)-6], "", "fruits", 2)
				} else if r.URL.Path[len(r.URL.Path)-5:] == "taiko" {
					htmlbuilder.CreateStateHTML(w, r.URL.Path[8:len(r.URL.Path)-6], "", "taiko", 2)
				} else if r.URL.Path[len(r.URL.Path)-5:] == "North" && r.URL.Path[8:len(r.URL.Path)-6] == "California" {
					htmlbuilder.CreateStateHTML(w, "California", "North", "all", 2)
				} else if r.URL.Path[len(r.URL.Path)-5:] == "South" && r.URL.Path[8:len(r.URL.Path)-6] == "California" {
					htmlbuilder.CreateStateHTML(w, "California", "South", "all", 2)
				} else {
					htmlbuilder.CreateStateHTML(w, r.URL.Path[8:], "", "all", 1)
				}
			}
		} else {
			http.ServeFile(w, r, "web/html/"+r.URL.Path[1:]+".html")
		}
	})

	http.HandleFunc("/update", func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			http.Error(w, fmt.Sprintf("error parsing url %v", err), 500)
		}
		state := r.FormValue("state")
		advstate := r.FormValue("adv")
		bg := r.FormValue("bg")
		token := r.FormValue("id")
		user := api.GetMe(token)
		if user.ID != 0 {
			if bg == "true" || bg == "false" {
				player.NewSetUserBg(bg, strconv.Itoa(user.ID))
			}
			if validations.ValidateState(state) {
				player.NewSetUserState(state, strconv.Itoa(user.ID))
				player.NewSetUserAdvState(advstate, strconv.Itoa(user.ID))
			} else {
				fmt.Fprint(w, "<h2>Invalid State.</h2>")
			}
			idInt, _ := strconv.Atoi(strconv.Itoa(user.ID))
			user := player.NewGetUserById(idInt)
			fmt.Fprint(w, (htmlbuilder.CreateUser(user, 0)))
		} else {
			fmt.Fprint(w, "<h2>Submission Failed.</h2>")
		}

	})

	http.HandleFunc("/submitPlayer", func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			http.Error(w, fmt.Sprintf("error parsing url %v", err), 500)
		}
		state := r.FormValue("state")
		id := r.FormValue("id")
		idInt, _ := strconv.Atoi(id)
		if player.NewGetUserById(idInt).ID != idInt {
			if validations.ValidateState(state) {
				if player.NewCheckStateLock(idInt) {
					fmt.Fprint(w, "<h2>This user is locked!</h2>")
				} else {
					if updater.IsUpdating {
						fmt.Fprint(w, "<h2>Update in progress, please wait.</h2>")
					} else {
						createUserFromId(id, state)
						fmt.Fprint(w, "<h2>Success!</h2>")
					}
				}
			} else {
				fmt.Fprint(w, "<h2>Invalid State.</h2>")
			}
		} else {
			fmt.Fprint(w, "<h2>User Already Exists.</h2>")
		}

	})

	http.HandleFunc("/adminUpdate", func(w http.ResponseWriter, r *http.Request) {
		token := r.FormValue("id")
		user := api.GetMe(token)
		//event := api.GetRecent(user.ID, token)
		localUser := player.NewGetUserById(user.ID)
		//achievement.CheckCompletion(event)
		fmt.Fprint(w, htmlbuilder.CreateAdminPanel(localUser))
	})

	//Serves local webpage for testing
	if OnlyBot == false {
		if secret.IS_TESTING {
			errhttp := http.ListenAndServe(":8080", nil)
			if errhttp != nil {
				log.Fatal("Web server (HTTP): ", errhttp)
			}
		} else {
			//Serves the webpage to the internet
			errhttps := http.ListenAndServeTLS(":443", "certs/cert.pem", "certs/key.pem", nil)
			if errhttps != nil {
				log.Fatal("Web server (HTTPS): ", errhttps)
			}
		}
	}

}

func user(w http.ResponseWriter, r *http.Request) {
	token := api.GetUserToken(r.URL.Query().Get("code"))
	user := api.GetMe(token)

	if user.Username == "" {
		cookie, _ := r.Cookie("Token")
		user = api.GetMe(cookie.Value)
	} else {
		// Lasts same amount of time as an osu token
		expiration := time.Now().Add(365 * 24 * time.Hour)
		cookie := http.Cookie{Name: "Token", Value: token, Expires: expiration}
		http.SetCookie(w, &cookie)
	}

	var localUser player.User
	if player.NewCheckDuplicate(user.ID) {
		localUser = player.NewGetUserById(user.ID)
		player.NewOverwriteExistingUser(localUser, user)
	} else {
		player.NewWriteToUser(user)
		localUser = player.NewGetUserById(user.ID)
	}
	fmt.Fprint(w, htmlbuilder.BuildHTMLHeader(1, "Just "+localUser.Username))
	fmt.Fprint(w, htmlbuilder.BuildHTMLNavbar())
	fmt.Fprint(w, htmlbuilder.CreateUser(localUser, 0))
	fmt.Fprint(w, htmlbuilder.CreateOptions(localUser, token))
	//if localUser.Admin {
	//fmt.Fprint(w, htmlbuilder.CreateAdminPanel(localUser))
	//}
	fmt.Fprint(w, htmlbuilder.BuildHTMLFooter())
}

func createUserFromId(id string, state string) {
	clientToken := api.GetClientToken()
	var newUser player.User = api.GetUser(id, clientToken)
	newUser.State = state
	player.NewWriteToUser(newUser)
}
