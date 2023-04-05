# Make Index for [Smeagol-Wiki](https://github.com/AustinWise/smeagol)

Got a nice directory tree full of markdown files you want Smeagol-Wiki to use? Want them indexed without having to do it yourself? This is the program for you.

## To Use:

I recommend stopping smeagol-wiki if it's already running. I haven't *seen* anything bad happen if I reindex the wiki while the server is running, but it seems unwise anyway.

```bash
$ make_index [root of your smeagol-wiki directory tree]
$ ~/.cargo/bin/smeagol-wiki --host 0.0.0.0 --fs [root of your smeagol-wiki directory tree]
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

## How it Works

- Function main() Location: main.go 
  Validates the path you passed make_index. Once the path you passed make_index has been validated, main.go calls the gen_index function in gen_index.go and passes it the validated path.

- Function gen_index(*path* string) Location: gen_index.go
  Calls gen_index_preflight with the path gen_index was passed. 
  
  If gen_index_preflight returns true, we create a new Index.md file, and iterate through the directory at the end of our path.
  
  If we hit an entry that, itself, is another directory, call gen_index as a go-routine (thread) with our path concatinated with the name of this new directory. Then create an anchor link in our Index.md to the Index.md in that new directory, and assume that by the time our new thread completes, there will *be* an Index.md there. 
  
  For any other type of file (determined by the dot-extention) we call one of several format_X_link functions to generate the exact text of the link, and write that to Index.md
  
  When we run out of entries in the directory at the end of our path, we return.

- Function format_X_link(*file_name* string, [*file_type* string]) Location: their own .go files.
  These functions generate the actual link text for each file type. Presently, only format_img_link requires the file_type string. In the future, I may standardize that API. 
  
  format_img_link is also special in that it does some logic to check to see whether you've named a given image nsfw_whatever, or just nsfw. If you have, it generates an anchor link, so you have to click to open the image. Otherwise, image links in an Index.md file open as images in your browser automatically.

- Function gen_index_preflight(*path* string) Location: gen_index.go
  If the Index.md file at the path given does not exist, or does exist but has an AUTOGEN tag in its YAML header, return true. Otherwise return false.

## Notes

- Q : Why Go? Why not Rust? A: I'm learning Go for another project because it has a mature graphics toolkit. Rust does not. Go is fast *enough* and safe *enough* for this project.

- Q: Why do you document your code so much? A: I'm still learning Go, and in six months, when I'm trying to fix your bug report, those comments will make it possible to recover my thought process when I wrote the code.

- Q: Can you add feature X? A: Maybe. I mostly wrote this for myself. I'm always willing to listen to cool ideas though.

- Q: I'm going to fork this so I can do it right! A: Feel free. If yours is better, I'll probably switch to it. :)

- 
