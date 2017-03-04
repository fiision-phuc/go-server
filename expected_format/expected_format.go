package expectedFormat

const (
	// Invalid param.
	InvalidParameter = "Expected \"Invalid '%s' parameter.\" but found \"%s\"."

	// Nil vs. non nil.
	Nil    = "Expected nil but found not nil."
	NotNil = "Expected not nil but found nil."

	// Panic.
	Panic = "Expected panic but never received."

	// Data types.
	BoolButFoundBool     = "Expected '%t' but found '%t'."
	NumberButFoundNumber = "Expected '%d' but found '%d'."
	StringButFoundString = "Expected '%s' but found '%s'."
)
