package app

import (
	"algogram/errores"
	TDAHash "algogram/hash"
	TDAHeap "algogram/heap"
	"math"
)

type usuario struct {
	nombre string
	id     int
	posts  TDAHeap.ColaPrioridad[*Post]
}

type listaUsuarios struct {
	lista TDAHash.Diccionario[string, Usuario]
}

func CrearListaDeUsuarios() ListaUsuarios {
	l := new(listaUsuarios)
	l.lista = TDAHash.CrearHash[string, Usuario]()
	return l
}

// Primitivas de ListaUsuarios

func (l *listaUsuarios) GuardarUsuario(nombre string, id int) Usuario {
	u := new(usuario)
	u.nombre = nombre
	u.id = id
	u.posts = TDAHeap.CrearHeap(u.func_afinidad)
	l.lista.Guardar(nombre, u)
	return u
}

func (l *listaUsuarios) BuscarUsuario(nombre string) (Usuario, error) {
	if !l.lista.Pertenece(nombre) {
		newError := new(errores.UsuarioInexistente)
		return nil, newError
	}
	return l.lista.Obtener(nombre), nil
}

func (l *listaUsuarios) GuardarPost(post Post) {
	for iter := l.lista.Iterador(); iter.HaySiguiente(); {
		_, usuario := iter.VerActual()
		usuario.AgregarPost(post)
		iter.Siguiente()
	}
}

// Primitivas de Usuario

func (u usuario) VerUsuario() (int, string) {
	return u.id, u.nombre
}

func (u *usuario) AgregarPost(post Post) {
	_, _, _, uid := post.VerPost()
	if uid != u.id {
		u.posts.Encolar(&post)
	}
}

func (u *usuario) VerProximoPost() (Post, error) {
	if u.posts.EstaVacia() {
		newError := new(errores.NoLoggeadoONoHayPosts)
		return nil, newError
	}
	return *u.posts.Desencolar(), nil
}

// Funciones y métodos auxiliares

func (u *usuario) func_afinidad(p1, p2 *Post) int {
	_, id1, _, uid1 := (*p1).VerPost()
	_, id2, _, uid2 := (*p2).VerPost()

	if int(math.Abs(float64(uid1)-float64(u.id))) == int(math.Abs(float64(uid2)-float64(u.id))) {
		return id2 - id1
	}
	return int(math.Abs(float64(uid2)-float64(u.id))) - int(math.Abs(float64(uid1)-float64(u.id)))
}
