package rc

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/chnsz/resource-cleanup/helper"
)

type ResourceQuery interface {
	QueryIds() []string
	GetName() string
}

func Clean(name string, ids []string) {
	workDir := helper.GetTmpDir()
	defer func() {
		_ = os.RemoveAll(workDir)
	}()

	if err := buildConfigFile(workDir); err != nil {
		log.Printf("[ERROR] failed to write main.tf: %s", err)
		return
	}

	// terraform init
	terraform(workDir, "init")
	if err := buildMainFile(workDir, name, ids); err != nil {
		log.Printf("[ERROR] failed to write main.tf: %s", err)
		return
	}
	for i, id := range ids {
		options := []string{
			"import",
			fmt.Sprintf("%s.test_%v", name, i),
			id,
		}
		terraform(workDir, options...)
	}
	_ = os.Remove(getMainTf(workDir))

	terraform(workDir, "apply", "-auto-approve")
}

func terraform(workDir string, options ...string) {
	cmd := exec.Command("terraform", options...)
	log.Printf("[INFO] run terraform %s", strings.Join(options, " "))
	cmd.Dir = workDir
	helper.Run(cmd, time.Minute*1)
}

func buildConfigFile(workDir string) error {
	authUrl := os.Getenv(helper.HwAuthUrl)
	if authUrl == "" {
		authUrl = "https://iam.myhuaweicloud.com:443/v3"
	}
	lines := make([]string, 0)
	lines = append(lines, fmt.Sprintf(`terraform {`))
	lines = append(lines, fmt.Sprintf(`  required_providers {`))
	lines = append(lines, fmt.Sprintf(`    huaweicloud = {`))
	lines = append(lines, fmt.Sprintf(`      source  = "huaweicloud/huaweicloud"`))
	lines = append(lines, fmt.Sprintf(`      version = ">= 1.63.0"`))
	lines = append(lines, fmt.Sprintf(`    }`))
	lines = append(lines, fmt.Sprintf(`  }`))
	lines = append(lines, fmt.Sprintf(`}`))
	lines = append(lines, fmt.Sprintf(``))
	lines = append(lines, fmt.Sprintf(`provider "huaweicloud" {`))
	lines = append(lines, fmt.Sprintf(`  auth_url   = "%s"`, authUrl))
	lines = append(lines, fmt.Sprintf(`  region     = "%s"`, os.Getenv(helper.HwRegionName)))
	lines = append(lines, fmt.Sprintf(`  access_key = "%s"`, os.Getenv(helper.HwAccessKey)))
	lines = append(lines, fmt.Sprintf(`  secret_key = "%s"`, os.Getenv(helper.HwSecretKey)))
	lines = append(lines, fmt.Sprintf(`  project_id = "%s"`, os.Getenv(helper.HwProjectId)))
	lines = append(lines, fmt.Sprintf(`}`))
	return helper.WriteToFile(fmt.Sprintf("%s%sconfig.tf", workDir, string(os.PathSeparator)), strings.Join(lines, "\n"))
}

func buildMainFile(workDir, name string, ids []string) error {
	lines := make([]string, 0)

	for i := range ids {
		lines = append(lines, fmt.Sprintf(`resource "%s" "test_%v" {`, name, i))
		lines = append(lines, `}`)
	}

	return helper.WriteToFile(getMainTf(workDir), strings.Join(lines, "\n"))
}

func getMainTf(workDir string) string {
	return fmt.Sprintf("%s%smain.tf", workDir, string(os.PathSeparator))
}
