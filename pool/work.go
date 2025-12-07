package pool

import (
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"designs.capital/dogepool/bitcoin"
	"designs.capital/dogepool/persistence"
)

// Main INPUT
func (p *PoolServer) fetchRpcBlockTemplatesAndCacheWork() error {
	var block *bitcoin.BitcoinBlock
	var err error
	template, auxBlocks, err := p.fetchAllBlockTemplatesFromRPC()
	if err != nil {
		// Switch nodes if we fail to get work
		err = p.CheckAndRecoverRPCs()
		if err != nil {
			return err
		}
		template, auxBlocks, err = p.fetchAllBlockTemplatesFromRPC()
		if err != nil {
			return err
		}
	}

	auxillary := p.config.BlockSignature

	if len(auxBlocks) > 0 {
		auxMerkleTree := bitcoin.BuildAuxChainMerkleTree(auxBlocks)

		for i := range auxBlocks {
			branch := auxMerkleTree.Branches[i]
			auxBlocks[i].MerkleIndex = branch.Index
			auxBlocks[i].MerkleBranch = branch.Siblings
		}

		if len(auxBlocks) == 1 {
			mergedPOW := auxBlocks[0].GetWork()
			auxillary = auxillary + hexStringToByteString(mergedPOW)
		} else {
			mergedPOW := auxBlocks[0].GetWorkWithMerkleRoot(auxMerkleTree.Root, auxMerkleTree.Size)
			auxillary = auxillary + hexStringToByteString(mergedPOW)
		}

		p.templates.AuxBlocks = auxBlocks
	}

	primaryName := p.config.GetPrimary()
	// TODO this is chain/bitcoin specific
	rewardPubScriptKey := p.GetPrimaryNode().RewardPubScriptKey
	extranonceByteReservationLength := 8

	var auxBlockPtr *bitcoin.AuxBlock
	if len(auxBlocks) > 0 {
		auxBlockPtr = &auxBlocks[0]
	}

	block, p.workCache, err = bitcoin.GenerateWork(&template, auxBlockPtr,
		primaryName, auxillary, rewardPubScriptKey,
		extranonceByteReservationLength)
	if err != nil {
		log.Print(err)
	}

	p.templates.BitcoinBlock = *block

	return nil
}

// Main OUTPUT
func (p *PoolServer) recieveWorkFromClient(share bitcoin.Work, client *stratumClient) error {
	primaryBlockTemplate := p.templates.GetPrimary()
	if primaryBlockTemplate.Template == nil {
		return errors.New("primary block template not yet set")
	}
	auxBlocks := p.templates.AuxBlocks

	var err error

	// TODO - this key and interface isn't very invertable..
	workerString := share[0].(string)
	workerStringParts := strings.Split(workerString, ".")
	if len(workerStringParts) < 2 {
		return errors.New("invalid miner address")
	}
	minerAddress := workerStringParts[0]
	rigID := workerStringParts[1]

	primaryBlockHeight := primaryBlockTemplate.Template.Height
	nonce := share[primaryBlockTemplate.NonceSubmissionSlot()].(string)
	extranonce2Slot, _ := primaryBlockTemplate.Extranonce2SubmissionSlot()
	extranonce2 := share[extranonce2Slot].(string)
	nonceTime := share[primaryBlockTemplate.NonceTimeSubmissionSlot()].(string)

	// TODO - validate input

	extranonce := client.extranonce1 + extranonce2

	_, err = primaryBlockTemplate.MakeHeader(extranonce, nonce, nonceTime)

	if err != nil {
		return err
	}

	result := validateAndWeighShare(&primaryBlockTemplate, auxBlocks, p.config.PoolDifficulty)

	if result.Status == shareInvalid {
		m := "❔ Invalid share for block %v from %v [%v] [%v]"
		m = fmt.Sprintf(m, primaryBlockHeight, client.ip, rigID, client.userAgent)
		return errors.New(m)
	}

	m := "Valid share for block %v from %v [%v]"
	m = fmt.Sprintf(m, primaryBlockHeight, client.ip, rigID)
	log.Println(m)

	blockTarget := bitcoin.Target(primaryBlockTemplate.Template.Target)
	blockDifficulty, _ := blockTarget.ToDifficulty()
	blockDifficulty = blockDifficulty * primaryBlockTemplate.ShareMultiplier()

	p.Lock()
	p.shareBuffer = append(p.shareBuffer, persistence.Share{
		PoolID:            p.config.PoolName,
		BlockHeight:       primaryBlockHeight,
		Miner:             minerAddress,
		Worker:            rigID,
		UserAgent:         client.userAgent,
		Difficulty:        result.ShareDifficulty,
		NetworkDifficulty: blockDifficulty,
		IpAddress:         client.ip,
		Created:           time.Now(),
	})
	p.Unlock()

	if result.Status == shareValid {
		return nil
	}

	submittedChains := make([]string, 0)

	for _, auxIndex := range result.AuxChainsMetTargets {
		auxBlock := auxBlocks[auxIndex]
		chainName := p.config.BlockChainOrder[auxIndex+1]

		log.Printf("Block candidate for %s at height %v from %v [%v]", chainName, auxBlock.Height, client.ip, rigID)

		err = p.submitAuxBlockForChain(primaryBlockTemplate, auxBlock, chainName)
		if err != nil {
			log.Printf("Failed to submit aux block for %s: %v", chainName, err)
			p.rpcManagers[chainName].CheckAndRecoverRPCs()
			err = p.submitAuxBlockForChain(primaryBlockTemplate, auxBlock, chainName)
		}

		if err == nil {
			auxTarget := bitcoin.Target(reverseHexBytes(auxBlock.Target))
			auxDifficulty, _ := auxTarget.ToDifficulty()
			auxDifficulty = auxDifficulty * bitcoin.GetChain(chainName).ShareMultiplier()

			found := persistence.Found{
				PoolID:                      p.config.PoolName,
				Chain:                       chainName,
				Created:                     time.Now(),
				Hash:                        auxBlock.Hash,
				NetworkDifficulty:           auxDifficulty,
				BlockHeight:                 uint(auxBlock.Height),
				TransactionConfirmationData: reverseHexBytes(auxBlock.CoinbaseHash),
				Status:                      persistence.StatusPending,
				Type:                        "Auxiliary",
				ConfirmationProgress:        0,
				Miner:                       minerAddress,
				Source:                      "",
			}

			err = persistence.Blocks.Insert(found)
			if err != nil {
				log.Println(err)
			} else {
				submittedChains = append(submittedChains, chainName)
				log.Printf("✅  Successful auxiliary block submission for %s at height %v from: %v [%v]", chainName, auxBlock.Height, client.ip, rigID)
			}
		}
	}

	if result.PrimaryMeetsTarget {
		log.Printf("Primary block candidate for %s at height %v from %v [%v]", p.config.GetPrimary(), primaryBlockHeight, client.ip, rigID)

		err = p.submitBlockToChain(primaryBlockTemplate)
		if err != nil {
			log.Printf("Failed to submit primary block: %v", err)
			p.rpcManagers[p.config.GetPrimary()].CheckAndRecoverRPCs()
			err = p.submitBlockToChain(primaryBlockTemplate)
		}

		if err == nil {
			found := persistence.Found{
				PoolID:               p.config.PoolName,
				Chain:                p.config.GetPrimary(),
				Created:              time.Now(),
				NetworkDifficulty:    blockDifficulty,
				BlockHeight:          primaryBlockHeight,
				Status:               persistence.StatusPending,
				Type:                 "Primary",
				ConfirmationProgress: 0,
				Miner:                minerAddress,
				Source:               "",
			}

			found.Hash, err = primaryBlockTemplate.HeaderHashed()
			if err != nil {
				log.Println(err)
			}
			found.TransactionConfirmationData, err = primaryBlockTemplate.CoinbaseHashed()
			if err != nil {
				log.Println(err)
			}

			err = persistence.Blocks.Insert(found)
			if err != nil {
				log.Println(err)
			} else {
				submittedChains = append(submittedChains, p.config.GetPrimary())
				log.Printf("✅  Successful primary block submission for %s at height %v from: %v [%v]", p.config.GetPrimary(), primaryBlockHeight, client.ip, rigID)
			}
		}
	}

	if len(submittedChains) > 0 {
		log.Printf("✅  Successfully submitted blocks to: %v", submittedChains)
	}

	return nil
}

func (pool *PoolServer) generateWorkFromCache(refresh bool) (bitcoin.Work, error) {
	work := append(pool.workCache, interface{}(refresh))

	return work, nil
}
