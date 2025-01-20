package routers

import (
    "gloriusaiapi/controllers"
    "github.com/gorilla/mux"
)

func SetupRouter() *mux.Router {
    router := mux.NewRouter()

    router.HandleFunc("/login", controllers.Login).Methods("POST")
    router.HandleFunc("/messages", controllers.CreateMessage).Methods("POST")
    router.HandleFunc("/getAllModels", controllers.GetAllModels).Methods("GET")
    router.HandleFunc("/setModel", controllers.SetModel).Methods("GET")

    return router
}
