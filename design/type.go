package design
import . "goa.design/goa/v3/dsl"
var StoredUsers = ResultType("application/vnd.cellar.stored-user", func() {
	Description("A StoredUser describes a user retrieved by the storage service.")
	Reference(User)
	TypeName("StoredUser")
	Attributes(func() {
		Attribute("id", UInt64, "ID is the unique id of the blog.")
		Attribute("name", String, "Name of person", func() {
			MaxLength(100)
		})
	})
	Required("id", "name")
})
var User = Type("User", func() {
	Description("User with id and name")
	Attribute("id", UInt64, "ID of a user")
	Attribute("name", String, "Name of person", func() {
		MaxLength(100)
	})
})