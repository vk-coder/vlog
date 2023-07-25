package vlog_test

import (
	"bytes"
	"os"
	"regexp"
	"testing"

	vlog "github.com/vkar/vlog/src"
)

func TestLogFile(t *testing.T) {
	f, err := os.CreateTemp("", "")
	if err != nil {
		t.Fatalf("could not create temporary file, %v", err)
	}
	filename := f.Name()

	defer os.Remove(filename)

	lg := vlog.GetLogger("file-logger", nil)
	lg.SetOutput(f)
	lg.Tracef("%s", "Hello!")

	f.Close()

	contents, err := os.ReadFile(filename)
	if err != nil {
		t.Fatalf("could not open temporary file for reading, %v", err)
	}

	validate(t, "file-logger .* trace Hello!", contents)
}

func validate(t *testing.T, pattern string, logStatement []byte) {
	matched, err := regexp.Match(pattern, logStatement)
	if err != nil {
		t.Fatalf("trace logging failed, %v", err)
	}

	if !matched {
		t.Fatalf("trace logging match failed, log-statement: %s", logStatement)
	}
}

func TestLog(t *testing.T) {
	type subTest struct {
		name    string
		f       func(lg *vlog.Logger)
		pattern string
	}
	samples := []subTest{
		{
			"Trace",
			func(lg *vlog.Logger) {
				lg.Trace("hello!")
			},
			".* trace hello!",
		},
		{
			"Tracef",
			func(lg *vlog.Logger) {
				lg.Tracef("%s", "hello!")
			},
			".* trace hello!",
		},
		{
			"Critical",
			func(lg *vlog.Logger) {
				lg.Critical("hello!")
			},
			".* critical hello!",
		},
		{
			"Criticalf",
			func(lg *vlog.Logger) {
				lg.Criticalf("%s", "hello!")
			},
			".* critical hello!",
		},
	}

	for _, smpl := range samples {
		t.Run(smpl.name, func(t *testing.T) {
			buf := new(bytes.Buffer)
			lg := vlog.GetLogger("lg", nil)
			lg.SetOutput(buf)
			smpl.f(lg)
			validate(t, smpl.pattern, buf.Bytes())
		})
	}
}

func TestLogLevel(t *testing.T) {
	type subTest struct {
		name    string
		f       func(lg *vlog.Logger)
		pattern string
	}
	samples := []subTest{
		{
			"Trace",
			func(lg *vlog.Logger) {
				lg.Trace("hello!")
			},
			"",
		},
		{
			"Debug",
			func(lg *vlog.Logger) {
				lg.Debugf("%s", "hello!")
			},
			"",
		},
		{
			"Info",
			func(lg *vlog.Logger) {
				lg.Infof("%s", "hello!")
			},
			".* info hello!",
		},
		{
			"Warn",
			func(lg *vlog.Logger) {
				lg.Warnf("%s", "hello!")
			},
			".* warn hello!",
		},
	}

	for _, smpl := range samples {
		t.Run(smpl.name, func(t *testing.T) {
			buf := new(bytes.Buffer)
			lg := vlog.GetLogger("lg", nil)
			lg.SetOutput(buf)
			lg.SetLevel(vlog.INFO)
			smpl.f(lg)
			validate(t, smpl.pattern, buf.Bytes())
		})
	}
}

func BenchmarkLogTracef(b *testing.B) {
	buf := new(bytes.Buffer)
	lg := vlog.GetLogger("lg", nil)
	lg.SetOutput(buf)
	for i := 0; i < b.N; i++ {
		lg.Tracef("%d %s", i+1, "item")
	}
}

func BenchmarkLogFile(b *testing.B) {
	f, err := os.CreateTemp("", "")
	if err != nil {
		b.Fatalf("could not create temporary file, %v", err)
	}
	filename := f.Name()

	defer f.Close()
	defer os.Remove(filename)

	lg := vlog.GetLogger("file-logger", nil)
	lg.SetOutput(f)

	for i := 0; i < b.N; i++ {
		lg.Tracef("%s", "Hello!")
	}
}
