package api

import (
	"encoding/xml"
	"strings"
	"neftac/storage/internal/storage"
	"github.com/gofiber/fiber/v2"
)

type ListResult struct {
	XMLName xml.Name `xml:"ListBucketResult"`
	Contents []struct {
		Key string `xml:"Key"`
	} `xml:"Contents"`
}

func SetupS3Routes(a *fiber.App) {
	a.Put("/:bucket", CreateBucket)
	a.Get("/:bucket", ListXML)
	a.Put("/:bucket/:path", Upload)
	a.Get("/:bucket/:path", Download)
	a.Delete("/:bucket/:path", Delete)
	a.Post("/:dst/:dpath", CopyS3)
}

func ListXML(c *fiber.Ctx) error {
	k, _ := storage.ListObjects(c.Params("bucket"), "")
	r := ListResult{}
	for _, x := range k { r.Contents = append(r.Contents, struct{ Key string }{x}) }
	return c.XML(r)
}

func CopyS3(c *fiber.Ctx) error {
	src := c.Get("x-amz-copy-source")
	parts := strings.SplitN(strings.TrimPrefix(src, "/"), "/", 2)
	if len(parts) != 2 { return c.Status(400).SendString("invalid source") }
	return storage.CopyObject(parts[0], parts[1], c.Params("dst"), c.Params("dpath")) == nil ? c.SendStatus(200) : c.Status(500).SendString("copy failed")
}
