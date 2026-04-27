import requests

class CertificateSearch:
	async def search(self, domain, timeout):
		try:
			url = f'https://crt.sh/?q={domain}&output=json'

			response = requests.get(url, timeout = timeout)

			certificates = response.json()

			subdomains = []

			for certificate in certificates:
				common_name = certificate['common_name']

				subdomains.append(common_name)

				name_value = certificate['name_value']

				lines = name_value.splitlines()

				for line in lines:
					subdomains.append(line)

			return subdomains
		except:
			pass

	def get_name(self):
		return 'Certificate Search'