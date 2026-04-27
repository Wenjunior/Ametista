import requests

class CertSpotter:
	async def search(self, domain, timeout):
		try:
			url = f'https://api.certspotter.com/v1/issuances?domain={domain}&include_subdomains=true&expand=dns_names'

			response = requests.get(url, timeout = timeout)

			certificates = response.json()

			subdomains = []

			for certificate in certificates:
				dns_names = certificate['dns_names']

				for subdomain in dns_names:
					subdomains.append(subdomain)

			return subdomains
		except:
			pass

	def get_name(self):
		return 'Cert Spotter'