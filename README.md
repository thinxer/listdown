listdown
========

Download a list of small files.

Each line of the list file looks like this:

    URL DEST_PATH

For example:

    http://example.com/1.png example.com/1.png
    http://example.net/2.png example.net/2.png

Files will never be partially written, and existing files will not be downloaded again. Since the program will buffer unfinished files in memory, it is not suitable for downloading large files.

Installation
------------

    go get github.com/thinxer/listdown

Usage
-----

Type `listdown -h`.
