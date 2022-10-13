package constants

const (
	// AuthSessionKey is the redis key of the auth session.
	AuthSessionKey = "auth:%s:session:%s"
	// ListAuthPermissionsKey is the redis key of the list of auth permissions.
	ListAuthPermissionsKey = "auth:permissions"
)
