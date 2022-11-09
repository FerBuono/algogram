package app

import (
	TDA_ABB "algogram/abb"
	"algogram/errores"
	TDAHash "algogram/hash"
	"strings"
)

type post struct {
	texto   string
	id      int
	uid     int
	usuario string
	likes   TDA_ABB.DiccionarioOrdenado[string, *string]
}

type listaPosts struct {
	diccPosts TDAHash.Diccionario[int, Post]
}

func CrearListaDePosts() ListaPosts {
	l := new(listaPosts)
	l.diccPosts = TDAHash.CrearHash[int, Post]()
	return l
}

// Primitivas de ListaPosts

func (l *listaPosts) GuardarPost(texto string, uid int, usuario string) Post {
	p := new(post)
	p.texto, p.uid, p.usuario = texto, uid, usuario
	p.id = l.diccPosts.Cantidad()
	p.likes = TDA_ABB.CrearABB[string, *string](func(a, b string) int { return strings.Compare(a, b) })
	l.diccPosts.Guardar(p.id, p)
	return p
}

func (l *listaPosts) LikearPost(id int, usuario string) error {
	if !l.diccPosts.Pertenece(id) {
		newError := new(errores.NoLoggeadoOPostInexistente)
		return newError
	}
	post := l.diccPosts.Obtener(id)
	post.Likear(usuario)
	l.diccPosts.Guardar(id, post)
	return nil
}

func (l listaPosts) MostrarLikes(id int) ([]string, error) {
	usuarios := []string{}
	if !l.diccPosts.Pertenece(id) {
		newError := new(errores.PostInexistenteOSinLikes)
		return usuarios, newError
	}
	post := l.diccPosts.Obtener(id)
	if post.Likes() == 0 {
		newError := new(errores.PostInexistenteOSinLikes)
		return usuarios, newError
	}
	usuarios = post.MostrarLikes()
	return usuarios, nil
}

// Primitivas de Post

func (p *post) VerPost() (string, int, string, int) {
	return p.texto, p.id, p.usuario, p.uid
}

func (p *post) MostrarLikes() []string {
	usuarios := []string{}
	for iter := p.likes.Iterador(); iter.HaySiguiente(); {
		usuario, _ := iter.VerActual()
		usuarios = append(usuarios, usuario)
		iter.Siguiente()
	}
	return usuarios
}

func (p *post) Likes() int {
	return p.likes.Cantidad()
}

func (p *post) Likear(usuario string) {
	p.likes.Guardar(usuario, nil)
}
