package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"
	"sync"
)

func createImage(rect image.Rectangle) (created *image.NRGBA) {
	pix := make([]uint8, rect.Dx()*rect.Dy()*4)

	created = &image.NRGBA{
		Pix:    pix,
		Stride: rect.Dx() * 4,
		Rect:   rect,
	}
	return
}

func addPixColorImage(iter int, max_Iter int, i int, y int, img *image.NRGBA) {

	var gradC color.NRGBA
	var n uint8 = uint8(((255) / (max_Iter)) * iter) // transformation du nombre d'itération en un uint8
	gradC.R = n
	gradC.G = n
	gradC.B = 1 / n
	gradC.A = 255
	img.SetNRGBA(i, y, gradC) // On ajoute une couleur au pixel i,y
}

func mandelbrot_Color(n_Max int, cx float32, cy float32) int {

	var xn, yn, tmpx, tmpy float32 = 0, 0, 0, 0
	var n int = 0

	for ((xn*xn + yn*yn) < 4) && (n < n_Max) { // itération pour mandelbrot
		tmpx = xn
		tmpy = yn

		xn = tmpx*tmpx - tmpy*tmpy + cx
		yn = 2*tmpx*tmpy + cy
		n++
	}

	return n
}

func worker(wg *sync.WaitGroup, id int, max_Iter int, cx float32, cy float32, i int, y int, img *image.NRGBA) {
	defer wg.Done()

	//fmt.Printf("Worker %v: Started\n", id)

	addPixColorImage(mandelbrot_Color(max_Iter, cx, cy), max_Iter, i, y, img)

	//fmt.Printf("Worker %v: Finished\n", id)
}

func main() {
	fmt.Println("Mandelbrot By LEONARD Matthias 20308 ECAM Student")

	// déclarations taille image et paramètres
	const img_Wide, img_Tall, max_Iter int = 5000, 5000, 180
	const x_Min, x_Max, y_Min, y_Max float32 = -2, +0.5, -1.25, +1.25 // domaine de mandelbrot
	var cx, cy float32 = 0, 0

	// création image
	rect := image.Rect(0, 0, img_Wide, img_Tall)
	img := createImage(rect)

	//
	var wg sync.WaitGroup
	var maxWorkers = 5

	// Le loadbalancer est implémenté de manière discrète il fait tourner les itérations de madelbrot et le découpage est basé sur les lignes de l'image.

	for i := 0; i < img_Wide; i++ { // Wide
		for y := 0; y < img_Tall; y = y + maxWorkers { // Tall

			for z := 0; z < maxWorkers; z++ {

				cx = ((float32(i)*(x_Max-x_Min))/(float32(img_Wide)+x_Min) - 2)
				cy = ((float32(y+z)*(y_Min-y_Max))/(float32(img_Tall)+y_Max) + 1.25) // Ajout du +z pour prendre en compte le nombre de worker et pas refaire les mêmes calculs

				//fmt.Println("Adding worker", z)
				wg.Add(1)
				go worker(&wg, z, max_Iter, cx, cy, i, y+z, img)
			}

			//fmt.Printf("Waiting for %d workers to finish\n", maxWorkers)
			wg.Wait()
			//fmt.Println("All Workers Completed")

		}
	}

	// Enregistrement de l'image

	outputFile, err := os.Create("Test5000x5000n180.png")
	if err != nil {
	}
	png.Encode(outputFile, img)
	outputFile.Close()

}
