# DarkMap
An EAC Mapper written in Go-lang

# Usage
You need a VDM, I recommend Xeroxz's one

# How it works
When the EasyAntiCheat driver is initialized, it walks through each loaded driver's read-only sections with MmCopyMemory to ensure that malicious patches have not taken place. But, EasyAntiCheat has a slight oversight resulting in certain drivers, known as session drivers, to not be accounted for during these initial scans. From my debugging, EasyAntiCheat entirely skips session drivers and does not make any attemps to ensure their integrity.

# Credits
Compiled-Code for the base
VDM-Visor for Golang rewrite
