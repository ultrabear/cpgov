# cpgov
A simple cli cpu governor control tool  
## Notice
This was made as a replacement for a bash script that did the same job, it is not a professional program by any means  
I placed this on my github as a convenience
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
