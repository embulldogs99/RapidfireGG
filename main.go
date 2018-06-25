package main

import(
  "net/http"
  	"html/template"
    "log"
    "fmt"
    "database/sql"
_ "github.com/lib/pq"
    _ "strconv"
    	"github.com/satori/go.uuid"
      	"os/exec"

)

type user struct {
  Email string
  Pass string
  Epicusername string
}
  //creates user database map variable
var dbu = map[string]user{} //user id, stores users
var dbs = map[string]string{} //session id, stores userids



func main() {
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


  s := &http.Server{

    Addr:    ":80",
    Handler: nil,
  }

  //pulls users from database


  //pulls users from database
	// for _, k := range dbuconnect() {
	// 	dbu[strings.Trim(k.email, " ")] = user{strings.Trim(k.email, " "), strings.Trim(k.pass, " "), k.o21, strings.Trim(k.role, " ")}
	// }

  dbusers, err := sql.Open("postgres", "postgres://postgres:rk@localhost:5432/postgres?sslmode=disable")
  if err != nil {log.Fatalf("Unable to connect to the database")}
  rowz, err := dbusers.Query("SELECT * FROM rfgg.members")
  if err != nil {log.Fatalf("Could not Scan User Data-select*")}
  for rowz.Next(){
    err := rowz.Scan(&email, &pass, &ppal, &wins, &losses, &heat, &refers, &memberflag,&credits,&grade,&epicusername,&gamertag)
    if err != nil {log.Fatalf("Could not Scan User Data-last step")}
    dbu[email] = user{email,pass,epicusername}
  }
  dbusers.Close()





  http.Handle("/favicon/", http.StripPrefix("/favicon/", http.FileServer(http.Dir("./favicon"))))
  http.Handle("/pics/", http.StripPrefix("/pics/", http.FileServer(http.Dir("./pics"))))
  http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("./css"))))
  http.Handle("/svg/", http.StripPrefix("/svg/", http.FileServer(http.Dir("./svg"))))
  http.HandleFunc("/", serve)
  http.HandleFunc("/signup", signup)
  http.HandleFunc("/signupform", signupform)
  http.HandleFunc("/waitingverification", waitingverification)
  http.HandleFunc("/login", login)
  http.HandleFunc("/logout", logout)
  http.HandleFunc("/profile", profile)
  http.HandleFunc("/tsignup", tsignup)
  http.HandleFunc("/tlaunch", tlaunch)
  http.HandleFunc("/tournaments", tournaments)
  http.HandleFunc("/freeweekly2", freeweekly)
  http.HandleFunc("/freeweekly3", freeweekly)
  http.HandleFunc("/freeweekly4", freeweekly)
  http.HandleFunc("/freeweekly5", freeweekly)
  http.HandleFunc("/freeweekly6", freeweekly)
  http.HandleFunc("/freeweekly7", freeweekly)
  http.HandleFunc("/freeweekly8", freeweekly)
  log.Fatal(s.ListenAndServe())

}


func serve(w http.ResponseWriter, r *http.Request){
  var tpl *template.Template
  tpl = template.Must(template.ParseFiles("main.gohtml","css/main.css","css/mcleod-reset.css"))
  tpl.Execute(w, nil)
}

func dbuconnect() []user {
	//opens conncetion to db for use

	dbusers, err := sql.Open("postgres", "postgres://postgres:rk@localhost:5432/postgres?sslmode=disable")
	if err != nil {log.Fatalf("Unable to connect to the database")}
	rows, err := dbusers.Query("SELECT * FROM rfgg.members")
	if err != nil {log.Fatal(err)}
	userslols := []user{}
	//cycles through the rows to grab the data by row
	for rows.Next() {
		userslol := user{}
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
		err := rows.Scan(&email, &pass, &ppal, &wins, &losses, &heat, &refers, &memberflag,&credits,&grade,&epicusername,&gamertag)
		if err != nil {
			log.Fatal(err)
		}
		// appends the rows
		userslols = append(userslols, userslol)
	}
	dbusers.Close()
	return userslols
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


func signup(w http.ResponseWriter, r *http.Request){
  if alreadyLoggedIn(r) {http.Redirect(w, r, "/profile", http.StatusSeeOther)}
  var tpl *template.Template
  tpl = template.Must(template.ParseFiles("signup.gohtml","css/main.css","css/mcleod-reset.css",))
  tpl.Execute(w, nil)
}

func signupform(w http.ResponseWriter, r *http.Request){
  if r.Method == http.MethodPost {
    email := r.FormValue("email")
    pass := r.FormValue("pass")
    epicusername := r.FormValue("epicusername")
    gamertag:= r.FormValue("gamertag")

    dbusers, err := sql.Open("postgres", "postgres://postgres:rk@localhost:5432/postgres?sslmode=disable")
    fmt.Println(email + " signed up with pass:" + pass)
    if err != nil {log.Fatalf("Unable to connect to the database")}
    sqlStatement := `INSERT INTO rfgg.members (email, pass, ppal, wins, losses, heat, refers, memberflag, credits, grade, epicusername, gamertag ) VALUES ($1, $2, true, 0, 0, 0, 0, 'y', 0, 0, $3, $4);`
    _, err = dbusers.Exec(sqlStatement, email,pass,epicusername,gamertag)
    if err != nil {print(err)}
    dbusers.Close()

    http.Redirect(w, r, "/waitingverification", http.StatusSeeOther)
    }
  var tpl *template.Template
  tpl = template.Must(template.ParseFiles("signupform.gohtml","css/main.css","css/mcleod-reset.css",))
  tpl.Execute(w, nil)

}



func login(w http.ResponseWriter, r *http.Request) {
	//if already logged in send to home page
	if alreadyLoggedIn(r) {
		http.Redirect(w, r, "/profile", http.StatusSeeOther)}

	//grab posted form information
	if r.Method == http.MethodPost {
		email := r.FormValue("email")
		pass := r.FormValue("pass")
		//defines u as dbu user

		u:= dbu[email]
		//pulls password from u and checks it with stored password
		if pass != u.Pass {
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
    http.Redirect(w, r, "/profile", http.StatusSeeOther)

	}else{	//html template
    var tpl *template.Template
    tpl = template.Must(template.ParseFiles("login.gohtml","css/main.css","css/mcleod-reset.css",))
    tpl.Execute(w, nil)}

}


func logout(w http.ResponseWriter, r *http.Request) {
	if !alreadyLoggedIn(r) {http.Redirect(w, r, "/login", http.StatusSeeOther)}
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


func tsignup(w http.ResponseWriter, r *http.Request){
  if !alreadyLoggedIn(r) {http.Redirect(w, r, "/signup", http.StatusSeeOther)}
  //provides user a cookie for some time and tracks login
  u := getUser(w, r)

  if r.Method == http.MethodPost {
  tournament := r.FormValue("tournament")
  starttime:= r.FormValue("starttime")
  roundnum:=1
  gametype:=r.FormValue("gametype")
  gamertag := r.FormValue("gamertag")
  epicusername := u.Epicusername
  email := u.Email
  wins := 0
  kills := 0
  matches:=0
  teamname:=r.FormValue("teamname")


  dbusers, err := sql.Open("postgres", "postgres://postgres:rk@localhost:5432/postgres?sslmode=disable")
  if err != nil {log.Fatalf("Unable to connect to the database")}
  sqlStatement := `INSERT INTO rfgg.tournaments (tournament,roundnum,gametype,gamertag,epicusername,wins,kills,matches,teamname,status,starttime) VALUES ($1, $2, $3,$4,$5,$6,$7,$8,$9,'open',$10);`
  _, err = dbusers.Exec(sqlStatement, tournament,roundnum,gametype,gamertag,epicusername,wins,kills,matches,teamname,starttime)
  if err != nil {
    fmt.Println("failed at user tournament entry step")
  }
  if epicusername != ""{
  sqlStatementer := `UPDATE rfgg.members SET epicusername=%s AND gamertag=%s WHERE email=%s VALUES ($1, $2, $3);`
  _, err = dbusers.Exec(sqlStatementer,epicusername,gamertag,email)
  if err != nil {fmt.Println("failed to update the members database or no update was needed")}
  fmt.Println(gamertag+"Signed up for a tournament")
  dbusers.Close()
  }else{
      http.Redirect(w, r, "/profile", http.StatusSeeOther)
  }
  http.Redirect(w, r, "/profile", http.StatusSeeOther)

  }
  var tpl *template.Template
  tpl = template.Must(template.ParseFiles("tsignup.gohtml","css/main.css","css/mcleod-reset.css"))
  tpl.Execute(w, u)
}



type Data struct{
  Email string
  Pass string
  Ppal bool
  Cwins int
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
  Matches int
  Teamname string
  Status string
  Kills int
  Starttime string

}





type Fortnitedata struct{
  Email string
  Pass string
  Ppal bool
  Cwins int
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
  Matches int
  Teamname string
  Status string
  Kills int
  Starttime string
  Last_updated sql.NullString
  Console sql.NullString
  Squadkill sql.NullFloat64
  Squadmatch sql.NullFloat64
  Squadkm sql.NullFloat64
  Duokill sql.NullFloat64
  Duomatch sql.NullFloat64
  Duokm sql.NullFloat64
  Solokill sql.NullFloat64
  Solomatch sql.NullFloat64
  Solokm sql.NullFloat64

}


func tlaunch(w http.ResponseWriter, r *http.Request){
	var tpl *template.Template
	tpl = template.Must(template.ParseFiles("tlaunch.gohtml","css/main.css","css/mcleod-reset.css"))
	tpl.Execute(w, nil)
	if r.Method == http.MethodPost {
		teamname:=r.FormValue("teamname")
		tournamentname:=r.FormValue("tournamentname")
		http.Redirect(w, r, "/freeweekly2", http.StatusSeeOther)
		tlaunchpython(teamname)
		}
}


func tlaunchpython( teamname string) {
	fmt.Println(teamname+" Launched a tournament")
	byteArray:=[]byte(teamname)
	
	cmd := exec.Command("python", "tlaunch.py")
	
	stdin, _ := cmd.StdinPipe()
		
	go func(){
	defer stdin.Close()
	if _,err :=stdin.Write(byteArray); err != nil{
		panic(err)
	}
}()

	out, _ := cmd.CombinedOutput()
	fmt.Println(string(out))
}








func profilepull(w http.ResponseWriter, r *http.Request) Fortnitedata{
  var email string
  var pass string
  var ppal bool
  var cwins int
  var wins int
  var losses int
  var heat int
  var refers int
  var memberflag string
  var credits int
  var grade int
  var epicusername string
  var gamertag string
  var gamertagt string
  var tournament string
  var roundnum int
  var gametype string
  var matches int
  var teamname string
  var status string
  var kills int
  var starttime string


  u := getUser(w, r)
  dbusers, _ := sql.Open("postgres", "postgres://postgres:rk@localhost:5432/postgres?sslmode=disable")
  _ = dbusers.QueryRow("SELECT * FROM rfgg.members WHERE email=$1 AND pass=$2",u.Email,u.Pass).Scan(&email, &pass, &ppal, &cwins, &losses, &heat, &refers, &memberflag,&credits,&grade,&epicusername,&gamertag)
  dbtourneys, _ := sql.Open("postgres", "postgres://postgres:rk@localhost:5432/postgres?sslmode=disable")
  err := dbtourneys.QueryRow("SELECT tournament, roundnum, gametype, epicusername, wins, kills, matches,teamname,status,gamertag,starttime FROM rfgg.tournaments WHERE epicusername=$1 AND status='open'",u.Epicusername).Scan(&tournament,&roundnum,&gametype,&epicusername,&wins,&kills,&matches,&teamname,&status,&gamertagt,&starttime)
  if err != nil{fmt.Println("failed to select from table")}

  var lastupdated sql.NullString
  var console sql.NullString
  var squadkill sql.NullFloat64
  var squadmatch sql.NullFloat64
  var squadkm sql.NullFloat64
  var duokill sql.NullFloat64
  var duomatch sql.NullFloat64
  var duokm sql.NullFloat64
  var solokill sql.NullFloat64
  var solomatch sql.NullFloat64
  var solokm sql.NullFloat64

  dbfortnite, _ :=sql.Open("postgres", "postgres://postgres:rk@localhost:5432/postgres?sslmode=disable")
  err = dbfortnite.QueryRow("SELECT last_updated,console,squadkill,squadmatch,squadkm,duokill,duomatch,duokm,solokill,solomatch,solokm FROM rfgg.fortniteplayerstats WHERE epicusername=$1 and console='xbl' ORDER BY squadmatch DESC LIMIT 1",u.Epicusername).Scan(&lastupdated,&console,&squadkill,&squadmatch,&squadkm,&duokill,&duomatch,&duokm,&solokill,&solomatch,&solokm)

  data:=Fortnitedata{email, pass, ppal, cwins, wins, losses, heat, refers, memberflag, credits, grade, epicusername, gamertagt, tournament, roundnum, gametype, matches,teamname,status, kills,starttime,lastupdated,console,squadkill,squadmatch,squadkm,duokill,duomatch,duokm,solokill,solomatch,solokm}
  fmt.Println(email + " logged on")
  dbusers.Close()
  dbtourneys.Close()
  return data
}


func profile(w http.ResponseWriter, r *http.Request){
  //are you already logged in?
	if !alreadyLoggedIn(r) {http.Redirect(w, r, "/login", http.StatusSeeOther)}
  //provides user a cookie for some time and tracks login
  data:=profilepull(w,r)
  var tpl *template.Template
  tpl = template.Must(template.ParseFiles("profile.gohtml","css/main.css","css/mcleod-reset.css"))
  tpl.Execute(w,data)
}

func tournaments(w http.ResponseWriter, r *http.Request){
  //are you already logged in?
	if !alreadyLoggedIn(r) {http.Redirect(w, r, "/login", http.StatusSeeOther)}
  //provides user a cookie for some time and tracks login
  u := getUser(w, r)
  if u.Email == "" {
    http.Error(w, "Please Unblock Cookies - They Help Our Website Run - and Login Again", http.StatusForbidden)
    return
  }
  if u.Email=="embulldogs99@yahoo.com"{
    var email string
    var pass string
    var ppal bool
    var cwins int
    var wins int
    var losses int
    var heat int
    var refers int
    var memberflag string
    var credits int
    var grade int
    var epicusername string
    var gamertagt string
    var tournament string
    var roundnum int
    var gametype string
    var matches int
    var teamname string
    var status string
    var kills int
    var starttime string

    dbtourneys, _ := sql.Open("postgres", "postgres://postgres:rk@localhost:5432/postgres?sslmode=disable")
    err := dbtourneys.QueryRow("SELECT * FROM rfgg.tournaments WHERE status='open'").Scan(&tournament,&roundnum,&gametype,&epicusername,&wins,&kills,&matches,&teamname,&status,&gamertagt,&starttime)
    if err != nil{fmt.Println("failed to select from table")}

    data:=Data{email, pass, ppal, cwins, wins, losses, heat, refers, memberflag, credits, grade, epicusername, gamertagt, tournament, roundnum, gametype, matches,teamname,status, kills,starttime}

    var tpl *template.Template
    tpl = template.Must(template.ParseFiles("tournaments.gohtml","css/main.css","css/mcleod-reset.css"))
    tpl.Execute(w,data)
    }

    var tpl *template.Template
    tpl = template.Must(template.ParseFiles("tournaments.gohtml","css/main.css","css/mcleod-reset.css"))
    tpl.Execute(w,nil)

}

type Tourn struct {
    Tournament string
    Roundnum string
    Gametype string
    Epicusername string
    Wins string
    Kills string
    Matches string
    Teamname string
    Status string
    Gamertag string
    Starttime string
    Rank int
  }



func freeweeklypull() []Tourn{
  dbtourneys, _ := sql.Open("postgres", "postgres://postgres:rk@localhost:5432/postgres?sslmode=disable")
  rowz, err := dbtourneys.Query("SELECT tournament, roundnum, gametype, epicusername, wins, kills, matches,teamname,status,gamertag,starttime, rank() over (order by kills desc) as rank FROM rfgg.tournaments WHERE tournament like 'freeweekly%' AND status='open'")
  if err != nil{fmt.Println("failed to select from tournaments table")}
  data := []Tourn{}
  for rowz.Next(){
    datas:=Tourn{}
    err=rowz.Scan(&datas.Tournament,&datas.Roundnum,&datas.Gametype,&datas.Epicusername,&datas.Wins,&datas.Kills,&datas.Matches,&datas.Teamname,&datas.Status,&datas.Gamertag,&datas.Starttime,&datas.Rank)
    if err != nil {log.Fatal(err)}
    data=append(data,datas)
    }
  dbtourneys.Close()
  return data
}

func freeweekly(w http.ResponseWriter, r *http.Request){

  var tpl *template.Template
  tpl = template.Must(template.ParseFiles("freeweekly.gohtml","css/main.css","css/mcleod-reset.css"))
  tpl.Execute(w,freeweeklypull())

}


func waitingverification(w http.ResponseWriter, r *http.Request){

  var tpl *template.Template
  tpl = template.Must(template.ParseFiles("waitingverification.gohtml","css/main.css","css/mcleod-reset.css",))
  tpl.Execute(w, nil)

}
