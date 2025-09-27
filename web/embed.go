package webassets

import "embed"

// Assets bundles HTML templates and static files for the checkout UI.
//
//go:embed templates/*.html static/*
var Assets embed.FS
