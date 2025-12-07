package bitcoin

import "regexp"

// Bellscoin
type Bellscoin struct{}

func (Bellscoin) ChainName() string { return "bellscoin" }
func (Bellscoin) CoinbaseDigest(coinbase string) (string, error) { return DoubleSha256(coinbase) }
func (Bellscoin) HeaderDigest(header string) (string, error) { return ScryptDigest(header) }
func (Bellscoin) ShareMultiplier() float64 { return 65536 }
func (Bellscoin) ValidMainnetAddress(address string) bool {
	return regexp.MustCompile("^(B)[a-km-zA-HJ-NP-Z1-9]{33,34}$").MatchString(address)
}
func (Bellscoin) ValidTestnetAddress(address string) bool {
	return regexp.MustCompile("^(n|2)[a-km-zA-HJ-NP-Z1-9]{33}$").MatchString(address)
}
func (Bellscoin) MinimumConfirmations() uint { return 251 }

// Pepecoin
type Pepecoin struct{}

func (Pepecoin) ChainName() string { return "pepecoin" }
func (Pepecoin) CoinbaseDigest(coinbase string) (string, error) { return DoubleSha256(coinbase) }
func (Pepecoin) HeaderDigest(header string) (string, error) { return ScryptDigest(header) }
func (Pepecoin) ShareMultiplier() float64 { return 65536 }
func (Pepecoin) ValidMainnetAddress(address string) bool {
	return regexp.MustCompile("^(P)[a-km-zA-HJ-NP-Z1-9]{33,34}$").MatchString(address)
}
func (Pepecoin) ValidTestnetAddress(address string) bool {
	return regexp.MustCompile("^(n|2)[a-km-zA-HJ-NP-Z1-9]{33}$").MatchString(address)
}
func (Pepecoin) MinimumConfirmations() uint { return 251 }

// Luckycoin
type Luckycoin struct{}

func (Luckycoin) ChainName() string { return "luckycoin" }
func (Luckycoin) CoinbaseDigest(coinbase string) (string, error) { return DoubleSha256(coinbase) }
func (Luckycoin) HeaderDigest(header string) (string, error) { return ScryptDigest(header) }
func (Luckycoin) ShareMultiplier() float64 { return 65536 }
func (Luckycoin) ValidMainnetAddress(address string) bool {
	return regexp.MustCompile("^(L)[a-km-zA-HJ-NP-Z1-9]{33,34}$").MatchString(address)
}
func (Luckycoin) ValidTestnetAddress(address string) bool {
	return regexp.MustCompile("^(n|2)[a-km-zA-HJ-NP-Z1-9]{33}$").MatchString(address)
}
func (Luckycoin) MinimumConfirmations() uint { return 120 }

// Junkcoin
type Junkcoin struct{}

func (Junkcoin) ChainName() string { return "junkcoin" }
func (Junkcoin) CoinbaseDigest(coinbase string) (string, error) { return DoubleSha256(coinbase) }
func (Junkcoin) HeaderDigest(header string) (string, error) { return ScryptDigest(header) }
func (Junkcoin) ShareMultiplier() float64 { return 65536 }
func (Junkcoin) ValidMainnetAddress(address string) bool {
	return regexp.MustCompile("^(J)[a-km-zA-HJ-NP-Z1-9]{33,34}$").MatchString(address)
}
func (Junkcoin) ValidTestnetAddress(address string) bool {
	return regexp.MustCompile("^(n|2)[a-km-zA-HJ-NP-Z1-9]{33}$").MatchString(address)
}
func (Junkcoin) MinimumConfirmations() uint { return 120 }

// Dingocoin
type Dingocoin struct{}

func (Dingocoin) ChainName() string { return "dingocoin" }
func (Dingocoin) CoinbaseDigest(coinbase string) (string, error) { return DoubleSha256(coinbase) }
func (Dingocoin) HeaderDigest(header string) (string, error) { return ScryptDigest(header) }
func (Dingocoin) ShareMultiplier() float64 { return 65536 }
func (Dingocoin) ValidMainnetAddress(address string) bool {
	return regexp.MustCompile("^(D)[a-km-zA-HJ-NP-Z1-9]{33,34}$").MatchString(address)
}
func (Dingocoin) ValidTestnetAddress(address string) bool {
	return regexp.MustCompile("^(n|2)[a-km-zA-HJ-NP-Z1-9]{33}$").MatchString(address)
}
func (Dingocoin) MinimumConfirmations() uint { return 251 }

// Dogmcoin
type Dogmcoin struct{}

func (Dogmcoin) ChainName() string { return "dogmcoin" }
func (Dogmcoin) CoinbaseDigest(coinbase string) (string, error) { return DoubleSha256(coinbase) }
func (Dogmcoin) HeaderDigest(header string) (string, error) { return ScryptDigest(header) }
func (Dogmcoin) ShareMultiplier() float64 { return 65536 }
func (Dogmcoin) ValidMainnetAddress(address string) bool {
	return regexp.MustCompile("^(D)[a-km-zA-HJ-NP-Z1-9]{33,34}$").MatchString(address)
}
func (Dogmcoin) ValidTestnetAddress(address string) bool {
	return regexp.MustCompile("^(n|2)[a-km-zA-HJ-NP-Z1-9]{33}$").MatchString(address)
}
func (Dogmcoin) MinimumConfirmations() uint { return 251 }

// Craftcoin
type Craftcoin struct{}

func (Craftcoin) ChainName() string { return "craftcoin" }
func (Craftcoin) CoinbaseDigest(coinbase string) (string, error) { return DoubleSha256(coinbase) }
func (Craftcoin) HeaderDigest(header string) (string, error) { return ScryptDigest(header) }
func (Craftcoin) ShareMultiplier() float64 { return 65536 }
func (Craftcoin) ValidMainnetAddress(address string) bool {
	return regexp.MustCompile("^(C)[a-km-zA-HJ-NP-Z1-9]{33,34}$").MatchString(address)
}
func (Craftcoin) ValidTestnetAddress(address string) bool {
	return regexp.MustCompile("^(n|2)[a-km-zA-HJ-NP-Z1-9]{33}$").MatchString(address)
}
func (Craftcoin) MinimumConfirmations() uint { return 120 }

// Newyorkcoin
type Newyorkcoin struct{}

func (Newyorkcoin) ChainName() string { return "newyorkcoin" }
func (Newyorkcoin) CoinbaseDigest(coinbase string) (string, error) { return DoubleSha256(coinbase) }
func (Newyorkcoin) HeaderDigest(header string) (string, error) { return ScryptDigest(header) }
func (Newyorkcoin) ShareMultiplier() float64 { return 65536 }
func (Newyorkcoin) ValidMainnetAddress(address string) bool {
	return regexp.MustCompile("^(N|R)[a-km-zA-HJ-NP-Z1-9]{33,34}$").MatchString(address)
}
func (Newyorkcoin) ValidTestnetAddress(address string) bool {
	return regexp.MustCompile("^(n|2)[a-km-zA-HJ-NP-Z1-9]{33}$").MatchString(address)
}
func (Newyorkcoin) MinimumConfirmations() uint { return 120 }

// Earthcoin
type Earthcoin struct{}

func (Earthcoin) ChainName() string { return "earthcoin" }
func (Earthcoin) CoinbaseDigest(coinbase string) (string, error) { return DoubleSha256(coinbase) }
func (Earthcoin) HeaderDigest(header string) (string, error) { return ScryptDigest(header) }
func (Earthcoin) ShareMultiplier() float64 { return 65536 }
func (Earthcoin) ValidMainnetAddress(address string) bool {
	return regexp.MustCompile("^(E|e)[a-km-zA-HJ-NP-Z1-9]{33,34}$").MatchString(address)
}
func (Earthcoin) ValidTestnetAddress(address string) bool {
	return regexp.MustCompile("^(n|2)[a-km-zA-HJ-NP-Z1-9]{33}$").MatchString(address)
}
func (Earthcoin) MinimumConfirmations() uint { return 120 }

// Worldcoin
type Worldcoin struct{}

func (Worldcoin) ChainName() string { return "worldcoin" }
func (Worldcoin) CoinbaseDigest(coinbase string) (string, error) { return DoubleSha256(coinbase) }
func (Worldcoin) HeaderDigest(header string) (string, error) { return ScryptDigest(header) }
func (Worldcoin) ShareMultiplier() float64 { return 65536 }
func (Worldcoin) ValidMainnetAddress(address string) bool {
	return regexp.MustCompile("^(W)[a-km-zA-HJ-NP-Z1-9]{33,34}$").MatchString(address)
}
func (Worldcoin) ValidTestnetAddress(address string) bool {
	return regexp.MustCompile("^(n|2)[a-km-zA-HJ-NP-Z1-9]{33}$").MatchString(address)
}
func (Worldcoin) MinimumConfirmations() uint { return 120 }

// Shibacoin
type Shibacoin struct{}

func (Shibacoin) ChainName() string { return "shibacoin" }
func (Shibacoin) CoinbaseDigest(coinbase string) (string, error) { return DoubleSha256(coinbase) }
func (Shibacoin) HeaderDigest(header string) (string, error) { return ScryptDigest(header) }
func (Shibacoin) ShareMultiplier() float64 { return 65536 }
func (Shibacoin) ValidMainnetAddress(address string) bool {
	return regexp.MustCompile("^(S)[a-km-zA-HJ-NP-Z1-9]{33,34}$").MatchString(address)
}
func (Shibacoin) ValidTestnetAddress(address string) bool {
	return regexp.MustCompile("^(n|2)[a-km-zA-HJ-NP-Z1-9]{33}$").MatchString(address)
}
func (Shibacoin) MinimumConfirmations() uint { return 251 }

// Beerscoin
type Beerscoin struct{}

func (Beerscoin) ChainName() string { return "beerscoin" }
func (Beerscoin) CoinbaseDigest(coinbase string) (string, error) { return DoubleSha256(coinbase) }
func (Beerscoin) HeaderDigest(header string) (string, error) { return ScryptDigest(header) }
func (Beerscoin) ShareMultiplier() float64 { return 65536 }
func (Beerscoin) ValidMainnetAddress(address string) bool {
	return regexp.MustCompile("^(B)[a-km-zA-HJ-NP-Z1-9]{33,34}$").MatchString(address)
}
func (Beerscoin) ValidTestnetAddress(address string) bool {
	return regexp.MustCompile("^(n|2)[a-km-zA-HJ-NP-Z1-9]{33}$").MatchString(address)
}
func (Beerscoin) MinimumConfirmations() uint { return 120 }

// Dogecoinev
type Dogecoinev struct{}

func (Dogecoinev) ChainName() string { return "dogecoinev" }
func (Dogecoinev) CoinbaseDigest(coinbase string) (string, error) { return DoubleSha256(coinbase) }
func (Dogecoinev) HeaderDigest(header string) (string, error) { return ScryptDigest(header) }
func (Dogecoinev) ShareMultiplier() float64 { return 65536 }
func (Dogecoinev) ValidMainnetAddress(address string) bool {
	return regexp.MustCompile("^(D)[a-km-zA-HJ-NP-Z1-9]{33,34}$").MatchString(address)
}
func (Dogecoinev) ValidTestnetAddress(address string) bool {
	return regexp.MustCompile("^(n|2)[a-km-zA-HJ-NP-Z1-9]{33}$").MatchString(address)
}
func (Dogecoinev) MinimumConfirmations() uint { return 251 }

// Bonkcoin
type Bonkcoin struct{}

func (Bonkcoin) ChainName() string { return "bonkcoin" }
func (Bonkcoin) CoinbaseDigest(coinbase string) (string, error) { return DoubleSha256(coinbase) }
func (Bonkcoin) HeaderDigest(header string) (string, error) { return ScryptDigest(header) }
func (Bonkcoin) ShareMultiplier() float64 { return 65536 }
func (Bonkcoin) ValidMainnetAddress(address string) bool {
	return regexp.MustCompile("^(B)[a-km-zA-HJ-NP-Z1-9]{33,34}$").MatchString(address)
}
func (Bonkcoin) ValidTestnetAddress(address string) bool {
	return regexp.MustCompile("^(n|2)[a-km-zA-HJ-NP-Z1-9]{33}$").MatchString(address)
}
func (Bonkcoin) MinimumConfirmations() uint { return 120 }

// Flincoin
type Flincoin struct{}

func (Flincoin) ChainName() string { return "flincoin" }
func (Flincoin) CoinbaseDigest(coinbase string) (string, error) { return DoubleSha256(coinbase) }
func (Flincoin) HeaderDigest(header string) (string, error) { return ScryptDigest(header) }
func (Flincoin) ShareMultiplier() float64 { return 65536 }
func (Flincoin) ValidMainnetAddress(address string) bool {
	return regexp.MustCompile("^(F)[a-km-zA-HJ-NP-Z1-9]{33,34}$").MatchString(address)
}
func (Flincoin) ValidTestnetAddress(address string) bool {
	return regexp.MustCompile("^(n|2)[a-km-zA-HJ-NP-Z1-9]{33}$").MatchString(address)
}
func (Flincoin) MinimumConfirmations() uint { return 120 }

// Marscoin
type Marscoin struct{}

func (Marscoin) ChainName() string { return "marscoin" }
func (Marscoin) CoinbaseDigest(coinbase string) (string, error) { return DoubleSha256(coinbase) }
func (Marscoin) HeaderDigest(header string) (string, error) { return ScryptDigest(header) }
func (Marscoin) ShareMultiplier() float64 { return 65536 }
func (Marscoin) ValidMainnetAddress(address string) bool {
	return regexp.MustCompile("^(M)[a-km-zA-HJ-NP-Z1-9]{33,34}$").MatchString(address)
}
func (Marscoin) ValidTestnetAddress(address string) bool {
	return regexp.MustCompile("^(n|2)[a-km-zA-HJ-NP-Z1-9]{33}$").MatchString(address)
}
func (Marscoin) MinimumConfirmations() uint { return 120 }

// BBQcoin
type BBQcoin struct{}

func (BBQcoin) ChainName() string { return "bbqcoin" }
func (BBQcoin) CoinbaseDigest(coinbase string) (string, error) { return DoubleSha256(coinbase) }
func (BBQcoin) HeaderDigest(header string) (string, error) { return ScryptDigest(header) }
func (BBQcoin) ShareMultiplier() float64 { return 65536 }
func (BBQcoin) ValidMainnetAddress(address string) bool {
	return regexp.MustCompile("^(b)[a-km-zA-HJ-NP-Z1-9]{33,34}$").MatchString(address)
}
func (BBQcoin) ValidTestnetAddress(address string) bool {
	return regexp.MustCompile("^(n|2)[a-km-zA-HJ-NP-Z1-9]{33}$").MatchString(address)
}
func (BBQcoin) MinimumConfirmations() uint { return 120 }

// Goldcoin
type Goldcoin struct{}

func (Goldcoin) ChainName() string { return "goldcoin" }
func (Goldcoin) CoinbaseDigest(coinbase string) (string, error) { return DoubleSha256(coinbase) }
func (Goldcoin) HeaderDigest(header string) (string, error) { return ScryptDigest(header) }
func (Goldcoin) ShareMultiplier() float64 { return 65536 }
func (Goldcoin) ValidMainnetAddress(address string) bool {
	return regexp.MustCompile("^(G)[a-km-zA-HJ-NP-Z1-9]{33,34}$").MatchString(address)
}
func (Goldcoin) ValidTestnetAddress(address string) bool {
	return regexp.MustCompile("^(n|2)[a-km-zA-HJ-NP-Z1-9]{33}$").MatchString(address)
}
func (Goldcoin) MinimumConfirmations() uint { return 120 }

// Catcoin
type Catcoin struct{}

func (Catcoin) ChainName() string { return "catcoin" }
func (Catcoin) CoinbaseDigest(coinbase string) (string, error) { return DoubleSha256(coinbase) }
func (Catcoin) HeaderDigest(header string) (string, error) { return ScryptDigest(header) }
func (Catcoin) ShareMultiplier() float64 { return 65536 }
func (Catcoin) ValidMainnetAddress(address string) bool {
	return regexp.MustCompile("^(9)[a-km-zA-HJ-NP-Z1-9]{33,34}$").MatchString(address)
}
func (Catcoin) ValidTestnetAddress(address string) bool {
	return regexp.MustCompile("^(n|2)[a-km-zA-HJ-NP-Z1-9]{33}$").MatchString(address)
}
func (Catcoin) MinimumConfirmations() uint { return 120 }

// Cyberyen
type Cyberyen struct{}

func (Cyberyen) ChainName() string { return "cyberyen" }
func (Cyberyen) CoinbaseDigest(coinbase string) (string, error) { return DoubleSha256(coinbase) }
func (Cyberyen) HeaderDigest(header string) (string, error) { return ScryptDigest(header) }
func (Cyberyen) ShareMultiplier() float64 { return 65536 }
func (Cyberyen) ValidMainnetAddress(address string) bool {
	return regexp.MustCompile("^(C)[a-km-zA-HJ-NP-Z1-9]{33,34}$").MatchString(address)
}
func (Cyberyen) ValidTestnetAddress(address string) bool {
	return regexp.MustCompile("^(n|2)[a-km-zA-HJ-NP-Z1-9]{33}$").MatchString(address)
}
func (Cyberyen) MinimumConfirmations() uint { return 120 }

// Infinitecoin
type Infinitecoin struct{}

func (Infinitecoin) ChainName() string { return "infinitecoin" }
func (Infinitecoin) CoinbaseDigest(coinbase string) (string, error) { return DoubleSha256(coinbase) }
func (Infinitecoin) HeaderDigest(header string) (string, error) { return ScryptDigest(header) }
func (Infinitecoin) ShareMultiplier() float64 { return 65536 }
func (Infinitecoin) ValidMainnetAddress(address string) bool {
	return regexp.MustCompile("^(i)[a-km-zA-HJ-NP-Z1-9]{33,34}$").MatchString(address)
}
func (Infinitecoin) ValidTestnetAddress(address string) bool {
	return regexp.MustCompile("^(n|2)[a-km-zA-HJ-NP-Z1-9]{33}$").MatchString(address)
}
func (Infinitecoin) MinimumConfirmations() uint { return 120 }

// IBithub
type IBithub struct{}

func (IBithub) ChainName() string { return "ibithub" }
func (IBithub) CoinbaseDigest(coinbase string) (string, error) { return DoubleSha256(coinbase) }
func (IBithub) HeaderDigest(header string) (string, error) { return ScryptDigest(header) }
func (IBithub) ShareMultiplier() float64 { return 65536 }
func (IBithub) ValidMainnetAddress(address string) bool {
	return regexp.MustCompile("^(I)[a-km-zA-HJ-NP-Z1-9]{33,34}$").MatchString(address)
}
func (IBithub) ValidTestnetAddress(address string) bool {
	return regexp.MustCompile("^(n|2)[a-km-zA-HJ-NP-Z1-9]{33}$").MatchString(address)
}
func (IBithub) MinimumConfirmations() uint { return 120 }

// Newenglandcoin
type Newenglandcoin struct{}

func (Newenglandcoin) ChainName() string { return "newenglandcoin" }
func (Newenglandcoin) CoinbaseDigest(coinbase string) (string, error) {
	return DoubleSha256(coinbase)
}
func (Newenglandcoin) HeaderDigest(header string) (string, error) { return ScryptDigest(header) }
func (Newenglandcoin) ShareMultiplier() float64                   { return 65536 }
func (Newenglandcoin) ValidMainnetAddress(address string) bool {
	return regexp.MustCompile("^(N)[a-km-zA-HJ-NP-Z1-9]{33,34}$").MatchString(address)
}
func (Newenglandcoin) ValidTestnetAddress(address string) bool {
	return regexp.MustCompile("^(n|2)[a-km-zA-HJ-NP-Z1-9]{33}$").MatchString(address)
}
func (Newenglandcoin) MinimumConfirmations() uint { return 120 }

// Bitbar
type Bitbar struct{}

func (Bitbar) ChainName() string { return "bitbar" }
func (Bitbar) CoinbaseDigest(coinbase string) (string, error) { return DoubleSha256(coinbase) }
func (Bitbar) HeaderDigest(header string) (string, error) { return ScryptDigest(header) }
func (Bitbar) ShareMultiplier() float64 { return 65536 }
func (Bitbar) ValidMainnetAddress(address string) bool {
	return regexp.MustCompile("^(B)[a-km-zA-HJ-NP-Z1-9]{33,34}$").MatchString(address)
}
func (Bitbar) ValidTestnetAddress(address string) bool {
	return regexp.MustCompile("^(n|2)[a-km-zA-HJ-NP-Z1-9]{33}$").MatchString(address)
}
func (Bitbar) MinimumConfirmations() uint { return 120 }

// Ferrite
type Ferrite struct{}

func (Ferrite) ChainName() string { return "ferrite" }
func (Ferrite) CoinbaseDigest(coinbase string) (string, error) { return DoubleSha256(coinbase) }
func (Ferrite) HeaderDigest(header string) (string, error) { return ScryptDigest(header) }
func (Ferrite) ShareMultiplier() float64 { return 65536 }
func (Ferrite) ValidMainnetAddress(address string) bool {
	return regexp.MustCompile("^(F)[a-km-zA-HJ-NP-Z1-9]{33,34}$|^(fe1)[0-9A-Za-z]{39}$").MatchString(address)
}
func (Ferrite) ValidTestnetAddress(address string) bool {
	return regexp.MustCompile("^(n|2)[a-km-zA-HJ-NP-Z1-9]{33}$").MatchString(address)
}
func (Ferrite) MinimumConfirmations() uint { return 120 }

// Flopcoin
type Flopcoin struct{}

func (Flopcoin) ChainName() string { return "flopcoin" }
func (Flopcoin) CoinbaseDigest(coinbase string) (string, error) { return DoubleSha256(coinbase) }
func (Flopcoin) HeaderDigest(header string) (string, error) { return ScryptDigest(header) }
func (Flopcoin) ShareMultiplier() float64 { return 65536 }
func (Flopcoin) ValidMainnetAddress(address string) bool {
	return regexp.MustCompile("^(F)[a-km-zA-HJ-NP-Z1-9]{33,34}$").MatchString(address)
}
func (Flopcoin) ValidTestnetAddress(address string) bool {
	return regexp.MustCompile("^(n|2)[a-km-zA-HJ-NP-Z1-9]{33}$").MatchString(address)
}
func (Flopcoin) MinimumConfirmations() uint { return 120 }

// Stohncoin
type Stohncoin struct{}

func (Stohncoin) ChainName() string { return "stohncoin" }
func (Stohncoin) CoinbaseDigest(coinbase string) (string, error) { return DoubleSha256(coinbase) }
func (Stohncoin) HeaderDigest(header string) (string, error) { return ScryptDigest(header) }
func (Stohncoin) ShareMultiplier() float64 { return 65536 }
func (Stohncoin) ValidMainnetAddress(address string) bool {
	return regexp.MustCompile("^(S)[a-km-zA-HJ-NP-Z1-9]{33,34}$").MatchString(address)
}
func (Stohncoin) ValidTestnetAddress(address string) bool {
	return regexp.MustCompile("^(n|2)[a-km-zA-HJ-NP-Z1-9]{33}$").MatchString(address)
}
func (Stohncoin) MinimumConfirmations() uint { return 120 }

// Sorachancoin
type Sorachancoin struct{}

func (Sorachancoin) ChainName() string { return "sorachancoin" }
func (Sorachancoin) CoinbaseDigest(coinbase string) (string, error) { return DoubleSha256(coinbase) }
func (Sorachancoin) HeaderDigest(header string) (string, error) { return ScryptDigest(header) }
func (Sorachancoin) ShareMultiplier() float64 { return 65536 }
func (Sorachancoin) ValidMainnetAddress(address string) bool {
	return regexp.MustCompile("^(S)[a-km-zA-HJ-NP-Z1-9]{33,34}$").MatchString(address)
}
func (Sorachancoin) ValidTestnetAddress(address string) bool {
	return regexp.MustCompile("^(n|2)[a-km-zA-HJ-NP-Z1-9]{33}$").MatchString(address)
}
func (Sorachancoin) MinimumConfirmations() uint { return 120 }

// Mooncoin
type Mooncoin struct{}

func (Mooncoin) ChainName() string { return "mooncoin" }
func (Mooncoin) CoinbaseDigest(coinbase string) (string, error) { return DoubleSha256(coinbase) }
func (Mooncoin) HeaderDigest(header string) (string, error) { return ScryptDigest(header) }
func (Mooncoin) ShareMultiplier() float64 { return 65536 }
func (Mooncoin) ValidMainnetAddress(address string) bool {
	return regexp.MustCompile("^(M|2)[a-km-zA-HJ-NP-Z1-9]{33,34}$").MatchString(address)
}
func (Mooncoin) ValidTestnetAddress(address string) bool {
	return regexp.MustCompile("^(n|2)[a-km-zA-HJ-NP-Z1-9]{33}$").MatchString(address)
}
func (Mooncoin) MinimumConfirmations() uint { return 120 }

// Fairbrix
type Fairbrix struct{}

func (Fairbrix) ChainName() string { return "fairbrix" }
func (Fairbrix) CoinbaseDigest(coinbase string) (string, error) { return DoubleSha256(coinbase) }
func (Fairbrix) HeaderDigest(header string) (string, error) { return ScryptDigest(header) }
func (Fairbrix) ShareMultiplier() float64 { return 65536 }
func (Fairbrix) ValidMainnetAddress(address string) bool {
	return regexp.MustCompile("^(f)[a-km-zA-HJ-NP-Z1-9]{33,34}$").MatchString(address)
}
func (Fairbrix) ValidTestnetAddress(address string) bool {
	return regexp.MustCompile("^(n|2)[a-km-zA-HJ-NP-Z1-9]{33}$").MatchString(address)
}
func (Fairbrix) MinimumConfirmations() uint { return 120 }

// Lebowskiscoin
type Lebowskiscoin struct{}

func (Lebowskiscoin) ChainName() string { return "lebowskiscoin" }
func (Lebowskiscoin) CoinbaseDigest(coinbase string) (string, error) {
	return DoubleSha256(coinbase)
}
func (Lebowskiscoin) HeaderDigest(header string) (string, error) { return ScryptDigest(header) }
func (Lebowskiscoin) ShareMultiplier() float64                   { return 65536 }
func (Lebowskiscoin) ValidMainnetAddress(address string) bool {
	return regexp.MustCompile("^(L)[a-km-zA-HJ-NP-Z1-9]{33,34}$").MatchString(address)
}
func (Lebowskiscoin) ValidTestnetAddress(address string) bool {
	return regexp.MustCompile("^(n|2)[a-km-zA-HJ-NP-Z1-9]{33}$").MatchString(address)
}
func (Lebowskiscoin) MinimumConfirmations() uint { return 120 }

// Bit
type Bit struct{}

func (Bit) ChainName() string { return "bit" }
func (Bit) CoinbaseDigest(coinbase string) (string, error) { return DoubleSha256(coinbase) }
func (Bit) HeaderDigest(header string) (string, error) { return ScryptDigest(header) }
func (Bit) ShareMultiplier() float64 { return 65536 }
func (Bit) ValidMainnetAddress(address string) bool {
	return regexp.MustCompile("^(B)[a-km-zA-HJ-NP-Z1-9]{33,34}$").MatchString(address)
}
func (Bit) ValidTestnetAddress(address string) bool {
	return regexp.MustCompile("^(n|2)[a-km-zA-HJ-NP-Z1-9]{33}$").MatchString(address)
}
func (Bit) MinimumConfirmations() uint { return 120 }

// Trumpow
type Trumpow struct{}

func (Trumpow) ChainName() string { return "trumpow" }
func (Trumpow) CoinbaseDigest(coinbase string) (string, error) { return DoubleSha256(coinbase) }
func (Trumpow) HeaderDigest(header string) (string, error) { return ScryptDigest(header) }
func (Trumpow) ShareMultiplier() float64 { return 65536 }
func (Trumpow) ValidMainnetAddress(address string) bool {
	return regexp.MustCompile("^(T)[a-km-zA-HJ-NP-Z1-9]{33,34}$").MatchString(address)
}
func (Trumpow) ValidTestnetAddress(address string) bool {
	return regexp.MustCompile("^(n|2)[a-km-zA-HJ-NP-Z1-9]{33}$").MatchString(address)
}
func (Trumpow) MinimumConfirmations() uint { return 120 }

// Mydogecoin
type Mydogecoin struct{}

func (Mydogecoin) ChainName() string { return "mydogecoin" }
func (Mydogecoin) CoinbaseDigest(coinbase string) (string, error) { return DoubleSha256(coinbase) }
func (Mydogecoin) HeaderDigest(header string) (string, error) { return ScryptDigest(header) }
func (Mydogecoin) ShareMultiplier() float64 { return 65536 }
func (Mydogecoin) ValidMainnetAddress(address string) bool {
	return regexp.MustCompile("^(D|M)[a-km-zA-HJ-NP-Z1-9]{33,34}$").MatchString(address)
}
func (Mydogecoin) ValidTestnetAddress(address string) bool {
	return regexp.MustCompile("^(n|2)[a-km-zA-HJ-NP-Z1-9]{33}$").MatchString(address)
}
func (Mydogecoin) MinimumConfirmations() uint { return 251 }
