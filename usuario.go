package main

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

type Usuario interface {

	// LeerUsuario devuelve el nombre de usuario
	NombreUsuario() string

	// Uid devuelve el id del usuario
	Uid() int

	// CrearPost crea un nuevo post en el feed
	CrearPost(texto string) *Post

	// AgregarPost agrega un nuevo post a su lista de posts
	AgregarPost(post *Post)

	// VerProximoPost devuelve el siguiente post en el feed de acuerdo a la funcion de afinidad
	VerProximoPost() (*Post, error)
}

type listaUsuarios struct {
	lista TDAHash.Diccionario[string, Usuario]
}

type ListaUsuarios interface {

	// GuardarUsuario agrega un nuevo usuario a la lista
	GuardarUsuario(usuario Usuario)

	// BuscarUsuario busca el usuario pedido por nombre en la lista
	BuscarUsuario(nombre string) (Usuario, error)

	// Iterador devuelve un IterDiccionario para la lista de usuarios
	Iterador() TDAHash.IterDiccionario[string, Usuario]
}

func CrearListaDeUsuarios() ListaUsuarios {
	l := new(listaUsuarios)
	l.lista = TDAHash.CrearHash[string, Usuario]()
	return l
}

func (l *listaUsuarios) GuardarUsuario(usuario Usuario) {
	l.lista.Guardar(usuario.NombreUsuario(), usuario)
}

func (l *listaUsuarios) BuscarUsuario(nombre string) (Usuario, error) {
	if !l.lista.Pertenece(nombre) {
		newError := new(errores.UsuarioInexistente)
		return nil, newError
	}
	return l.lista.Obtener(nombre), nil
}

func (l *listaUsuarios) Iterador() TDAHash.IterDiccionario[string, Usuario] {
	return l.lista.Iterador()
}

func CrearUsuario(nombre string, id int) Usuario {
	u := new(usuario)
	u.nombre = nombre
	u.id = id
	u.posts = TDAHeap.CrearHeap(u.func_afinidad)
	return u
}

func (u usuario) NombreUsuario() string {
	return u.nombre
}

func (u usuario) Uid() int {
	return u.id
}

func (u *usuario) CrearPost(texto string) *Post {
	post := new(Post)
	post.texto, post.uid, post.usuario = texto, u.id, u.nombre
	return post
}

func (u *usuario) AgregarPost(post *Post) {
	if post.uid != u.id {
		u.posts.Encolar(post)
	}
}

func (u *usuario) VerProximoPost() (*Post, error) {
	if u.posts.EstaVacia() {
		newError := new(errores.NoLoggeadoONoHayPosts)
		return nil, newError
	}
	return u.posts.Desencolar(), nil
}

func (u *usuario) func_afinidad(p1, p2 *Post) int {
	if int(math.Abs(float64(p1.uid)-float64(u.id))) == int(math.Abs(float64(p2.uid)-float64(u.id))) {
		return p2.id - p1.id
	}
	return int(math.Abs(float64(p2.uid)-float64(u.id))) - int(math.Abs(float64(p1.uid)-float64(u.id)))
}
