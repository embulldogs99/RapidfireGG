package main

import(
  "net/http"
  	"html/template"
    "log"
    "fmt"
    "io/ioutil"
    "database/sql"
_ "github.com/lib/pq"



)


func main() {

  s := &http.Server{

    Addr:    ":80",
    Handler: nil,
  }

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
  log.Fatal(s.ListenAndServe())




/*type user struct {
	email string
	pass  string
	ppal  bool
	wins  int
  losses int
  heat int
  refers int
  memberflag int
  credits float64
  gread int
}*/

/*var dbs = map[string]string{} //session id, stores userids*/


}


func dbusignup(e string,p string) {

	dbusers, err := sql.Open("postgres", "postgres://postgres:rk@localhost:5432/postgres?sslmode=disable")
	if err != nil {
		log.Fatalf("Unable to connect to the database")
	}
  sqlStatement := `INSERT INTO rfgg.members (email, pass, ppal, wins, losses, heat, refers, memberflag, credits, grade ) VALUES ($1, $2, true, 0, 0, 0, 0, 'Y', 0, 0);`
  _, err = dbusers.Exec(sqlStatement, e,p)
  if err != nil {
    panic(err)
  }
  dbusers.Close()
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

func login(w http.ResponseWriter, r *http.Request){
  var tpl *template.Template
  tpl = template.Must(template.ParseFiles("login.gohtml","css/main.css","css/mcleod-reset.css",))
  tpl.Execute(w, nil)
}



func profile(w http.ResponseWriter, r *http.Request){
  if r.Method == http.MethodPost {
    var emailcheck string
    var passcheck string
    emailcheck = r.FormValue("email")
    passcheck = r.FormValue("pass")
    fmt.Println(emailcheck)
    dbusers, err := sql.Open("postgres", "postgres://postgres:rk@localhost:5432/postgres?sslmode=disable")
    fmt.Println("Connected")
  	if err != nil {
      log.Fatalf("Unable to connect to the database")
    }

    var email string
    var pass string
    var ppal string
    var wins int
    var losses int
    var heat int
    var refers int
    var memberflag string
    var credits float64
    var grade int

    err = dbusers.QueryRow("SELECT * FROM rfgg.members WHERE email=$1 AND pass=$2",emailcheck,passcheck).Scan(&email, &pass, &ppal, &wins, &losses, &heat, &refers, &memberflag, &credits, &grade)

    switch{
    case err == sql.ErrNoRows:
      log.Printf("No user with that ID.")
      http.Redirect(w, r, "/login", http.StatusSeeOther)
    case err != nil:
      log.Fatal(err)
    default:
      var tpl *template.Template
      tpl = template.Must(template.ParseFiles("profile.gohtml","css/main.css","css/mcleod-reset.css",))
      tpl.Execute(w, nil)
      fmt.Println(email)
      fmt.Println(wins)
      fmt.Println(pass)
      fmt.Println(ppal)
      fmt.Println(losses)
      fmt.Println(heat)
      fmt.Println(memberflag)
      fmt.Println(refers)
      fmt.Println(credits)
      fmt.Println(grade)
      }
  }
}


func waitingregister(w http.ResponseWriter, r *http.Request){

  if r.Method == http.MethodPost {
    email := r.FormValue("email")
    pass := r.FormValue("pass")

    fmt.Println(email + " signed up with pass:" + pass)
    dbusignup(email,pass)

    err := ioutil.WriteFile("test.txt", []byte(email+":"+pass), 0666)
    if err != nil {
        log.Fatal(err)
    }
    http.Redirect(w, r, "/waitingregister", http.StatusSeeOther)
    }

  var tpl *template.Template
  tpl = template.Must(template.ParseFiles("waitingverification.gohtml","css/main.css","css/mcleod-reset.css",))
  tpl.Execute(w, nil)

}
