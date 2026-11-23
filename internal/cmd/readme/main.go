// Package main regenerates the parts of README.md that are derived from the
// code: the count of typed tasks and the provider support matrix.
//
// Run it from the repository root:
//
//	go run ./internal/cmd/readme
package main

import (
	"errors"
	"fmt"
	"io/fs"
	"log/slog"
	"os"
	"reflect"
	"slices"
	"strconv"
	"strings"

	"github.com/aarock1234/unicap/v3"
	"github.com/aarock1234/unicap/v3/provider/anticaptcha"
	"github.com/aarock1234/unicap/v3/provider/capsolver"
	"github.com/aarock1234/unicap/v3/provider/twocaptcha"
	"github.com/aarock1234/unicap/v3/tasks"
)

const (
	// readmePath is the generated document, relative to the working directory.
	readmePath = "README.md"

	// readmeMode is the permission the document is written with.
	readmeMode fs.FileMode = 0o644

	// placeholderKey satisfies the providers' constructors. No request is
	// ever sent with it.
	placeholderKey = "readme"
)

const (
	// supportedMark fills the cell of a provider that solves the task.
	supportedMark = "✓"

	// unsupportedMark fills the cell of a provider that rejects the task.
	unsupportedMark = "–"

	// footnoteMark follows the name of a typed task no provider solves. The
	// footnote it points to is hand-written below the matrix.
	footnoteMark = "¹"

	// taskHeader titles the first column of the matrix.
	taskHeader = "Task"

	// taskAlignment leaves the task column with the default alignment.
	taskAlignment = "---"

	// providerAlignment centers a provider column.
	providerAlignment = ":---:"
)

// errMarkerMissing reports that the document lacks a marker of a generated
// region.
var errMarkerMissing = errors.New("marker missing")

// column is one provider column of the support matrix.
type column struct {
	title    string
	supports func(unicap.TaskType) bool
}

// row is one task row of the support matrix. It holds one support flag per
// column, in column order.
type row struct {
	label     string
	supported []bool
}

// region is a generated span of the document, delimited by the HTML comments
// "<!-- name:start -->" and "<!-- name:end -->".
type region struct {
	name    string
	content string
}

func main() {
	if err := run(); err != nil {
		slog.Error("readme generation failed", slog.Any("error", err))
		os.Exit(1)
	}
}

// run rewrites the generated regions of the README in the working directory.
func run() error {
	columns, err := newColumns()
	if err != nil {
		return err
	}

	typedTasks := tasks.All()

	regions := []region{
		{
			name:    "task-count",
			content: strconv.Itoa(len(typedTasks)),
		},
		{
			name:    "support-matrix",
			content: "\n" + renderMatrix(columns, newRows(columns, typedTasks)),
		},
	}

	document, err := os.ReadFile(readmePath)
	if err != nil {
		return fmt.Errorf("reading readme: %w", err)
	}

	generated, err := applyRegions(string(document), regions)
	if err != nil {
		return fmt.Errorf("generating %s: %w", readmePath, err)
	}

	if err := os.WriteFile(readmePath, []byte(generated), readmeMode); err != nil {
		return fmt.Errorf("writing readme: %w", err)
	}

	return nil
}

// newColumns creates the built-in providers and returns their columns in
// display order. Constructing a provider sends no request.
func newColumns() ([]column, error) {
	capSolver, err := capsolver.New(placeholderKey)
	if err != nil {
		return nil, fmt.Errorf("creating capsolver provider: %w", err)
	}

	twoCaptcha, err := twocaptcha.New(placeholderKey)
	if err != nil {
		return nil, fmt.Errorf("creating 2captcha provider: %w", err)
	}

	antiCaptcha, err := anticaptcha.New(placeholderKey)
	if err != nil {
		return nil, fmt.Errorf("creating anti-captcha provider: %w", err)
	}

	return []column{
		{
			title:    "CapSolver",
			supports: capSolver.Supports,
		},
		{
			title:    "2Captcha",
			supports: twoCaptcha.Supports,
		},
		{
			title:    "Anti-Captcha",
			supports: antiCaptcha.Supports,
		},
	}, nil
}

// newRows returns one row per typed task, followed by the RawTask row.
func newRows(columns []column, typedTasks []unicap.Task) []row {
	rows := make([]row, 0, len(typedTasks)+1)

	for _, task := range typedTasks {
		supported := supportFlags(columns, task)

		rows = append(rows, row{
			label:     typedLabel(typeName(task), supported),
			supported: supported,
		})
	}

	rawTask := &tasks.RawTask{}

	return append(rows, row{
		label:     code(typeName(rawTask)) + " (passthrough)",
		supported: supportFlags(columns, rawTask),
	})
}

// supportFlags reports, for each column in order, whether its provider solves
// the task.
func supportFlags(columns []column, task unicap.Task) []bool {
	supported := make([]bool, len(columns))

	for index, column := range columns {
		supported[index] = column.supports(task.Type())
	}

	return supported
}

// typeName returns the Go type name of a task given as a pointer to its
// struct.
func typeName(task unicap.Task) string {
	return reflect.TypeOf(task).Elem().Name()
}

// typedLabel returns the first cell of a typed task's row, with the footnote
// mark when no provider solves the task.
func typedLabel(name string, supported []bool) string {
	if slices.Contains(supported, true) {
		return code(name)
	}

	return code(name) + footnoteMark
}

// code formats an identifier as inline Markdown code.
func code(identifier string) string {
	return "`" + identifier + "`"
}

// renderMatrix renders the support matrix as a GitHub-flavoured Markdown
// table without column padding. Every line ends with a newline.
func renderMatrix(columns []column, rows []row) string {
	header := []string{taskHeader}
	alignment := []string{taskAlignment}

	for _, column := range columns {
		header = append(header, column.title)
		alignment = append(alignment, providerAlignment)
	}

	var table strings.Builder

	table.WriteString(renderTableLine(header))
	table.WriteString(renderTableLine(alignment))

	for _, row := range rows {
		table.WriteString(renderTableLine(row.cells()))
	}

	return table.String()
}

// cells returns the row's table cells: its label, then one mark per column.
func (r row) cells() []string {
	cells := []string{r.label}

	for _, supported := range r.supported {
		cells = append(cells, supportMark(supported))
	}

	return cells
}

// supportMark returns the cell text for a support flag.
func supportMark(supported bool) string {
	if supported {
		return supportedMark
	}

	return unsupportedMark
}

// renderTableLine joins cells into one Markdown table line.
func renderTableLine(cells []string) string {
	return "| " + strings.Join(cells, " | ") + " |\n"
}

// applyRegions replaces the content of every region in the document.
func applyRegions(document string, regions []region) (string, error) {
	for _, region := range regions {
		start := fmt.Sprintf("<!-- %s:start -->", region.name)
		end := fmt.Sprintf("<!-- %s:end -->", region.name)

		replaced, err := replaceBetween(document, start, end, region.content)
		if err != nil {
			return "", fmt.Errorf("replacing %s region: %w", region.name, err)
		}

		document = replaced
	}

	return document, nil
}

// replaceBetween replaces the text between the first start marker and the
// first end marker after it. The markers themselves are kept.
func replaceBetween(document, start, end, replacement string) (string, error) {
	before, rest, found := strings.Cut(document, start)
	if !found {
		return "", fmt.Errorf("%w: %s", errMarkerMissing, start)
	}

	_, after, found := strings.Cut(rest, end)
	if !found {
		return "", fmt.Errorf("%w: %s", errMarkerMissing, end)
	}

	return before + start + replacement + end + after, nil
}
