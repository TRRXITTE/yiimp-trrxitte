package bitcoin

const BitcoinMinConfirmations = 102

type Blockchain interface {
	ChainName() string
	CoinbaseDigest(coinbase string) (string, error)
	HeaderDigest(header string) (string, error)
	ShareMultiplier() float64
	MinimumConfirmations() uint

	ValidMainnetAddress(address string) bool
	ValidTestnetAddress(address string) bool
}

func GetChain(chainName string) Blockchain {
	switch chainName {
	case "dogecoin":
		return Dogecoin{}
	case "litecoin":
		return Litecoin{}
	case "bellscoin":
		return Bellscoin{}
	case "pepecoin":
		return Pepecoin{}
	case "luckycoin":
		return Luckycoin{}
	case "junkcoin":
		return Junkcoin{}
	case "dingocoin":
		return Dingocoin{}
	case "dogmcoin":
		return Dogmcoin{}
	case "craftcoin":
		return Craftcoin{}
	case "newyorkcoin":
		return Newyorkcoin{}
	case "earthcoin":
		return Earthcoin{}
	case "worldcoin":
		return Worldcoin{}
	case "shibacoin":
		return Shibacoin{}
	case "beerscoin":
		return Beerscoin{}
	case "dogecoinev":
		return Dogecoinev{}
	case "bonkcoin":
		return Bonkcoin{}
	case "flincoin":
		return Flincoin{}
	case "marscoin":
		return Marscoin{}
	case "bbqcoin":
		return BBQcoin{}
	case "goldcoin":
		return Goldcoin{}
	case "catcoin":
		return Catcoin{}
	case "cyberyen":
		return Cyberyen{}
	case "infinitecoin":
		return Infinitecoin{}
	case "ibithub":
		return IBithub{}
	case "newenglandcoin":
		return Newenglandcoin{}
	case "bitbar":
		return Bitbar{}
	case "ferrite":
		return Ferrite{}
	case "flopcoin":
		return Flopcoin{}
	case "stohncoin":
		return Stohncoin{}
	case "sorachancoin":
		return Sorachancoin{}
	case "mooncoin":
		return Mooncoin{}
	case "fairbrix":
		return Fairbrix{}
	case "lebowskiscoin":
		return Lebowskiscoin{}
	case "bit":
		return Bit{}
	case "trumpow":
		return Trumpow{}
	case "mydogecoin":
		return Mydogecoin{}
	default:
		panic("Unknown blockchain: " + chainName)
	}
}
