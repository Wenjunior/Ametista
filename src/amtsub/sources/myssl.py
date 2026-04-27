import requests

class MySSL:
	async def search(self, domain, timeout):
		try:
			url = 'https://myssl.com/api/v1/discover_sub_domain?domain=' + domain

			response = requests.get(url, timeout = timeout)

			certificates = response.json()['data']

			subdomains = []

			for certificate in certificates:
				subdomain = certificate['domain']

				subdomains.append(subdomain)

			return subdomains
		except:
			pass

	def get_name(self):
		return 'MySSL'