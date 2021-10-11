package main

import (
  "crypto/ed25519"
  "encoding/hex"
  "fmt"
  "os"

  "github.com/bwmarrin/discordgo"
  "github.com/gofiber/fiber/v2"
)

var discordPublicKey ed25519.PublicKey

func main() {

  publicKey, err := hex.DecodeString(os.Getenv("DISCORD_PUBLIC_KEY"))
  if err != nil {
    panic(err)
  }
  discordPublicKey = ed25519.PublicKey(publicKey)

  app := fiber.New()

  app.Get("/healthy", func (c *fiber.Ctx) error {
    c.JSON(map[string]string {
      "healthy": "yus!",
    })

    return nil
  })

  interactionGroup := app.Group("/interactions", interactionMiddleware)
  interactionGroup.Post("/", interactionHandler)

  app.Listen(":3000")
}

func interactionMiddleware(c *fiber.Ctx) error {
  if !verifyInteraction(c, discordPublicKey) {
    c.Status(401).JSON(map[string]string {
      "Begone": "foul demon!",
    })
    return nil
  }

  return c.Next()
}

func interactionHandler(c *fiber.Ctx) error {

  interaction := new(discordgo.Interaction)
  if err := c.BodyParser(interaction); err != nil {
    return err
  }
  fmt.Println(interaction)
  c.JSON(map[string]interface{} {
    "type": discordgo.InteractionResponsePong,
  })
  return nil
}

func verifyInteraction(c *fiber.Ctx, key ed25519.PublicKey) bool {

  signatureStr := c.Get("X-Signature-Ed25519", "")
  if signatureStr == "" {
    return false
  }
  signature, err := hex.DecodeString(signatureStr)
  if err != nil {
    return false
  }
  if len(signature) != ed25519.SignatureSize {
    return false
  }

  timestamp := c.Get("X-Signature-Timestamp", "")
  if timestamp == "" {
    return false
  }

  return ed25519.Verify(key, append([]byte(timestamp), c.Body()...), signature)
}