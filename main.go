package main

import(
  "net/http"
  	"html/template"
    "log"
    "fmt"
    "io/ioutil"
    "database/sql"
_ "github.com/lib/pq"
    _ "strconv"
    	"github.com/satori/go.uuid"

)

  type user struct {
    email string
    epicusername string
    pass  string
  }
  //creates user database map variable
var dbu = map[string]user{} //user id, stores users
var dbs = map[string]string{} //session id, stores userids



func main() {

  s := &http.Server{

    Addr:    ":80",
    Handler: nil,
  }


  var email string
  var pass string
  var epicusername string
  var ppal bool
  var wins int
  var losses int
  var heat int
  var refers int
  var memberflag string
  var credits int
  var grade int
  var gamertag string
  //pulls users from database
  dbusers, err := sql.Open("postgres", "postgres://postgres:rk@localhost:5432/postgres?sslmode=disable")
  if err != nil {log.Fatalf("Unable to connect to the database")}
  err = dbusers.QueryRow("SELECT * FROM rfgg.members ").Scan(&email, &pass, &ppal, &wins, &losses, &heat, &refers, &memberflag,&credits,&grade,&epicusername,&gamertag)

  if err != nil {log.Fatalf("Could not Scan User Data")}

  dbu[email] = user{email,epicusername,pass}

  http.Handle("/favicon/", http.StripPrefix("/favicon/", http.FileServer(http.Dir("./favicon"))))
  http.Handle("/pics/", http.StripPrefix("/pics/", http.FileServer(http.Dir("./pics"))))
  http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("./css"))))
  http.Handle("/svg/", http.StripPrefix("/svg/", http.FileServer(http.Dir("./svg"))))
  http.HandleFunc("/", serve)
  http.HandleFunc("/verify", verify)
  http.HandleFunc("/weeklyregister", weeklyregister)
  http.HandleFunc("/waitingregister", waitingregister)
  http.HandleFunc("/login", login)
  http.HandleFunc("/profile", profile)
  http.HandleFunc("/tournamentsignup", tournamentsignup)
  log.Fatal(s.ListenAndServe())



}

func alreadyLoggedIn(req *http.Request) bool {
	c, err := req.Cookie("session")
	if err != nil {
		return false
	}
  email := dbs[c.Value]
	_, ok := dbu[email]
	return ok
}

func logout(w http.ResponseWriter, r *http.Request) {
	if !alreadyLoggedIn(r) {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
	}

	c, _ := r.Cookie("session")
	//delete the session
	delete(dbs, c.Value)
	//remove the cookie
	c = &http.Cookie{
		Name:  "session",
		Value: "",
		//max avge value of less than 0 means delete the cookie now
		MaxAge: -1,
	}
	http.SetCookie(w, c)
	http.Redirect(w, r, "/login", http.StatusSeeOther)

}

func getUser(w http.ResponseWriter, r *http.Request) user {
	//gets cookie
	c, err := r.Cookie("session")
	if err != nil {
		sID, _ := uuid.NewV4()
		c = &http.Cookie{
			Name:  "session",
			Value: sID.String(),
		}
	}
	//sets max age of cookie (time available to be logged in) and creates a cookie
	const cookieLength int = 14400
	c.MaxAge = cookieLength
	http.SetCookie(w, c)

	//if user already exists, get user
	var u user
	if email, ok := dbs[c.Value]; ok {
		u = dbu[email]
	}
	return u

}

func tournamentsignup(w http.ResponseWriter, r *http.Request){
  if r.Method == http.MethodPost {
  tournament := "freeweekly1"
  roundnum:=1
  gametype:=r.FormValue("gametype")
  gamertag := r.FormValue("gamertag")
  epicusername := r.FormValue("epicusername")
  email := r.FormValue("email")
  wins := 0
  kills := 0
  matches:=0
  teamname:=r.FormValue("teamname")

  dbusers, err := sql.Open("postgres", "postgres://postgres:rk@localhost:5432/postgres?sslmode=disable")
  if err != nil {log.Fatalf("Unable to connect to the database")}
  sqlStatement := `INSERT INTO rfgg.tournaments (tournament,roundnum,gametype,gamertag,epicusername,wins,kills,matches,teamname) VALUES ($1, $2, $3,$4,$5,$6,$7,$8,$9);`
  _, err = dbusers.Exec(sqlStatement, tournament,roundnum,gametype,gamertag,epicusername,wins,kills,matches,teamname)
  if err != nil {panic(err)}
  sqlStatementer := `UPDATE rfgg.members SET epicusername=%s AND gamertag=%s WHERE email=%s VALUES ($1, $2, $3);`
  _, err = dbusers.Exec(sqlStatementer, epicusername,gamertag,email)
  if err != nil {panic(err)}
  fmt.Printf(gamertag+"Signed up for a tournament")
  dbusers.Close()
  }
  http.Redirect(w, r, "/profile", http.StatusSeeOther)
}


func serve(w http.ResponseWriter, r *http.Request){
  var tpl *template.Template
  tpl = template.Must(template.ParseFiles("main.gohtml","css/main.css","css/mcleod-reset.css"))
  tpl.Execute(w, nil)
}

func verify(w http.ResponseWriter, r *http.Request){
  var tpl *template.Template
  tpl = template.Must(template.ParseFiles("verification.gohtml","css/main.css","css/mcleod-reset.css",))
  tpl.Execute(w, nil)
}

func weeklyregister(w http.ResponseWriter, r *http.Request){
  var tpl *template.Template
  tpl = template.Must(template.ParseFiles("tregistration.gohtml","css/main.css","css/mcleod-reset.css",))
  tpl.Execute(w, nil)
}


func login(w http.ResponseWriter, r *http.Request) {
	//if already logged in send to home page
	if alreadyLoggedIn(r) {
		http.Redirect(w, r, "/profile", http.StatusSeeOther)
		return
	}

	//grab posted form information
	if r.Method == http.MethodPost {
		email := r.FormValue("email")
		pass := r.FormValue("pass")

		//defines u as dbu user info (email,pass) then matches form email with stored email
		u, ok := dbu[email]
		if !ok {
			http.Error(w, "Username and/or password not found", http.StatusForbidden)
			return
		}
		//pulls password from u and checks it with stored password
		if pass != u.pass {
			http.Error(w, "Username and/or password not found", http.StatusForbidden)
			return
		}
		//create new session (cookie) to identify user
		sID, _ := uuid.NewV4()
		c := &http.Cookie{
			Name:  "session",
			Value: sID.String(),
		}
		http.SetCookie(w, c)
		dbs[c.Value] = email
		http.Redirect(w, r, "/login", http.StatusSeeOther)

	}

	//html template
  var tpl *template.Template
  tpl = template.Must(template.ParseFiles("login.gohtml","css/main.css","css/mcleod-reset.css",))
  tpl.Execute(w, nil)

}




func profile(w http.ResponseWriter, r *http.Request){
  //are you already logged in?
	if !alreadyLoggedIn(r) {http.Redirect(w, r, "/signup", http.StatusSeeOther)}
  //provides user a cookie for some time and tracks login
  u := getUser(w, r)
  if u.email == "" {
    http.Error(w, "Please Unblock Cookies - They Help Our Website Run - and Login Again", http.StatusForbidden)
    return
  }

  if r.Method == http.MethodPost {
    var emailcheck string
    var passcheck string
    emailcheck = r.FormValue("email")
    passcheck = r.FormValue("pass")
    dbusers, err := sql.Open("postgres", "postgres://postgres:rk@localhost:5432/postgres?sslmode=disable")
  	if err != nil {
      log.Fatalf("Unable to connect to the database")
    }

    var email string
    var pass string
    var ppal bool
    var wins int
    var losses int
    var heat int
    var refers int
    var memberflag string
    var credits int
    var grade int
    var epicusername string
    var gamertag string
    var tournament string
    var roundnum int
    var gametype string
    var kills int

    type Data struct{
      Email string
      Pass string
      Ppal bool
      Wins int
      Losses int
      Heat int
      Refers int
      Memberflag string
      Credits int
      Grade int
      Epicusername string
      Gamertag string
      Tournament string
      Roundnum int
      Gametype string
      Kills int

    }

    err = dbusers.QueryRow("SELECT * FROM rfgg.members WHERE email=$1 AND pass=$2 AND memberflag=$3",emailcheck,passcheck,"y").Scan(&email, &pass, &ppal, &wins, &losses, &heat, &refers, &memberflag,&credits,&grade,&epicusername,&gamertag)
    _ = dbusers.QueryRow("SELECT * FROM rfgg.tournaments WHERE epicusername=$1",epicusername).Scan(&tournament,&roundnum,&gametype,&gamertag,&epicusername,&kills)

    data:=Data{email, pass, ppal, wins, losses, heat, refers, memberflag, credits, grade, epicusername, gamertag, tournament, roundnum, gametype, kills}
    switch{
    case err == sql.ErrNoRows:
      log.Printf("No user with that ID.")
      http.Redirect(w, r, "/login", http.StatusSeeOther)
    case err != nil:
      log.Fatal(err)
    default:
      fmt.Println(email + " logged on")
      var tpl *template.Template
      tpl = template.Must(template.ParseFiles("profile.gohtml","css/main.css","css/mcleod-reset.css"))

      tpl.Execute(w,data)

      }
  }

    http.Redirect(w, r, "/login", http.StatusSeeOther)
}


func waitingregister(w http.ResponseWriter, r *http.Request){
  //create new session (cookie) to identify user
  sID, _ := uuid.NewV4()
  c := &http.Cookie{
    Name:  "session",
    Value: sID.String(),
  }

  if r.Method == http.MethodPost {
    email := r.FormValue("email")
    pass := r.FormValue("pass")
    epicusername := r.FormValue("epicusername")
    gamertag:= r.FormValue("gamertag")

    http.SetCookie(w, c)
    dbs[c.Value] = email

  	dbusers, err := sql.Open("postgres", "postgres://postgres:rk@localhost:5432/postgres?sslmode=disable")
    fmt.Println(email + " signed up with pass:" + pass)
  	if err != nil {log.Fatalf("Unable to connect to the database")}
    sqlStatement := `INSERT INTO rfgg.members (email, pass, ppal, wins, losses, heat, refers, memberflag, credits, grade, epicusername, gamertag ) VALUES ($1, $2, true, 0, 0, 0, 0, 'y', 0, 0, $3, $4);`
    _, err = dbusers.Exec(sqlStatement, email,pass,epicusername,gamertag)
    if err != nil {http.Redirect(w, r, "/verify", http.StatusSeeOther)}
    dbusers.Close()


    err = ioutil.WriteFile("test.txt", []byte(email+":"+pass), 0666)
    if err != nil {
        log.Fatal(err)
    }
    http.Redirect(w, r, "/waitingregister", http.StatusSeeOther)
    }

  var tpl *template.Template
  tpl = template.Must(template.ParseFiles("waitingverification.gohtml","css/main.css","css/mcleod-reset.css",))
  tpl.Execute(w, nil)

}
