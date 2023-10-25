package Search

import (
	"sync"

	"github.com/notnil/chess"
)

var pieceValues = map[chess.PieceType]int{
	chess.Pawn:   100,
	chess.Knight: 320,
	chess.Bishop: 330,
	chess.Rook:   500,
	chess.Queen:  900,
	chess.King:   0,
}

// NewTranspositionTable initializes a new transposition table.
func NewTranspositionTable(size int) *TranspositionTable {
	return &TranspositionTable{
		entries: make(map[uint64]TranspositionTableEntry, size),
	}
}

const (
	ExactScore int = iota
	LowerBound
	UpperBound

	EndgameFactor   = 250
	EarlyGameFactor = 100
)

// TranspositionTableEntry represents an entry in the transposition table.
type TranspositionTableEntry struct {
	HashKey   uint64
	Depth     int
	Score     int
	ScoreType int
	BestMove  *chess.Move // Add this field to store the best move
}

// TranspositionTable is a simple transposition table implementation.
type TranspositionTable struct {
	entries map[uint64]TranspositionTableEntry
	mutex   sync.Mutex
}

// Store stores an entry in the transposition table.
func (tt *TranspositionTable) Store(key uint64, entry TranspositionTableEntry) {
	tt.mutex.Lock()
	defer tt.mutex.Unlock()
	tt.entries[key] = entry
}

// Lookup looks up an entry in the transposition table.
func (tt *TranspositionTable) Lookup(key uint64) (TranspositionTableEntry, bool) {
	tt.mutex.Lock()
	defer tt.mutex.Unlock()
	entry, found := tt.entries[key]
	return entry, found
}

func ClearTranspositionTable(tt *TranspositionTable) {
	tt.mutex.Lock()
	defer tt.mutex.Unlock()
	tt.entries = make(map[uint64]TranspositionTableEntry)
}

func HashPosition(position *chess.Position) uint64 {
	var hashKey uint64

	// Iterate through the squares on the board
	for sq := chess.A1; sq <= chess.H8; sq++ {
		piece := position.Board().Piece(sq)

		if piece != chess.NoPiece {
			// Combine the piece type and square into a single value and XOR it to the hash key
			pieceValue := uint64(piece.Type()) << 3
			squareValue := uint64(sq)
			hashKey ^= pieceValue | squareValue
		}
	}

	return hashKey
}

var ATTACK_TABLE = [8][8]int{
	{0, 0, 0, 0, 0, 0, 0, 0},
	{0, 0, 0, 0, 0, 0, 0, 0},
	{0, 0, 0, 0, 0, 0, 0, 0},
	{0, 0, 0, 0, 0, 0, 0, 0},
	{0, 0, 0, 0, 0, 0, 0, 0},
	{0, 0, 0, 0, 0, 0, 0, 0},
	{0, 0, 0, 0, 0, 0, 0, 0},
	{0, 0, 0, 0, 0, 0, 0, 0},
}

var PAWN_TABLE = [8][8]int{
	{0, 0, 0, 0, 0, 0, 0, 0},
	{50, 50, 50, 50, 50, 50, 50, 50},
	{10, 10, 20, 30, 30, 20, 10, 10},
	{5, 5, 10, 25, 25, 10, 5, 5},
	{0, 0, 0, 20, 20, 0, 0, 0},
	{5, -5, -10, 0, 0, -10, -5, 5},
	{5, 10, 10, -20, -20, 10, 10, 5},
	{0, 0, 0, 0, 0, 0, 0, 0},
}

// KNIGHT_TABLE
var KNIGHT_TABLE = [8][8]int{
	{-50, -40, -30, -30, -30, -30, -40, -50},
	{-40, -20, 0, 0, 0, 0, -20, -40},
	{-30, 0, 10, 15, 15, 10, 0, -30},
	{-30, 5, 15, 20, 20, 15, 5, -30},
	{-30, 0, 15, 20, 20, 15, 0, -30},
	{-30, 5, 10, 15, 15, 10, 5, -30},
	{-40, -20, 0, 5, 5, 0, -20, -40},
	{-50, -40, -30, -30, -30, -30, -40, -50},
}

// BISHOPS_TABLE
var BISHOPS_TABLE = [8][8]int{
	{-20, -10, -10, -10, -10, -10, -10, -20},
	{-10, 0, 0, 0, 0, 0, 0, -10},
	{-10, 0, 5, 10, 10, 5, 0, -10},
	{-10, 5, 5, 10, 10, 5, 5, -10},
	{-10, 0, 10, 10, 10, 10, 0, -10},
	{-10, 10, 10, 10, 10, 10, 10, -10},
	{-10, 5, 0, 0, 0, 0, 5, -10},
	{-20, -10, -10, -10, -10, -10, -10, -20},
}

// ROOKS_TABLE
var ROOKS_TABLE = [8][8]int{
	{0, 0, 0, 0, 0, 0, 0, 0},
	{5, 10, 10, 10, 10, 10, 10, 5},
	{-5, 0, 0, 0, 0, 0, 0, -5},
	{-5, 0, 0, 0, 0, 0, 0, -5},
	{-5, 0, 0, 0, 0, 0, 0, -5},
	{-5, 0, 0, 0, 0, 0, 0, -5},
	{-5, 0, 0, 0, 0, 0, 0, -5},
	{0, 0, 0, 5, 5, 0, 0, 0},
}

// QUEENS_TABLE
var QUEENS_TABLE = [8][8]int{
	{-20, -10, -10, -5, -5, -10, -10, -20},
	{-10, 0, 0, 0, 0, 0, 0, -10},
	{-10, 0, 5, 5, 5, 5, 0, -10},
	{-5, 0, 5, 5, 5, 5, 0, -5},
	{0, 0, 5, 5, 5, 5, 0, -5},
	{-10, 5, 5, 5, 5, 5, 0, -10},
	{-10, 0, 5, 0, 0, 0, 0, -10},
	{-20, -10, -10, -5, -5, -10, -10, -20},
}

// KINGS_TABLE
var KINGS_TABLE = [8][8]int{
	{-30, -40, -40, -50, -50, -40, -40, -30},
	{-30, -40, -40, -50, -50, -40, -40, -30},
	{-30, -40, -40, -50, -50, -40, -40, -30},
	{-30, -40, -40, -50, -50, -40, -40, -30},
	{-20, -30, -30, -40, -40, -30, -30, -20},
	{-10, -20, -20, -20, -20, -20, -20, -10},
	{20, 20, 0, 0, 0, 0, 20, 20},
	{20, 30, 10, 0, 0, 10, 30, 20},
}

func Eval(board *chess.Board, game *chess.Game, tt *TranspositionTable, depth int, movesPlayed []string) int {
	NoOfPieces := 0

	WhiteScore := 0
	BlackScore := 0

	hashKey := HashPosition(game.Position())

	if entry, found := tt.Lookup(hashKey); found {
		if entry.Depth >= depth {
			return entry.Score
		}
	}

	// Iterate through the squares on the board
	for sq := chess.A1; sq <= chess.H8; sq++ {
		piece := board.Piece(sq)

		if piece != chess.NoPiece {
			// Get the value of the piece from the pieceValues map
			value := pieceValues[piece.Type()]
			NoOfPieces++
			// Use PST to find value
			if piece.Color() == chess.Black {
				switch piece.Type() {
				case chess.Pawn:
					if NoOfPieces <= 10 {
						WhiteScore += (PAWN_TABLE[sq/8][sq%8] * -1) + value + EndgameFactor
					} else {
						WhiteScore += (PAWN_TABLE[sq/8][sq%8] * -1) + value
					}
				case chess.Knight:
					WhiteScore += (KNIGHT_TABLE[sq/8][sq%8] * -1) + value
				case chess.Bishop:
					WhiteScore += (BISHOPS_TABLE[sq/8][sq%8] * -1) + value
				case chess.Rook:
					WhiteScore += (ROOKS_TABLE[sq/8][sq%8] * -1) + value
				case chess.Queen:
					WhiteScore += (QUEENS_TABLE[sq/8][sq%8] * -1) + value
				}
			} else {
				switch piece.Type() {
				case chess.Pawn:
					if NoOfPieces <= 10 {
						BlackScore += (PAWN_TABLE[sq/8][sq%8]) + value + EndgameFactor
					} else {
						BlackScore += (PAWN_TABLE[sq/8][sq%8]) + value
					}
				case chess.Knight:
					BlackScore += (KNIGHT_TABLE[sq/8][sq%8]) + value
				case chess.Bishop:
					BlackScore += (BISHOPS_TABLE[sq/8][sq%8]) + value
				case chess.Rook:
					BlackScore += (ROOKS_TABLE[sq/8][sq%8]) + value
				case chess.Queen:
					BlackScore += (QUEENS_TABLE[sq/8][sq%8]) + value
					// case chess.King:
					// 	BlackScore += (KINGS_TABLE[sq/8][sq%8]) + value
				}
			}
		}
	}

	WhiteMobility, BlackMobility := calculateMobility(board)
	safety := CalculateCaptures(game)
	FinalScore := (WhiteScore + WhiteMobility) - (BlackScore + BlackMobility) + safety

	tt.Store(hashKey, TranspositionTableEntry{
		HashKey:   hashKey,
		Depth:     depth,
		Score:     FinalScore,
		ScoreType: ExactScore,
	})

	return FinalScore
}

func findCaptures(game *chess.Game) []*chess.Move {
	ValMoves := game.ValidMoves()
	var captures []*chess.Move
	for _, move := range ValMoves {
		if move.HasTag(chess.Capture) {
			captures = append(captures, move)
		}
	}
	return captures
}

func calculateMobility(board *chess.Board) (int, int) {
	WhiteMobility := 0
	BlackMobility := 0

	for sq := chess.A1; sq <= chess.H8; sq++ {
		piece := board.Piece(sq)

		if piece != chess.NoPiece {
			mobility := calculatePieceMobility(board, sq)
			if piece.Color() == chess.White {
				WhiteMobility += mobility
			} else {
				BlackMobility += mobility
			}
		}
	}

	return WhiteMobility, BlackMobility
}

func calculatePieceMobility(board *chess.Board, square chess.Square) int {
	mobility := 0

	piece := board.Piece(square)
	color := piece.Color()

	switch piece.Type() {
	case chess.Pawn:
		mobility = calculatePawnMobility(board, square, color)
	case chess.Knight:
		mobility = calculateKnightMobility(board, square)
	case chess.Bishop:
		mobility = calculateBishopMobility(board, square)
	case chess.Rook:
		mobility = calculateRookMobility(board, square)
	case chess.Queen:
		mobility = calculateQueenMobility(board, square)
	case chess.King:
		mobility = calculateKingMobility(board, square)
	}

	return mobility
}

func calculatePawnMobility(board *chess.Board, square chess.Square, color chess.Color) int {
	mobility := 0

	// Define the squares that a pawn can move to based on its color
	var targetSquares []chess.Square
	if color == chess.White {
		// For White pawns
		targetSquares = []chess.Square{square + 8, square + 16}
	} else {
		// For Black pawns
		targetSquares = []chess.Square{square - 8, square - 16}
	}

	// Check if the target squares are valid and unoccupied
	for _, targetSquare := range targetSquares {
		if targetSquare >= chess.A1 && targetSquare <= chess.H8 && board.Piece(targetSquare) == chess.NoPiece {
			mobility++
		}
	}

	return mobility
}

func calculateKnightMobility(board *chess.Board, square chess.Square) int {
	mobility := 0

	// Define the possible knight moves
	possibleMoves := []chess.Square{
		square + 15, square + 17,
		square + 10, square + 6,
		square - 15, square - 17,
		square - 10, square - 6,
	}

	// Check if the possible moves are within the board bounds and unoccupied
	for _, move := range possibleMoves {
		if move >= chess.A1 && move <= chess.H8 && board.Piece(move) == chess.NoPiece {
			mobility++
		}
	}

	return mobility
}

// Function to calculate mobility for a Rook
func calculateRookMobility(board *chess.Board, square chess.Square) int {
	mobility := 0

	// Define the possible moves for a rook (up, down, left, and right)
	moves := [4]chess.Square{chess.D1, chess.D2, chess.D3, chess.D4}

	for _, move := range moves {
		targetSquare := square + move
		for targetSquare >= chess.A1 && targetSquare <= chess.H8 {
			if board.Piece(targetSquare) == chess.NoPiece {
				mobility++
			} else {
				break
			}
			targetSquare += move
		}
	}

	return mobility
}

// Function to calculate mobility for a Queen (combining Rook and Bishop mobility)
func calculateQueenMobility(board *chess.Board, square chess.Square) int {
	mobility := calculateRookMobility(board, square) + calculateBishopMobility(board, square)

	return mobility
}

// Function to calculate mobility for a King
func calculateKingMobility(board *chess.Board, square chess.Square) int {
	mobility := 0

	// Define the possible moves for a king
	possibleMoves := [8]chess.Square{
		square + 1, square - 1,
		square + 7, square - 7,
		square + 8, square - 8,
		square + 9, square - 9,
	}

	// Check if the possible moves are within the board bounds and unoccupied
	for _, move := range possibleMoves {
		if move >= chess.A1 && move <= chess.H8 && board.Piece(move) == chess.NoPiece {
			mobility++
		}
	}

	return mobility
}

func calculateBishopMobility(board *chess.Board, square chess.Square) int {
	mobility := 0

	// Define the possible diagonal moves for a bishop
	moves := [4]chess.Square{chess.D1, chess.D3, chess.D5, chess.D7}

	for _, move := range moves {
		targetSquare := square + move
		for targetSquare >= chess.A1 && targetSquare <= chess.H8 {
			if board.Piece(targetSquare) == chess.NoPiece {
				mobility++
			} else {
				break
			}
			targetSquare += move
		}
	}

	return mobility
}

func CalculateCaptures(game *chess.Game) int {
	whiteSafety := 0
	blackSafety := 0
	EvaluateCapture := 0

	attackingPieces := make(map[chess.Color][]*chess.Piece)
	captures := findCaptures(game)

	for _, capture := range captures {
		fromSquare := capture.S1()
		targetSquare := capture.S2()

		capturingPiece := game.Position().Board().Piece(fromSquare)
		capturedPiece := game.Position().Board().Piece(targetSquare)

		// Handle captured pieces based on their type and color
		if capturedPiece.Color() == game.Position().Turn() {
			// Piece of the current player's color was captured
			if capturedPiece.Color() == chess.Black {
				blackSafety -= pieceValues[capturedPiece.Type()]
			} else {
				whiteSafety -= pieceValues[capturedPiece.Type()]
			}

			if !IsPieceProtected(game, *capture, targetSquare) {
				if capturedPiece.Type() == chess.Queen {
					EvaluateCapture += 900
				} else {
					EvaluateCapture += pieceValues[capturedPiece.Type()] - pieceValues[capturingPiece.Type()]
				}
			}
			whiteSafety += IsFriendlyProtected(game, &capturingPiece, &capturedPiece, *capture, targetSquare)
		} else {
			// Piece of the opposing player's color was captured
			if capturedPiece.Color() == chess.Black {
				blackSafety -= pieceValues[capturedPiece.Type()]
			} else {
				whiteSafety -= pieceValues[capturedPiece.Type()]
			}

			EvaluateCapture += pieceValues[capturedPiece.Type()]

			if IsPieceProtected(game, *capture, targetSquare) {
				EvaluateCapture += pieceValues[capturedPiece.Type()]
			}

			// Handle protection and multiple attackers
			whiteSafety += IsFriendlyProtected(game, &capturingPiece, &capturedPiece, *capture, targetSquare)
			attackingPieces[chess.White] = append(attackingPieces[chess.White], &capturingPiece)
		}
	}

	// Handle safety due to multiple attackers
	if len(attackingPieces[chess.White]) > 1 {
		whiteSafety += (len(attackingPieces[chess.White]) - 1) * 50
	}
	if len(attackingPieces[chess.Black]) > 1 {
		whiteSafety -= (len(attackingPieces[chess.Black]) - 1) * 50
	}

	return whiteSafety - blackSafety + EvaluateCapture
}

func IsPieceProtected(game *chess.Game, move chess.Move, targetSquare chess.Square) bool {
	// Clone the game to simulate the move.
	i := 0
	copyGame := game.Clone()
	copyGame.Move(&move)

	// Find all captures after the move.
	captures := findCaptures(copyGame)

	// Check if the target square appears in the list of captures.
	for _, capture := range captures {
		// Access row and col using the targetSquare methods.
		ATTACK_TABLE[targetSquare.Rank()][targetSquare.File()] -= 200
		if capture.S2() == targetSquare {
			i += 1
		}
	}
	return i > 0
}

func IsFriendlyProtected(game *chess.Game, capturingPiece *chess.Piece, targetPiece *chess.Piece, move chess.Move, targetSquare chess.Square) int {
	if IsPieceProtected(game, move, targetSquare) {
		if capturingPiece.Color() == targetPiece.Color() {
			return pieceValues[targetPiece.Type()] - pieceValues[capturingPiece.Type()]
		}
		return -pieceValues[targetPiece.Type()] + pieceValues[capturingPiece.Type()]
	} else {
		if targetPiece.Type() == chess.Queen {
			return 0
		}
		return -pieceValues[targetPiece.Type()]
	}
}

func IsAttackingPieceProtected(game *chess.Game, move chess.Move) int {
	fromSquare := move.S1()
	targetSquare := move.S2()

	attackingPiece := game.Position().Board().Piece(fromSquare)
	targetedPiece := game.Position().Board().Piece(targetSquare)
	if IsPieceProtected(game, move, targetSquare) {
		return -pieceValues[attackingPiece.Type()] + pieceValues[targetedPiece.Type()]
	}
	// Check if the piece on the target square is protected
	return pieceValues[targetedPiece.Type()]
}
