import chess
import chess.pgn

# Initialize a chess board
board = chess.Board()

# List of moves from your input
move_list = ["e4", "e6", "d4", "d5", "Nc3", "Bb4", "Ne2", "dxe4", "e1d2", "d8d5", "a2a4", "d5e5", "d4e5", "b4c3", "d2c3", "b8d7", "d1d7", "e8d7", "c3d4", "e4e3", "f2e3", "a8b8", "a1b1", "b8a8", "b1a1", "a8b8", "a1b1", "b8a8", "a1b1", "a8b8", "b1a1", "a8b8", "h1g1", "a8b8", "b1a1", "a8b8", "h1g1", "a8b8", "b1a1", "a8b8", "h1g1", "a8b8", "b1a1", "a8b8", "h1g1", "a8b8", "b1a1", "a8b8", "a2a1", "g8f6", "c1d2", "h8d8", "d2c1", "a8b8", "e2c3", "d8e8", "c1d2", "e8g8", "f1d3", "a7a5", "a2a1", "c7c5", "d4c5", "g8f8", "d2c1", "f8d8", "d3e4", "f6d5", "c1d2", "d8f8", "a1a2", "d5f6", "g1a1", "f8d8", "a1g1", "d8f8", "g1a1", "f8d8", "a1g1", "d8f8", "g1a1", "f8d8", "a1g1", "d8f8", "g1a1", "f8d8", "a2a1", "d8f8", "a1a3", "f8d8", "g1a1", "d8f8", "a1g1", "f8d8", "g1a1", "d8f8", "a1g1", "f8d8", "g1a1", "d8f8", "a1g1", "f8d8", "a1a2", "f8d8", "d2c1", "f6h5", "e4g6", "h7g6", "c5d4", "d7e7", "c1d2", "h8h5", "a2a1", "h5h6", "a1c1", "g2e1", "c1a1", "h6h5", "a1c1", "e1c2", "c1b1", "h5h6", "b1c1", "h6h5", "c1b1", "h5h6", "b1c1", "h6h5", "c1b1", "h5h6", "b1c1", "h6h5", "c1b1", "h5h6", "b1c1", "h6h7", "c1b1", "h7h5", "b1c1", "h5h6", "c1b1", "h6h5", "b1c1", "h5h6", "c1b1", "h6h5", "b1c1", "h5h6", "c1b1", "h6h5", "b1c1", "h5h6", "c1b1", "h6h5", "b1c1", "h5h6", "b1c1", "h6h7", "c1b1", "h7h8", "c1b1", "h8e8", "d2c1", "c8d7", "a3a1", "b8c8"]

# Create a PGN game
pgn_game = chess.pgn.Game()

# Iterate through the moves and add them to the game
for move in move_list:
    move_obj = chess.Move.from_uci(move)
    pgn_game.add_main_variation(move_obj)

# Print the game in PGN format
print(pgn_game)
