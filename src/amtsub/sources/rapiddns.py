import requests
from bs4 import BeautifulSoup

class RapidDNS:
	async def search(self, domain, timeout):
		try:
			url = 'https://rapiddns.io/subdomain/' + domain

			response = requests.get(url, timeout = timeout)

			body = response.text

			soup = BeautifulSoup(body, 'html.parser')

			tds = soup.find_all('td')

			subdomains = []

			for td in tds:
				subdomain = td.text

				subdomains.append(subdomain)

			anchors = soup.find_all('a')

			for anchor in anchors:
				href = anchor.href

				if href:
					subdomain = href.split('/')[2].split('#')[0]

					subdomains.append(subdomain)

			return subdomains
		except:
			pass

	def get_name(self):
		return 'RapidDNS'