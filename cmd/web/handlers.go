package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/julienschmidt/httprouter"
	"twitbox.vedantkugaonkar.net/internal/model"
	"twitbox.vedantkugaonkar.net/internal/validator"
)

type snippetCreateForm struct {
	Title               string `form:"title"`
	Content             string `form:"content"`
	Expires             int    `form:"expires"`
	validator.Validator `form:"-"`
}

func (app *application) home(w http.ResponseWriter, req *http.Request) {
	twits, err := app.twits.Latest()
	if err != nil {
		app.serverError(w, err)
		return
	}
	data := app.newTemplateData(req)
	data.Twits = twits
	app.renderer(w,req, "home.tmpl.html", http.StatusOK, data)
}
func (app *application) twitView(w http.ResponseWriter, req *http.Request) {
	//when httprouter parses a req, any named parameters are stored in the
	//req context so here ParseFromContext is used to retrieve/get the id
	params := httprouter.ParamsFromContext(req.Context())
	//ByName can then be used to get the actual value
	id, err := strconv.Atoi(params.ByName("id"))
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}
	twit, err := app.twits.Get(id)
	if err != nil {
		if errors.Is(err, model.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
		return
	}

	data := app.newTemplateData(req)
	data.Twit = twit

	app.renderer(w,req, "view.tmpl.html", http.StatusOK, data)

}
func (app *application) twitCreate(w http.ResponseWriter, req *http.Request) {
	data := app.newTemplateData(req)
        data.CurrentYear = time.Now().Year()
	data.Form = &snippetCreateForm{
		Expires: 365,
	}
	app.renderer(w,req, "create.tmpl.html", http.StatusCreated, data)

}
func (app *application) twitCreatePost(w http.ResponseWriter, req *http.Request) {
        app.infoLog.Printf("POST /twit/create - Content-Type: %s", req.Header.Get("Content-Type"))
	app.infoLog.Printf("POST /twit/create - Form values before decode: %v", req.PostForm)
	var scF snippetCreateForm
	err := app.decodePostForm(req, &scF)
	if err != nil {
                app.errorLog.Printf("Form decode error: %v", err)
		app.clientError(w, http.StatusBadRequest)
		return
	}
        app.infoLog.Printf("POST /twit/create - Decoded form: Title=%s, Content=%s, Expires=%d", scF.Title, scF.Content, scF.Expires)

	scF.CheckField(validator.NotBlank(scF.Title), "title", "This field cannot be blank!")
	scF.CheckField(validator.MaxChars(scF.Title, 100), "title", "This field cannot be more than 100 characters")
	scF.CheckField(validator.PermittedInt(scF.Expires, 1, 7, 365), "expires", "This field must include 1,7 or 365")
	scF.CheckField(validator.NotBlank(scF.Content), "content", "This field cannot be empty!")
	if !scF.Valid() {
		data := app.newTemplateData(req)
		data.Form = scF
		app.renderer(w,req, "create.tmpl.html", http.StatusUnprocessableEntity, data)
		return
	}
	id, err := app.twits.Insert(scF.Title, scF.Content, scF.Expires)
	if err != nil {
                
		app.serverError(w, err)
		return
	}
	app.sessionManager.Put(req.Context(), "flash", "Snippet successfully created!")
	//redirecting user to the relevant page for the twit
	http.Redirect(w, req, fmt.Sprintf("/twit/view/%d", id), http.StatusSeeOther)
}

type userCreateForm struct {
	Name                string `form:"name"`
	Email               string `form:"email"`
	Password            string `form:"password"`
	validator.Validator `form:"-"`
}

func (app *application) userSignup(w http.ResponseWriter, req *http.Request) {
	data := app.newTemplateData(req)
	data.Form = userCreateForm{}
	app.renderer(w,req, "signup.tmpl.html", http.StatusCreated, data)
}
func (app *application) userSignupPost(w http.ResponseWriter, req *http.Request) {
	var uCF userCreateForm
	err := app.decodePostForm(req, &uCF)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	uCF.CheckField(validator.NotBlank(uCF.Name), "name", "This field cannot be blank!")
	uCF.CheckField(validator.NotBlank(uCF.Email), "email", "This field cannot be empty!")
	uCF.CheckField(validator.Matches(uCF.Email, validator.EmailRX), "email", "This field must be a valid email id")
	uCF.CheckField(validator.NotBlank(uCF.Password), "password", "This field cannot be empty!")
	uCF.CheckField(validator.MinChars(uCF.Password, 8), "password", "This field must be atleast 8 characters long!")
	if !uCF.Valid() {
		data := app.newTemplateData(req)
		data.Form = uCF
		app.renderer(w,req, "signup.tmpl.html", http.StatusUnprocessableEntity, data)
		return
	}
	err = app.users.Insert(uCF.Name, uCF.Email, uCF.Password)
	if err != nil {
		if errors.Is(err, model.ErrDuplicateEmail) {
			uCF.AddFieldErrors("email", "Email address already in use!")
			data := app.newTemplateData(req)
			data.Form = uCF
			app.renderer(w,req, "signup.tmpl.html", http.StatusUnprocessableEntity, data)
		} else {
			app.serverError(w, err)
		}
		return
	}
	app.sessionManager.Put(req.Context(), "flash", "Your signup was successful.Please log in")
	http.Redirect(w, req, "/user/login", http.StatusSeeOther)
}

type UserLoginForm struct {
	Email    string `form:"email"`
	Password string `form:"password"`
	validator.Validator
}

func (app *application) userLogin(w http.ResponseWriter, req *http.Request) {
	data := app.newTemplateData(req)
	data.Form = UserLoginForm{}
	app.renderer(w,req, "login.tmpl.html", http.StatusOK, data)
}
func (app *application) userLoginPost(w http.ResponseWriter, req *http.Request) {
	var uLF UserLoginForm
	err := app.decodePostForm(req, &uLF)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	uLF.CheckField(validator.NotBlank(uLF.Email), "email", "This field cannot be blank!")
	uLF.CheckField(validator.Matches(uLF.Email, validator.EmailRX), "email", "This field must be a valid email address")
	uLF.CheckField(validator.NotBlank(uLF.Password), "password", "This field cannot be blank!")
	if !uLF.Valid() {
		data := app.newTemplateData(req)
		data.Form = uLF
		app.renderer(w,req, "login.tmpl.html", http.StatusUnprocessableEntity, data)
		return
	}
	id, err := app.users.Authenticate(uLF.Email, uLF.Password)
	if err != nil {
		if errors.Is(err, model.ErrInvalidCredentials) {
			uLF.AddNonFieldError("Email or password is incorrect!")
			data := app.newTemplateData(req)
			data.Form = uLF
			app.renderer(w,req, "login.tmpl.html", http.StatusUnprocessableEntity, data)
		} else {
			app.serverError(w, err)
		}
		return
	}
	err = app.sessionManager.RenewToken(req.Context())
	if err != nil {
		app.serverError(w, err)
		return
	}
	app.sessionManager.Put(req.Context(), "authenticatedUserId", id)
	http.Redirect(w, req, "/twit/create", http.StatusSeeOther)

}
func (app *application) userLogOut(w http.ResponseWriter, req *http.Request) {
	err := app.sessionManager.RenewToken(req.Context())
	if err != nil {
		app.serverError(w, err)
		return
	}
	app.sessionManager.Remove(req.Context(), "authenticatedUserId")
	app.sessionManager.Put(req.Context(), "flash", "You've been logged out successfully!")
	http.Redirect(w, req, "/", http.StatusSeeOther)

}
func (app *application) about(w http.ResponseWriter, req *http.Request) {
	data := app.newTemplateData(req)
	app.renderer(w, req, "about.tmpl.html", http.StatusOK, data)
}
func (app *application) accountView(w http.ResponseWriter, req *http.Request) {
	userID := app.sessionManager.GetInt(req.Context(), "authenticatedUserId")

	user, err := app.users.Get(userID)
	if err != nil {
		if errors.Is(err, model.ErrNoRecord) {
			http.Redirect(w, req, "/user/login", http.StatusSeeOther)
		} else {
			app.serverError(w, err)
		}
		return
	}

	data := app.newTemplateData(req)
	data.User = user

	app.renderer(w, req , "account.tmpl.html",http.StatusOK, data)
}

