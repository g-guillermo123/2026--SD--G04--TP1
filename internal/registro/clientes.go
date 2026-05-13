package registro

import (
	"net"
	"sync"
)

// RegistroClientes mantiene el listado de conexiones activas de forma segura
type RegistroClientes struct {
	// TODO 1: agregar un campo sync.RWMutex para proteger el mapa .-HECHO
	mu       sync.RWMutex // Declaración del mutex para proteger el acceso al mapa de clientes
	clientes map[string]net.Conn
}

// NuevoRegistro crea un registro vacío
func NuevoRegistro() *RegistroClientes {
	// TODO 2: inicializar el mapa de clientes .-HECHO
	// Inicialización del mapa de clientes
	r := &RegistroClientes{
		clientes: make(map[string]net.Conn),
	}
	return r
}

// Agregar añade un cliente al registro
func (r *RegistroClientes) Agregar(nombre string, conexion net.Conn) {
	// TODO 3: bloquear para escritura, agregar al mapa, desbloquear .-HECHO
	// Bloqueo para escritura, agregar al mapa y desbloquear
	r.mu.Lock()                   // Bloqueo exclusivo para escritura
	r.clientes[nombre] = conexion // Agregar el cliente al mapa
	r.mu.Unlock()                 // Desbloqueo después de agregar
}

// Eliminar remueve un cliente del registro
func (r *RegistroClientes) Eliminar(nombre string) {
	// TODO 4: bloquear para escritura, eliminar del mapa, desbloquear .-HECHO
	// Bloqueo para escritura, eliminar del mapa y desbloquear
	r.mu.Lock()                // Bloqueo exclusivo para escritura
	delete(r.clientes, nombre) // Eliminación del registro
	r.mu.Unlock()              // Desbloqueo después de eliminar
}

// ObtenerConexiones devuelve una copia de todas las conexiones activas
func (r *RegistroClientes) ObtenerConexiones() []net.Conn {
	// TODO 5: bloquear para lectura, copiar conexiones a un slice, desbloquear .-HECHO
	r.mu.RLock() // Bloqueo para lectura
	// Crear un slice para almacenar las conexiones activas
	conexiones := make([]net.Conn, 0, len(r.clientes))
	for _, conn := range r.clientes {
		conexiones = append(conexiones, conn) // Copiar cada conexión al slice
	}
	r.mu.RUnlock() // Desbloqueo después de copiar
	return conexiones
}

// Cantidad devuelve el número de clientes conectados
func (r *RegistroClientes) Cantidad() int {
	// TODO 6: bloquear para lectura, retornar len del mapa, desbloquear .-HECHO
	r.mu.RLock()                // Bloqueo para lectura
	cantidad := len(r.clientes) // Obtener la cantidad de clientes conectados
	r.mu.RUnlock()              // Desbloqueo después de obtener la cantidad
	return cantidad
}

// Nombres devuelve un slice con los nombres de los clientes
func (r *RegistroClientes) Nombres() []string {
	// TODO 7: bloquear para lectura, copiar nombres a un slice, desbloquear .-HECHO
	r.mu.RLock() // Bloqueo para lectura
	nombres := make([]string, 0, len(r.clientes))
	for nombre := range r.clientes {
		nombres = append(nombres, nombre)
	}
	r.mu.RUnlock() // Desbloqueo después de copiar
	return nombres
}

// ObtenerClientes devuelve una copia del mapa de clientes activos
func (r *RegistroClientes) ObtenerClientes() map[string]net.Conn {
	r.mu.RLock()
	defer r.mu.RUnlock()

	copia := make(map[string]net.Conn)

	for nombre, conexion := range r.clientes {
		copia[nombre] = conexion
	}

	return copia
}
