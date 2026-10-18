package web

import (
	"strings"
	"testing"
)

func TestVideoGenPrewarmsPlayerAudioRouting(t *testing.T) {
	videoJS, err := templateFS.ReadFile("templates/static/js/videogen.js")
	if err != nil {
		t.Fatalf("read videogen.js: %v", err)
	}
	appJS, err := templateFS.ReadFile("templates/static/js/app.js")
	if err != nil {
		t.Fatalf("read app.js: %v", err)
	}

	// Embedded sources may be checked out with CRLF line endings on Windows.
	videoSrc := normalizeLineEndings(string(videoJS))
	appSrc := normalizeLineEndings(string(appJS))

	for _, want := range []string{
		"prepareAudioRouting: function(options = {})",
		"if (!this._audioRoutingConnected)",
		"window.VideoGen.prepareAudioRouting();",
	} {
		if !strings.Contains(videoSrc, want) {
			t.Fatalf("videogen.js missing shared audio routing token %q", want)
		}
	}

	for _, want := range []string{
		`document.addEventListener("pointerdown", resumeVideoGenAudioRouting`,
		`document.addEventListener("keydown", resumeVideoGenAudioRouting`,
		`ap.on("play", () => {
    resumeVideoGenAudioRouting();`,
	} {
		if !strings.Contains(appSrc, want) {
			t.Fatalf("app.js missing early audio resume token %q", want)
		}
	}
}

func normalizeLineEndings(s string) string {
	return strings.ReplaceAll(s, "\r\n", "\n")
}
