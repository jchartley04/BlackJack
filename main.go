package main

import (
	"fmt" //Formatting input / output operations

	"strings" //String manipulation

	"github.com/santosh/gophercises/deck" //Imports deck package to handle cards and decks of cards
)

type Hand []deck.Card //Represents a player or dealer's hand, which is represented by a slice of cards from the deck package

func (h Hand) String() string {
	strs := make([]string, len(h)) //Creates a slice of strings the same length as the hand (number of cards)
	for i := range h {             //Loops through the cards currently in the hand
		strs[i] = h[i].String() //Converts each card to its string representation. Ex: King of Spades, Two of Hearts, etc.
	}
	return strings.Join(strs, ", ") //Outputs the cards in the hand
}

func (h Hand) DealerString() string {
	return h[0].String() + ", **HIDDEN**" //Returns the first card as a string, hides the dealers second card
}

func (h Hand) Score() int {
	minScore := h.MinScore() //Calculates the min score of the hand (Aces counted as 1)
	if minScore > 11 {       //If the current min score is more than 11, return the minScore
		return minScore
	}
	for _, c := range h { //Loops over all the cards in the deck
		if c.Rank == deck.Ace {
			// ace is currently worth 1, and we are changing it to be worth 11
			// 11 - 1 = 10
			return minScore + 10
		}
	}
	return minScore
}

func (h Hand) MinScore() int {
	score := 0
	for _, c := range h { //Loops through each card in the player or dealers hand
		score += min(int(c.Rank), 10) //Converts Jack, Q, and K cards to 10.
	}
	return score //Returns the minimum score in the hand
}

func min(a, b int) int { //Min function used in MinScore() func
	if a < b {
		return a
	}
	return b
} //returns lowest of 2 int

func Shuffle(gs GameState) GameState {
	ret := clone(gs)                                //Clones the current gameState
	ret.Deck = deck.New(deck.Deck(3), deck.Shuffle) //Creates a new deck with 3 sets of cards. Shuffles this 3 set deck.
	return ret                                      //Returns the shuffled deck
}

func Deal(gs GameState) GameState {
	ret := clone(gs)              //Clones the current game state
	ret.Player = make(Hand, 0, 5) //Players hand capacity is set to 5 cards
	ret.Dealer = make(Hand, 0, 5) //Dealers hand capacity is set to 5 cards
	var card deck.Card            //card var holds each drawn card
	for i := 0; i < 2; i++ {      //Deals two cards to both the player and the dealer
		card, ret.Deck = draw(ret.Deck)       //Draws first card from the deck
		ret.Player = append(ret.Player, card) //Gives player the first card in the deck
		card, ret.Deck = draw(ret.Deck)       //Draws second card
		ret.Dealer = append(ret.Dealer, card) //Deals second card to the dealer
	}
	ret.State = StatePlayerTurn //Sets the game state to indicate its the players turn to act
	return ret                  //Returns the updated gameS state with the dealt cards and updated state
}

func Stand(gs GameState) GameState {
	ret := clone(gs) //Clones current game state
	ret.State++      //Advances the game state (player to dealer, dealer to hand over) (SEEN BELOW)
	return ret       //Returns the updated gameState
}

func Hit(gs GameState) GameState {
	ret := clone(gs)            //Clones current gameState
	hand := ret.CurrentPlayer() //Gets current player's hand
	//This can change depending on the game state. Either the player or the dealer
	var card deck.Card              //card holds value of each drawn card
	card, ret.Deck = draw(ret.Deck) //Draws card and updates the cards in the cloned state
	*hand = append(*hand, card)     //Adds the drawn card to the current players/dealers hand
	if hand.Score() > 21 {
		return Stand(ret) //Automatically stand if the score is above 21. The player or dealer busted
	}
	return ret //Retuyrns the updated game state
}

func EndHand(gs GameState) GameState {
	ret := clone(gs)                                         //Clone current game state
	pScore, dScore := ret.Player.Score(), ret.Dealer.Score() //Calculates the scores for both player and dealer
	fmt.Println("==FINAL HANDS==")                           //Print out final hand in output
	fmt.Println("Player:", ret.Player, "\nScore:", pScore)   //Print player score
	fmt.Println("Dealer:", ret.Dealer, "\nScore:", dScore)   //Print dealer score
	switch {
	case pScore > 21: //Score over 21, you lose
		fmt.Println("You busted")
	case dScore > 21: //Dealer score over 21, dealer loses
		fmt.Println("Dealer busted")
	case pScore > dScore: //You win if your score is greater than the dealer's score
		fmt.Println("You win!")
	case dScore > pScore: //If dealer score is above player, you lose
		fmt.Println("You lose")
	case dScore == pScore: //Draw if the player score and dealer score is the same
		fmt.Println("Draw")
	}
	fmt.Println()

	//Reset the player's and dealer's hands to nil(empty) for the next round
	ret.Player = nil
	ret.Dealer = nil
	return ret //Return the updated game state
}

func main() {
	var gs GameState

	gs = Shuffle(gs) //Shuffles the deck of cards

	//Play 3 rounds of black jack
	for i := 0; i < 3; i++ {
		gs = Deal(gs) //At the start of each round deal cards to player and dealer

		var input string

		//Pplayers turn:
		for gs.State == StatePlayerTurn {
			//Display the hands of player & dealer
			fmt.Println("Player:", gs.Player)
			fmt.Println("Dealer:", gs.Dealer.DealerString())
			fmt.Println("What will you do? (h)it, (s)tand")

			//User input whether to hit or stand
			fmt.Scanf("%s\n", &input)
			switch input {
			case "h":
				gs = Hit(gs)
			case "s":
				gs = Stand(gs)
			default:
				fmt.Println("Invalid option:", input)
			}
		}
		//Dealers turn
		for gs.State == StateDealerTurn {
			//Dealer hits if their score is 16 or less or a soft 17 (With an Ace)
			if gs.Dealer.Score() <= 16 || (gs.Dealer.Score() == 17 && gs.Dealer.MinScore() != 17) {
				gs = Hit(gs)
			} else {
				gs = Stand(gs)
			}
		}

		gs = EndHand(gs)
	}
}

// Removes the top card from the deck and it returns it along with the remaining deck
func draw(cards []deck.Card) (deck.Card, []deck.Card) {
	return cards[0], cards[1:]
}

// State represents the different stages of the game.
type State int8

const (
	StatePlayerTurn State = iota //0; Players turn
	StateDealerTurn              //1; Dealers turn
	StateHandOver                //2; Hand over
)

// Gamestate holds the current state of the game, including the deck, player, and dealer hands
type GameState struct {
	Deck   []deck.Card //Deck of cards being used in the game
	State  State       //The current state (players turn, dealers, or hand over)
	Player Hand        //The players hand
	Dealer Hand        //Dealers hand
}

// CurrentPLayer returns the hand of the player whose turn it is
func (gs *GameState) CurrentPlayer() *Hand {
	switch gs.State {
	case StatePlayerTurn: //return the players hand during their turn
		return &gs.Player
	case StateDealerTurn: //Return the dealers hand during their turn
		return &gs.Dealer
	default:
		panic("it isn't currently any player's turn")
	}
}

// clone creates a copy of the current game state to avoid modifying the original state.
func clone(gs GameState) GameState {
	ret := GameState{
		Deck:   make([]deck.Card, len(gs.Deck)), //Copy the deck
		State:  gs.State,                        //Copy the current state
		Player: make(Hand, len(gs.Player)),      //Copy the current players hand
		Dealer: make(Hand, len(gs.Dealer)),      //Copy the current dealers hand
	}
	copy(ret.Deck, gs.Deck)     //copies elements from gs.Deck slice to slice ret.Deck
	copy(ret.Player, gs.Player) //Copies elements from gs.player to ret.player
	copy(ret.Dealer, gs.Dealer)
	return ret
}
