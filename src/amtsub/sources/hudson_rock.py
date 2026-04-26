import requests

class HudsonRock:
	async def search(self, domain, timeout):
		try:
			url = 'https://cavalier.hudsonrock.com/api/json/v2/osint-tools/urls-by-domain?domain=' + domain

			response = requests.get(url, timeout = timeout)

			data = response.json()['data']

			employees_urls = data['employees_urls']

			clients_urls = data['clients_urls']

			compromised_urls = employees_urls + clients_urls

			subdomains = []

			for compromised_url in compromised_urls:
				compromised_url = compromised_url['url']

				subdomain = compromised_url.split('/')[2]

				subdomains.append(subdomain)

			return subdomains
		except:
			pass

	def get_name(self):
		return 'Hudson Rock'