package paca

import (
	"bytes"
	"fmt"
	"os"
	"strings"

	"golang.org/x/net/html"
)

var (
	errCond = "\tassertNoErr(t, %s)\n"
)

type ideCase struct {
	Title string
	Cmds  []*selCmd
}

func (c *ideCase) TestCode() string {
	output := bytes.Buffer{}
	output.WriteString("package seltest\n\nimport(\n\t\"testing\"\n\t\"time\"\n\n\t\"github.com/sclevine/agouti\"\n)\n")
	output.WriteString("var _ = 42 * time.Second\n\n")
	output.WriteString(fmt.Sprintf("func %s(t *testing.T, page *agouti.Page){\n", Camelize(c.Title)))
	for _, cmd := range c.Cmds {
		output.WriteString(cmd.Code())
	}
	output.WriteString("}\n")

	return output.String()
}

type selCmd struct {
	Action string
	Target string
	Value  string
}

func (c *selCmd) Code() string {
	output := &bytes.Buffer{}
	switch c.Action {
	case "open":
		output.WriteString(fmt.Sprintf(errCond, fmt.Sprintf("page.Navigate(TargetHost + \"%s\")", c.Target)))
	case "store":
		output.WriteString(fmt.Sprintf("\t%s %s := \"%s\"\n", `//`, c.Value, c.Target))
	case "type":
		find := fmt.Sprintf(`page.Find("input[%s]")`, c.Target)
		output.WriteString(fmt.Sprintf(errCond, fmt.Sprintf("%s.Fill(\"%s\")", find, c.Value)))
	case "click":
		if strings.HasPrefix(c.Target, "//") {
			output.WriteString(fmt.Sprintf(errCond, fmt.Sprintf("page.FindByXPath(\"%s\").Click()", c.Target)))
		} else {
			parts := strings.Split(c.Target, "=")
			if len(parts) == 2 {
				switch parts[0] {
				case "css":
					output.WriteString(fmt.Sprintf(errCond, fmt.Sprintf("page.Find(\"%s\").Click()", parts[1])))
				case "id":
					output.WriteString(fmt.Sprintf(errCond, fmt.Sprintf("page.Find(\"#%s\").Click()", parts[1])))
				case "link":
					output.WriteString(fmt.Sprintf(errCond, fmt.Sprintf("page.FindByLink(\"%s\").Click()", parts[1])))
				default:
					output.WriteString("\t" + `// click ` + c.Target + "\n")
				}
			}
		}
	case "pause":
		output.WriteString("\ttime.Sleep(" + c.Target + " * time.Millisecond)\n")
	default:
		output.WriteString(fmt.Sprintf("\t%s %s | %s | %s\n", `//`, c.Action, c.Target, c.Value))
	}
	return output.String()
}

func IDEConverter(path string) (*ideCase, error) {
	// open file
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	// parse HTML
	doc, err := html.Parse(f)
	if err != nil {
		return nil, err
	}

	// extract case
	c := &ideCase{}
	htmlParser(c, doc)
	return c, nil
}

func htmlParser(c *ideCase, n *html.Node) {
	if n.Type == html.ElementNode {
		switch n.Data {
		case "title":
			if content := n.FirstChild; content != nil {
				c.Title = content.Data
			}
		case "tr":
			trParser(c, n)
		default:
			//fmt.Println("default:", n.Data)
		}
	}
	for child := n.FirstChild; child != nil; child = child.NextSibling {
		htmlParser(c, child)
	}

}

func trParser(c *ideCase, tr *html.Node) {
	tds := []*html.Node{}
	for c1 := tr.FirstChild; c1 != nil; c1 = c1.NextSibling {
		if c1.Type == html.ElementNode && c1.Data == "td" {
			tds = append(tds, c1)
		}
	}
	if len(tds) != 3 {
		return
	}
	if tds[0].FirstChild == nil {
		return
	}
	cmd := &selCmd{Action: tds[0].FirstChild.Data}
	if tds[1].FirstChild != nil {
		cmd.Target = tds[1].FirstChild.Data
	}
	if tds[2].FirstChild != nil {
		cmd.Value = tds[2].FirstChild.Data
	}
	c.Cmds = append(c.Cmds, cmd)
}

func Camelize(str string) string {
	output := &bytes.Buffer{}
	if strings.Contains(str, " ") {
		for _, el := range strings.Split(str, " ") {
			output.WriteString(strings.Title(el))
		}
	}
	if strings.Contains(str, "_") {
		for _, el := range strings.Split(str, "_") {
			output.WriteString(strings.Title(el))
		}
	}

	return output.String()
}

func HelperFileContent() string {
	return `package seltest

import(
	"testing"
	"time"
	"github.com/sclevine/agouti"
)

var TargetHost = "https://google.com" // EDIT ME
var _ = 42 * time.Second

func assertNoErr(t *testing.T, err error) {
	_, thisFile, thisLine, _ := runtime.Caller(1)
	undo := strings.Repeat("\x08", len(fmt.Sprintf("%s:%d:      ", filepath.Base(thisFile), thisLine)))
	_, file, line, _ := runtime.Caller(1)

	if err != nil {
		t.Fatalf("%s%s:%d - %v\n", undo, filepath.Base(file), line, err)
	}
}

func TestFirstScenario(t *testing.T) {
	driver := agouti.Selenium()
	if err := driver.Start(); err != nil {
		t.Fatal(err)
	}
	page, err := driver.NewPage(agouti.Browser("chrome"))
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		page.Destroy()
		driver.Stop()
	}()
	if err = page.Navigate(TargetHost); err != nil {
		t.Fatal("can't navigate to target host", err)
	}
	// call testFunction passing t and page
}
`
}
