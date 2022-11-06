package main

import (
	"algogram/errores"
	TDAHash "algogram/hash"
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func abrirArchivo(archivo string) *os.File {
	file, err := os.Open(archivo)
	if err != nil {
		newError := new(errores.ErrorLeerArchivo)
		fmt.Fprintln(os.Stdout, newError.Error())
		os.Exit(0)
	}
	return file
}

func guardarUsuarios(usuarios *os.File) TDAHash.Diccionario[string, Usuario] {
	listaUsuarios := TDAHash.CrearHash[string, Usuario]()
	id := 0
	scannerUsuarios := bufio.NewScanner(usuarios)
	for scannerUsuarios.Scan() {
		nombre := scannerUsuarios.Text()
		postsOrdenados := []*Post{}
		usuario := CrearUsuario(nombre, id, postsOrdenados)
		listaUsuarios.Guardar(nombre, usuario)
		id++
	}
	return listaUsuarios
}

func main() {
	var newError error
	var args = os.Args[1:]
	if len(args) != 1 {
		newError = new(errores.ErrorParametros)
		fmt.Fprintln(os.Stdout, newError.Error())
		os.Exit(0)
	}

	usuarios := abrirArchivo(args[0])
	listaUsuarios := guardarUsuarios(usuarios)
	listaPosts := CrearListaDePosts(func(p1, p2 Post) int { return p1.uid - p2.uid })

	var usuarioLoggeado Usuario

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {

		input := strings.Split(scanner.Text(), " ")
		action := input[0]

		switch action {

		case "login":
			if usuarioLoggeado != nil {
				newError = new(errores.UsuarioLoggeado)
				fmt.Fprintln(os.Stdout, newError.Error())
				break
			}

			usuario := input[1]
			if !listaUsuarios.Pertenece(usuario) {
				newError = new(errores.UsuarioInexistente)
				fmt.Fprintln(os.Stdout, newError.Error())
				break
			}
			usuarioLoggeado = listaUsuarios.Obtener(usuario)
			listaPosts.OrdenarPosts(usuarioLoggeado.Uid())
			fmt.Println("Hola", usuarioLoggeado.LeerUsuario())

		case "logout":
			if usuarioLoggeado == nil {
				newError = new(errores.NoLoggeado)
				fmt.Fprintln(os.Stdout, newError.Error())
				break
			}
			usuarioLoggeado = nil
			fmt.Println("Adios")

		case "publicar":
			if usuarioLoggeado == nil {
				newError = new(errores.NoLoggeado)
				fmt.Fprintln(os.Stdout, newError.Error())
				break
			}
			texto := strings.Join(input[1:], " ")
			post := usuarioLoggeado.CrearPost(texto)
			listaPosts.GuardarPost(post)
			fmt.Println("Post publicado")

		case "ver_siguiente_feed":
			if usuarioLoggeado == nil {
				newError = new(errores.NoLoggeadoONoHayPosts)
				fmt.Fprintln(os.Stdout, newError.Error())
				break
			}
			post, err := listaPosts.VerProximo(usuarioLoggeado.Uid())
			if err != nil {
				fmt.Fprintln(os.Stdout, err.Error())
				break
			}
			fmt.Println("Post ID", post.id)
			fmt.Println(post.usuario, "dijo:", post.texto)
			fmt.Println("Likes:", post.likes.Cantidad())

		case "likear_post":
			id, err := strconv.Atoi(input[1])
			if err != nil {
				newError = new(errores.ErrorParametros)
				fmt.Fprintln(os.Stdout, newError.Error())
				break
			}
			err = listaPosts.LikearPost(id, usuarioLoggeado.LeerUsuario())
			if err != nil {
				fmt.Fprintln(os.Stdout, err.Error())
				break
			}
			fmt.Println("Post likeado")

		case "mostrar_likes":
			id, err := strconv.Atoi(input[1])
			if err != nil {
				newError = new(errores.ErrorParametros)
				fmt.Fprintln(os.Stdout, newError.Error())
				break
			}
			likes, err := listaPosts.MostrarLikes(id)
			if err != nil {
				fmt.Fprintln(os.Stdout, err.Error())
				break
			}
			if len(likes) == 1 {
				fmt.Println("El post tiene", len(likes), "like:")
			} else {
				fmt.Println("El post tiene", len(likes), "likes:")
			}
			for _, usuario := range likes {
				fmt.Println("\t", usuario)
			}

		default:
			fmt.Fprintln(os.Stdout, "Comando incorrecto")

		}

	}
}
