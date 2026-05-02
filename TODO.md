## Sub

- Checar se o domínio raiz é valido
- Adicionar novas fontes:
	- URLscan - `https://urlscan.io/api/v1/search?q={}`
	- Wayback Machine - `http://web.archive.org/cdx/search/cdx?url=*.{}/*&output=text&fl=original&collapse=urlkey`
	- CommonCrawl - `https://index.commoncrawl.org/collinfo.json`
		- Adicione `?url=.{}` ao final da URL encontrada em `cdx-api`

## Scan

- Adicionar suporte a CIDR e range de IPs
- Checar e caso necessário aumentar o ulimit antes de iniciar a varredura
- Adicionar a opção de salvar os resultados