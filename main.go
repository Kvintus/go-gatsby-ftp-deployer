package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os/exec"

	"github.com/gofiber/fiber"
)

type Ftp struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type Project struct {
	Repo string `json:"repo"`
	Ftp  Ftp    `json:"ftp"`
}

func main() {
	file, _ := ioutil.ReadFile("setup.json")
	var projects map[string]Project
	json.Unmarshal([]byte(file), &projects)
	app := fiber.New()

	app.Get("/:project", func(c *fiber.Ctx) {
		buildCommand := exec.Command("yarn", "build")
		buildCommand.Dir = "/home/kvintus/Documents/programming/wave/ezermann/client"
		res, _ := buildCommand.Output()
		fmt.Println(string(res))
		c.Send(projects[c.Params("project")].Repo)
	})
	app.Listen(3000)

}
