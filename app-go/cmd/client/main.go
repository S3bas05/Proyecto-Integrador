package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"imageflow/app-go/internal/client"
	"imageflow/app-go/internal/models"
)

func main() {
	serverURL := os.Getenv("IMAGEFLOW_SERVER_URL")
	if serverURL == "" {
		serverURL = "http://localhost:8080"
	}

	api := client.NewAPIClient(serverURL)
	fmt.Println("====================================")
	fmt.Println("       IMAGEFLOW - CLIENTE GO")
	fmt.Println("====================================")
	fmt.Println("Servidor:", serverURL)
	if err := api.Health(); err != nil {
		fmt.Println("No fue posible conectar con ImageFlow:", err)
		fmt.Println("Verifica que el servidor esté iniciado y que la URL sea correcta.")
		return
	}
	fmt.Println("Conexión establecida correctamente.")

	reader := bufio.NewReader(os.Stdin)
	for {
		showMenu()
		option := readLine(reader, "Selecciona una opción: ")
		fmt.Println()

		switch option {
		case "1":
			uploadImage(reader, api)
		case "2":
			listImages(api, "", "")
		case "3":
			name := readLine(reader, "Nombre o parte del nombre (Enter para omitir): ")
			category := readLine(reader, "Categoría (Enter para omitir): ")
			listImages(api, name, category)
		case "4":
			id := readLine(reader, "ID de la imagen: ")
			destination := readLine(reader, "Ruta donde guardar la imagen (ej. descargas/foto.png): ")
			if id == "" || destination == "" {
				fmt.Println("El ID y la ruta de destino son obligatorios.")
				break
			}
			if err := api.Download(id, destination); err != nil {
				fmt.Println("Error:", err)
			} else {
				fmt.Println("Imagen descargada en:", destination)
			}
		case "5":
			id := readLine(reader, "ID de la imagen que deseas eliminar: ")
			if id == "" {
				fmt.Println("Debes ingresar un ID.")
				break
			}
			confirmation := strings.ToLower(readLine(reader, "Escribe S para confirmar: "))
			if confirmation != "s" {
				fmt.Println("Operación cancelada.")
				break
			}
			if err := api.Delete(id); err != nil {
				fmt.Println("Error:", err)
			} else {
				fmt.Println("Imagen eliminada correctamente.")
			}
		case "0":
			fmt.Println("Cerrando ImageFlow. Hasta luego.")
			return
		default:
			fmt.Println("Opción no válida. Selecciona un número del menú.")
		}
		fmt.Println()
	}
}

func showMenu() {
	fmt.Println("----------- MENÚ PRINCIPAL -----------")
	fmt.Println("1. Cargar imagen")
	fmt.Println("2. Consultar todas las imágenes")
	fmt.Println("3. Buscar por nombre o categoría")
	fmt.Println("4. Descargar imagen")
	fmt.Println("5. Eliminar imagen")
	fmt.Println("0. Salir")
}

func readLine(reader *bufio.Reader, prompt string) string {
	fmt.Print(prompt)
	value, _ := reader.ReadString('\n')
	return strings.TrimSpace(value)
}

func uploadImage(reader *bufio.Reader, api *client.APIClient) {
	path := readLine(reader, "Ruta de la imagen que deseas cargar: ")
	category := readLine(reader, "Categoría de la imagen: ")
	if path == "" {
		fmt.Println("La ruta del archivo es obligatoria.")
		return
	}

	metadata, err := api.Upload(path, category)
	if err != nil {
		fmt.Println("No se pudo cargar la imagen:", err)
		return
	}
	fmt.Println("Imagen cargada correctamente.")
	printImage(metadata)
}

func listImages(api *client.APIClient, name, category string) {
	images, err := api.ListImages(name, category)
	if err != nil {
		fmt.Println("No se pudieron consultar las imágenes:", err)
		return
	}
	if len(images) == 0 {
		fmt.Println("No se encontraron imágenes.")
		return
	}
	fmt.Printf("Se encontraron %d imagen(es):\n", len(images))
	for _, image := range images {
		printImage(image)
	}
}

func printImage(image models.Metadata) {
	fmt.Println("--------------------------------------")
	fmt.Println("ID:         ", image.ID)
	fmt.Println("Nombre:     ", image.Name)
	fmt.Println("Categoría:  ", image.Category)
	fmt.Println("Formato:    ", image.ContentType)
	fmt.Printf("Tamaño:      %d bytes\n", image.Size)
	fmt.Println("Fecha:      ", image.CreatedAt.Format("2006-01-02 15:04:05"))
}
