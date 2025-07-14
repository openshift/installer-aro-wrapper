package main

// Copyright (c) Microsoft Corporation.
// Licensed under the Apache License 2.0.

// See https://github.com/openshift/installer/blob/master/cmd/openshift-install/log.go

import (
	"io"
	"strings"

	"github.com/sirupsen/logrus"
)

type fileHook struct {
	file      io.Writer
	formatter logrus.Formatter
	level     logrus.Level

	truncateAtNewLine bool
}

func newFileHook(file io.Writer, level logrus.Level, formatter logrus.Formatter) *fileHook {
	return &fileHook{
		file:      file,
		formatter: formatter,
		level:     level,
	}
}

func newFileHookWithNewlineTruncate(file io.Writer, level logrus.Level, formatter logrus.Formatter) *fileHook {
	f := newFileHook(file, level, formatter)
	f.truncateAtNewLine = true
	return f
}

func (h fileHook) Levels() []logrus.Level {
	var levels []logrus.Level
	for _, level := range logrus.AllLevels {
		if level <= h.level {
			levels = append(levels, level)
		}
	}

	return levels
}

func (h *fileHook) Fire(entry *logrus.Entry) error {
	// logrus reuses the same entry for each invocation of hooks.
	// so we need to make sure we leave them message field as we received.
	orig := entry.Message
	defer func() { entry.Message = orig }()

	msgs := []string{orig}
	if h.truncateAtNewLine {
		msgs = strings.Split(orig, "\n")
	}

	for _, msg := range msgs {
		// this makes it easier to call format on entry
		// easy without creating a new one for each split message.
		entry.Message = msg
		line, err := h.formatter.Format(entry)
		if err != nil {
			return err
		}

		if _, err := h.file.Write(line); err != nil {
			return err
		}
	}

	return nil
}
