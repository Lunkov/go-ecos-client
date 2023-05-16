package messages

import (
)

type NFTToken struct {
  ID              uint64  `json:"id"`
  AssetType       string  `json:"asset_type"`
  TokenType       string  `json:"token_type"`
  
  Name            string  `json:"name"`
  Description     string  `json:"description"`
  Symbol          string  `json:"symbol"`
  TokenUri        string  `json:"TokenUri" validate:"string,max=2000"`
  
  Owner           string  `json:"Owner,omitempty" validate:"string"`
  
  IsBurned        bool   `json:"IsBurned" validate:"bool"`
	BurnedBy        string `json:"BurnedBy,omitempty" validate:"string"`
	BurnedDate      string `json:"BurnedDate,omitempty" validate:"string"`
}

// https://docs.kaleido.io/kaleido-services/digital-assets/
// https://docs.oracle.com/en/cloud/paas/blockchain-cloud/usingoci/scaffolded-go-nft-project.html
