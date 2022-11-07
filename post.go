package main

import (
	TDA_ABB "algogram/abb"
	"algogram/errores"
	TDAHash "algogram/hash"
	"strings"
)

type Post struct {
	texto   string
	id      int
	uid     int
	usuario string
	likes   TDA_ABB.DiccionarioOrdenado[string, *string]
}

type posts struct {
	diccPosts TDAHash.Diccionario[int, *Post]
}

type Posts interface {
	// GuardarPost agrega un post a la lista
	GuardarPost(p *Post)

	// LikearPost le suma una unidad a los likes del post buscado por id
	LikearPost(id int, usuario string) error

	// MostrarLikes devuelve un arreglo con los usuarios que likearon un post
	MostrarLikes(id int) ([]string, error)
}

func CrearListaDePosts(criterio func(Post, Post) int) Posts {
	p := new(posts)
	p.diccPosts = TDAHash.CrearHash[int, *Post]()
	return p
}

func (p *posts) GuardarPost(post *Post) {
	post.id = p.diccPosts.Cantidad()
	post.likes = TDA_ABB.CrearABB[string, *string](func(a, b string) int { return strings.Compare(a, b) })
	p.diccPosts.Guardar(post.id, post)
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
	if post.likes.Cantidad() == 0 {
		newError := new(errores.PostInexistenteOSinLikes)
		return usuarios, newError
	}
	for iter := post.likes.Iterador(); iter.HaySiguiente(); {
		usuario, _ := iter.VerActual()
		usuarios = append(usuarios, usuario)
		iter.Siguiente()
	}
	return usuarios, nil
}
