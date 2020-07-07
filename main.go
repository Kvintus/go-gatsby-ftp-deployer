package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os/exec"

	"github.com/gofiber/fiber"
)

type Ftp struct {
	Host	string `json:"host"`
	Username    string `json:"username"`
	Password string `json:"password"`
}

type Project struct {
	PathToProject string `json:"pathToProject"`
	Repo string `json:"repo"`
	Ftp  Ftp    `json:"ftp"`
}

func gatsbyClean(path string) error {
	fmt.Println("The cleaning has started")
	cleanCommand := exec.Command("gatsby", "clean")
	cleanCommand.Dir = path
	defer fmt.Println("Cleaning finished")
	return cleanCommand.Run();
}

func gatsbyBuild(path string) string {
	fmt.Println("The build process has started...")
	buildCommand := exec.Command("yarn", "build")
	buildCommand.Dir = path
		res, _ := buildCommand.Output()
	resString := string(res)
	fmt.Println(resString);
	fmt.Println("The project has been successfully built!")
	return resString
}

func upload(ftp Ftp, path string) error {
	fmt.Println("Starting an upload to a server...")
	scpPath := fmt.Sprintf("%s@%s:/%s/web", ftp.Username, ftp.Host, ftp.Host)
	uploadCommand := exec.Command("sshpass", "-p", ftp.Password, "scp", "-rp", "public/.", scpPath)
	uploadCommand.Dir = path
	resp, err :=  uploadCommand.Output()
	if (err != nil) {
		fmt.Println(err)
	}
	fmt.Println(string(resp))
	fmt.Println("Upload finished!")
	return err;
}

func main() {
	file, _ := ioutil.ReadFile("setup.json")
	var projects map[string]Project
	json.Unmarshal([]byte(file), &projects)
	app := fiber.New()

	app.Get("/:project", func(c *fiber.Ctx) {
		project := projects[c.Params("project")]

		gatsbyClean(project.PathToProject)
		gatsbyBuild(project.PathToProject)
		upload(project.Ftp, project.PathToProject)
		
		c.Status(200).Send("ok")
	})
	app.Listen(3000)
}
