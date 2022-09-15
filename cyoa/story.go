// cyoa is a package for building Choose Your Own Adventure
// stories that can be rendered via the resulting http.Handler
package cyoa

import (
	"encoding/json"
	"errors"
	"html/template"
	"io"
	"log"
	"net/http"
	"strings"
)

// HandlerOptions are used with the NewHandler function to
// configure the returned http.Handler.
type HandlerOption func(h *handler)

// Story represents a Choose Your Own Adventure story.
// Each key is the name of a story chapter and
// each value is a Chapter.
type Story map[string]Chapter

// Chapter represents a CYOA story chapter. Each chapter
// includes its title, the paragraphs it is composed of,
// and options available for the reader to take at the
// end of the chapter. If the options are empty it is
// assumed that you have reached the end of that particular
// story path.
type Chapter struct {
	Title      string   `json:"title"`
	Paragraphs []string `json:"story"`
	Options    []Option `json:"options"`
}

type handler struct {
	story    Story
	template *template.Template
	pathFn   func(r *http.Request) string
}

// Option represents a choice offered at the end of a story
// chapter. Text is the visible text to end users, while
// the Chapter field is the key to a chapter.
type Option struct {
	Text    string `json:"text"`
	Chapter string `json:"arc"`
}

var (
	errRuntime = errors.New("invalid memory address or nil pointer dereference")

	tmpl     *template.Template
	tmplPath = "template/cyoaweb/template.html"
)

func init() {
	tmpl = template.Must(template.ParseFiles(tmplPath))
}

func (h handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := h.pathFn(r)
	if chapter, ok := h.story[path]; ok {
		if h.template == nil {
			log.Fatal(errRuntime)
		}

		err := h.template.Execute(w, chapter)
		if err != nil {
			log.Printf("%v", err)
			http.Error(w, "Something went wrong...", http.StatusInternalServerError)
		}
		return
	}
	http.Error(w, "Chapter not found.", http.StatusNotFound)

}

func defaultPathFn(r *http.Request) string {
	path := strings.TrimSpace(r.URL.Path)
	if path == "" || path == "/" {
		path = "/intro"
	}
	return path[1:]
}

// JsonToStory will decode a story using the incoming reader
// and the encoding/json package. It is assumed that the
// provided reader has the story stored in JSON.
func JsonToStory(r io.Reader) (Story, error) {
	d := json.NewDecoder(r)
	var story Story
	if err := d.Decode(&story); err != nil {
		return nil, err
	}
	return story, nil
}

// NewHandler will construct an http.Handler that will render
// the provided story with some configurations (optional).
func NewHandler(s Story, opts ...HandlerOption) http.Handler {
	h := handler{s, tmpl, defaultPathFn}
	for _, opt := range opts {
		if opt == nil {
			log.Fatal(errRuntime)
		}
		opt(&h)
	}
	return h
}

// SetPathParser is an option to provide a custom function
// for processing the story chapter from the incoming request.
func SetPathParser(fn func(r *http.Request) string) HandlerOption {
	return func(h *handler) {
		h.pathFn = fn
	}
}

// SetTemplate is an option to provide a custom template to
// be used when rendering stories.
func SetTemplate(t *template.Template) HandlerOption {
	return func(h *handler) {
		h.template = t
	}
}
