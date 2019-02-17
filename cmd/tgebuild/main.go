package main

import (
	"errors"
	"flag"
	"fmt"
	"go/build"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"runtime"

	decentcopy "github.com/hugocarreira/go-decent-copy"
	"github.com/otiai10/copy"
)

type Builder struct {
	target      string
	devMode     bool
	cwd         string
	packagePath string
	distPath    string
	programName string
	goPath      string
	tgeRootPath string
}

func determineGoVersion() error {
	gobin, err := exec.LookPath("go")
	if err != nil {
		return errors.New("go not found")
	}
	goVersionOut, err := exec.Command(gobin, "version").CombinedOutput()
	if err != nil {
		return fmt.Errorf("'go version' failed: %v, %s", err, goVersionOut)
	}
	var minor int
	if _, err := fmt.Sscanf(string(goVersionOut), "go version go1.%d", &minor); err != nil {
		// Ignore unknown versions; it's probably a devel version.
		return nil
	}
	if minor < 11 {
		return errors.New("Go 1.11 or newer is required")
	}
	return nil
}

func findTGERootPath() (string, error) {
	gopath := os.Getenv("GOPATH")
	if gopath == "" {
		gopath = build.Default.GOPATH
	}
	var tgeRootPath string
	err := filepath.Walk(gopath, func(p string, info os.FileInfo, err error) error {
		if !info.IsDir() && info.Name() == "tge.marker" {
			tgeRootPath = path.Dir(p)
		}
		return nil
	})

	if err != nil {
		return "", fmt.Errorf("Failed to find TGE root path: %s", err)
	}

	if tgeRootPath == "" {
		return "", fmt.Errorf("Failed to find TGE root path: not found in GOPATH")
	}

	return tgeRootPath, nil
}

func (b *Builder) init() error {
	err := determineGoVersion()
	if err != nil {
		return err
	}

	b.cwd, _ = os.Getwd()

	if !path.IsAbs(b.packagePath) {
		b.packagePath = path.Join(b.cwd, b.packagePath)
	}

	if _, err = os.Stat(path.Join(b.packagePath, "go.mod")); os.IsNotExist(err) {
		fmt.Println("WARNING:\n   > 'go.mod' not found in package path, tgebuild will not be able to retrieve your dependencies if needed.")
	}

	b.programName = path.Base(b.packagePath)

	if err = os.Chdir(b.packagePath); err != nil {
		return err
	}

	if b.tgeRootPath, err = findTGERootPath(); err != nil {
		return err
	}

	b.distPath = path.Join(b.packagePath, "dist", b.target)

	if !b.devMode {
		if err = os.RemoveAll(b.distPath); err != nil {
			return err
		}
	}

	if _, err = os.Stat(b.distPath); os.IsNotExist(err) {
		if err = os.MkdirAll(b.distPath, os.ModeDir|0777); err != nil {
			return err
		}
	}

	b.goPath = os.Getenv("GOPATH")
	if b.goPath == "" {
		b.goPath = build.Default.GOPATH
	}

	return nil
}

func (b *Builder) copyResources() error {
	resourcesInPath := path.Join(b.packagePath, b.target)
	boolFirstCopy := false
	var err error
	if _, err = os.Stat(resourcesInPath); os.IsNotExist(err) {
		boolFirstCopy = true
		if err = os.MkdirAll(resourcesInPath, os.ModeDir|0777); err != nil {
			return err
		}
		if err = copy.Copy(path.Join(b.tgeRootPath, b.target), resourcesInPath); err != nil {
			return err
		}
		fmt.Printf("NOTICE:\n   > './%s' folder has been added to your project and can be used to customize your build (see content for details)\n", b.target)
	}
	if boolFirstCopy || !b.devMode {
		if err = copy.Copy(resourcesInPath, b.distPath); err != nil {
			return err
		}
	}
	return nil
}

func (b *Builder) buildDesktop(packagePath string) error {
	b.target = runtime.GOOS
	b.packagePath = packagePath
	err := b.init()
	if err != nil {
		return err
	}

	if runtime.GOOS == "windows" {
		b.programName = fmt.Sprintf("%s.exe", b.programName)
	}

	cmd := exec.Command("go", "build", "-o", path.Join(b.distPath, b.programName))
	cmd.Env = append(os.Environ())

	if err := cmd.Run(); err != nil {
		return err
	}

	if err := b.copyResources(); err != nil {
		return err
	}

	return nil
}

func (b *Builder) buildBrowser(packagePath string) error {
	b.target = "browser"
	b.packagePath = packagePath
	err := b.init()
	if err != nil {
		return err
	}

	b.programName = "main.wasm"

	cmd := exec.Command("go", "build", "-o", path.Join(b.distPath, b.programName))
	cmd.Env = append(os.Environ(),
		"GOOS=js",
		"GOARCH=wasm",
	)
	if err := cmd.Run(); err != nil {
		return err
	}

	if err := b.copyResources(); err != nil {
		return err
	}

	return nil
}

func (b *Builder) buildAndroid(packagePath string) error {
	b.target = "android"
	b.packagePath = packagePath
	err := b.init()
	if err != nil {
		return err
	}

	gomobilebin, err := exec.LookPath("gomobile")
	if err != nil {
		gomobilebin = path.Join(b.goPath, "bin", "gomobile")
		if _, err = os.Stat(gomobilebin); os.IsNotExist(err) {
			fmt.Println("NOTICE:\n   > installing gomobile in your workspace...")
			cmd := exec.Command("go", "get", "golang.org/x/mobile/cmd/gomobile")
			cmd.Env = append(os.Environ())
			if err := cmd.Run(); err != nil {
				return err
			}
		}
	}

	if _, err = os.Stat(path.Join(b.goPath, "pkg", "gomobile", "ndk-toolchains")); os.IsNotExist(err) {
		androidNDKPath := os.Getenv("ANDROID_NDK")
		if androidNDKPath == "" {
			fmt.Println("ERROR:\n   > ANDROID_NDK is not set (should be $ANDROID_HOME/ndk-bundle), see https://developer.android.com/ndk/guides/.")
			return fmt.Errorf("cannot initialize gomobile")
		}

		fmt.Println("NOTICE:\n   > initializing gomobile...")
		cmd := exec.Command("gomobile", "init", "-ndk", androidNDKPath)
		cmd.Env = append(os.Environ())
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("cannot initialize gomobile: %s", err)
		}
	}

	if _, err = os.Stat(path.Join(b.packagePath, b.target, "AndroidManifest.xml")); os.IsNotExist(err) {
		if err = decentcopy.Copy(path.Join(b.tgeRootPath, b.target, "AndroidManifest.xml"), path.Join(b.packagePath, "AndroidManifest.xml")); err != nil {
			return err
		}
	} else {
		if err = decentcopy.Copy(path.Join(b.packagePath, b.target, "AndroidManifest.xml"), path.Join(b.packagePath, "AndroidManifest.xml")); err != nil {
			return err
		}
	}

	if b.devMode {
		b.programName = fmt.Sprintf("%s.apk", b.programName)

		cmd := exec.Command(gomobilebin, "build", "-target=android", "-o", path.Join(b.distPath, b.programName))
		cmd.Env = append(os.Environ())
		if err := cmd.Run(); err != nil {
			return err
		}
	} else {
		for _, t := range []string{"arm", "386", "amd64", "arm64"} {
			cmd := exec.Command(gomobilebin, "build", fmt.Sprintf("-target=android/%s", t), "-o", path.Join(b.distPath, fmt.Sprintf("%s-%s.apk", b.programName, t)))
			cmd.Env = append(os.Environ())
			if err := cmd.Run(); err != nil {
				return err
			}
		}

	}

	if err := b.copyResources(); err != nil {
		return err
	}

	os.Remove(path.Join(b.packagePath, "AndroidManifest.xml"))
	os.Remove(path.Join(b.distPath, "AndroidManifest.xml"))

	return nil
}

func (b *Builder) buildIOS(packagePath string, bundleID string) error {
	b.target = "ios"
	b.packagePath = packagePath
	err := b.init()
	if err != nil {
		return err
	}

	gomobilebin, err := exec.LookPath("gomobile")
	if err != nil {
		gomobilebin = path.Join(b.goPath, "bin", "gomobile")
		if _, err = os.Stat(gomobilebin); os.IsNotExist(err) {
			fmt.Println("NOTICE:\n   > installing gomobile in your workspace...")
			cmd := exec.Command("go", "get", "golang.org/x/mobile/cmd/gomobile")
			cmd.Env = append(os.Environ())
			if err := cmd.Run(); err != nil {
				return err
			}
		}
	}

	b.programName = fmt.Sprintf("%s.app", b.programName)

	cmd := exec.Command(gomobilebin, "build", "-target=ios", "-v", "-o", path.Join(b.distPath, b.programName))
	cmd.Env = append(os.Environ())
	if err := cmd.Run(); err != nil {
		return err
	}

	if err := b.copyResources(); err != nil {
		return err
	}

	return nil
}

func main() {
	targetFlag := flag.String("t", "desktop", "build target : desktop, android, ios, browser")
	devModeFlag := flag.Bool("d", false, "Dev mode, skip clean, assets copy & arch split (faster)")
	bundleIDFlag := flag.String("b", "", "IOS only, bundleId to use for app")
	flag.Parse()

	if len(flag.Args()) == 0 {
		fmt.Println(usage)
		return
	}

	var err error
	builder := &Builder{}
	builder.devMode = *devModeFlag
	switch *targetFlag {
	case "desktop":
		err = builder.buildDesktop(flag.Args()[0])
	case "browser":
		err = builder.buildBrowser(flag.Args()[0])
	case "android":
		err = builder.buildAndroid(flag.Args()[0])
	case "ios":
		if *bundleIDFlag == "" {
			fmt.Println("ERROR: Missing bunldeId for IOS")
			return
		}
		err = builder.buildIOS(flag.Args()[0], *bundleIDFlag)
	default:
		fmt.Printf("ERROR: Unsupported target '%s'\n", *targetFlag)
		flag.Usage()
		return
	}

	if err != nil {
		os.RemoveAll(builder.distPath)
		fmt.Printf("ERROR: %s\n", err)
		return
	}

	fmt.Printf("Build SUCCEED \\o/\n   > %s\n", builder.distPath)
}

var usage = `TGE command line tool builds and packages TGE applications.

To install:
	$ go get github.com/thommil/tge/cmd/tgebuild
	
Usage:
	tgebuild [-d] [-t target] [-b bundleId] package
	
Use 'tgebuild -h' for arguments details.`

/*
Build compiles and encodes the app named by the import path.

The named package must define a main function.

The -target flag takes a target system name, either android (the default) or ios.

For -target android, if an AndroidManifest.xml is defined in the package directory, it is added to the APK output. Otherwise, a default manifest is generated. By default, this builds a fat APK for all supported instruction sets (arm, 386, amd64, arm64). A subset of instruction sets can be selected by specifying target type with the architecture name. E.g. -target=android/arm,android/386.

For -target ios, gomobile must be run on an OS X machine with Xcode installed.

If the package directory contains an assets subdirectory, its contents are copied into the output.

Flag -iosversion sets the minimal version of the iOS SDK to compile against. The default version is 6.1.

The -bundleid flag is required for -target ios and sets the bundle ID to use with the app.

The -o flag specifies the output file name. If not specified, the output file name depends on the package built.

The -v flag provides verbose output, including the list of packages built.

The build flags -a, -i, -n, -x, -gcflags, -ldflags, -tags, and -work are shared with the build command. For documentation, see 'go help build'.

*/
