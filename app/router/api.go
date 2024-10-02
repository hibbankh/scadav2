package router

import (
	"app/controllers"
	"app/mqtt"
	utils "framework/utils/common"
	"framework/utils/session/jwt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
)

type RouterHandler struct {
	mux      *chi.Mux
	internal *chi.Mux
	jwt      *jwt.JWTStruct
}

type ApiVersion struct {
	Version  string
	ApiRoute []Route
}

type Route struct {
	Pattern string
	Handler http.Handler
}

func (h *RouterHandler) Routes(ApiVersion ...ApiVersion) {

	for _, apiVersion := range ApiVersion {
		if len(apiVersion.ApiRoute) > 0 && apiVersion.Version != "" {
			// Register Routes
			h.mux.Route(apiVersion.Version, func(r chi.Router) {
				if os.Getenv("ENV") == "dev" {
					cors := cors.New(cors.Options{
						// AllowedOrigins: []string{"https://foo.com"}, // Use this to allow specific origin hosts
						AllowedOrigins: []string{"*"},
						// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
						AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
						AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
						ExposedHeaders:   []string{"Link"},
						AllowCredentials: true,
						MaxAge:           300, // Maximum value not ignored by any of major browsers
					})

					r.Use(cors.Handler)
					r.Use(h.jwt.AuthMiddleware)
				}
				for _, route := range apiVersion.ApiRoute {
					r.Mount(route.Pattern, route.Handler)
				}
			})
		}
	}
}

var HandlerInstance *RouterHandler

func NewRouter() *RouterHandler {
	//why we use chi mux instead of Gorilla
	//sebab chimux lgi ringan and laju instead of gorilla which is more features technically
	//tpi chi ni handle static files lgi senang berbanding dengan gorilla
	mux := chi.NewRouter()
	internal := chi.NewRouter()

	// auth
	auth := jwt.NewAuth()

	if HandlerInstance == nil {
		return &RouterHandler{
			mux:      mux,
			internal: internal,
			jwt:      auth,
		}
	}
	return HandlerInstance
}

func (h *RouterHandler) Run() {
	// Init internal and exposed port
	OUTPORT := utils.GetEnv("OUT_PORT", "8080")
	INTERNALPORT := utils.GetEnv("IN_PORT", "3000")

	log.Println("Public Port: ", OUTPORT)
	go http.ListenAndServe(":"+OUTPORT, h.mux)

	log.Println("Internal Port: ", INTERNALPORT)
	http.ListenAndServe(":"+INTERNALPORT, h.internal)
}

func Register() []ApiVersion {

	var apiVersion []ApiVersion
	route_v1 := []Route{
		// Write route here
		{Pattern: "/alarm", Handler: controllers.AlarmRoute()},
		{Pattern: "/incinerator", Handler: controllers.IncineratorRoute()},
		{Pattern: "/setting", Handler: controllers.SettingRoute()},
	}
	apiVersion = append(apiVersion, ApiVersion{
		Version:  "/api/v1",
		ApiRoute: route_v1,
	})

	apiVersion = append(apiVersion, ApiVersion{
		Version: "/api/v2",
		ApiRoute: []Route{
			{Pattern: "/incinerator", Handler: controllers.IncineratorRouteV2()},
		},
	})

	return apiVersion
}

/*
	Register MQTT topics
*/
func Topics() {
	mqtt.GetMqttClient().Subscribe("afes/iot/scada/reading_logs", 0, mqtt.MsgHandlerInstReadingLog)
	mqtt.GetMqttClient().Subscribe("afes/iot/scada/alarm_logs", 0, mqtt.MsgHandlerAlarmLog)
}

/*
	Start webserver
*/
func Start() {
	Topics()
	server := NewRouter()
	route := Register()
	server.Routes(route...)
	server.Run()
}
