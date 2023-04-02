# Make Index for Smeagol-Wiki

Got a nice directory tree full of markdown files you want Smeagol-Wiki to use? Want them indexed without having to do it yourself? This is the program for you.

## To Use:

```bash
$ ./make_index [root of your smeagol-wiki directory tree]
```

This will recursively dive your Smeagol-Wiki directory tree and create Index.md files in each directory.

## Features:

make_index is multithreaded, using Go's goroutines function. Depending on the size of your directory tree, this may speed things up dramatically.

make_index will refrain from clobbering user-generated Index.md files, or in fact any Index.md file which does not contain the YAML tag: AUTOGEN.

make_index understands the following filetypes:

- Markdown files, with the extention .md

- Directories

- Jpeg, Gif, and PNG files with the extensions .jpg .gif or .png

It will ignore any other files it finds.

## Idiosyncracies:

- make_index is very, very chatty. It's also multithreaded, so the log output may make no sense at all.

- make_index will choke on any filename with spaces in it. I've found the Linux command "detox" prior to running make_index solves this in a suitably automatic fashion.

- make_index will ignore any directory or file whose name is not capitalized. This is so your wiki can have other stuff, like the smeagol.toml file, or perhaps even a code directory, and not have those directories indexed.

- When make_index indexes your image directory, it will put an anchor link rather than an image link to any image that begins with nsfw_ . There are some images you probably don't want popping onto your browser in the coffee shop.

## Building:

```bash
$ git clone https://github.com/jrstrick/Make_Index-for-Smeagol-Wiki.git
$ cd  Make_Index-for-Smeagol-Wiki
$ go build
```

## Notes

- Q : Why Go? Why not Rust? A: I'm learning Go for another project because it has a mature graphics toolkit. Rust does not. Go is fast *enough* and safe *enough* for this project.

- Q: Why do you document your code so much? A: I'm still learning Go, and in six months, when I'm trying to fix your bug report, those comments will make it possible to recover my thought process when I wrote the code.

- Q: Can you add feature X? A: Maybe. I mostly wrote this for myself. I'm always willing to listen to cool ideas though.

- Q: I'm going to fork this so I can do it right! A: Feel free. If yours is better, I'll probably switch to it. :)
