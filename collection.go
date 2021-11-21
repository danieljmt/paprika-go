package paprika

import (
	"archive/zip"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
)

type Collection struct {
	recipes []Recipe
}

func NewCollection() *Collection {
	return &Collection{}
}

func LoadCollection(file string) (*Collection, error) {
	c := &Collection{}
	if err := c.LoadRecipes(file); err != nil {
		return nil, err
	}
	return c, nil
}

func (c *Collection) Add(r ...Recipe) {
	c.recipes = append(c.recipes, r...)
}

func (c *Collection) Remove(rs ...Recipe) {
	for i := 0; i < len(c.recipes); i++ {
		if contains(rs, c.recipes[i]) {
			c.recipes = append(c.recipes[:i], c.recipes[i+1:]...)
			i--
		}
	}
}

func contains(rs []Recipe, r Recipe) bool {
	for _, rsr := range rs {
		if rsr.UID == r.UID {
			return true
		}
	}
	return false
}

func (c *Collection) Get(rm ...RecipeMatcher) {
}

func (c *Collection) LoadRecipes(file ...string) error {
	for _, f := range file {
		if err := c.loadRecipes(f); err != nil {
			return err
		}
	}
	return nil
}

func (c *Collection) loadRecipes(file string) error {
	zr, err := zip.OpenReader(file)
	if err != nil {
		return fmt.Errorf("error reading collection file `%s`: %w", file, err)
	}
	for _, f := range zr.File {
		if !f.FileInfo().IsDir() {
			rc, err := f.Open()
			if err != nil {
				return fmt.Errorf("error opening recipe file `%s`: %w", f.Name, err)
			}
			defer rc.Close()
			r, err := decompress(rc)
			if err != nil {
				return fmt.Errorf("error decompressing recipe `%s`: %w", f.Name, err)
			}
			c.recipes = append(c.recipes, *r)
		}
	}
	zr.Close()
	return nil
}

func (c *Collection) Export(dest string) error {
	archive, err := os.Create(filepath.Join(dest, "Recipes.paprikarecipes"))
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := archive.Close(); err != nil {
			log.Println("failed to close archive writer: %w", err)
		}
	}()
	zw := zip.NewWriter(archive)
	defer func() {
		if err := zw.Close(); err != nil {
			log.Println("failed to close zip writer: %w", err)
		}
	}()

	for _, r := range c.recipes {
		data, err := r.compress()
		if err != nil {
			return fmt.Errorf("failed compressing recipe: %w", err)
		}

		illegalChars := regexp.MustCompile("[/\\?%*:|\"<>]")
		fileName := illegalChars.ReplaceAllString(r.Name, "_")

		iow, err := zw.Create(filepath.Clean(fmt.Sprintf("%s.paprikarecipe", fileName)))
		if err != nil {
			return fmt.Errorf("failed creating recipe in archive: %w", err)
		}
		if _, err := iow.Write(data); err != nil {
			return fmt.Errorf("failed writing recipe to archive: %w", err)
		}
	}

	return nil
}
