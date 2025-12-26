package log

import (
	"log/slog"
	"os"

	"github.com/lmittmann/tint"
	"github.com/mattn/go-colorable"
	"github.com/mattn/go-isatty"
)

func init() {
	w := os.Stdout
	slog.SetDefault(slog.New(
		tint.NewHandler(colorable.NewColorable(w), &tint.Options{
			TimeFormat: "Mon, Jan 2 2006, 3:04:05 pm MST",
			NoColor:    !isatty.IsTerminal(w.Fd()),
			Level:      slog.LevelDebug,
		}),
	))
}
