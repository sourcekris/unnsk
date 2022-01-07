## unnsk - extractor for NaShrinK files

A proprietary format by Nashsoft Systems of Bangalore, India. These files 
are DCL imploded with some basic header.

#### Requires

To build:
 - golang 1.16

To install binary release:
 - See releases for builds for most major OS and architectures.

#### Usage

```shell
usage: unnsk [-h] -e FILENAME [-d PATH]

Extract NaShrinK files.

optional arguments:
  -h, --help            show this help message and exit
  -e FILENAME, --extract FILENAME
                        The NSK file to extract.
  -d PATH, --destination PATH
                        An optional output folder.

Please file bugs on the GitHub Issues. Thanks
```

### Format Info
Reverse engineering of the format I found:

- 3 byte file signature: "NSK"
- 4 byte integer - Compressed data size
- 5 bytes unknown, usually:
- '\x20\x??\x??\x26\x54'
- 4 byte integer - Uncompressed data size
- 1 byte filename length
- n bytes Filename
- remainder - DCL Imploded payload

### References
Archive format information: http://fileformats.archiveteam.org/wiki/NaShrinK
DOS Archiver: https://www.sac.sk/download/pack/nsk50.zip

#### Author

- Kris Hunt (@ctfkris)