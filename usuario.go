package main

import (
	TDAHeap "algogram/heap"
	"math"
)

type usuario struct {
	nombre         string
	id             int
	postsOrdenados TDAHeap.ColaPrioridad[*Post]
}

type Usuario interface {

	// LeerUsuario devuelve el nombre de usuario
	LeerUsuario() string

	// Uid devuelve el id del usuario
	Uid() int

	// CrearPost crea un nuevo post en el feed
	CrearPost(texto string) *Post
}

func CrearUsuario(nombre string, id int, posts []*Post) Usuario {
	u := new(usuario)
	u.nombre = nombre
	u.id = id
	u.postsOrdenados = TDAHeap.CrearHeapArr(posts, u.func_afinidad)
	return u
}

func (u usuario) LeerUsuario() string {
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

func (u *usuario) func_afinidad(p1, p2 *Post) int {
	if int(math.Abs(float64(p1.uid)-float64(u.id))) == int(math.Abs(float64(p2.uid)-float64(u.id))) {
		return p2.id - p1.id
	}
	return int(math.Abs(float64(p2.uid)-float64(u.id))) - int(math.Abs(float64(p1.uid)-float64(u.id)))
}
