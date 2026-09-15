package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os/exec"
	"runtime"
	"strings"
	"time"

	"io.github.com/conv-pfx/internal/assets"
	p "io.github.com/conv-pfx/internal/collections/pix"
	v "io.github.com/conv-pfx/internal/collections/van"
	"io.github.com/conv-pfx/internal/convert"
	"io.github.com/conv-pfx/internal/utils"
)

type Page struct {
	Success   bool   `json:"sucesso"`
	Error     string `json:"erro,omitempty"`
	NomeCert  string `json:"nomeCert,omitempty"`
	CertPEM   string `json:"certPem,omitempty"`
	NomeChave string `json:"nomeChave,omitempty"`
	KeyPEM    string `json:"keyPem,omitempty"`
}

func formatarNomeArquivo(nomeInput, prefixo string) string {

	nomeLimpo := strings.TrimSpace(nomeInput)

	nomeLimpo = strings.TrimSuffix(nomeLimpo, ".pem")
	nomeLimpo = strings.TrimSuffix(nomeLimpo, ".crt")
	nomeLimpo = strings.TrimSuffix(nomeLimpo, ".key")

	if nomeLimpo == "" {
		return prefixo + ".pem"
	}

	return fmt.Sprintf("%s%s.pem", prefixo, nomeLimpo)
}

func abrirNavegador(url string) {
	time.Sleep(300 * time.Millisecond)
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("cmd", "/c", "start", "chrome", url)
	case "darwin":
		cmd = exec.Command("open", "-a", "Google Chrome", url)
	default:
		cmd = exec.Command("google-chrome", url)
	}

	_ = cmd.Start()

}

func responseJSON(w http.ResponseWriter, res Page) {
	w.Header().Set("content-type", "application/json")
	_ = json.NewEncoder(w).Encode(res)
}

var createCollections = map[string]func() ([]byte, error){
	"sicoob:van":    func() ([]byte, error) { return v.SicoobVanCollection(v.SicoobVanConfig{}) },
	"sicoob:pix":    func() ([]byte, error) { return p.SicoobPixCollection(&p.SicoobPixConfig{}) },
	"santander:van": func() ([]byte, error) { return v.SantanderVanCollection(v.SantanderVanConfig{}) },
	"santander:pix": func() ([]byte, error) { return p.SantanderPixCollection(p.SantanderPixConfig{}) },
	"bradesco:van":  func() ([]byte, error) { return v.BradescoVanCollection(v.BradescoVanConfig{}) },
	"bradesco:pix":  func() ([]byte, error) { return p.BradescoPixCollection(p.BradescoPixConfig{}) },
	"sicredi:van":   func() ([]byte, error) { return v.SicrediVanCollection(&v.SicrediVanConfig{}) },
	"sicredi:pix":   func() ([]byte, error) { return p.SicrediPixCollection(&p.SicrediPixConfig{}) },
	"bb:van":        func() ([]byte, error) { return v.BancoDoBrasilVanCollection(v.BancoDoBrasilVanConfig{}) },
	"bb:pix":        func() ([]byte, error) { return p.BancoDoBrasilPixCollection(p.BancoDoBrasilPixConfig{}) },
}

func pathConstruct(r *http.Request) string {
	banco := r.PathValue("banco")
	produto := r.URL.Query().Get("produto")
	if produto == "" {
		return banco
	}

	return banco + ":" + produto
}

func handlerFieldsCollection(w http.ResponseWriter, r *http.Request) {
	gen, ok := createCollections[pathConstruct(r)]

	if !ok {
		http.NotFound(w, r)
		return
	}

	data, err := gen()
	if err != nil {
		responseJSON(w, Page{Success: false, Error: err.Error()})
		return
	}

	fields, err := utils.ExtractField(data)
	if err != nil {
		responseJSON(w, Page{Success: false, Error: err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(fields)
}

func handlerGenCollection(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}

	key := pathConstruct(r)

	gen, ok := createCollections[key]
	if !ok {
		http.NotFound(w, r)
	}

	var values map[string]string

	err := json.NewDecoder(r.Body).Decode(&values)
	if err != nil {
		responseJSON(w, Page{Success: false, Error: err.Error()})
		return
	}

	data, err := gen()
	if err != nil {
		responseJSON(w, Page{Success: false, Error: err.Error()})
		return
	}

	data = utils.CompletePlaceholders(data, values)
	fileName := strings.ReplaceAll(key, ":", "-") + ".json"

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Disposition", `attachment; filename=`+fileName+`"`)
	w.Write(data)
}

func main() {
	http.Handle("/", http.FileServer(assets.GetCoreFile()))

	http.HandleFunc("/convert-pem", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		if err := r.ParseMultipartForm(10 << 20); err != nil {
			responseJSON(w, Page{Success: false, Error: "O arquivo excede o limite permitido 10MB"})
			return
		}

		senha := r.FormValue("password")

		file, _, err := r.FormFile("pfxfile")
		if err != nil {
			responseJSON(w, Page{Success: false, Error: "Falha ao carregar o arquivo pfx"})
			return
		}

		defer file.Close()

		pfxByte, err := io.ReadAll(file)
		if err != nil {
			responseJSON(w, Page{Success: false, Error: "Falha ao ler o conteúdo do pfx"})
			return
		}

		keyPEM, certPEM, err := convert.ConvertPFX(pfxByte, senha)
		if err != nil {
			responseJSON(w, Page{Success: false, Error: err.Error()})
			return
		}

		nomePub := formatarNomeArquivo(r.FormValue("nomeCert"), "ChavePublica")
		nomePriv := formatarNomeArquivo(r.FormValue("nomeChave"), "ChavePrivada")

		responseJSON(w, Page{
			Success:   true,
			NomeCert:  nomePub,
			CertPEM:   string(certPEM),
			NomeChave: nomePriv,
			KeyPEM:    string(keyPEM),
		})

	})

	http.HandleFunc("/convert-next", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		if err := r.ParseMultipartForm(10 << 20); err != nil {
			responseJSON(w, Page{Success: false, Error: "O arquivo excede o limite permitido 10MB"})
			return
		}

		senha := r.FormValue("password")
		novaSenha := r.FormValue("novaSenha")

		if novaSenha == "" {
			novaSenha = senha
		}

		file, _, err := r.FormFile("pfxfile")
		if err != nil {
			responseJSON(w, Page{Success: false, Error: "Falha ao carregar o arquivo pfx"})
			return
		}

		defer file.Close()

		certPFX, err := io.ReadAll(file)
		if err != nil {
			responseJSON(w, Page{Success: false, Error: "Falha ao ler o conteúdo o pfx"})
			return
		}

		newPfx, err := convert.ConvertPfxNext(certPFX, senha, novaSenha)
		if err != nil {
			responseJSON(w, Page{Success: false, Error: err.Error()})
			return
		}

		nomeSaida := formatarNomeArquivo(r.FormValue("nomeSaida"), "certificadonovo")
		nomeSaida = strings.TrimSuffix(nomeSaida, ".pfx")

		w.Header().Set("Content-Type", "application/pkcs-12")
		w.Header().Set("Content-Disposition", `attachment, filename="`+nomeSaida+`"`)
		w.Write(newPfx)

	})

	http.HandleFunc("GET /collections/{bancos}/campos", handlerFieldsCollection)
	http.HandleFunc("POST /collections/{banco}", handlerGenCollection)

	go abrirNavegador("http://localhost:8080")

	fmt.Println("Servidor inicializado em http://localhost:8080")

	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Printf("Erro ao iniciar o servidor")
	}
}
