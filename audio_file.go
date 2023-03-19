package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

type AudioFile struct {
	Path       string   `json:"path"`
	Format     string   `json:"format"`
	SampleRate int      `json:"sampleRate,omitempty"`
	Tracks     []*Track `json:"tracks,omitempty"`
}

type Tag struct {
	Name  string
	Value string
}

func NewAudio(path string) (af *AudioFile, err error) {
	af = &AudioFile{Path: path, Format: strings.ToLower(filepath.Ext(path)[1:])}

	cmd := exec.Command(
		"ffprobe",
		"-loglevel", "error",
		"-of", "compact=nk=1:p=0",
		"-show_entries", "stream=sample_rate",
		"-select_streams", "a",
		af.Path,
	)
	var out strings.Builder
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr
	if err = cmd.Run(); err == nil {
		af.SampleRate, err = strconv.Atoi(strings.TrimSpace(out.String()))
	} else {
		err = fmt.Errorf("failed to get sample rate: %s", err)
	}

	return
}

func (af *AudioFile) Extract(t *Track, filename string) (err error) {
	// ffmpeg options
	args := []string{"-loglevel", "error", "-y", "-i", af.Path, "-map_metadata", "-1"}
	tags := []Tag{
		{"composer", t.Composer},
		{"artist", t.Artist()},
		{"performer", t.Performer},
		{"album", t.Album},
		{"title", t.Title},
		{"genre", t.Genre},
		{"date", t.Date},
	}

	var diskNumber string
	if t.DiskNumber != 0 {
		diskNumber = strconv.Itoa(t.DiskNumber)
	}
	var totalDisks string
	if t.TotalDisks != 0 {
		totalDisks = strconv.Itoa(t.TotalDisks)
	}

	switch *format {
	case "flac":
		tags = append(tags,
			Tag{"tracknumber", strconv.Itoa(t.Number)},
			Tag{"tracktotal", strconv.Itoa(*t.TotalTracks)},
			Tag{"discnumber", diskNumber},
			Tag{"totaldiscs", totalDisks},
		)
		if af.SampleRate > 192000 {
			args = append(args, "-ar", "192000")
		}

	case "ogg":
		tags = append(tags,
			Tag{"tracknumber", strconv.Itoa(t.Number)},
			Tag{"discnumber", diskNumber},
			Tag{"totaldiscs", totalDisks},
		)
		if af.SampleRate > 192000 {
			args = append(args, "-ar", "192000")
		}

	case "mp3":
		tags = append(tags,
			Tag{"track", fmt.Sprintf("%d/%d", t.Number, t.TotalTracks)},
		)
		if diskNumber != "" && totalDisks != "" {
			tags = append(tags, Tag{"disc", fmt.Sprintf("%s/%s", diskNumber, totalDisks)})
		}
		args = append(args, "-write_id3v2", "1", "-id3v2_version", "3", "-qscale:a", "3")
	}

	args = append(args, ffmpegArgs...)

	for _, t := range tags {
		if t.Value != "" {
			args = append(args, "-metadata", fmt.Sprintf("%s=%s", t.Name, t.Value))
		}
	}

	atrim := fmt.Sprintf("atrim=start_sample=%d", t.StartAtSample)
	if t.EndAtSample != 0 {
		atrim = fmt.Sprintf("%s:end_sample=%d", atrim, t.EndAtSample)
	}
	args = append(args, "-af", atrim, filename)

	cmd := exec.Command("ffmpeg", args...)
	cmd.Stderr = os.Stderr

	return cmd.Run()
}
