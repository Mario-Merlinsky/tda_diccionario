package diccionario

import (
	"fmt"
)

type estado int

const (
	VACIO estado = iota
	OCUPADO
	BORRADO
	CAPACIDAD_INICIAL     = 13
	FACTOR_REDIMENSION    = 2
	FACTOR_CARGA_INFERIOR = 0.20
	FACTOR_CARGA_SUPERIOR = 0.70
	VALOR_INICIAL         = 0
	ERROR_CLAVE           = "La clave no pertenece al diccionario"
	ERROR_ITERADOR        = "El iterador termino de iterar"
)

type celdaHash[K comparable, V any] struct {
	clave  K
	valor  V
	estado estado
}

type hashCerrado[K comparable, V any] struct {
	tabla    []celdaHash[K, V]
	cantidad int
	tam      int
	borrados int
}

type iteradorHash[K comparable, V any] struct {
	hash          *hashCerrado[K, V]
	indice_actual int
}

func CrearHash[K comparable, V any]() Diccionario[K, V] {
	return &hashCerrado[K, V]{tabla: crearTabla[K, V](CAPACIDAD_INICIAL), tam: CAPACIDAD_INICIAL}
}

func (hash *hashCerrado[K, V]) Guardar(clave K, dato V) {
	if (float32(hash.cantidad)+float32(hash.borrados))/float32(hash.tam) >= FACTOR_CARGA_SUPERIOR {
		hash.redimension(hash.tam * FACTOR_REDIMENSION)
	}
	indice_definitivo, estado := buscarPosicionValida[K, V](hash.tabla, clave, hash.tam)

	if estado == VACIO {
		hash.tabla[indice_definitivo].clave = clave
		hash.tabla[indice_definitivo].estado = OCUPADO
		hash.cantidad++
	}
	hash.tabla[indice_definitivo].valor = dato
}

func (hash hashCerrado[K, V]) Pertenece(clave K) bool {
	_, estado := buscarPosicionValida[K, V](hash.tabla, clave, hash.tam)
	return !(estado == VACIO)
}

func (hash hashCerrado[K, V]) Obtener(clave K) V {

	indice_definitivo, estado := buscarPosicionValida[K, V](hash.tabla, clave, hash.tam)

	if estado == VACIO {
		panic(ERROR_CLAVE)
	}

	return hash.tabla[indice_definitivo].valor
}

func (hash *hashCerrado[K, V]) Borrar(clave K) V {
	if float32(hash.cantidad)/float32(hash.tam) <= FACTOR_CARGA_INFERIOR && hash.tam > CAPACIDAD_INICIAL {
		hash.redimension(hash.tam / FACTOR_REDIMENSION)
	}
	indice_definitivo, estado := buscarPosicionValida[K, V](hash.tabla, clave, hash.tam)

	if estado == VACIO {
		panic(ERROR_CLAVE)
	}
	hash.tabla[indice_definitivo].estado = BORRADO
	hash.borrados++
	hash.cantidad--
	return hash.tabla[indice_definitivo].valor
}

func (hash hashCerrado[K, V]) Cantidad() int {
	return hash.cantidad
}

func (hash hashCerrado[K, V]) Iterar(visitar func(clave K, dato V) bool) {
	for _, celda := range hash.tabla {
		if celda.estado == OCUPADO && !visitar(celda.clave, celda.valor) {
			break
		}
	}
}

func (hash *hashCerrado[K, V]) Iterador() IterDiccionario[K, V] {
	iter := &iteradorHash[K, V]{hash: hash, indice_actual: VALOR_INICIAL}
	iter.encontrarPosicionOcupada()
	return iter
}

func (iter *iteradorHash[K, V]) HaySiguiente() bool {
	return !(iter.indice_actual == iter.hash.tam)
}

func (iter iteradorHash[K, V]) VerActual() (K, V) {
	if !iter.HaySiguiente() {
		panic(ERROR_ITERADOR)
	}

	return iter.hash.tabla[iter.indice_actual].clave, iter.hash.tabla[iter.indice_actual].valor
}

func (iter *iteradorHash[K, V]) Siguiente() {
	if !iter.HaySiguiente() {
		panic(ERROR_ITERADOR)
	}
	iter.indice_actual++
	iter.encontrarPosicionOcupada()
}

func (iter *iteradorHash[K, V]) encontrarPosicionOcupada() {
	for iter.indice_actual < iter.hash.tam && iter.hash.tabla[iter.indice_actual].estado != OCUPADO {
		iter.indice_actual++
	}
}

func hashing[K comparable](clave K) uint64 {
	return HashBytes64(convertirABytes[K](clave))
}

func convertirABytes[K comparable](clave K) []byte {
	return []byte(fmt.Sprintf("%v", clave))
}

func buscarPosicionValida[K comparable, V any](tabla []celdaHash[K, V], clave K, tam int) (int, estado) {
	indice := int(hashing(clave) % uint64(tam))
	for i := 0; indice+i < tam; i++ {
		if (tabla[indice+i].estado == OCUPADO && tabla[indice+i].clave == clave) || tabla[indice+i].estado == VACIO {
			return indice + i, tabla[indice+i].estado
		}
	}
	return buscarPosicionValida[K, V](tabla, clave, indice)
}

func (hash *hashCerrado[K, V]) redimension(nuevo_largo int) {
	tabla_nueva := crearTabla[K, V](nuevo_largo)
	for _, celda := range hash.tabla {
		if celda.estado == OCUPADO {
			indice_definitivo, _ := buscarPosicionValida[K, V](tabla_nueva, celda.clave, nuevo_largo)

			tabla_nueva[indice_definitivo].clave = celda.clave
			tabla_nueva[indice_definitivo].valor = celda.valor
			tabla_nueva[indice_definitivo].estado = OCUPADO
		}
	}
	hash.tabla = tabla_nueva
	hash.tam = nuevo_largo
	hash.borrados = VALOR_INICIAL
}

func crearTabla[K comparable, V any](largo int) []celdaHash[K, V] {
	return make([]celdaHash[K, V], largo)
}
