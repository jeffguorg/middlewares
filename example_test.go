package middlewares_test

import (
	"github.com/go-chi/chi"
	"github.com/jeffguorg/middlewares"
	"net/http"
)

func ExampleUseMiddleware() {
	router := chi.NewRouter()

	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		name := middlewares.Parameter(r, "name").(string)
		w.Write([]byte("Hello, " + name + "!"))
	})
	router.Use(middlewares.RequireParametersInQuery("name"))

	http.ListenAndServe(":8888", router)
}
