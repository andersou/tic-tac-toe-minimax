package main

import "fmt"

type Jogador int8

const (
	Jogador1 Jogador = 1
	Jogador2 Jogador = 2
	NaoJogou Jogador = 0
)

type Estado struct {
	matriz       [3][3]Jogador
	jogadorAtual Jogador
	minimax      int8
	terminou     bool
}

func (e *Estado) FimDeJogo() (bool, Jogador) {
	m := e.matriz
	for i := 0; i < 3; i++ {
		//verifica linha
		for j := Jogador1; j <= Jogador2; j++ {
			//linha i
			terminouLinha := m[i][0] == j && m[i][1] == j && m[i][2] == j

			//coluna i
			terminouColuna := m[0][i] == j && m[1][i] == j && m[2][i] == j
			if terminouLinha || terminouColuna {
				return true, j
			}
		}
	}
	//verificar diagonais
	for j := Jogador1; j <= Jogador2; j++ {
		//diagonal 1
		diagonal1 := m[0][0] == j && m[1][1] == j && m[2][2] == j

		//diagonal2
		diagonal2 := m[0][2] == j && m[1][1] == j && m[2][0] == j
		if diagonal1 || diagonal2 {
			return true, j
		}
	}
	//se ainda há possibilidade
	for _, lin := range e.matriz {
		for _, cel := range lin {
			if cel == NaoJogou {
				return false, NaoJogou
			}
		}
	}
	//empate
	return true, NaoJogou
}

type Nodo struct {
	filhos   []*Nodo
	conteudo interface{}
}

func jogarJogoDaVelha(e *Estado) (possiveisEstados []*Estado) {
	proxJogador := Jogador2
	if e.jogadorAtual == Jogador2 {
		proxJogador = Jogador1
	}
	for l, lin := range e.matriz {
		for c, col := range lin {
			if col == NaoJogou {
				novoEstado := *e
				novoEstado.jogadorAtual = proxJogador
				novoEstado.matriz[l][c] = e.jogadorAtual
				possiveisEstados = append(possiveisEstados, &novoEstado)
			}
		}
	}
	return
}

func nodoFromEstados(estados []*Estado) (nodos []*Nodo) {
	for _, e := range estados {
		n := &Nodo{conteudo: e}
		nodos = append(nodos, n)
	}
	return
}

func constroiArvore(n *Nodo) (placar [3]int) {
	possiveisEstados := jogarJogoDaVelha(n.conteudo.(*Estado))
	n.filhos = nodoFromEstados(possiveisEstados)
	for _, filho := range n.filhos {
		estado := filho.conteudo.(*Estado)
		if terminou, quemVence := estado.FimDeJogo(); !terminou {
			placarFilho := constroiArvore(filho)
			for i := range placarFilho {
				placar[i] += placarFilho[i]
			}
		} else {
			estado.terminou = true
			if quemVence == Jogador1 {
				estado.minimax = 1
			} else if quemVence == Jogador2 {
				estado.minimax = -1
			}
			placar[quemVence]++
		}
	}
	return
}

func qtdeNodos(n *Nodo) int {
	qtde := 0
	if n.filhos != nil {
		qtde += len(n.filhos)

		for _, filho := range n.filhos {
			qtde += qtdeNodos(filho)
		}
	}
	return qtde
}

func calculaMinimax(nodo *Nodo) int8 {
	if nodo.filhos == nil {
		return nodo.conteudo.(*Estado).minimax
	}
	jogadorAtual := nodo.conteudo.(*Estado).jogadorAtual
	var valores []int8
	for _, filho := range nodo.filhos {
		estado := filho.conteudo.(*Estado)
		estado.minimax = calculaMinimax(filho)
		valores = append(valores, estado.minimax)
	}

	minimax := valores[0]
	for _, val := range valores {
		if jogadorAtual == Jogador1 {
			if val > minimax {
				minimax = val
			}
		} else {
			if jogadorAtual == Jogador2 {
				if val < minimax {
					minimax = val
				}
			}
		}
	}
	return minimax
}
func main() {
	//primeiro nodo
	n := Nodo{}
	n.conteudo = &Estado{jogadorAtual: Jogador1}
	placar := constroiArvore(&n)
	fmt.Printf("Quantidade de estados totais \t%d \nQuantidade de Empates \t\t%d\nQuantidade de Vitória J1 \t%d \nQuantidade de Vitória J2 \t%d \n", qtdeNodos(&n), placar[0], placar[1], placar[2])
	fmt.Printf("Valor Minimax  \t\t\t%d\n", calculaMinimax(&n))

}
