package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"

	"sd-broadcast/pkg/protocolo"
)

func main() {
	direccionServidor := os.Getenv("SERVIDOR")
	if direccionServidor == "" {
		direccionServidor = "localhost:4000"
	}

	nombre := os.Getenv("NOMBRE")
	if nombre == "" {
		fmt.Print("Ingrese su nombre: ")
		lector := bufio.NewReader(os.Stdin)
		nombreBytes, _, _ := lector.ReadLine()
		nombre = string(nombreBytes)
	}

	// TODO 20: conectar al servidor usando net.Dial("tcp", direccionServidor)
	conexion, err := net.Dial("tcp", direccionServidor)
	if err != nil {
		log.Fatalf("No se pudo conectar al servidor %s: %v", direccionServidor, err)
	}
	defer conexion.Close()

	log.Printf("Conectado al servidor %s como %s", direccionServidor, nombre)

	// TODO 21: enviar mensaje de identificación con protocolo.Codificar
	mensajeIdentificacion := protocolo.NuevoMensaje(nombre, "Cliente identificado", "identificacion")

	err = protocolo.Codificar(conexion, mensajeIdentificacion)
	if err != nil {
		log.Fatalf("No se pudo enviar la identificación: %v", err)
	}

	// TODO 22: iniciar una goroutine que escuche mensajes del servidor en paralelo
	go recibirMensajes(conexion)

	// TODO 23: en el hilo principal, leer líneas de stdin y enviar mensajes de tipo "broadcast"
	lector := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("> ")

		linea, err := lector.ReadString('\n')
		if err != nil {
			log.Printf("Error leyendo entrada: %v", err)
			break
		}

		linea = strings.TrimSpace(linea)

		if linea == "" {
			continue
		}

		if linea == "/salir" {
			log.Println("Saliendo del cliente...")
			break
		}

		mensaje := protocolo.NuevoMensaje(nombre, linea, "broadcast")

		err = protocolo.Codificar(conexion, mensaje)
		if err != nil {
			log.Printf("Error enviando mensaje: %v", err)
			break
		}
	}

	log.Println("Cliente finalizado")
}

// recibirMensajes lee continuamente desde la conexión e imprime en consola
func recibirMensajes(conexion net.Conn) {
	for {
		mensaje, err := protocolo.Decodificar(conexion)
		if err != nil {
			fmt.Println("Desconectado del servidor")
			return
		}

		fmt.Printf("\n[%s] %s - %s\n> ",
			mensaje.Emisor,
			mensaje.Contenido,
			mensaje.Timestamp.Format("15:04:05"),
		)
	}
}
