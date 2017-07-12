package main

import (
	"flag"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/unixpickle/anyvec"
	"github.com/unixpickle/essentials"
	"github.com/unixpickle/serializer"
	"github.com/unixpickle/wordembed/glove"
)

const NumWordResults = 50

type Server struct {
	Addr     string
	AssetDir string

	Embedding *glove.Embedding
}

func main() {
	var server Server
	var embedFile string
	flag.StringVar(&server.Addr, "addr", ":8083", "address to listen on")
	flag.StringVar(&server.AssetDir, "assets", "assets", "web asset directory")
	flag.StringVar(&embedFile, "embedding", "../embedding_out", "embedding file")
	flag.Parse()

	log.Println("Loading embedding...")
	if err := serializer.LoadAny(embedFile, &server.Embedding); err != nil {
		essentials.Die(err)
	}

	http.Handle("/", http.FileServer(http.Dir(server.AssetDir)))
	http.HandleFunc("/word", server.ServeWord)

	log.Println("Listening at " + server.Addr + "...")
	http.ListenAndServe(server.Addr, nil)
}

func (s *Server) ServeWord(w http.ResponseWriter, r *http.Request) {
	query := strings.ToLower(r.FormValue("word"))
	if !s.Embedding.Tokens.Contains(query) {
		s.serveTemplate(w, "notfound", map[string]string{"word": query})
		return
	}
	vec := s.Embedding.Embed(query)
	s.serveTemplate(w, "word", s.wordInfo(query, vec))
}

func (s *Server) serveTemplate(w http.ResponseWriter, name string, obj interface{}) {
	path := filepath.Join(s.AssetDir, "templates", name+".html")
	t, err := template.ParseFiles(path)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := t.Execute(w, obj); err != nil {
		log.Println(err)
	}
}

func (s *Server) wordInfo(word string, vec anyvec.Vector) interface{} {
	ids, correlations := s.Embedding.Lookup(vec, NumWordResults)

	var matches []map[string]interface{}
	for i, id := range ids {
		matchVec := s.Embedding.EmbedID(id)
		matchWord := s.Embedding.Tokens.Token(id)
		if matchWord == word || matchWord == "" {
			continue
		}
		matches = append(matches, map[string]interface{}{
			"id":          id,
			"vec":         matchVec.Data(),
			"word":        matchWord,
			"correlation": correlations[i],
		})
	}
	return map[string]interface{}{
		"vec":     vec.Data(),
		"word":    word,
		"matches": matches,
	}
}
