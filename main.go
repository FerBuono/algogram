package main

import (
	"algogram/errores"
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

func guardarUsuarios(usuarios *os.File) ListaUsuarios {
	listaUsuarios := CrearListaDeUsuarios()
	id := 0
	scannerUsuarios := bufio.NewScanner(usuarios)
	for scannerUsuarios.Scan() {
		nombre := scannerUsuarios.Text()
		usuario := CrearUsuario(nombre, id)
		listaUsuarios.GuardarUsuario(usuario)
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
			usuario := strings.Join(input[1:], " ")
			usuarioLoggeado, newError = listaUsuarios.BuscarUsuario(usuario)
			if newError != nil {
				fmt.Fprintln(os.Stdout, newError.Error())
				break
			}
			fmt.Println("Hola", usuarioLoggeado.NombreUsuario())

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
			for iter := listaUsuarios.Iterador(); iter.HaySiguiente(); {
				_, usuario := iter.VerActual()
				usuario.AgregarPost(post)
				iter.Siguiente()
			}
			fmt.Println("Post publicado")

		case "ver_siguiente_feed":
			if usuarioLoggeado == nil {
				newError = new(errores.NoLoggeadoONoHayPosts)
				fmt.Fprintln(os.Stdout, newError.Error())
				break
			}
			post, newError := usuarioLoggeado.VerProximoPost()
			if newError != nil {
				fmt.Fprintln(os.Stdout, newError.Error())
				break
			}
			fmt.Println("Post ID", post.id)
			fmt.Println(post.usuario, "dijo:", post.texto)
			fmt.Println("Likes:", post.likes.Cantidad())

		case "likear_post":
			if usuarioLoggeado == nil {
				newError = new(errores.NoLoggeadoOPostInexistente)
				fmt.Fprintln(os.Stdout, newError.Error())
				break
			}
			id, err := strconv.Atoi(input[1])
			if err != nil {
				newError = new(errores.ErrorParametros)
				fmt.Fprintln(os.Stdout, newError.Error())
				break
			}
			newError = listaPosts.LikearPost(id, usuarioLoggeado.NombreUsuario())
			if newError != nil {
				fmt.Fprintln(os.Stdout, newError.Error())
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
			likes, newError := listaPosts.MostrarLikes(id)
			if newError != nil || len(likes) == 0 {
				fmt.Fprintln(os.Stdout, newError.Error())
				break
			}
			fmt.Println("El post tiene", len(likes), "likes:")
			for _, usuario := range likes {
				fmt.Printf("\t%s\n", usuario)
			}

		default:
			fmt.Fprintln(os.Stdout, "Comando incorrecto")

		}

	}
}
