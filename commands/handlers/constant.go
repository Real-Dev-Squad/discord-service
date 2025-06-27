package handlers

const(
	DiscordService = "Discord Service"
	DiscordAvatarBaseURL = "https://cdn.discordapp.com/avatars"
	VerificationString = "Please verify your discord account by clicking the link below 👇"
	VerificationNote = "By granting authorization, you agree to permit us to manage your server nickname displayed ONLY in the Real Dev Squad server and to sync your joining data with your user account on our platform."
)

type RequestHeader struct {
	Service string
	Authorization string
	ContentType string
}

var DefaultHeaders = RequestHeader{
	Service: "x-service-name",
	Authorization: "Authorization",
	ContentType: "Content-Type",
}