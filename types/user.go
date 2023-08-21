package types

// bson:"_id" json:"id,omitempty": These are struct tags. Struct tags provide metadata about the field. In this case:
// bson:"_id": This tag specifies how the field should be mapped when using MongoDB's BSON format. 
// The field will be stored in the BSON document with the key _id.

// json:"id,omitempty": This tag specifies how the field should be marshaled/unmarshaled to/from JSON. 
// The field will be named id in the JSON representation. The omitempty option indicates that if the 
// field's value is empty (in this case, an empty string), it should be omitted from the JSON output.

type User struct {
	ID 		  string `bson:"_id" json:"id,omitempty"`
	FirstName string `bson:"firstName" json:"firstName"`
	LastName  string `bson:"lastName" json:"lastName"`
}