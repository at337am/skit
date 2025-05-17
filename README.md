## What Is This

A bunch of small command-line scripts I wrote in Go to make my daily workflow easier.

## How to Use It

Make sure your system has the following installed:

- Go (version 1.16 or higher)
- [just](https://github.com/casey/just) - a command runner
- [ffmpeg](https://ffmpeg.org) - required by some media-related scripts

To build and install all commonly used scripts, simply run:

```bash
just install-all
```

## Features

Quickly clean and convert videos in bulk using the following tools:

```bash
$ xfixer ./
$ vid2mp4 ./
$ fmp4 ./
$ tmrn -d ./
```

