# cpgov
A simple cli cpu governor control tool

# Installation
Normal install
Usage: `sudo cpgov` recommended for multi user installs
```bash
sudo make install
```
Setuid install
Usage: `cpgov` not recommended for multi user installs
```bash
sudo make install-setuid
```
# Usage
List governors availible and current governor
```bash
cpgov
```
Set a certain cpu governor for all cpus
```bash
cpgov <governor>
```
