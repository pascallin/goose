package goose

const tagName = "goose"

const (
	// field level tags
	indexTag      = "index"
	primaryKeyTag = "primary"
	createdAtTag  = "createdAt"
	updatedAtTag  = "updatedAt"
	deletedAtTag  = "deletedAt"
	// row level tags
	refTag       = "ref"
	forignKeyTag = "forignKey"
	populateTag  = "populate"
)
