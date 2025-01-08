package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/urfave/cli/v2"
)

var downloadFlags = []cli.Flag{}

var DownloadCmd = cli.Command{
	Name:                   "download",
	Aliases:                []string{"d"},
	Usage:                  "Download datasets",
	UsageText:              "download [opt..] [directory]",
	UseShortOptionHandling: true,
	Flags:                  downloadFlags,
	Action: func(cCtx *cli.Context) error {
		if cCtx.NArg() == 0 {
			fmt.Println("\n  a directory is needed!\n")
			cli.ShowSubcommandHelpAndExit(cCtx, 1)
		}

		return Download(cCtx)
	},
}

func Download(cli *cli.Context) error {
	folder := cli.Args().First()
	configDir := filepath.Join(folder, "domains")
	if _, err := os.Stat(configDir); os.IsNotExist(err) {
		if err = os.MkdirAll(configDir, 0750); err != nil {
			fmt.Println(err)
		}
	}

	DownloadSuffix(configDir)

	return nil
}

func DownloadSuffix(dirname string) {
	fmt.Println("Downloading public suffix...")
	url := "https://publicsuffix.org/list/public_suffix_list.dat"
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)

	lines := strings.Split(string(body), "\n")
	var buffer bytes.Buffer

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "// ===BEGIN PRIVATE DOMAINS") {
			break
		}
		if line != "" && !strings.HasPrefix(line, "//") {
			line = strings.Replace(line, "*.", "", 1)
			line = strings.Replace(line, "!", "", 1)
			buffer.WriteString(line)
			buffer.WriteString("\n")
		}
	}

	if err := os.WriteFile(filepath.Join(dirname, "suffix.lst"), buffer.Bytes(), 0666); err != nil {
		fmt.Println(err)
	}

}
