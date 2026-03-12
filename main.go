package main

import (
"fmt"
"os"
"os/exec"

"github.com/harshpatel5940/nestop/internal"
)

// Version is set at build time via ldflags
var (
	Version   = "dev"
	BuildTime = "unknown"
)

func main() {
if len(os.Args) > 1 {
	switch os.Args[1] {
	case "--version", "-v":
		fmt.Printf("nestop %s (built %s)\n", Version, BuildTime)
		return
	case "--help", "-h":
		fmt.Println("nestop — Interactive NestJS project scaffolder")
		fmt.Println()
		fmt.Println("Usage: nestop [flags]")
		fmt.Println()
		fmt.Println("Flags:")
		fmt.Println("  --version, -v   Print version")
		fmt.Println("  --help,    -h   Show this help")
		return
	}
}

fmt.Println("🪺 nestop — NestJS Starter Scaffolder")
fmt.Println()

config, err := internal.RunForm()
if err != nil {
fmt.Fprintf(os.Stderr, "Error: %v\n", err)
os.Exit(1)
}

fmt.Println()
fmt.Printf("📦 Generating project: %s\n", config.ProjectName)

if err := internal.Generate(config); err != nil {
fmt.Fprintf(os.Stderr, "❌ Generation failed: %v\n", err)
os.Exit(1)
}

fmt.Println("✅ Project files generated!")

if config.InitGit {
fmt.Println("⚡ Initializing git repository...")
gitInit := exec.Command("git", "init", config.ProjectName)
gitInit.Stdout = os.Stdout
gitInit.Stderr = os.Stderr
if err := gitInit.Run(); err != nil {
fmt.Fprintf(os.Stderr, "⚠️  git init failed: %v\n", err)
} else {
// Stage all files
gitAdd := exec.Command("git", "-C", config.ProjectName, "add", "-A")
gitAdd.Run()
gitCommit := exec.Command("git", "-C", config.ProjectName, "commit", "-m", "chore: initial scaffold")
gitCommit.Stdout = os.Stdout
gitCommit.Stderr = os.Stderr
gitCommit.Run()
}
}

if config.InstallDeps {
fmt.Printf("📥 Installing dependencies with %s...\n", config.PackageManager)
pm := string(config.PackageManager)
install := exec.Command(pm, "install")
install.Dir = config.ProjectName
install.Stdout = os.Stdout
install.Stderr = os.Stderr
if err := install.Run(); err != nil {
fmt.Fprintf(os.Stderr, "⚠️  %s install failed: %v\n", pm, err)
fmt.Printf("Run manually: cd %s && %s install\n", config.ProjectName, pm)
} else {
fmt.Println("✅ Dependencies installed!")
}
}

fmt.Println()
fmt.Println("🎉 Done! Next steps:")
fmt.Printf("  cd %s\n", config.ProjectName)
if !config.InstallDeps {
if config.UsesBun() {
fmt.Println("  bun install")
} else {
fmt.Println("  pnpm install")
}
}
if config.UsesBun() {
fmt.Println("  bun run start:dev")
} else {
fmt.Println("  pnpm start:dev")
}
fmt.Println()
fmt.Println("📖 Swagger docs: http://localhost:3000/api/swagger")
}
