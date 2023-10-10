package main

import (
	"fmt"
	"log"
	"os"
	"unicode"

	"github.com/KEINOS/go-pallet/pallet"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
	"github.com/google/uuid"
	"github.com/vincent-petithory/dataurl"
)

func Handler() *fiber.App {
	app := fiber.New(fiber.Config{
		Views: html.New("./views", ".html"),
	})

	app.Get("/", func(c *fiber.Ctx) error {
		return c.Render("index", fiber.Map{})
	})

	app.Get("/api/oekaki", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).JSON(&fiber.Map{
			"image": LatestOekaki.Image,
		})
	})

	app.Get("/api/oekaki/all", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).JSON(&fiber.Map{
			"success": true,
			"result":  DB.Find(&OekakiStore{}),
		})
	})

	app.Post("/api/oekaki", func(c *fiber.Ctx) error {
		// Check if the input is valid as expected JSON format
		r := new(OekakiRequest)
		if err := c.BodyParser(r); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
				"success": false,
				"error":   "Bad JSON format",
			})
		}

		// Check if next answer is not empty
		if r.NextAnswer == "" {
			return c.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
				"success": false,
				"error":   "Next answer is not given",
			})
		}

		// Check if the answer is entered (skip if this is first one)
		if r.Answer == "" && LatestOekaki.Answer != "" {
			return c.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
				"success": false,
				"error":   "Your answer is not given",
			})
		}

		// Check if the answer & next answer is all hiragana
		for _, r := range r.Answer {
			if !unicode.In(r, unicode.Hiragana) {
				return c.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
					"success": false,
					"error":   "Your answer contains non hiragana character(s)",
				})
			}
		}
		for _, r := range r.NextAnswer {
			if !unicode.In(r, unicode.Hiragana) {
				return c.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
					"success": false,
					"error":   "Your next answer contains non hiragana character(s)",
				})
			}
		}

		// Check if the input is valid as data URI
		image_data, err := dataurl.DecodeString(r.Image)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
				"success": false,
				"error":   "Bad data URI",
			})
		}

		// Check input image size for prevention from spamming
		image_file_name := fmt.Sprintf("%v.png", uuid.New().String())
		image_file, err := os.Create(image_file_name)
		if err != nil {
			log.Println(err)
		}
		defer image_file.Close()
		image_file.Write(image_data.Data)

		// Save the size of this image
		fi, _ := image_file.Stat()
		image_size := fi.Size()

		// Count the color used in this image
		pal, err := pallet.Load(image_file_name)
		if err != nil {
			log.Println(err)
		}
		pixinfo := pallet.ByOccurrence(pal)

		// Reject if only one color is used
		if len(pixinfo) < 2 {
			return c.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
				"success": false,
				"error":   "Nothing is drew",
			})
		}

		colors := pixinfo[0:2]

		// Delete temporary file
		err = os.Remove(image_file_name)
		if err != nil {
			log.Printf("Failed to delete temporary file %v", image_file_name)
		}

		// Reject the image with more than 2 colors
		color_ok := true
		for _, v := range colors {
			if !(v.R == 255 && v.G == 255 && v.B == 255) {
				if !(v.R == 0 && v.G == 0 && v.B == 0) {
					color_ok = false
					break
				}
			}
		}
		if !(color_ok) {
			return c.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
				"success": false,
				"error":   "Some colors are in use other than pure black and white (in RGB)",
			})
		}

		// Reject the image smaller than 20,000 bytes
		// (20,000 bytes is nearly equivalent to all white for 500x500 png image)
		if image_size < 20000 {
			return c.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
				"success": false,
				"error":   "Too small input image",
			})
		}

		// Handling in-memory oekaki result
		new := OekakiStore{
			Id:         uuid.NewString(),
			Answer:     LatestOekaki.Answer,
			UserAnswer: r.Answer,
			Image:      r.Image,
		}
		DB.Create(&new)
		prev := LatestOekaki.Answer
		LatestOekaki.Answer = r.NextAnswer
		LatestOekaki.Image = r.Image

		return c.Status(fiber.StatusOK).JSON(&fiber.Map{
			"success": true,
			"correct": prev == r.Answer,
		})
	})

	return app
}
