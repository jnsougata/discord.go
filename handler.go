package main

import (
	"crypto/ed25519"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"

	"github.com/gin-gonic/gin"
)

func handler(ctx *gin.Context, client *Client) {
	signatureHex := ctx.GetHeader("X-Signature-Ed25519")
	timestamp := ctx.GetHeader("X-Signature-Timestamp")
	if signatureHex == "" || timestamp == "" {
		ctx.JSON(401, gin.H{"message": "Unauthorized"})
		return
	}
	body, _ := io.ReadAll(ctx.Request.Body)
	message := []byte(timestamp)
	message = append(message, body...)
	signatureBytes, _ := hex.DecodeString(signatureHex)
	publicKeyBytes, _ := hex.DecodeString(client.PublicKey)
	ok := ed25519.Verify(publicKeyBytes, message, signatureBytes)
	if !ok {
		ctx.JSON(401, gin.H{"message": "Unauthorized"})
		return
	}
	var interaction Interaction
	interaction.Client = client
	interaction.Context = ctx
	_ = json.Unmarshal(body, &interaction)
	switch interaction.Type {
	case InteractionTypePing:
		ctx.JSON(200, gin.H{"type": 1})
	case InteractionTypeApplicationCommand:
		key := fmt.Sprintf("%s:%d", interaction.Data.Name, interaction.Data.Type)
		globalHandlerMap[key](&interaction)
	case InteractionTypeMessageComponent:
		key := fmt.Sprintf("%s:%d", interaction.Data.CustomId, interaction.Data.ComponentType)
		globalHandlerMap[key](&interaction)
	case InteractionTypeModalSubmit:
		globalHandlerMap[interaction.Data.CustomId](&interaction)
	default:
		// TODO: handle unknown interaction type
	}
}
