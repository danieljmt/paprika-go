package paprika

import (
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/google/uuid"
)

type Recipe struct {
	UID             string   `json:"uid"`
	Created         string   `json:"created"`
	Hash            string   `json:"hash"`
	Name            string   `json:"name"`
	Description     string   `json:"description"`
	Ingredients     string   `json:"ingredients"`
	Directions      string   `json:"directions"`
	Notes           string   `json:"notes"`
	NutritionalInfo string   `json:"nutritional_info"`
	PrepTime        string   `json:"prep_time"`
	CookTime        string   `json:"cook_time"`
	TotalTime       string   `json:"total_time"`
	Difficulty      string   `json:"difficulty"`
	Servings        string   `json:"servings"`
	Rating          int      `json:"rating"`
	Source          string   `json:"source"`
	SourceURL       string   `json:"source_url"`
	Photo           string   `json:"photo"`
	PhotoLarge      string   `json:"photo_large"`
	PhotoHash       string   `json:"photo_hash"`
	ImageURL        string   `json:"image_url"`
	Categories      []string `json:"categories"`
	PhotoData       string   `json:"photo_data"`
}

// NewRecipe creates a valid paprika recipe
func NewRecipe(name, ingredients, directions string) *Recipe {
	uid := uuid.NewMD5(uuid.Nil, []byte(name)).String()
	r := &Recipe{
		UID:         uid,
		Created:     time.Now().Format("2006-01-02 15:04:05"), // 2020-08-31 20:51:49
		Name:        name,
		Ingredients: ingredients,
		Directions:  directions,
	}
	return r
}

// SetImageURL sets the recipe image from an image URL
func (r *Recipe) SetImageURL(url string) error {
	if r.UID == "" {
		return errors.New("recipe must have a UID")
	}
	resp, err := http.Get(url) // #nosec G107 needs to be using a variable
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	dat, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	imParts := strings.Split(url, ".")
	r.Photo = fmt.Sprintf("%s.%s", r.UID, imParts[len(imParts)-1])
	r.ImageURL = url
	r.PhotoData = base64.StdEncoding.EncodeToString(dat)
	return nil
}

// SetImage sets the recipe image from a local image file
func (r *Recipe) SetImage(path string) error {
	if r.UID == "" {
		return errors.New("recipe must have a UID")
	}
	dat, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	imParts := strings.Split(path, ".")
	r.Photo = fmt.Sprintf("%s.%s", r.UID, imParts[len(imParts)-1])
	r.PhotoData = base64.StdEncoding.EncodeToString(dat)
	return nil
}
