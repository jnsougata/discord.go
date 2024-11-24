# discord.go
Discord HTTP Interaction API Wrapper for Serverless Applications

## Work In Progress

### Quick Start
```go
package main

import (
	"fmt"
	"log"
	"math/rand"

	"github.com/gin-gonic/gin"
)

func main() {
	app := App(
		&AppState{
			Port:          8080,
			Path:          "...",
			DiscordToken:  "...",
			PublicKey:     "...",
			ApplicationId: "...",
			ReleaseMode:   gin.ReleaseMode,
		})
	app.Preload(deleteButton)
	app.AddCommands(echo)
	if os.Getenv("ENV") == "DEVEL" {
		app.Sync()
	}
	if err := app.Run(); err != nil {
		log.Fatal(fmt.Errorf("failed to run app: %w", err))
	}
}

var deleteButton = Button(ButtonConfig{
	Label:    "Delete",
	Style:    ButtonStyleDanger,
	CustomId: "delete",
	Handler: func(interaction *Interaction) {
		interaction.Message.Delete(interaction.Client)
	},
})

var echo = SlashCommand(
	SlashCommandConfig{
		Name:        "echo",
		Description: "echoes your input",
		Options: Options(StringOption(OptionConfigString{
			Name:        "input",
			Description: "The input to echo",
			Required:    true,
			MinLength:   3,
		})),
		Handler: func(interaction *Interaction) {
			var options struct {
				Input string `json:"input"`
			}
			interaction.Bind(&options)
			interaction.Response(MessageOptions{
				Content:    options.Input,
				Components: Componenets(Row(deleteButton)),
			})
		},
	},
)
```
