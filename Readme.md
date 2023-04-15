# Make Index for [Smeagol(wiki)](https://github.com/AustinWise/smeagol)

This program exists because I am lazy. I had an entire directory tree full of markdown files, and as a novelist writing in markdown, I generate more. A lot. I wanted to be able to keep track of all my characters, events, and the chapters of the novels in a wiki so I don't go crazy looking up how old my main character was three novels ago in the spring, and what color her cat is. After much searching, I settled on [smeagol](https://github.com/AustinWise/smeagol).

Only one little problem. At that time, smeagol's default response, in light of a missing index file (say, Index.md, for argument) was to do nothing. I wasn't about to write all those dozens of index files by hand, especially in light of how much they change. This program is the result.

This program will go through a directory tree and generate an index file (Index.md, by default) in each subdirectory it finds containing links to all the files in that subdirectory. It will then put links to each index file in each directory in the parent directory's index file, producing a navigable wiki tree.

## BREAKING CHANGES

### Parameter system has changed

In the interests of better parameter handling, I switched to Go's flags package instead of parsing the parameters myself the way we did back in the stone age. This did break one feature. Before, you could type 

```bash
make_index /home/jim/mywiki
```

and make_index would treat the parameter /home/jim/mywiki as the root of the wiki you're trying to index. This no longer works. If you pass make_index parameters this way, they will be ignored.

The new way is to type

```bash
make_index -wiki_root /home/jim/wiki
```

This makes it much easier to add additional parameters. Which is why I changed it in the first place.

## MORE BREAKING CHANGES

If you don't have a smeagol.toml file in your wiki root *and* you don't specify an index_file, make_index's new default index file name is README.md, to stay in sync with smeagol's defaults. If you were using the previous version with smeagol, you already have that smeagol.toml file so smeagol could find the indexes you created. Now you have a choice. But it might break things.

## To Use:

I recommend stopping smeagol-wiki if it's already running. I haven't *seen* anything bad happen if I reindex the wiki while the server is running, but it seems unwise anyway.

```bash
$ make_index [OPTIONS]
```

## Option flags:

These flags are all optional, for flexability. By default, make_index will run in the current directory, read from smeagol.toml if it's there, and write Index.md files.

```
-config_file <filename.toml> default: smeagol.toml, ignored if not found.
-index_file <filename.md> default: README.md or whatever is in smeagol.toml
-wiki_root <path name> default: current directory
```

## Examples

```bash
$./make_index -wiki_root ~/home/jim/my_wiki -config_file smeagol.toml -index_file Index.md
//Here, we explicitly tell make_index where the wiki root is, what config file to load when it gets there, and //what file name to use when generating the index files.

//or

$./make_index -wiki_root ~/home/jim/my_wiki
//Here, we just tell make_index what the root of the wiki is, and accept smeagol.toml as our config file. If it's there, it will set the index file name to the same thing smeagol is using. Otherwise it will revert to my default, which is Index.md. Yes, I'm entirely aware this is different from smeagol's 
```

## Features:

make_index's behavior can be changed through the use of flags. It can use a custom config file, or read your smeagol.toml file by default. It can use different names for the index files. And you can call it from some other directory and pass it the root of your wiki. See Option Flags, above.

make_index is multithreaded, using go-routines. Depending on the size of your directory tree, this may speed things up dramatically.

make_index will refrain from clobbering user-generated index files, or in fact any index file which does not contain the YAML tag: AUTOGEN.

make_index understands the following filetypes:

- Markdown files, with the extention .md

- Directories

- Jpeg, Gif, and PNG files with the extensions .jpg .gif or .png

It will ignore any other files it finds.

## Idiosyncracies:

- make_index is very, very chatty. It's also multithreaded, so the log output may make no sense at all.

- make_index will choke on any filename with spaces in it. I've found the Linux command "detox" prior to running make_index solves this in a suitably automatic fashion.

- make_index will ignore any directory whose name is not capitalized. This is so your wiki can have other stuff, like a code directory, without having that directory indexed.

- make_index will ignore any file type it does not explicitly understand. This means you can have things like smeagol.toml where it needs to be and not have it in an index.

- When make_index indexes your image directory, it will put an anchor link rather than an image link to any image that begins with nsfw_ . There are some images you probably don't want popping onto your browser in the coffee shop.

- make_index is dumb about file types. It goes strictly off the dot extension, so any file with .jpg will be linked as an image, any file with a .md extentson will be linked as a markdown file, and any file with an extension it doesn't know will be ignored. I should probably replace this logic with proper mime typing.

## Building:

```bash
$ git clone https://github.com/jrstrick/Make_Index-for-Smeagol-Wiki.git
$ cd  Make_Index-for-Smeagol-Wiki
$ go build
$ cp ./make_index [somewhere on your search path for binaries]
$ cp ./smeagol.toml [root directory of your smeagol-wiki]
```

## Notes

- Q : Why Go? Why not Rust? A: I'm learning Go for another project because it has a mature graphics toolkit. Rust does not. Go is fast *enough* and safe *enough* for this project.

- Q: Why do you document your code so much? A: I'm still learning Go, and in six months, when I'm trying to fix your bug report, those comments will make it possible to recover my thought process when I wrote the code.

- Q: Can you add feature X? A: Maybe. I mostly wrote this for myself. I'm always willing to listen to cool ideas though.

- Q: I'm going to fork this so I can do it right! A: Feel free. If yours is better, I'll probably switch to it. :)

- Q: Why go-routines? This should be a proper recursion! A: Because I could. :) Also because recursion chews up stack memory, and why have all those cores sitting idle while one thread does all the work? Also, end state testing in actual recursion seemed harder.
