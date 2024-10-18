package library

import (
	"fmt"
	"os"
	"path"
	"strings"

	"mime"

	"github.com/dhowden/tag"
	"github.com/nktknshn/gomusiclibrary/cmd/cli"
	"github.com/spf13/cobra"
)

var cmdLibraryTags = cobra.Command{
	Use:  "tags",
	RunE: tags,
}

func tags(cmd *cobra.Command, args []string) error {
	db := cli.GetDatabaseMust()

	list, err := db.FilesList()

	if err != nil {
		return err
	}

	for _, dbf := range list {
		ext := path.Ext(dbf.Path)

		if ext == "" {
			fmt.Println("Skipping", dbf.Path)
			continue
		}

		if ext == ".m3u" || ext == ".wav" {
			fmt.Println("Skipping", dbf.Path)
			return nil
		}

		m := mime.TypeByExtension(ext)

		if !strings.HasPrefix(m, "audio/") {
			fmt.Println("Skipping", dbf.Path)
			continue
		}

		f, err := os.Open(dbf.Path)

		if err != nil {
			return err
		}

		defer f.Close()

		t, err := tag.ReadFrom(f)

		if err != nil {
			// return fmt.Errorf("read tags from %s: %w", dbf.Path, err)
			fmt.Printf("Error reading tags from %s: %s\n", dbf.Path, err)
			continue
		}

		// fmt.Println(t.Raw())

		// var file_tags []models.AudioTag
		for name, value := range t.Raw() {
			fmt.Printf("%s: %t %s\n", name, value, value)
		}
		// for _, tag := range t.Raw() {
		// 	file_tags = append(file_tags, models.AudioTag{
		// 		FileID:    dbf.ID,
		// 		Name:      tag.Name,
		// 		Value:     tag.Value,
		// 		CreatedAt: time.Now(),
		// 	})
		// }
	}

	return nil
}
