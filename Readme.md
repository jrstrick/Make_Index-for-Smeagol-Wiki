# Make Index for [Smeagol(wiki)](https://github.com/AustinWise/smeagol)

This program exists because I am lazy. I had an entire directory tree full of markdown files, and as a novelist writing in markdown, I generate more. A lot. I wanted to be able to keep track of all my characters, events, and the chapters of the novels in a wiki so I don't go crazy looking up how old my main character was three novels ago in the spring, and what color her cat is. After much searching, I settled on [smeagol](https://github.com/AustinWise/smeagol).

Only one little problem. At that time, smeagol's default response, in light of a missing index file (say, Index.md, for argument) was to do nothing. I wasn't about to write all those dozens of index files by hand, especially in light of how much they change. This program is the result.

This program will go through a directory tree and generate an index file (Index.md, by default) in each subdirectory it finds containing links to all the files in that subdirectory. It will then put links to each index file in each directory in the parent directory's index file, producing a navigable wiki tree.

## BREAKING CHANGES for V0.2

Why breaking changes? On chatting with AustinWise, author of smeagol, he pointed out that this utility could be useful for other things. In order to make make_index more flexible for that purpose, and at the same time better integrated with smeagol, I changed how parameter handling was done and added several parameters. Since my userbase appears to be one (me), this seemed like a good point to make this change.

1. ### The Parameter System Has Changed
   
   In the interests of better parameter handling, I switched to Go's flags package instead of parsing the parameters myself the way we did back in the stone age. This did break one feature. Before, you could type
   
   \$``make_index /home/jim/mywiki``
   
   and make_index would treat the parameter /home/jim/mywiki as the root of the wiki you're trying to index. This no longer works. If you pass make_index parameters this way, they will be ignored.
   
   The new way is to type
   \$``make_index -wiki_root /home/jim/wiki``

2. ### The Default Index File Name Has Changed.
   
   By default, make_index now produces README.md files. You can change this, either by passing the -index_file flag on the command line, or by having a smeagol.toml file in your wiki root that sets the index file name —aka index-page—to some other name.
   
   If you were using the previous version of make_index with smeagol, you already have that smeagol.toml file specifying Index.md as the index file name so smeagol could find the indexes you created. This will continue to work.

## To Use:

$ ``make_index [OPTIONS]``

### Option flags:

These flags are all optional. By default, make_index will run in the current directory, read from smeagol.toml if it's there, and write README.md files. unless otherwise specified in smeagol.toml.

```
-config_file <filename.toml> default: smeagol.toml, ignored if not found.
-index_file <filename.md> default: README.md or whatever is in smeagol.toml
-wiki_root <path name> default: current directory
-debug
```

- -config_file specifies the name of the config file to look for in the wiki_root directory.
- -index_file specifies what make_index should use for its index files, and also what name existing index files will have.
- -wiki_root specifies the root directory for the wiki. Almost certainly the same place you're pointing smeagol itself.

### Examples

- $``./make_index -wiki_root ~/home/jim/my_wiki -config_file notsmeagol.toml -index_file Index.md -debug``
  
  Here, we explicitly tell make_index where the wiki root is,what config file to load when it gets there, and what file name to use when generating the index files, and turn on the very, very chatty debug.

- $``./make_index -wiki_root ~/home/jim/my_wiki``
  
  Here, we just tell make_index what the root of the wiki is, and accept smeagol.toml as our config file, and enjoy the peace and quiet without debug.
  
  If smeagol.toml is present in the wiki root directory, make_index will set the index file name to the same thing smeagol is using. Otherwise it will revert to the default used by smeagol: README.md

## Features:

make_index's behavior can be changed through the use of flags. It can use a custom config file, or read your smeagol.toml file by default. It can use different names for the index files. And you can call it from some other directory and pass it the root of your wiki. See Option Flags, above.

make_index is multithreaded, using go-routines. Depending on the size of your directory tree, this may speed things up dramatically.

make_index will refrain from clobbering user-generated index files, or in fact any index file which does not contain the YAML tag: AUTOGEN.

make_index understands the following filetypes:

- Markdown files, with the extention .md
- Directories
- Jpeg, Gif, and PNG files with the extensions .jpg, .jpeg .gif or .png

It will ignore any other files it finds.

## Idiosyncracies:

- make_index will choke on any filename with spaces in it. I've found the Linux command "detox" prior to running make_index solves this in a suitably automatic fashion.
- make_index will ignore any directory whose name is not capitalized. This is so your wiki can have other stuff, like a code directory, without having that directory indexed.
- make_index will ignore any file type it does not explicitly understand. This means you can have things like smeagol.toml where they need to be and not have them indexed.
- When make_index indexes your image directory, it will put an anchor link rather than an image link to any image that begins with nsfw . There are some images you probably don't want popping onto your browser in the coffee shop.
- make_index is dumb about file types. It goes strictly off the dot extension, so any file with .jpg will be linked as an image, any file with a .md extentson will be linked as a markdown file, and any file with an extension it doesn't know will be ignored.

## Building:

```bash
git clone https://github.com/jrstrick/Make_Index-for-Smeagol-Wiki.git
cd  Make_Index-for-Smeagol-Wiki
go build
cp ./make_index [somewhere on your search path for binaries]
cp ./smeagol.toml [root directory of your smeagol-wiki]
```

## Notes

- Q: Wait, where did all the log/debug spam go?
  A: Pass the flag -debug when calling make_index, and it will be as chatty as you remember.
- Q : Why Go? Why not Rust?
  A: I'm learning Go for another project because it has a mature graphics toolkit. Rust does not. Go is fast *enough* and safe *enough* for this project.
- Q: Why do you document your code so much?
  A: I'm still learning Go, and in six months, when I'm trying to fix your bug report, those comments will make it possible to recover my thought process when I wrote the code.
- Q: Can you add feature X?
  A: Maybe. I mostly wrote this for myself. I'm always willing to listen to cool ideas though. I won't take it personally if you fork this project. I might even switch to your fork if it's better.
- Q: Why go-routines? This should be a proper recursion!
  A: Because I could. :) Also because recursion chews up stack memory, and why have all those cores sitting idle while one thread does all the work? Also, end state testing in actual recursion seemed harder.
