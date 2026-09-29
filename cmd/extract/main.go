// Command extract preserves the original disk transfers and module byte spans.
package main

import (
	"crypto/sha256"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/olivierh59500/go-spaceballs/internal/source"
)

func main() {
	disk := flag.String("disk", "", "verified unpacked ADF path (required)")
	directory := flag.String("out", "assets/raw", "destination for extracted source data")
	flag.Parse()
	if *disk == "" {
		log.Fatal("-disk is required")
	}
	if err := run(*disk, *directory); err != nil {
		log.Fatal(err)
	}
}

func run(path, directory string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	if err := source.ValidateDisk(data); err != nil {
		return err
	}
	if err := os.MkdirAll(directory, 0755); err != nil {
		return err
	}
	type record struct {
		source.Region
		Bytes  int
		SHA256 string
	}
	var records []record
	for _, region := range source.Regions {
		bank, err := source.ReadRegion(data, region)
		if err != nil {
			return err
		}
		if err := os.WriteFile(filepath.Join(directory, region.Name+".bin"), bank, 0644); err != nil {
			return err
		}
		records = append(records, record{region, len(bank), fmt.Sprintf("%x", sha256.Sum256(bank))})
	}
	for name, offset := range map[string]int{"loader": 0x6e00, "state-of-the-art": 0x22226} {
		module, err := source.ReadModule(data, offset)
		if err != nil {
			return err
		}
		if err := os.WriteFile(filepath.Join(directory, name+".mod"), module, 0644); err != nil {
			return err
		}
	}
	manifest, err := json.MarshalIndent(records, "", "  ")
	if err != nil {
		return err
	}
	if err := os.WriteFile(filepath.Join(directory, "transfers.json"), append(manifest, '\n'), 0644); err != nil {
		return err
	}
	fmt.Printf("Extracted %d authored transfers and both verified MOD soundtracks.\n", len(records))
	return nil
}
