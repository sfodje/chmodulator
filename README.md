# chmodulator

A command-line script for deciphering linux permisions

## Installation
Download relevant executable from [here](https://github.com/sfodje/chmodulator/releases)

## Usage
```bash
$ chmodulator 0655
-rw-r-xr-x
0655
Owner: rw
Group: rx
Other: rx
```

```bash
$ chmodulator -rwxr--r--
-rwxr--r--
0744
Owner: rwx
Group: r
Other: r
```

```bash
$ chmodulator -owner=rwx -group=r-- -other=r--
-rwxr--r--
0744
Owner: rwx
Group: r
Other: r
```
