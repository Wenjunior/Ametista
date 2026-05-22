## Sub

- Checar se o domínio raiz é valido
- Adicionar novas fontes:
	- URLscan - `https://urlscan.io/api/v1/search?q={}`
	- Wayback Machine - `http://web.archive.org/cdx/search/cdx?url=*.{}/*&output=text&fl=original&collapse=urlkey`
	- CommonCrawl - `https://index.commoncrawl.org/collinfo.json`
		- Adicione `?url=.{}` ao final da URL encontrada em `cdx-api`

## Scan

- Checar se os números estão entre 0 e 65536
- Embaralhar a ordem das portas antes de escanear
- Adicionar suporte a CIDR e range de IPs
- Adicionar as técnicas TCP SYN (stealth) scan e UDP scan

## Probe

- Analisar o corpo da resposta para descobrir o content type
- Adicionar novas opções:
	- Mostrar o response time
	- Habilitar detecção de tecnologias
- Adicionar matchers:
	- `-mt` - response time
	- `-mc` - status code
	- `-ml` - content length
	- `-ms` - string
	- `-mr` - regex
- Adicionar filtros:
	- `-ft` - response time
	- `-fc` - status code
	- `-fl` - content length
	- `-fs` - string
	- `-fr` - regex