package main

import (
	"bufio"
	"fmt"
	_ "image/jpeg"
	_ "image/png"
	"log"
	"path/filepath"
	"sort"
	"strconv"

	"github.com/disintegration/imageorient"
	"github.com/nfnt/resize"
	"github.com/pkg/errors"
	"github.com/urfave/cli/v2"

	"image/jpeg"
	"os"
	"strings"
	"time"
)

func Convert(inputPath, outputPath string, resizeImg bool) error {
	inputFile, err := os.Open(inputPath)
	if err != nil {
		return errors.Wrap(err, "Open input file failed")
	}
	defer inputFile.Close()

	imgSrc, _, err := imageorient.Decode(inputFile)
	if err != nil {
		return errors.Wrap(err, "png.Decode failed")
	}

	outputFile, err := os.Create(outputPath)
	if err != nil {
		return errors.Wrap(err, "os.Create(outputPath) failed")
	}
	defer outputFile.Close()

	if resizeImg && (imgSrc.Bounds().Max.X > MaxWidth || imgSrc.Bounds().Max.Y > MaxHeight) {
		var newSizeX, newSizeY uint
		if float64(imgSrc.Bounds().Max.X)/float64(MaxWidth) > float64(imgSrc.Bounds().Max.Y)/float64(MaxHeight) {
			newSizeX = uint(MaxWidth)
		} else {
			newSizeY = uint(MaxHeight)
		}
		fmt.Printf("Resize %s : width * height: %d * %d  ->  %d * %d   preserve ratio\n", outputPath, imgSrc.Bounds().Max.X, imgSrc.Bounds().Max.Y, newSizeX, newSizeY)

		imgSrc = resize.Resize(newSizeX, newSizeY, imgSrc, resize.Lanczos3)
	}

	var opt jpeg.Options
	opt.Quality = 100
	if err := jpeg.Encode(outputFile, imgSrc, &opt); err != nil {
		return errors.Wrap(err, "jpeg.Encode failed")
	}
	return nil
}

func ConvertJPG(inputPath, outputPath string, resizeImg bool) error {
	inputFile, err := os.Open(inputPath)
	if err != nil {
		return errors.Wrap(err, "Open input file failed")
	}
	defer inputFile.Close()

	imgSrc, _, err := imageorient.Decode(inputFile)
	if err != nil {
		return errors.Wrap(err, "png.Decode failed")
	}

	outputFile, err := os.Create(outputPath)
	if err != nil {
		return errors.Wrap(err, "os.Create(outputPath) failed")
	}
	defer outputFile.Close()

	if resizeImg && (imgSrc.Bounds().Max.X > MaxWidth || imgSrc.Bounds().Max.Y > MaxHeight) {
		var newSizeX, newSizeY uint
		if float64(imgSrc.Bounds().Max.X)/float64(MaxWidth) > float64(imgSrc.Bounds().Max.Y)/float64(MaxHeight) {
			newSizeX = uint(MaxWidth)
		} else {
			newSizeY = uint(MaxHeight)
		}
		fmt.Printf("Resize %s : width * height: %d * %d  ->  %d * %d   preserve ratio\n", outputPath, imgSrc.Bounds().Max.X, imgSrc.Bounds().Max.Y, newSizeX, newSizeY)

		imgSrc = resize.Resize(newSizeX, newSizeY, imgSrc, resize.Lanczos3)
	}

	var opt jpeg.Options
	opt.Quality = 100
	if err := jpeg.Encode(outputFile, imgSrc, &opt); err != nil {
		return errors.Wrap(err, "jpeg.Encode failed")
	}

	return nil
}

func readFiles(ext string) []string {
	var result []string
	err := filepath.Walk("./", func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() && !strings.HasPrefix(strings.ToLower(path), "output") && strings.Contains(strings.ToLower(path), ext) {
			result = append(result, path)
		}
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}

	return result
}

func wait() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("please enter any key to exit...")
	_, _ = reader.ReadString('\n')
}

func mkdir() {
	if _, err := os.Stat("./output"); os.IsNotExist(err) {
		os.MkdirAll("./output", 0777)
	}
}

var (
	MaxWidth  = 1280
	MaxHeight = 720
)

func main() {
	app := cli.NewApp()
	app.Usage = `Transforme un png en jpg ou lit des jpg pour les sauver non entrelacés. Sauve les fichiers dans le répertoire output avec sous directory.
		Si aucune commande est donnée alors les images sont testées sur leur taille et celles qui sont au dessus de 1280x720 (ou 720x1280 si portrait) sont listées.`
	app.Name = "NixPlayJpgTools"
	app.Compiled = time.Now()
	app.Version = "1.0.0"
	app.Commands = []*cli.Command{
		{
			Name: "png",
			Usage: `convertit un png en jpg et l'écrit dans le répertoire output
		NixPlayJpgTools.exe png --resize va convertir des images png en jpg et les retailler si besoin
		NixPlayJpgTools.exe --portrait png --resize va faire comme ci-dessus mais va utiliser une taille 720x1280`,
			Flags: []cli.Flag{
				&cli.BoolFlag{
					Name:  "resize",
					Usage: "Resize l'image pour 1280x720",
				},
			},
			Action: func(c *cli.Context) error {
				if c.Bool("portrait") {
					MaxHeight, MaxWidth = MaxWidth, MaxHeight
				}
				succeedFileNum := 0
				files := readFiles(".png")
				mkdir()
				for index, name := range files {
					i := strconv.Itoa(index + 1)
					outputName := strings.Replace(strings.Replace(name, ".png", ".jpg", 1), ".PNG", ".jpg", 1)
					os.MkdirAll(filepath.Dir("./output/"+name), os.ModePerm)
					if err := Convert(name, "./output/"+outputName, c.Bool("resize")); err != nil {
						fmt.Println(i + ". " + name + ": Failed!!!")
						fmt.Println(err)
					} else {
						fmt.Println(i + ". " + name + ": Done")
						succeedFileNum++
					}
				}
				fmt.Println(strconv.Itoa(succeedFileNum) + "/" + strconv.Itoa(len(files)) + " done")
				if !c.Bool("nowait") {
					wait()
				}

				return nil
			},
		},
		{
			Name: "jpg",
			Usage: `relit des jpg et les sauve non entrelacés dans le répertoire output
		NixPlayJpgTools.exe jpg --resize va convertir des images jpg en jpg non entrelacés et les retailler si besoin, le fichier de sortie sera en .jpg
		NixPlayJpgTools.exe --portrait jpg --resize va faire comme ci-dessus mais va utiliser une taille 720x1280`,
			Flags: []cli.Flag{
				&cli.BoolFlag{
					Name:  "resize",
					Usage: fmt.Sprintf("Resize l'image pour %d*%d", MaxWidth, MaxHeight),
				},
			},
			Action: func(c *cli.Context) error {
				if c.Bool("portrait") {
					MaxHeight, MaxWidth = MaxWidth, MaxHeight
				}
				succeedFileNum := 0
				files := readFiles(".jpg")
				files2 := readFiles(".jpeg")
				files = append(files, files2...)
				mkdir()
				for index, name := range files {
					i := strconv.Itoa(index + 1)
					os.MkdirAll(filepath.Dir("./output/"+name), os.ModePerm)
					newName := "./output/" + strings.Replace(name, filepath.Ext(name), ".jpg", 1)
					if err := ConvertJPG(name, newName, c.Bool("resize")); err != nil {
						fmt.Println(i + ". " + name + ": Failed!!!")
						fmt.Println(err)
					} else {
						fmt.Println(i + ". " + name + ": Done")
						succeedFileNum++
					}
				}
				fmt.Println(strconv.Itoa(succeedFileNum) + "/" + strconv.Itoa(len(files)) + " done")
				if !c.Bool("nowait") {
					wait()
				}
				return nil
			},
		},
	}

	app.Flags = []cli.Flag{
		&cli.BoolFlag{
			Name:  "nowait",
			Usage: "Ne pas attendre une touche après process",
		},
		&cli.BoolFlag{
			Name:  "portrait",
			Usage: "Va utiliser 702x1280 pour le resize si demandé",
		},
	}

	sort.Sort(cli.CommandsByName(app.Commands))
	sort.Sort(cli.FlagsByName(app.Flags))

	app.Action = func(c *cli.Context) error {
		if c.Bool("portrait") {
			MaxHeight, MaxWidth = MaxWidth, MaxHeight
		}
		files := readFiles(".jpg")
		files2 := readFiles(".jpeg")
		files3 := readFiles(".png")
		files = append(files, files2...)
		files = append(files, files3...)
		countFilesNotGoodSize := 0

		for _, f := range files {
			reader, err := os.Open(f)
			if err != nil {
				log.Fatal(err)
			}

			imagef, _, err := imageorient.Decode(reader)
			reader.Close()
			if err != nil {
				fmt.Fprintf(os.Stderr, "%s: %v\n", f, err)
			} else {
				if imagef.Bounds().Max.X > MaxWidth || imagef.Bounds().Max.Y > MaxHeight {
					fmt.Printf("%s  width * height: %d * %d\n", f, imagef.Bounds().Max.X, imagef.Bounds().Max.Y)
					countFilesNotGoodSize++
				}
			}
		}
		fmt.Printf("\nNombre de fichiers avec une taille plus grande que %d*%d: %d\n", MaxWidth, MaxHeight, countFilesNotGoodSize)

		return nil
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
