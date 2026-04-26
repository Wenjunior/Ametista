import requests

class AnubisDB:
	async def search(self, domain, timeout):
		try:
			url = 'https://anubisdb.com/anubis/subdomains/' + domain

			response = requests.get(url, timeout = timeout)

			subdomains = response.json()

			return subdomains
		except:
			pass

	def get_name(self):
		return 'AnubisDB'