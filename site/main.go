package main

import (
	"encoding/csv"
	"flag"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
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

	EmbeddingPath string
	Embedding     *glove.Embedding
}

func main() {
	var server Server
	flag.StringVar(&server.Addr, "addr", ":8083", "address to listen on")
	flag.StringVar(&server.AssetDir, "assets", "assets", "web asset directory")
	flag.StringVar(&server.EmbeddingPath, "embedding", "../embedding_out", "embedding file")
	flag.Parse()

	log.Println("Loading embedding...")
	if err := serializer.LoadAny(server.EmbeddingPath, &server.Embedding); err != nil {
		essentials.Die(err)
	}

	http.Handle("/", http.FileServer(http.Dir(server.AssetDir)))
	http.HandleFunc("/word", server.ServeWord)
	http.HandleFunc("/download_csv", server.ServeDownloadCSV)
	http.HandleFunc("/download_raw", server.ServeDownloadRaw)

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

func (s *Server) ServeDownloadCSV(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/csv")
	w.Header().Set("Content-Disposition", "attachment; filename=embeddings.csv;")

	csvWriter := csv.NewWriter(w)
	defer csvWriter.Flush()

	for id, token := range s.Embedding.Tokens {
		vec := s.Embedding.EmbedID(id)
		record := make([]string, vec.Len()+1)
		record[0] = token
		for i, comp := range vec.Data().([]float32) {
			record[i+1] = strconv.FormatFloat(float64(comp), 'f', -1, 32)
		}
		if err := csvWriter.Write(record); err != nil {
			return
		}
	}
}

func (s *Server) ServeDownloadRaw(w http.ResponseWriter, r *http.Request) {
	f, err := os.Open(s.EmbeddingPath)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer f.Close()

	stats, err := os.Stat(s.EmbeddingPath)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/octet-stream")
	http.ServeContent(w, r, "embedding_out", stats.ModTime(), f)
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
