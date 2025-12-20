package cardingester
import (
	"time"
)

type Bulk_Option_Response struct {
	Object  string `json:"object"`
	HasMore bool   `json:"has_more"`
	Data    []Bulk_Entry `json:"data"`
}

type Bulk_Entry struct {
	Object          string    `json:"object"`
	ID              string    `json:"id"`
	Type            string    `json:"type"`
	UpdatedAt       time.Time `json:"updated_at"`
	URI             string    `json:"uri"`
	Name            string    `json:"name"`
	Description     string    `json:"description"`
	Size            int       `json:"size"`
	DownloadURI     string    `json:"download_uri"`
	ContentType     string    `json:"content_type"`
	ContentEncoding string    `json:"content_encoding"`
}

type MtgCard struct {
	Object        string `json:"object"`
	ID            string `json:"id"`
	OracleID      string `json:"oracle_id"`
	MultiverseIds []int  `json:"multiverse_ids"`
	MtgoID        int    `json:"mtgo_id"`
	ArenaID       int    `json:"arena_id"`
	TcgplayerID   int    `json:"tcgplayer_id"`
	CardmarketID  int    `json:"cardmarket_id"`
	Name          string `json:"name"`
	Lang          string `json:"lang"`
	ReleasedAt    string `json:"released_at"`
	URI           string `json:"uri"`
	ScryfallURI   string `json:"scryfall_uri"`
	Layout        string `json:"layout"`
	HighresImage  bool   `json:"highres_image"`
	ImageStatus   string `json:"image_status"`
	ImageUris     struct {
		Small      string `json:"small"`
		Normal     string `json:"normal"`
		Large      string `json:"large"`
		Png        string `json:"png"`
		ArtCrop    string `json:"art_crop"`
		BorderCrop string `json:"border_crop"`
	} `json:"image_uris"`
	ManaCost      string        `json:"mana_cost"`
	Cmc           float64       `json:"cmc"`
	TypeLine      string        `json:"type_line"`
	OracleText    string        `json:"oracle_text"`
	Colors        []string `json:"colors"`
	ColorIdentity []string      `json:"color_identity"`
	Keywords      []string `json:"keywords"`
	ProducedMana  []string      `json:"produced_mana"`
	Legalities    struct {
		Standard        string `json:"standard"`
		Future          string `json:"future"`
		Historic        string `json:"historic"`
		Timeless        string `json:"timeless"`
		Gladiator       string `json:"gladiator"`
		Pioneer         string `json:"pioneer"`
		Modern          string `json:"modern"`
		Legacy          string `json:"legacy"`
		Pauper          string `json:"pauper"`
		Vintage         string `json:"vintage"`
		Penny           string `json:"penny"`
		Commander       string `json:"commander"`
		Oathbreaker     string `json:"oathbreaker"`
		Standardbrawl   string `json:"standardbrawl"`
		Brawl           string `json:"brawl"`
		Alchemy         string `json:"alchemy"`
		Paupercommander string `json:"paupercommander"`
		Duel            string `json:"duel"`
		Oldschool       string `json:"oldschool"`
		Premodern       string `json:"premodern"`
		Predh           string `json:"predh"`
	} `json:"legalities"`
	Games           []string `json:"games"`
	Reserved        bool     `json:"reserved"`
	GameChanger     bool     `json:"game_changer"`
	Foil            bool     `json:"foil"`
	Nonfoil         bool     `json:"nonfoil"`
	Finishes        []string `json:"finishes"`
	Oversized       bool     `json:"oversized"`
	Promo           bool     `json:"promo"`
	Reprint         bool     `json:"reprint"`
	Variation       bool     `json:"variation"`
	SetID           string   `json:"set_id"`
	Set             string   `json:"set"`
	SetName         string   `json:"set_name"`
	SetType         string   `json:"set_type"`
	SetURI          string   `json:"set_uri"`
	SetSearchURI    string   `json:"set_search_uri"`
	ScryfallSetURI  string   `json:"scryfall_set_uri"`
	RulingsURI      string   `json:"rulings_uri"`
	PrintsSearchURI string   `json:"prints_search_uri"`
	CollectorNumber string   `json:"collector_number"`
	Digital         bool     `json:"digital"`
	Rarity          string   `json:"rarity"`
	CardBackID      string   `json:"card_back_id"`
	Artist          string   `json:"artist"`
	ArtistIds       []string `json:"artist_ids"`
	IllustrationID  string   `json:"illustration_id"`
	BorderColor     string   `json:"border_color"`
	Frame           string   `json:"frame"`
	FullArt         bool     `json:"full_art"`
	Textless        bool     `json:"textless"`
	Booster         bool     `json:"booster"`
	StorySpotlight  bool     `json:"story_spotlight"`
	Prices          struct {
		Usd       string      `json:"usd"`
		UsdFoil   string      `json:"usd_foil"`
		UsdEtched string `json:"usd_etched"`
		Eur       string      `json:"eur"`
		EurFoil   string      `json:"eur_foil"`
		Tix       string      `json:"tix"`
	} `json:"prices"`
	RelatedUris struct {
		Gatherer                  string `json:"gatherer"`
		TcgplayerInfiniteArticles string `json:"tcgplayer_infinite_articles"`
		TcgplayerInfiniteDecks    string `json:"tcgplayer_infinite_decks"`
		Edhrec                    string `json:"edhrec"`
	} `json:"related_uris"`
	PurchaseUris struct {
		Tcgplayer   string `json:"tcgplayer"`
		Cardmarket  string `json:"cardmarket"`
		Cardhoarder string `json:"cardhoarder"`
	} `json:"purchase_uris"`
}
