package Search

import (
	"sort"
	"time"

	"github.com/notnil/chess"
)

var nodesVisited int
var BestMove *chess.Move

func QuiescenceSearch(alpha, beta int, game *chess.Game, tt *TranspositionTable, depth int, movesPlayed []string, moves []*chess.Move) int {
	// Calculate the stand-pat score based on your current evaluation function
	standPatScore := Eval(game.Position().Board(), game, tt, depth, movesPlayed)
	// Compare the stand-pat score with beta
	if standPatScore >= beta {
		return beta
	}

	// If the stand-pat score is greater than alpha, update alpha
	if alpha < standPatScore {
		alpha = standPatScore
	}

	// Check all moves to see if they are captures or checks
	ValMoves := game.ValidMoves()

	for _, move := range ValMoves {
		// Check if the move is a capture or check
		isCapture := move.HasTag(chess.Capture)
		isCheck := move.HasTag(chess.Check)

		// If it's a capture or check, apply the QFilterMoves function to assess its value
		if isCapture || isCheck {
			moveValue := QFilterMove(move, game)

			if moveValue == -1 {
				continue
			}

			Copy := game.Clone()
			Copy.Move(move)
			score := -QuiescenceSearch(-beta, -alpha, Copy, tt, depth-1, movesPlayed, moves)
			score += valueOfPieceThreatened(move, game)

			if score >= beta {
				return beta
			}
			if score > alpha {
				alpha = score
			}
		}
	}

	return alpha
}

// Function for the alpha-beta search
func NegaMaxAlphabeta(game *chess.Game, depth, alpha, beta int, MaximisingPlayer bool, tt *TranspositionTable, movesPlayed []string) (int, *chess.Move, int) {
	nodesVisited++

	hashKey := HashPosition(game.Position())
	entry, found := tt.Lookup(hashKey)

	if found && entry.Depth >= depth {
		if entry.ScoreType == ExactScore {
			return entry.Score, entry.BestMove, 1
		}
		if entry.ScoreType == LowerBound && entry.Score >= beta {
			return entry.Score, entry.BestMove, 1
		}
		if entry.ScoreType == UpperBound && entry.Score <= alpha {
			return entry.Score, entry.BestMove, 1
		}
	}

	ValMoves := game.ValidMoves()
	OrderedMoves := OrderMoves(game, ValMoves)

	if game.Outcome() != chess.NoOutcome || depth == 0 {
		return QuiescenceSearch(alpha, beta, game, tt, depth-1, movesPlayed, OrderedMoves), nil, 1
	}

	visitedNodes := 1
	var BestMove *chess.Move

	if MaximisingPlayer {
		MaxEval := -9999

		for _, move := range OrderedMoves {
			Copy := game.Clone()
			Copy.Move(move)
			eval, _, visited := NegaMaxAlphabeta(Copy, depth-1, alpha, beta, false, tt, movesPlayed)

			visitedNodes += visited

			if eval > MaxEval {
				MaxEval = eval
				BestMove = move
			}

			alpha = max(alpha, eval)
			if beta <= alpha {
				break
			}
		}

		var scoreType int
		if MaxEval <= alpha {
			scoreType = UpperBound
		} else if MaxEval >= beta {
			scoreType = LowerBound
		} else {
			scoreType = ExactScore
		}

		tt.Store(hashKey, TranspositionTableEntry{
			HashKey:   hashKey,
			Depth:     depth,
			Score:     MaxEval,
			ScoreType: scoreType,
			BestMove:  BestMove,
		})

		return MaxEval, BestMove, visitedNodes
	} else {
		MinEval := 9999

		for _, move := range OrderedMoves {
			Copy := game.Clone()
			Copy.Move(move)
			eval, _, visited := NegaMaxAlphabeta(Copy, depth-1, alpha, beta, true, tt, movesPlayed)
			visitedNodes += visited

			if eval < MinEval {
				MinEval = eval
				BestMove = move
			}

			beta = min(beta, eval)
			if beta <= alpha {
				break
			}
		}

		var scoreType int
		if MinEval <= alpha {
			scoreType = UpperBound
		} else if MinEval >= beta {
			scoreType = LowerBound
		} else {
			scoreType = ExactScore
		}

		tt.Store(hashKey, TranspositionTableEntry{
			HashKey:   hashKey,
			Depth:     depth,
			Score:     MinEval,
			ScoreType: scoreType,
			BestMove:  BestMove,
		})

		return MinEval, BestMove, visitedNodes
	}
}

func IterativeDeepening(game *chess.Game, maxDepth int, tt *TranspositionTable, movesPlayed []string, timeout time.Duration) (int, *chess.Move, int) {
	nodesVisited = 0
	bestMove := (*chess.Move)(nil)
	bestScore := -9999
	startTime := time.Now()

	for depth := 1; depth <= maxDepth; depth++ {
		elapsedTime := time.Since(startTime)
		if elapsedTime >= timeout {
			break
		}

		score, move, visited := NegaMaxAlphabeta(game, depth, -9999, 9999, true, tt, movesPlayed)
		nodesVisited += visited
		if move != nil {
			bestMove = move
			bestScore = score
		}

		elapsedTime = time.Since(startTime)
		if elapsedTime >= timeout {
			break
		}
	}

	return bestScore, bestMove, nodesVisited
}

func OrderMoves(game *chess.Game, moves []*chess.Move) []*chess.Move {
	sort.SliceStable(moves, func(i, j int) bool {
		priorityI := movePriority(game, moves[i])
		priorityJ := movePriority(game, moves[j])

		// Compare moves based on their priority
		if priorityI != priorityJ {
			return priorityI > priorityJ
		}
		return false
	})

	return moves
}

func movePriority(game *chess.Game, move *chess.Move) int {
	if game.Method() == chess.Checkmate {
		return 4
	}
	if move.Promo() != chess.NoPieceType {
		return 3
	}
	if move.HasTag(chess.QueenSideCastle) || move.HasTag(chess.KingSideCastle) {
		return 1
	}
	if move.HasTag(chess.Capture) {
		// Evaluate the capture and capture value
		capturedValue := QFilterMove(move, game)

		if capturedValue > 0 {
			// Positive values mean it's a good trade
			return 2
		} else if capturedValue < 0 {
			// Negative values mean it's a bad trade
			return -1
		}
	}
	srcSquare := move.S1()
	piece := game.Position().Board().Piece(srcSquare)
	if piece.Type() == chess.King {

		return -1
	}

	Copy := game.Clone()
	Copy.Move(move)
	if Copy.Outcome() == chess.Draw {
		return -2
	}
	return 0
}

func QFilterMove(move *chess.Move, game *chess.Game) int {
	fromSquare := move.S1()
	toSquare := move.S2()
	movingPiece := game.Position().Board().Piece(fromSquare)
	targetPiece := game.Position().Board().Piece(toSquare)

	movingPieceValue := pieceValues[movingPiece.Type()]
	targetPieceValue := pieceValues[targetPiece.Type()]

	if movingPieceValue <= targetPieceValue {
		return 1
	}
	return 0

}

func valueOfPieceThreatened(move *chess.Move, game *chess.Game) int {
	fromSquare := move.S1()
	toSquare := move.S2()
	ThreateningPiece := game.Position().Board().Piece(fromSquare)
	ThreatenedPiece := game.Position().Board().Piece(toSquare)

	movingPieceValue := pieceValues[ThreateningPiece.Type()]
	targetPieceValue := pieceValues[ThreatenedPiece.Type()]

	if (movingPieceValue > 100) && (targetPieceValue > 100) {
		return PenaltyValue
	}

	return 0
}
