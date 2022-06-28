# cnmap

> I created cmap while studying for the OSCP to quickly scan hosts and build working directories and files. 
> It became tedius to do this manually and write out common scans like nikto, gobuster, etc to kick off 
> enumeration. The scanner does not replace nmap or anything else, it is just used to quickly get started with
> a new hackthebox, tryhackme, vulnhun, proving grounds, or whatever else type box. 

## Features

- Scan a single, or multiple hosts
- TCP Syn Scan on each host
- Scans each open port for http service
- Creates a working directory
- Creates a README file and writes ports and useful scan command strings
- Future state -  maybe kick off scans for each host

## Usage

Switches:

	 -h		sing lehost - (format: ip,hostname)
	 -H  	hosts file, one per line (format: ip,hostname)
	 -o		create working directory and create README for each host

Examples:

- cmap -h ip,hostname -o 
- cmap -H hosts.txt -o

Notes:

- if no hostname is supplied the IP is used
- output dir is hardcoded to ./
- template path is currently hardcoded to /opt/cmap/template.md
- the template uses nunjucks-like syntax for placeholders 

## TODO

- add error checking (lol)
- run scans and write output files