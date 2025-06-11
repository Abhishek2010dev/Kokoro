package kokoro

type (
	// H is a shorthand for a JSON object (map[string]any).
	// Useful for quick response construction: ctx.JSON(H{"msg": "ok"})
	H = map[string]any

	// Payload is an alias for H, used for semantically naming API responses.
	// Example: return ctx.JSON(Payload{"user": user})
	Payload = H
)
