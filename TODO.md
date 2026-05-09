## Sub

- Checar se o domínio raiz é valido
- Adicionar novas fontes:
	- URLscan - `https://urlscan.io/api/v1/search?q={}`
	- Wayback Machine - `http://web.archive.org/cdx/search/cdx?url=*.{}/*&output=text&fl=original&collapse=urlkey`
	- CommonCrawl - `https://index.commoncrawl.org/collinfo.json`
		- Adicione `?url=.{}` ao final da URL encontrada em `cdx-api`

## Scan

- Embaralhar a ordem das portas antes de escanear
- Adicionar suporte a CIDR e range de IPs
- Adicionar as técnicas TCP SYN (stealth) scan e UDP scan

## Geral

- Organizar utils.go
- Analisar o código com um linter