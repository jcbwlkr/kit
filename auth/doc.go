// Package auth package provides API's for managing users who will be accessing
// our web API's and applications. This includes all the CRUD related support
// for users and authentication. For now the user system is simple since users
// are pre-registered using our cli tool `kit`. In the future this can be
// expanded to provide more UI based support.
//
// Users
//
// A user is an entity that can be authenticated on the system and granted rights
// to the API's and applications that we build. From an authentication standpoint
// several fields from a User document are important:
//
// PublicID  : This is the users public identifier and can be shared with the world.
// It provides a unique id for each user. It is used to lookup users from the database.
// This is a randomlu generated UUID.
//
// PrivateID : This is the users private identifier and must not be shared with
// the world. It is used in conjunction with the user supplied password to create
// an encrypted password. To authenticate with a password you need the users
// password and this private id. This is a randomlu generated UUID.
//
// Password  : This is a hased value based on a user provided string and the
// user's private identifier. These values are combined and encrypted to create
// a hash value that is stored in the user document as the password.
//
// A user can be looked up by their public identifier or by email address. The
// email address can be used as a username and an index in the collection for
// users but make this field unique.
//
// Sessions
//
// A session is a document in the database tied to a user via their PublicID.
// Sessions provide a level of security for web tokens by giving them an expiration
// date and a lookup point for the user accessing the API. The SessionID is what
// is used to look up the user performing the web call. Though the user's PublicID
// is not private, it is not used directly.
//
// Web Tokens
//
// Access to the different web service API's using kit require sending a web token
// on every request. HTTP Basic Authorization is being used:
//
// 		Authorization: Basic WebToken
//
// A web token is a value that is not stored in the database but can be consistently
// generated by having a user document and a session identifier. It is made up of
// two parts, a SessionID and a Token which are concatinated together and then base64
// encoded for use over HTTP:
//
// 	base64Encode(SessionID:Token)
//
// The SessionID is a randomly generated UUID and is recorded in the sessions
// collection. It is tied to a user via their PublicID and has an expiration date:
//
// 		SessionID   string
// 		PublicID    string
// 		DateExpires time.Time
// 		DateCreated time.Time
//
// The Token is generated by using the following fields from the user document:
//
// 		Password     string
// 		PublicID     string
// 		PrivateID    string
// 		Email        string
//
// The PublicID, PrivateID, Email fields are used to create a Salt that is
// combined with the Password to create a signed SHA256 hash value. This is a
// value that can be consistenly re-created when all the same values are present.
// If any of the fields used in this token change, the token will be invalidated.
//
// Web Token Authentication
//
// To make things as secure as possible, database lookups are performed as part
// of web token authentication. The user must keep their token secure.
//
// Here are the steps to web token authentication:
//
// 	* Decode the web token and break it into its parts of SessionID and Token.
// 	* Retrieve the session document for the SessionID and validate it has not expired.
// 	* Retrieve the user document from the PublicID in the session document.
// 	* Validate the Token is valid by generating a new token from the retrieved user document.
//
// If any of these steps fail, authorization fails.
//
// User Login
//
// This has not been coded yet.
package auth
