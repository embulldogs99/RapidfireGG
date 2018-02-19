package main

import(
  "net/http"
  	"html/template"
    "log"
    "fmt"
    "io/ioutil"
    "database/sql"
_ "github.com/lib/pq"
    "strconv"
    "text/template"

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


}


func dbusignup(e string,p string) {

	dbusers, err := sql.Open("postgres", "postgres://postgres:rk@localhost:5432/postgres?sslmode=disable")
	if err != nil {
		log.Fatalf("Unable to connect to the database")
	}
  sqlStatement := `INSERT INTO rfgg.members (email, pass, ppal, wins, losses, heat, refers, memberflag, credits, grade ) VALUES ($1, $2, true, 0, 0, 0, 0, 'y', 0, 0);`
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



	/*bks := make([]Data, 0)
	//cycles through the rows to grab the data by row
	for rows.Next() {
		bk := Data{}
		err := rows.Scan(&bk.Email, &bk.Pass, &bk.Ppal, &bk.Wins,&bk.Heat, &bk.Losses, &bk.Refers, &bk.Memberflag, &bk.Credits, &bk.Grade)

		if err != nil {
			http.Error(w, http.StatusText(500), 500)
			log.Fatal(err)
		}
		// appends the rows
		bks = append(bks, bk)

	}
	db.Close()
	//returns the databse values for use in another function
	return bks
}
*/

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
    var ppal bool
    var wins int
    var losses int
    var heat int
    var refers int
    var memberflag string
    var credits int
    var grade int

    err = dbusers.QueryRow("SELECT * FROM rfgg.members WHERE email=$1 AND pass=$2 AND memberflag=$3",emailcheck,passcheck,"y").Scan(&email, &pass, &ppal, &wins, &losses, &heat, &refers, &memberflag, &credits, &grade)
    data:="Wins:"+strconv.Itoa(wins)+" "+"Heat:"+strconv.Itoa(heat)+" "+"Referalls:"+strconv.Itoa(refers)+" "+"Credits:"+strconv.Itoa(credits)+"Grade:"+strconv.Itoa(grade)

    switch{
    case err == sql.ErrNoRows:
      log.Printf("No user with that ID.")
      http.Redirect(w, r, "/login", http.StatusSeeOther)
    case err != nil:
      log.Fatal(err)
    default:
      var tpl *template.Template
      tpl = template.Parse("profile.gohtml","css/main.css","css/mcleod-reset.css")
      tpl.Execute(w,data)
      fmt.Println(email+"has logged in")

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
