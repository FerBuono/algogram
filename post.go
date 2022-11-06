package main

import (
	"algogram/errores"
	TDAHash "algogram/hash"
	TDAHeap "algogram/heap"
	"math"
)

type Post struct {
	texto   string
	id      int
	uid     int
	usuario string
	likes   TDAHash.Diccionario[string, *string]
}

type posts struct {
	lista     []*Post
	diccPosts TDAHash.Diccionario[int, *Post]
	heapPosts TDAHeap.ColaPrioridad[*Post]
}

type Posts interface {
	// GuardarPost agrega un post a la lista
	GuardarPost(p *Post)

	// OrdenarPosts ordena los posts para poder mostrarlos según la afinidad de los usuarios
	OrdenarPosts(uid int)

	// VerProximo devuelve el post correspondiente al usuario con más afinidad
	VerProximo(uid int) (*Post, error)

	// LikearPost le suma una unidad a los likes del post buscado por id
	LikearPost(id int, usuario string) error

	// MostrarLikes devuelve un arreglo con los usuarios que likearon un post
	MostrarLikes(id int) ([]string, error)
}

func CrearListaDePosts(criterio func(Post, Post) int) Posts {
	p := new(posts)
	p.lista = []*Post{}
	p.diccPosts = TDAHash.CrearHash[int, *Post]()
	return p
}

func (p *posts) GuardarPost(post *Post) {
	post.id = len(p.lista)
	post.likes = TDAHash.CrearHash[string, *string]()
	p.lista = append(p.lista, post)
	p.diccPosts.Guardar(post.id, post)
}

func (p *posts) OrdenarPosts(uid int) {
	p.heapPosts = TDAHeap.CrearHeapArr(p.lista, func(p1, p2 *Post) int {
		if int(math.Abs(float64(p1.uid)-float64(uid))) == int(math.Abs(float64(p2.uid)-float64(uid))) {
			return p2.id - p1.id
		}
		return int(math.Abs(float64(p2.uid)-float64(uid))) - int(math.Abs(float64(p1.uid)-float64(uid)))
	})
}

func (p *posts) VerProximo(uid int) (*Post, error) {
	if p.heapPosts.EstaVacia() {
		newError := new(errores.NoLoggeadoONoHayPosts)
		return nil, newError
	}
	if p.heapPosts.VerMax().uid == uid {
		p.heapPosts.Desencolar()
		return p.VerProximo(uid)
	}
	return p.heapPosts.Desencolar(), nil
}

func (p *posts) LikearPost(id int, usuario string) error {
	if !p.diccPosts.Pertenece(id) {
		newError := new(errores.NoLoggeadoOPostInexistente)
		return newError
	}
	post := p.diccPosts.Obtener(id)
	post.likes.Guardar(usuario, nil)
	p.diccPosts.Guardar(id, post)
	return nil
}

func (p posts) MostrarLikes(id int) ([]string, error) {
	usuarios := []string{}
	if !p.diccPosts.Pertenece(id) {
		newError := new(errores.PostInexistenteOSinLikes)
		return usuarios, newError
	}
	post := p.diccPosts.Obtener(id)
	for iter := post.likes.Iterador(); iter.HaySiguiente(); {
		usuario, _ := iter.VerActual()
		usuarios = append(usuarios, usuario)
		iter.Siguiente()
	}
	return usuarios, nil
}
