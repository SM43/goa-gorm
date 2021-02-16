package design
import (
	. "goa.design/goa/v3/dsl"
	cors "goa.design/plugins/v3/cors/dsl"
)
var _ = API("user", func() {
	Title("User Service")
	Description("Service to perform CRUD operations using goa")
    Server("user", func() {
        Host("localhost", func() {
            URI("http://localhost:8080")
        })
    })
})
var _ = Service("user", func() {
	Description("The user service gives user details.")
	cors.Origin("/.*local.thaha.xyz/", func() {
		cors.Headers("X-Shared-Secret")
		cors.Methods("GET", "POST")
		cors.Expose("X-Time", "X-Api-Version")
		cors.MaxAge(100)
		cors.Credentials()
	})
	Error("db_error", ErrorResult, "Unable to process db request")
	//Method to post new user
	Method("add", func() {
		Description("Add new user and return its ID.")
		Payload(User)
		Result(String)
		HTTP(func() {
			POST("/")
			Response(StatusCreated)
			Response("db_error", StatusInternalServerError)
		})
	})
	//Method to get all existing users
	Method("list", func() {
		Description("List all users")
		Result(ArrayOf(StoredUsers))
		HTTP(func() {
			GET("/")
			Response(StatusOK)
			Response("db_error", StatusInternalServerError)
		})
	})
	Files("/openapi.json", "./gen/http/openapi.json")
})