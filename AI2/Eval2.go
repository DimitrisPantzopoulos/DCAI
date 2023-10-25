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
	chess.Queen:  1000,
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

const PenaltyValue = -200 // You can adjust this value as needed

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
	WhiteScore := 0
	BlackScore := 0

	hashKey := HashPosition(game.Position())

	if entry, found := tt.Lookup(hashKey); found {
		if entry.Depth >= depth {
			return entry.Score
		}
	}

	WhiteMobility, BlackMobility := calculateMobility(board)

	// Iterate through the squares on the board
	for sq := chess.A1; sq <= chess.H8; sq++ {
		piece := board.Piece(sq)

		if piece != chess.NoPiece {
			// Get the value of the piece from the pieceValues map
			value := pieceValues[piece.Type()]

			// Use PST to find value
			if piece.Color() == chess.White {
				switch piece.Type() {
				case chess.Pawn:
					WhiteScore += (PAWN_TABLE[sq/8][sq%8] * -1) + value
				case chess.Knight:
					WhiteScore += (KNIGHT_TABLE[sq/8][sq%8] * -1) + value
				case chess.Bishop:
					WhiteScore += (BISHOPS_TABLE[sq/8][sq%8] * -1) + value
				case chess.Rook:
					WhiteScore += (ROOKS_TABLE[sq/8][sq%8] * -1) + value
				case chess.Queen:
					WhiteScore += (QUEENS_TABLE[sq/8][sq%8] * -1) + value
				case chess.King:
					WhiteScore += (KINGS_TABLE[sq/8][sq%8] * -1) + value
				}
			} else {
				switch piece.Type() {
				case chess.Pawn:
					BlackScore += (PAWN_TABLE[sq/8][sq%8]) + value
				case chess.Knight:
					BlackScore += (KNIGHT_TABLE[sq/8][sq%8]) + value
				case chess.Bishop:
					BlackScore += (BISHOPS_TABLE[sq/8][sq%8]) + value
				case chess.Rook:
					BlackScore += (ROOKS_TABLE[sq/8][sq%8]) + value
				case chess.Queen:
					BlackScore += (QUEENS_TABLE[sq/8][sq%8]) + value
				case chess.King:
					BlackScore += (KINGS_TABLE[sq/8][sq%8]) + value
				}
			}
		}
	}
	FinalScore := (WhiteScore + WhiteMobility) - (BlackScore + BlackMobility)

	tt.Store(hashKey, TranspositionTableEntry{
		HashKey:   hashKey,
		Depth:     depth,
		Score:     FinalScore,
		ScoreType: ExactScore,
	})

	return FinalScore
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
