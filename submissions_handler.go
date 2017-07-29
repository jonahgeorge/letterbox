package main

import (
	"encoding/json"
	"github.com/flosch/pongo2"
	"github.com/gorilla/mux"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	"net/http"
)

func (app *Application) SubmissionsCreateHandler(w http.ResponseWriter, r *http.Request) {
	uuid := mux.Vars(r)["uuid"]
	r.ParseForm()

	form, err := NewFormsRepository(app.db).FindByUuid(uuid)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	json, err := json.MarshalIndent(r.Form, "", " ")
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	_, err = NewSubmissionsRepository(app.db).Create(form.id, string(json))
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	// TODO Send email
	user, _ := NewUsersRepository(app.db).FindById(form.userId)

	plainTextTemplate := `
Plaintext New Form Submission

{{ json }}
	`

	htmlTemplate := `
<h1>New Form Submission</h1> 

<br>

<pre>{{ json }}</pre>
	`

	c := pongo2.Context{
		"json": string(json),
	}

	htmlContent, _ := pongo2.FromString(htmlTemplate)
	plainTextContent, _ := pongo2.FromString(plainTextTemplate)

	plainText, _ := plainTextContent.Execute(c)
	html, _ := htmlContent.Execute(c)

	message := mail.NewSingleEmail(
		mail.NewEmail("Letterdrop Team", "team@letterdrop.herokuapp.com"),
		"New Form Submission",
		mail.NewEmail(user.name, user.email),
		plainText,
		html,
	)

	app.emailClient.Send(message)
	http.Redirect(w, r, "/", 302)
}
