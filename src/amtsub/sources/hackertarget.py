import requests

class HackerTarget:
	async def search(self, domain, timeout):
		try:
			url = 'https://api.hackertarget.com/hostsearch/?q=' + domain

			response = requests.get(url, timeout = timeout)

			body = response.text

			lines = body.splitlines()

			subdomains = []

			for line in lines:
				subdomain = line.split(',')[0]

				subdomains.append(subdomain)

			return subdomains
		except:
			pass

	def get_name(self):
		return 'HackerTarget'