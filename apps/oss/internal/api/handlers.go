package api

import (
	"encoding/json"
	"neftac/storage/internal/db"
	"neftac/storage/internal/storage"
	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(a *fiber.App) {
	v1 := a.Group("/v1").Use(Auth)
	v1.Put("/buckets/:bucket", CreateBucket)
	v1.Get("/buckets", ListBuckets)
	v1.Put("/buckets/:bucket/objects/:path", Policy("write"), Upload)
	v1.Get("/buckets/:bucket/objects/:path", Policy("read"), Download)
	v1.Delete("/buckets/:bucket/objects/:path", Policy("write"), Delete)
	v1.Get("/buckets/:bucket/objects", List)
	v1.Post("/sign-upload", SignUpload)
	v1.Post("/copy", Copy)
	v1.Post("/rename", Rename)
	v1.Post("/move", Move)
}

func CreateBucket(c *fiber.Ctx) error { return db.DB.Create(&db.Bucket{Name: c.Params("bucket")}).Error == nil ? c.JSON(fiber.Map{"ok": true}) : c.Status(500).JSON(fiber.Map{"error": "failed"}) }
func ListBuckets(c *fiber.Ctx) error { var b []db.Bucket; db.DB.Find(&b); n := []string{}; for _, x := range b { n = append(n, x.Name) }; return c.JSON(n) }
func Upload(c *fiber.Ctx) error { e, _, err := storage.WriteObject(c.Params("bucket"), c.Params("path"), c.Context().RequestBodyStream()); return err == nil ? c.JSON(fiber.Map{"etag": e}) : c.Status(500).JSON(fiber.Map{"error": err.Error()}) }
func Download(c *fiber.Ctx) error { f, err := storage.ReadObject(c.Params("bucket"), c.Params("path")); if err != nil { return c.Status(404).JSON(fiber.Map{"error": "not found"}) }; defer f.Close(); return c.SendStream(f) }
func Delete(c *fiber.Ctx) error { return storage.DeleteObject(c.Params("bucket"), c.Params("path")) == nil ? c.SendStatus(204) : c.Status(500).JSON(fiber.Map{"error": "failed"}) }
func List(c *fiber.Ctx) error { k, _ := storage.ListObjects(c.Params("bucket"), c.Query("prefix")); return c.JSON(k) }

func SignUpload(c *fiber.Ctx) error {
	var r struct{ Bucket, Key string; Expires int }
	json.NewDecoder(c.Context().RequestBodyStream()).Decode(&r)
	url := "https://s3.airsoko.com/" + r.Bucket + "/" + r.Key + "?X-Neftac-Signed=1"
	return c.JSON(fiber.Map{"url": url, "expires_in": r.Expires})
}

func Copy(c *fiber.Ctx) error {
	var r struct{ SrcB, SrcK, DstB, DstK string }
	json.NewDecoder(c.Context().RequestBodyStream()).Decode(&r)
	return storage.CopyObject(r.SrcB, r.SrcK, r.DstB, r.DstK) == nil ? c.JSON(fiber.Map{"ok": true}) : c.Status(500).JSON(fiber.Map{"error": "copy failed"})
}

func Rename(c *fiber.Ctx) error { Copy(c); Delete(c); return c.SendStatus(200) }
func Move(c *fiber.Ctx) error { Copy(c); Delete(c); return c.SendStatus(200) }
