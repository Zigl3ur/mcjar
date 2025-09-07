package backup

import (
	"archive/zip"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/Zigl3ur/mcli/internal/utils"
	"github.com/Zigl3ur/mcli/internal/utils/loader"
	"github.com/spf13/cobra"
)

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "backup",
		Short:   "Backup a server folders to a zip archive",
		Long:    "Backup a server folders to a zip archive",
		PreRunE: validate,
		RunE:    execute,
	}

	cmd.Flags().StringArrayP("folders", "f", []string{}, "the folder(s) to backup")
	cmd.Flags().StringP("destination", "d", ".", "the folder destination where the backup will be stored")

	//nolint:errcheck
	cmd.MarkFlagRequired("folders")
	//nolint:errcheck
	cmd.MarkFlagDirname("folders")
	//nolint:errcheck
	cmd.MarkFlagDirname("destination")

	cmd.Flags().SortFlags = false

	return cmd
}

func validate(cmd *cobra.Command, args []string) error {
	folders, _ := cmd.Flags().GetStringArray("folders")

	for _, folder := range folders {
		if _, err := os.Stat(folder); err != nil {
			return err
		}
	}

	if !cmd.Flag("destination").Changed {
		wd, err := os.Getwd()

		if err != nil {
			return errors.New("failed to get working dir")
		}

		//nolint:errcheck
		cmd.Flags().Set("destination", wd)
	}

	return nil
}

func execute(cmd *cobra.Command, args []string) error {
	folders, _ := cmd.Flags().GetStringArray("folders")
	destination, _ := cmd.Flags().GetString("destination")

	loader.Start("Creating Backup")

	if err := utils.CheckDir(destination); err != nil {
		loader.Stop()
		return err
	}
	now := time.Now().Format("2006-01-02_15-04-05")
	archiveName := fmt.Sprintf("Backup-%s.zip", now)
	archivePath := filepath.Join(destination, archiveName)
	archive, err := os.Create(archivePath)
	if err != nil {
		loader.Stop()
		return err
	}

	//nolint:errcheck
	defer archive.Close()

	writer := zip.NewWriter(archive)

	for _, elt := range folders {
		basePath := filepath.Dir(elt)
		if err := addToArchive(writer, basePath, elt); err != nil {
			loader.Stop()
			return err
		}
	}

	//nolint:errcheck
	writer.Close()
	loader.Stop()
	fmt.Printf("Successfully created backup at %s\n", archivePath)

	return nil
}

func addToArchive(w *zip.Writer, basePath, currentPath string) error {
	file, err := os.Stat(currentPath)
	if err != nil {
		return nil
	}

	relPath, err := filepath.Rel(basePath, currentPath)
	if err != nil {
		return err
	}

	loader.UpdateMessage(fmt.Sprintf("Adding %s to archive", filepath.Dir(relPath)))

	if file.IsDir() {
		if relPath != "." {
			dirPath := relPath + string(os.PathSeparator)
			_, err := w.Create(dirPath)
			if err != nil {
				return err
			}
		}

		files, err := os.ReadDir(currentPath)
		if err != nil {
			return err
		}

		for _, file := range files {
			if err = addToArchive(w, basePath, filepath.Join(currentPath, file.Name())); err != nil {
				return err
			}
		}
	} else {
		file, err := os.Open(currentPath)
		if err != nil {
			return err
		}
		//nolint:errcheck
		defer file.Close()

		wzip, err := w.Create(relPath)
		if err != nil {
			return err
		}

		if _, err := io.Copy(wzip, file); err != nil {
			return err
		}
	}

	return nil
}
