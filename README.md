# unflac

A command line tool for fast frame accurate audio image + cue sheet splitting.

This project is started mostly out of frustration over supporting
[split2flac](https://github.com/ftrvxmtrx/split2flac) with all the
external dependencies and their quirks.

## Installation/upgrade and running

You need [Go](https://golang.org/) installed.

```
go install git.sr.ht/~ft/unflac@latest
~/go/bin/unflac -h
```

## Dependencies

 * ffmpeg (ffprobe command is also used, part of ffmpeg package)

Yeah, that's it.

## Some useful options

 * Set a specific sample format and rate for output files:

   `unflac -F -sample_fmt -F s16 -F -ar -F 192000 ...`

   See `ffmpeg -sample_fmts`. `man ffmpeg` contains a lot more.

## Output file naming

To set a custom output file naming use `-n` option (run `unflac -h` to see the default value).
The format of the argument is described [here](https://golang.org/pkg/text/template).

```
Elem - a function that replaces invalid file path characters with ones that look (almost) the same but valid

.Input.TrackNumberFmt - track number printf format (it's either "%02d" or "%03d" depending on the number of tracks)
.Input.Composer       - composer (can be empty)
.Input.Performer      - performer (can be empty)
.Input.SongWriter     - song writer (can be empty)
.Input.Title          - title, that's album name in general (can be empty)
.Input.Genre          - genre
.Input.Date           - date
.Input.TotalTracks    - total number of tracks
.Input.TotalDisks     - total number of disks, 0 if there is only one disk
.Input.Artist         - a special handy field that is selected based on Composer/Performer/SongWriter; is never empty

.Track.Number         - track number
.Track.DiskNumber     - disk number, 0 if there is only one disk
.Track.Composer       - composer (can be empty)
.Track.Performer      - performer (can be empty)
.Track.SongWriter     - song writer (can be empty)
.Track.Title          - track title
.Track.Artist         - a special handy field that is selected based on Composer/Performer/SongWriter; is never empty
```

## TODO

 * copying cue, log, images over to the destination
 * clean up on errors
 * replay gain
 * support multiple artist/performer/composer?

## DONE

 * basic splitting
 * tagging
 * split to other formats (ogg/vorbis, mp3)
 * parallel extraction
 * support more input formats (wavpack, wave, mac)
 * get the artist/performer/composer right
 * multi-disk cue sheets support
 * detect cue sheet encoding and convert to utf8
 * an option to format the output path in a specific way

## NOTES

 * using id3 v2.3 tags because it seems basically no software reads id3 v2.4
