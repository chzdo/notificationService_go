package main

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"notification_service/internals/appHandlers"
	"notification_service/internals/appValidators"
	"notification_service/internals/mailing"
	"notification_service/internals/models"
	pushnotification "notification_service/internals/pushNotification"
	"notification_service/internals/responses"
	"notification_service/internals/services"
	"notification_service/internals/socket"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	socketio "github.com/googollee/go-socket.io"
	"github.com/mailgun/mailgun-go/v4"

	"github.com/tbalthazar/onesignal-go"
)

// func corsMiddleware(next http.Handler) (http.Handler, socketio.Server) {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		allowHeaders := "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization"

// 		w.Header().Set("Content-Type", "application/json")
// 		w.Header().Set("Access-Control-Allow-Origin", "*")
// 		w.Header().Set("Access-Control-Allow-Methods", "POST, PUT, PATCH, GET, DELETE")
// 		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
// 		w.Header().Set("Access-Control-Allow-Credentials", "true")
// 		w.Header().Set("Access-Control-Allow-Headers", allowHeaders)

// 		next.ServeHTTP(w, r)
// 	})
// }

func (app *application) Routes() (http.Handler, socketio.Server) {

	model, err := models.RegisterModel(app.logs)
	mailgun := mailgun.NewMailgun(os.Getenv("MAILGUN_DOMAIN"), os.Getenv("MAILGUN_API_KEY"))
	push := onesignal.NewClient(nil)

	push.AppKey = os.Getenv("ONE_SIGNAL_APP_KEY")

	if err != nil {
		app.logs.ErrorLogs.Panicln(err)
	}
	handle := &appHandlers.Handler{
		Handler: chi.NewMux(),
		Services: &services.Services{
			Logs:      *app.logs,
			Models:    model,
			Validator: *appValidators.Validator,
			Socket:    socket.Websocket,
			Mailer: mailing.Mailer{
				Driver: mailgun,
			},
			Push: pushnotification.PushNotification{
				Driver: *push,
			},
		},
		Responses: &responses.ResponseFunctions{
			Logs: app.logs,
		},
	}

	handle.Handler.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"*"},
	}))
	workDir, _ := os.Getwd()
	filesDir := http.Dir(filepath.Join(workDir, "static/"))

	handle.Handler.Use(func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			rw.Header().Set("Content-Type", "application/json")
			h.ServeHTTP(rw, r)
		})

	})
	FileServer(handle.Handler, "/", filesDir)

	handle.Handler.Handle("/socket.io/", &socket.Websocket.Socket)
	handle.TriggerRoutes()
	handle.OrganizationSettingsRoutes()
	handle.UserMobileSettingsRoutes()
	handle.OrganizationNotificationRoutes()
	handle.UserSeenNotificationRoutes()

	handle.Handler.NotFound(func(response http.ResponseWriter, request *http.Request) {

		handle.Responses.ErrorRespond(response, &responses.ErrorResponse{
			Status:  http.StatusNotFound,
			Success: false,
			Message: http.StatusText(http.StatusNotFound),
		})

	})

	return handle.Handler, socket.Websocket.Socket
}

func FileServer(r chi.Router, path string, root http.FileSystem) {

	if strings.ContainsAny(path, "{}*") {
		panic("File Server does not permit any URL Parameters")
	}

	if path != "/" && path[len(path)-1] != '/' {
		r.Get(path, http.RedirectHandler(path+"/", 301).ServeHTTP)
		path += "/"
	}

	path += "*"

	r.Get(path, func(w http.ResponseWriter, r *http.Request) {
		rectx := chi.RouteContext(r.Context())
		pathprefix := strings.TrimSuffix(rectx.RoutePattern(), "/*")
		w.Header().Set("Content-Type", "text2/html")
		fs := http.StripPrefix(pathprefix, http.FileServer(root))
		fs.ServeHTTP(w, r)
	})
}
