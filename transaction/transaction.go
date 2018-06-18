package transaction

type Transaction struct {
	ID []byte
	Vin []TXInput
	Vout []TXOutput
}


type TXInput struct {
	
}

type TXOutput struct {
	Value int
	ScriptPubKey string
}
