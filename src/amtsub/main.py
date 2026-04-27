#!/usr/bin/env python3
import os
import re
import sys
import asyncio
import argparse
from sys import stderr
from colorama import Fore, Style
from argparse import HelpFormatter, ArgumentParser

from sources.myssl import MySSL
from sources.rapiddns import RapidDNS
from sources.anubis_db import AnubisDB
from sources.hudson_rock import HudsonRock
from sources.hackertarget import HackerTarget
from sources.certificate_search import CertificateSearch

DEFAULT_TIMEOUT = 10

class CapitalisedHelpFormatter(HelpFormatter):
	def add_usage(self, usage, actions, groups, prefix = None):
		if not prefix:
			prefix = 'Usage: '

		return super(CapitalisedHelpFormatter, self).add_usage(usage, actions, groups, prefix)

def parse_arguments():
	parser = ArgumentParser(add_help = False, formatter_class = CapitalisedHelpFormatter)

	parser._optionals.title = 'Options'

	parser.add_argument('-h', action = 'help', help = argparse.SUPPRESS)

	parser.add_argument('-d', nargs = '+', dest = 'domains', help = 'Target root domain names')

	parser.add_argument('-l', dest = 'filename', help = 'File containing a list of target root domain names')

	parser.add_argument('-t', type = int, dest = 'timeout', default = DEFAULT_TIMEOUT, help = f'Set how many seconds it should wait for a response (default is {DEFAULT_TIMEOUT})')

	parser.add_argument('-o', dest = 'output', help = 'File to write results to')

	return parser.parse_args()

def eprint(error, color = Fore.YELLOW):
	stderr.write(f'{color}{error}{Style.RESET_ALL}\n')

def panic(error):
	eprint(error, color = Fore.RED)

	sys.exit(1)

async def enumerate_subdomains(domain, timeout):
	print('Enumerating subdomains for ' + domain)

	sources = [
		MySSL(),
		RapidDNS(),
		AnubisDB(),
		HudsonRock(),
		HackerTarget(),
		CertificateSearch()
	]

	tasks = []

	for source in sources:
		task = source.search(domain, timeout)

		tasks.append(task)

	all_results = []

	results = await asyncio.gather(*tasks)

	for source, result in zip(sources, results):
		if result:
			all_results += result

			continue

		name = source.get_name()

		eprint('Could not search on ' + name)

	pattern = r'[0-9a-z-\.]+' + domain

	all_results = [result for result in all_results if re.fullmatch(pattern, result)]

	all_results = sorted(list(set(all_results)))

	return all_results

async def main():
	args = parse_arguments()

	target_domains = []

	domains = args.domains

	if domains:
		target_domains += domains

	filename = args.filename

	if filename:
		try:
			with open(filename, 'r') as file:
				for line in file:
					line = line.strip()

					if not line:
						continue

					target_domains.append(line)
		except PermissionError:
			panic(f'Could not read "{filename}": permission denied')
		except FileNotFoundError:
			panic(f'Could not read "{filename}": file not found')
		except IsADirectoryError:
			panic(f'Could not read "{filename}": is a directory')
		except NotADirectoryError:
			panic(f'Could not read "{filename}": invalid path')
		except UnicodeDecodeError:
			panic(f'Could not read "{filename}": not encoded in UTF-8')

	timeout = args.timeout

	results = []

	for target_domain in target_domains:
		found_subdomains = await enumerate_subdomains(target_domain.removeprefix('.'), timeout)

		for found_subdomain in found_subdomains:
			print(Fore.GREEN + found_subdomain + Style.RESET_ALL)

		results += found_subdomains

	print(f'{len(results)} subdomains was discovered')

	output = args.output

	if output:
		try:
			with open(output, 'w') as file:
				line_separator = os.linesep

				content = line_separator.join(results)

				file.write(content)
		except PermissionError:
			panic(f'Could not create "{output}": permission denied')
		except NotADirectoryError:
			panic(f'Could not create "{output}": invalid path')

try:
	asyncio.run(main())
except KeyboardInterrupt:
	pass