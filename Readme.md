# BankHub Core (pfx-conversor)

Ferramenta interna de suporte da MK Solutions para conversão de certificados digitais mTLS e utilidades de apoio ao atendimento, usada nas integrações com Banco do Brasil, Itaú, Sicoob, Sicredi, Santander e Bradesco.

É um binário Go único, sem dependências externas em tempo de execução: o frontend inteiro (HTML/CSS/JS) é embutido no próprio executável via `go:embed`, e nenhum dado é gravado em disco — toda conversão acontece em memória.

## Funcionalidades

### Conversores de certificado (por banco)

- **Banco do Brasil, Sicoob, Santander, Bradesco** — `.pfx`/`.p12` → chave pública + privada em PEM. Extrai também a cadeia de certificação completa (raiz + intermediários) quando embutida no arquivo, exigida pelo portal Developer do BB.
- **Re-encodar PFX** — gera um novo `.pfx` a partir de um existente, trocando a senha (compatibilidade com ferramentas que exigem outro formato de criptografia).
- **Sicredi** — único banco que não usa `.pfx`: o portal já entrega um `.cer` + `.key` separados, então aqui só normalizamos os dois para PEM.
- **Itaú** — não usa conversão de PFX; gera um CSR (Certificate Signing Request) + chave privada RSA 2048 novos, para o processo de credenciamento mTLS.

### Utilidades

- **Substituição de Códigos** — extrai a coluna de códigos de fatura de uma planilha (`.csv`/`.xlsx`), localizando a coluna automaticamente pelo cabeçalho. Resultado pode ser copiado sem aspas ou entre aspas simples, e opcionalmente já monta uma query SQL pronta (`select * from payment where external_reference in (...)`, com ou sem filtro de status).
- **Validador de Remessa** — valida a integridade estrutural de arquivos de remessa bancária CNAB 240/400 (tamanho de linha, sequência dos registros, estrutura de lotes), apontando a linha exata de cada inconsistência.

## Rodando localmente

```bash
go run ./cmd/main.go
```

O servidor sobe em `http://localhost:8080` e abre o navegador automaticamente.

## Gerando um executável distribuível

Como o frontend é embutido no binário, um único arquivo compilado já é suficiente — quem for rodar não precisa instalar Go nem nada além do próprio executável.

```bash
# Windows
GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build -ldflags="-s -w" -o pfx-conversor.exe ./cmd/main.go

# macOS (Apple Silicon)
GOOS=darwin GOARCH=arm64 CGO_ENABLED=0 go build -ldflags="-s -w" -o pfx-conversor-mac ./cmd/main.go

# Linux
GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -ldflags="-s -w" -o pfx-conversor-linux ./cmd/main.go
```

## Estrutura do projeto

```
cmd/main.go                      servidor HTTP, rotas e handlers
internal/assets/                 embute e serve o frontend (static/index.html)
internal/convert/
  pfx_to_pem.go                  conversor genérico .pfx → PEM (Sicoob/Santander/Bradesco)
  pfx_bb.go                      conversor do BB, com extração da cadeia de certificação
  pfx_next.go                    re-encodar .pfx com nova senha
  pem_sicredi.go                 normaliza .cer/.key do Sicredi para PEM
  csr_itau.go                    gera CSR + chave privada para o Itaú
internal/utils/planilha.go    extração de códigos de fatura de planilhas
internal/utils/remessa.go      validador estrutural de remessa CNAB 240/400
```

## Rotas HTTP

| Rota | Descrição |
|---|---|
| `POST /convert-pem` | Converte `.pfx` → PEM (genérico ou BB, conforme `banco`) |
| `POST /convert-sicredi-pem` | Converte `.cer` + `.key` do Sicredi → PEM |
| `POST /convert-next` | Re-encoda um `.pfx` com nova senha |
| `POST /convert-csr` | Gera CSR + chave privada para o Itaú |
| `POST /planilha/faturas` | Extrai códigos de fatura de uma planilha |
| `POST /remessa/validar` | Valida um arquivo de remessa CNAB 240/400 |

## Dependências

- [`golang.org/x/crypto`](https://pkg.go.dev/golang.org/x/crypto) e [`software.sslmate.com/src/go-pkcs12`](https://pkg.go.dev/software.sslmate.com/src/go-pkcs12) — leitura/escrita de `.pfx` (o pacote da sslmate é usado no lugar do `x/crypto/pkcs12` porque lida corretamente com arquivos que trazem a cadeia de certificação embutida, algo comum nos certificados bancários brasileiros).
- [`github.com/xuri/excelize`](https://pkg.go.dev/github.com/xuri/excelize/v2) — leitura de planilhas `.xlsx`.