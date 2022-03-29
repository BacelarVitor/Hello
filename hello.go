package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

const monitoramentos = 3
const delay = 5

func main() {

	exibeIntroducao()

	for {

		exibeMenu()

		comando := leComando()

		switch comando {
		case 1:
			iniciarMonitoramento()
		case 2:
			fmt.Println("Exibindo Logs...")
			imprimeLogs()
		case 0:
			fmt.Println("Saindo do programa")
			os.Exit(0)
		default:
			fmt.Println("Não conheço este comando")
			os.Exit(-1)
		}
	}

}

func exibeIntroducao() {

	nome := "Vitor"
	versao := 1.1
	fmt.Println("Olá sr(a).", nome)
	fmt.Println("Este programa está na versão", versao)
}

func exibeMenu() {

	fmt.Println("1- Iniciar Monitoramento")
	fmt.Println("2- Exibir Logs")
	fmt.Println("0- Sair do Programa")
}

func iniciarMonitoramento() {

	fmt.Println("Monitorando...")
	fmt.Println("")

	sites := leSitesDoArquivo()

	for i := 0; i < monitoramentos; i++ {

		for _, site := range sites {

			go testaSite(site)
		}

		time.Sleep(delay * time.Second)
		fmt.Println("")
	}

	fmt.Println("")
}

func leComando() int {

	var comandoLido int
	fmt.Scan(&comandoLido)
	fmt.Println("O valor da variável comando é:", comandoLido)

	return comandoLido
}

func testaSite(site string) {

	resp, err := http.Get(site)

	if err != nil {

		fmt.Println("Ocorreu um erro:", err, "func: testaSite")
	}

	if resp.StatusCode == 200 {

		fmt.Println("Site:", site, "foi carregado com sucesso!")
		registraLogs(site, true)
	} else {

		fmt.Println("Site:", site, "está com problemas. Status Code:", resp.StatusCode)
		registraLogs(site, false)
	}
}

func leSitesDoArquivo() []string {

	var sites []string

	arquivo, err := os.Open("sites.txt")

	if err != nil {

		fmt.Println("Ocorreu um erro:", err, "func: leSitesdoArquivo")
	}

	leitor := bufio.NewReader(arquivo)

	for {

		linha, err := leitor.ReadString('\n')
		linha = strings.TrimSpace(linha)

		sites = append(sites, linha)

		if err == io.EOF {

			break
		}
	}

	arquivo.Close()

	return sites
}

func registraLogs(site string, status bool) {

	arquivo, err := os.OpenFile("log.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)

	if err != nil {

		fmt.Println("Ocorreu um erro:", err, "func: registraLogs")
	}

	arquivo.WriteString(time.Now().Format("02/01/2006 15:04:05") + " - " + site + " - online: " + strconv.FormatBool(status) + "\n")

	arquivo.Close()
}

func imprimeLogs() {

	arquivo, err := ioutil.ReadFile("log.txt")

	if err != nil {

		fmt.Println("Ocorreu um erro:", err, "func: imprimeLogs")
	}

	fmt.Println(string(arquivo))
}
