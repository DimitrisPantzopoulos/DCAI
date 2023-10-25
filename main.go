package main

import (
	"fmt"
	"time"

	S "DCAI.com/packages/AI"
	V "DCAI.com/packages/AI2"
	"DCAI.com/packages/util"
	"github.com/notnil/chess"
)

func main() {
	originalGame := chess.NewGame()
	movesPlayed := []string{}

	InitialCap := (64 * 1048576) / 16
	tt := S.NewTranspositionTable(InitialCap)
	tt2 := V.NewTranspositionTable(InitialCap)

	match := true
	match2 := true

	for i := 0; i < 100; i++ {
		if originalGame.Outcome() == chess.NoOutcome {
			if i <= 7 && match {
				theoryName, nextMove := util.AIOpeningBook(movesPlayed)
				fmt.Println("Theory Move:", theoryName)
				fmt.Println("Next Move:", nextMove)

				if nextMove == "Still not found :(" {
					match = false
					score, bestMove, visitedNodes := S.IterativeDeepening(originalGame, 4, tt, movesPlayed, time.Second)
					moveStr := bestMove.String()
					movesPlayed = append(movesPlayed, moveStr)
					originalGame.Move(bestMove)
					fmt.Println("Move Safety Negamax || Score", score, "||Best Move", bestMove, "||visitedNodes:", visitedNodes)
					fmt.Println("Move Safety Negamax Current game position:")
				} else {
					originalGame.MoveStr(nextMove)
					movesPlayed = append(movesPlayed, nextMove)
				}
			} else {
				score, bestMove, visitedNodes := S.IterativeDeepening(originalGame, 4, tt, movesPlayed, time.Second)
				fmt.Println("Move Safety Negamax || Score", score, "||Best Move", bestMove, "||visitedNodes:", visitedNodes)
				originalGame.Move(bestMove)
				fmt.Println("Move Safety Negamax Current game position:")
			}

			fmt.Println(originalGame.Position().Board().Draw())

			if i <= 7 && match2 {
				theoryName, nextMove := util.AIOpeningBook(movesPlayed)
				fmt.Println("Theory Move:", theoryName)
				fmt.Println("Next Move:", nextMove)

				if nextMove == "Still not found :(" {
					match2 = false
					Vscore, VbestMove, VvisitedNodes := V.IterativeDeepening(originalGame, 4, tt2, movesPlayed, time.Second)
					moveStr := VbestMove.String()
					movesPlayed = append(movesPlayed, moveStr)

					movesPlayed = append(movesPlayed, VbestMove.String())
					fmt.Println("Negamax || Score", Vscore, "||Best Move", VbestMove, "||visitedNodes:", VvisitedNodes)
					originalGame.Move(VbestMove)
					fmt.Println("Negamax Current game position:")
				} else {
					originalGame.MoveStr(nextMove)
					movesPlayed = append(movesPlayed, nextMove)
				}
			} else {
				Vscore, VbestMove, VvisitedNodes := V.IterativeDeepening(originalGame, 4, tt2, movesPlayed, time.Second)
				fmt.Println("Negamax || Score", Vscore, "||Best Move", VbestMove, "||visitedNodes:", VvisitedNodes)
				originalGame.Move(VbestMove)
				fmt.Println("Negamax Current game position:")
			}
			fmt.Println(originalGame.Position().Board().Draw())
		} else {
			fmt.Println(originalGame.Outcome())
			break
		}
	}

	fmt.Printf("Game completed. %s by %s.\n", originalGame.Outcome(), originalGame.Method())
	fmt.Println(originalGame.String())
}
