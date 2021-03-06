package models

import (
	"fmt"
	"log"
	"strings"

	"github.com/jinzhu/gorm"
	"github.com/qor/l10n"
	"github.com/qor/media_library"
	"github.com/qor/publish"
	"github.com/qor/qor-example/db"
	"github.com/qor/slug"
	"github.com/qor/sorting"
	"github.com/qor/validations"
)

type Product struct {
	gorm.Model
	l10n.Locale
	publish.Status
	sorting.SortingDESC

	Name            string
	NameWithSlug    slug.Slug        `l10n:"sync"`
	Code            string           `l10n:"sync"`
	CategoryID      uint             `l10n:"sync"`
	Category        Category         `l10n:"sync"`
	Collections     []Collection     `l10n:"sync" gorm:"many2many:product_collections"`
	MadeCountry     string           `l10n:"sync"`
	Price           float32          `l10n:"sync"`
	Description     string           `sql:"size:2000"`
	ColorVariations []ColorVariation `l10n:"sync"`
	Enabled         bool
}

func (product Product) DefaultPath() string {
	defaultPath := "/"
	if len(product.ColorVariations) > 0 {
		defaultPath = fmt.Sprintf("/products/%s_%s", product.Code, product.ColorVariations[0].ColorCode)
	}
	return defaultPath
}

func (product Product) MainImageUrl() string {
	return product.ColorVariations[0].MainImageUrl()
}

func (product Product) Validate(db *gorm.DB) {
	if strings.TrimSpace(product.Name) == "" {
		db.AddError(validations.NewError(product, "Name", "Name can not be empty"))
	}

	if strings.TrimSpace(product.Code) == "" {
		db.AddError(validations.NewError(product, "Code", "Code can not be empty"))
	}
}

type ColorVariation struct {
	gorm.Model
	ProductID      uint
	Product        Product
	ColorID        uint
	Color          Color
	ColorCode      string
	Images         []ColorVariationImage
	SizeVariations []SizeVariation
}

type ColorVariationImage struct {
	gorm.Model
	ColorVariationID uint
	Image            ColorVariationImageStorage `sql:"type:varchar(4096)"`
}

type ColorVariationImageStorage struct{ media_library.FileSystem }

func (colorVariation ColorVariation) MainImageUrl() string {
	imageURL := "/images/default_product.png"
	if len(colorVariation.Images) > 0 {
		imageURL = colorVariation.Images[0].Image.URL()
	}
	return imageURL
}

func (ColorVariationImageStorage) GetSizes() map[string]media_library.Size {
	return map[string]media_library.Size{
		"small":  {Width: 320, Height: 320},
		"middle": {Width: 640, Height: 640},
		"big":    {Width: 1280, Height: 1280},
	}
}

type SizeVariation struct {
	gorm.Model
	ColorVariationID  uint
	ColorVariation    ColorVariation
	SizeID            uint
	Size              Size
	AvailableQuantity uint
}

func SizeVariations() []SizeVariation {
	sizeVariations := []SizeVariation{}
	if err := db.DB.Debug().Preload("ColorVariation.Color").Preload("ColorVariation.Product").Preload("Size").Find(&sizeVariations).Error; err != nil {
		log.Fatalf("query sizeVariations (%v) failure, got err %v", sizeVariations, err)
		return sizeVariations
	}
	return sizeVariations
}

func (sizeVariation SizeVariation) Stringify() string {
	if colorVariation := sizeVariation.ColorVariation; colorVariation.ID != 0 {
		product := colorVariation.Product
		return fmt.Sprintf("%s (%s-%s-%s)", product.Name, product.Code, colorVariation.Color.Code, sizeVariation.Size.Code)
	}
	return fmt.Sprint(sizeVariation.ID)
}
