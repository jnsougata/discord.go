package main

import (
	"fmt"
	"io"

	"github.com/gin-gonic/gin"
)

type AppState struct {
	Port          int
	Path          string
	DiscordToken  string
	PublicKey     string
	ApplicationId string
	ReleaseMode   string
	commands      []ApplicationCommand
	handlerMap    map[string]func(interaction *Interaction)
}

type app struct {
	AppState
	Engine *gin.Engine
	http   *HttpClient
}

func (a *app) Run() error {
	a.Engine.POST(a.Path, func(c *gin.Context) { handler(c, a) })
	var PORT = ":8080"
	if a.Port != 0 {
		PORT = fmt.Sprintf(":%d", a.Port)
	}
	return a.Engine.Run(PORT)
}

func (a *app) Sync() {
	resp, err := a.http.Sync(a.commands)
	if err != nil {
		panic(err)
	}
	body, _ := io.ReadAll(resp.Body)
	resp.Body.Close()
	fmt.Println(resp.StatusCode)
	fmt.Println(string(body))
}

func (a *app) AddCommands(commands ...ApplicationCommand) {
	for _, command := range commands {
		a.handlerMap[fmt.Sprintf("%s:%d", command.Name, command.Type)] = command.Handler
	}
	a.commands = append(a.commands, commands...)
}

func App(state *AppState) *app {
	gin.SetMode(state.ReleaseMode)
	state.handlerMap = map[string]func(interaction *Interaction){}
	return &app{*state, gin.Default(), NewHttpClient(state)}
}
