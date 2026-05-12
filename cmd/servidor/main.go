package main

import (
	"log"
	"net"
	"os"

	"sd-broadcast/internal/registro"
	"sd-broadcast/pkg/protocolo"
)

const puertoPorDefecto = "4000"

func main() {
	puerto := os.Getenv("PUERTO")
	if puerto == "" {
		puerto = puertoPorDefecto
	}

	escuchador, err := net.Listen("tcp", ":"+puerto)
	if err != nil {
		log.Fatalf("No se pudo iniciar el escuchador: %v", err)
	}
	defer escuchador.Close()

	log.Printf("Servidor de broadcast escuchando en :%s", puerto)

	// TODO 8: crear un RegistroClientes usando registro.NuevoRegistro() .-HECHO
	var registroClientes *registro.RegistroClientes
	registroClientes = registro.NuevoRegistro() // Creación del registro de clientes

	// TODO 9: iniciar goroutine para descubrimiento UDP (bonus) .-HECHO

	go iniciarDescubrimientoUDP(puerto)

	for {
		conexion, err := escuchador.Accept()
		if err != nil {
			log.Printf("Error al aceptar conexión: %v", err)
			continue
		}

		// TODO 10: en lugar de llamar directamente a manejarCliente,
		// lanzar una goroutine para atender la conexión concurrentemente
		go manejarCliente(conexion, registroClientes) // Lanzar una goroutine para manejar la conexión

		// manejarCliente(conexion, registroClientes) no utilizo esta llamada al manejador
	}
}

func manejarCliente(conexion net.Conn, registroClientes *registro.RegistroClientes) {
	defer conexion.Close()

	// TODO 11: leer el primer mensaje de identificación del cliente
	// Usar protocolo.Decodificar para obtener el nombre del emisor
	msgInicial, err := protocolo.Decodificar(conexion)
	if err != nil {
		return // Si no se identifica, cerrar
	}
	nombreCliente := msgInicial.Emisor // El protocolo debería traer quién es

	log.Printf("Cliente conectado: %s desde %s", nombreCliente, conexion.RemoteAddr())

	// TODO 12: agregar el cliente al registro usando registroClientes.Agregar(nombreCliente, conexion) .-HECHO
	registroClientes.Agregar(nombreCliente, conexion) // Agregar el cliente al registro

	// TODO 13: notificar a todos los demás clientes que "nombreCliente se unió" .-HECHO
	difundirMensaje(registroClientes, protocolo.NuevoMensaje("Sistema", nombreCliente+" se unió", "sistema"), nombreCliente)
	// Usar difundirMensaje excepto al emisor

	// TODO 14: defer para eliminar al cliente del registro al desconectar .-HECHO
	defer registroClientes.Eliminar(nombreCliente)
	defer difundirMensaje(registroClientes, protocolo.NuevoMensaje("Sistema", nombreCliente+" se desconectó", "sistema"), nombreCliente)

	// TODO 15: bucle para leer mensajes del cliente y reenviarlos a todos los demás .-HECHO
	// Usar protocolo.Decodificar en un for {}
	// Si el mensaje.Tipo es "broadcast", usar difundirMensaje
	// Si hay error en Decode, salir del bucle (cliente desconectado)
	for {
		mensaje, err := protocolo.Decodificar(conexion)
		if err != nil {
			log.Printf("Error al decodificar mensaje de %s: %v", nombreCliente, err)
			break
		}

		if mensaje.Tipo == "broadcast" {
			difundirMensaje(registroClientes, mensaje, nombreCliente)
		}
	}

	log.Printf("Cliente desconectado: %s", nombreCliente)
}

// difundirMensaje envía un mensaje a todos los clientes excepto al emisor indicado
func difundirMensaje(r *registro.RegistroClientes, msg protocolo.Mensaje, exceptoEmisor string) {
	// Obtener la slice de las conexiones y los nombres de los clientes registrados
	conexiones := r.ObtenerConexiones()
	nombres := r.Nombres()

	// Iteración
	for i, conn := range conexiones {
		// No se envía el mensaje al emisor
		if nombres[i] != exceptoEmisor {
			// Enviar el mensaje usando protocolo.Codificar
			err := protocolo.Codificar(conn, msg)
			if err != nil {
				log.Printf("Error enviando a %s: %v", nombres[i], err)
			}
		}
	}
}

// iniciarDescubrimientoUDP inicia un servidor UDP para responder a solicitudes de descubrimiento
func iniciarDescubrimientoUDP(puerto string) {
	addr, _ := net.ResolveUDPAddr("udp", ":9999") // Puerto estándar de descubrimiento
	conn, _ := net.ListenUDP("udp", addr)
	defer conn.Close()

	for {
		buf := make([]byte, 1024)
		_, remoteAddr, _ := conn.ReadFromUDP(buf)
		// Respondemos con el puerto TCP en el que estamos escuchando, el 4000
		conn.WriteToUDP([]byte(puerto), remoteAddr)
	}
}

func manejarSolicitudUDP(conn *net.UDPConn, addr *net.UDPAddr, puerto string) {

}
