package handlers

const(
	DISCORD_SERVICE = "Discord Service"
	DISCORD_AVATAR_BASE_URL = "https://cdn.discordapp.com/avatars"
	VERIFICATION_STRING = "Please verify your discord account by clicking the link below 👇"
	VERIFICATION_SUBSTRING = "By granting authorization, you agree to permit us to manage your server nickname displayed ONLY in the Real Dev Squad server and to sync your joining data with your user account on our platform."
)

type HEADER_KEYS struct {
	SERVICE string
	AUTHORIZATION string
	CONTENT_TYPE string
}

var HEADERS = HEADER_KEYS{
	SERVICE: "x-service-name",
	AUTHORIZATION: "Authorization",
	CONTENT_TYPE: "Content-Type",
}