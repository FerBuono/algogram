package main

import (
	"algogram/app"
	"algogram/errores"
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

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
	listaPosts := app.CrearListaDePosts()

	var usuarioLoggeado app.Usuario

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
			fmt.Println("Hola", usuario)

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
			uid, usuario := usuarioLoggeado.VerUsuario()
			post := listaPosts.GuardarPost(texto, uid, usuario)
			listaUsuarios.GuardarPost(post)
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
			texto, id, usuario, _ := post.VerPost()
			fmt.Println("Post ID", id)
			fmt.Println(usuario, "dijo:", texto)
			fmt.Println("Likes:", post.Likes())

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
			_, usuario := usuarioLoggeado.VerUsuario()
			newError = listaPosts.LikearPost(id, usuario)
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
