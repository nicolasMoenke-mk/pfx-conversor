package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os/exec"
	"runtime"
	"strings"
	"text/template"
	"time"

	"io.github.com/conv-pfx/convert"
)

type Page struct {
	Success   bool   `json:"sucesso"`
	Error     string `json:"erro,omitempty"`
	NomeCert  string `json:"nomeCert,omitempty"`
	CertPEM   string `json:"certPem,omitempty"`
	NomeChave string `json:"nomeChave,omitempty"`
	KeyPEM    string `json:"keyPem,omitempty"`
}

var page = template.Must(template.New("index").Parse(
	`
	<!DOCTYPE html>
<html lang="pt-BR">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Conversor de Certificados</title>
    <style>
        :root {
            --primary: #2563eb;
            --primary-hover: #1d4ed8;
            --bg: #f8fafc;
            --card-bg: #ffffff;
            --text: #0f172a;
            --text-muted: #64748b;
            --border: #e2e8f0;
            --error-bg: #fef2f2;
            --error-text: #991b1b;
            --success-bg: #f0fdf4;
            --success-text: #166534;
        }

        * { box-sizing: border-box; margin: 0; padding: 0; }
        body { font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, sans-serif; background-color: var(--bg); color: var(--text); padding: 40px 20px; display: flex; justify-content: center; }
        .container { background: var(--card-bg); width: 100%; max-width: 540px; padding: 32px; border-radius: 12px; border: 1px solid var(--border); box-shadow: 0 10px 15px -3px rgba(0,0,0,0.05); }
        h1 { font-size: 1.5rem; font-weight: 700; margin-bottom: 4px; }
        p.subtitle { font-size: 0.875rem; color: var(--text-muted); margin-bottom: 20px; }
        
        /* Menu de Abas */
        .tabs { display: flex; border-bottom: 2px solid var(--border); margin-bottom: 24px; gap: 8px; }
        .tab-btn { flex: 1; padding: 10px 12px; background: none; border: none; font-size: 0.875rem; font-weight: 600; color: var(--text-muted); cursor: pointer; border-bottom: 2px solid transparent; margin-bottom: -2px; transition: all 0.2s; }
        .tab-btn:hover { color: var(--primary); }
        .tab-btn.active { color: var(--primary); border-bottom-color: var(--primary); }

        .tab-content { display: none; }
        .tab-content.active { display: block; }

        .alert { padding: 12px 16px; border-radius: 8px; font-size: 0.875rem; margin-bottom: 20px; line-height: 1.4; display: none; white-space: pre-line; }
        .alert-error { background-color: var(--error-bg); color: var(--error-text); border: 1px solid #fecaca; }
        .alert-success { background-color: var(--success-bg); color: var(--success-text); border: 1px solid #bbf7d0; }

        .form-group { margin-bottom: 18px; }
        label { display: block; font-size: 0.875rem; font-weight: 600; margin-bottom: 6px; }
        input[type="password"], input[type="text"] { width: 100%; padding: 10px 12px; font-size: 0.9rem; border: 1px solid var(--border); border-radius: 6px; outline: none; }
        input:focus { border-color: var(--primary); }

        .file-wrapper { display: flex; align-items: center; gap: 8px; }
        input[type="file"] { flex: 1; padding: 8px; border: 1px solid var(--border); border-radius: 6px; font-size: 0.875rem; background: #fff; cursor: pointer; }
        .btn-remove { background: #ef4444; color: white; border: none; border-radius: 6px; padding: 10px 14px; cursor: pointer; display: none; font-weight: bold; }
        .btn-remove:hover { background: #dc2626; }

        .prefix-badge { font-size: 0.75rem; color: var(--primary); font-weight: bold; margin-top: 4px; }
        
        button[type="submit"] { width: 100%; padding: 12px; background-color: var(--primary); color: white; font-weight: 600; border: none; border-radius: 6px; cursor: pointer; font-size: 0.95rem; margin-top: 8px; }
        button[type="submit"]:hover { background-color: var(--primary-hover); }
        button:disabled { background-color: var(--text-muted); cursor: not-allowed; }
    </style>
</head>
<body>
    <div class="container">
        <h1>Conversor de Certificados</h1>
        <p class="subtitle">Selecione o tipo de conversão desejado abaixo.</p>

        <!-- Navegação entre Abas -->
        <div class="tabs">
            <button class="tab-btn active" onclick="alternarAba('pemTab', this)">PFX para PEMs Soltos</button>
            <button class="tab-btn" onclick="alternarAba('nextTab', this)">Compatível com Next.js (PFX)</button>
        </div>

        <div id="alertError" class="alert alert-error"></div>
        <div id="alertSuccess" class="alert alert-success"></div>

        <!-- ABA 1: CONVERTER PFX PARA PEM -->
        <div id="pemTab" class="tab-content active">
            <form id="formPem">
                <div class="form-group">
                    <label>1. Selecione o Certificado (.PFX / .P12)</label>
                    <div class="file-wrapper">
                        <input type="file" name="pfxfile" accept=".pfx,.p12" required onchange="validarArquivo(this, 'btnRemovePem')">
                        <button type="button" id="btnRemovePem" class="btn-remove" onclick="removerArquivo('formPem', 'btnRemovePem')">🗑️</button>
                    </div>
                </div>

                <div class="form-group">
                    <label>2. Senha do Certificado</label>
                    <input type="password" name="password" placeholder="Digite a senha" required>
                </div>

                <div class="form-group">
                    <label>3. Identificador do Certificado Público (Sufixo)</label>
                    <input type="text" name="nomeCert" placeholder="Ex: Producao">
                    <div class="prefix-badge">O arquivo será: ChavePublica_&lt;seu_texto&gt;.pem</div>
                </div>

                <div class="form-group">
                    <label>4. Identificador da Chave Privada (Sufixo)</label>
                    <input type="text" name="nomeChave" placeholder="Ex: Producao">
                    <div class="prefix-badge">O arquivo será: ChavePrivada_&lt;seu_texto&gt;.pem</div>
                </div>

                <button type="submit" class="btnSubmit">Converter & Baixar PEMs</button>
            </form>
        </div>

        <!-- ABA 2: COMPATÍVEL COM NEXT.JS (PFX RE-ENCODADO) -->
        <div id="nextTab" class="tab-content">
            <form id="formNext" action="/convert-next" method="POST" enctype="multipart/form-data">
                <div class="form-group">
                    <label>1. Selecione o Certificado Original (.PFX / .P12)</label>
                    <div class="file-wrapper">
                        <input type="file" name="pfxfile" accept=".pfx,.p12" required onchange="validarArquivo(this, 'btnRemoveNext')">
                        <button type="button" id="btnRemoveNext" class="btn-remove" onclick="removerArquivo('formNext', 'btnRemoveNext')">🗑️</button>
                    </div>
                </div>

                <div class="form-group">
                    <label>2. Senha Atual do Certificado</label>
                    <input type="password" name="password" placeholder="Digite a senha atual" required>
                </div>

                <div class="form-group">
                    <label>3. Nome do Novo Arquivo (Opcional)</label>
                    <input type="text" name="nomeSaida" placeholder="Ex: certificado_next">
                    <div class="prefix-badge">O arquivo final será baixado como .pfx</div>
                </div>

                <button type="submit" class="btnSubmit">Converter & Baixar PFX Novo</button>
            </form>
        </div>
    </div>

    <script>
        function alternarAba(abaId, elementoBotao) {
            document.querySelectorAll('.tab-content').forEach(tab => tab.classList.remove('active'));
            document.querySelectorAll('.tab-btn').forEach(btn => btn.classList.remove('active'));
            
            document.getElementById(abaId).classList.add('active');
            elementoBotao.classList.add('active');

            document.getElementById('alertError').style.display = 'none';
            document.getElementById('alertSuccess').style.display = 'none';
        }

        function validarArquivo(input, btnId) {
            const btnRemove = document.getElementById(btnId);
            if (input.files && input.files.length > 0) {
                btnRemove.style.display = 'inline-block';
            } else {
                btnRemove.style.display = 'none';
            }
        }

        function removerArquivo(formId, btnId) {
            const form = document.getElementById(formId);
            const input = form.querySelector('input[type="file"]');
            const btnRemove = document.getElementById(btnId);
            input.value = '';
            btnRemove.style.display = 'none';
        }

        function dispararDownload(conteudoTextual, nomeArquivo) {
            const blob = new Blob([conteudoTextual], { type: 'application/x-pem-file' });
            const url = window.URL.createObjectURL(blob);
            const a = document.createElement('a');
            a.style.display = 'none';
            a.href = url;
            a.download = nomeArquivo;
            document.body.appendChild(a);
            a.click();
            window.URL.revokeObjectURL(url);
            document.body.removeChild(a);
        }

        // Submissão Aba 1 (PEMs Soltos via JSON)
        document.getElementById('formPem').addEventListener('submit', async function(e) {
            e.preventDefault();

            const alertError = document.getElementById('alertError');
            const alertSuccess = document.getElementById('alertSuccess');
            const btnSubmit = this.querySelector('.btnSubmit');

            alertError.style.display = 'none';
            alertSuccess.style.display = 'none';
            btnSubmit.disabled = true;
            btnSubmit.innerText = 'Convertendo...';

            const formData = new FormData(this);

            try {
                const response = await fetch('/convert-pem', {
                    method: 'POST',
                    body: formData
                });

                if (!response.ok) throw new Error('Servidor retornou erro HTTP ' + response.status);

                const data = await response.json();

                if (!data.sucesso) {
                    alertError.innerText = data.erro;
                    alertError.style.display = 'block';
                } else {
                    dispararDownload(data.certPem, data.nomeCert);
                    setTimeout(() => dispararDownload(data.keyPem, data.nomeChave), 300);

                    alertSuccess.innerText = "Download dos arquivos iniciado:\n• " + data.nomeCert + "\n• " + data.nomeChave;
                    alertSuccess.style.display = 'block';
                }
            } catch (err) {
                alertError.innerText = "Erro ao se comunicar com a aplicação: " + err.message;
                alertError.style.display = 'block';
            } finally {
                btnSubmit.disabled = false;
                btnSubmit.innerText = 'Converter & Baixar PEMs';
            }
        });
    </script>
</body>
</html>`))

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

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		_ = page.Execute(w, nil)
	})

	http.HandleFunc("/convert-pfx", func(w http.ResponseWriter, r *http.Request) {
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
			responseJSON(w, Page{Success: false, Error: "ERROR"})
			return
		}

		senha := r.FormValue("password")
		novaSenha := r.FormValue("novaSenha")

		if novaSenha == "" {
			novaSenha = senha
		}

		file, _, err := r.FormFile("pfxfile")
		if err != nil {
			responseJSON(w, Page{Success: false, Error: "ERROR"})
			return
		}

		defer file.Close()

		certPFX, err := io.ReadAll(file)
		if err != nil {
			responseJSON(w, Page{Success: false, Error: "ERROR"})
			return
		}

		newPfx, err := convert.ConvertPfxNext(certPFX, senha, novaSenha)
		if err != nil {
			responseJSON(w, Page{Success: false, Error: "ERROR"})
			return
		}

		nomeSaida := formatarNomeArquivo(r.FormValue("nomeSaida"), "certificadonovo")
		nomeSaida = strings.TrimSuffix(nomeSaida, ".pfx")

		w.Header().Set("Content-Type", "application/pkcs-12")
		w.Header().Set("Content-Disposition", `attachment, filename="`+nomeSaida+`"`)
		w.Write(newPfx)

	})

	go abrirNavegador("http://localhost:8080")

	fmt.Println("Servidor inicializado em http://localhost:8080")

	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Printf("Erro ao iniciar o servidor")
	}
}
