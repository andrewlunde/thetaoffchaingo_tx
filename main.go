package tx

import (
	"fmt"

	"github.com/spf13/cobra"
)

// Common flags used in Tx sub commands.
var (
	chainIDFlag                  string
	fromFlag                     string
	toFlag                       string
	pathFlag                     string
	seqFlag                      uint64
	thetaAmountFlag              string
	tfuelAmountFlag              string
	gasAmountFlag                uint64
	feeFlag                      string
	resourceIDsFlag              []string
	resourceIDFlag               string
	durationFlag                 uint64
	reserveFundInTFuelFlag       string
	reserveCollateralInTFuelFlag string
	reserveSeqFlag               uint64
	paymentSeqFlag               uint64
	addressesFlag                []string
	percentagesFlag              []string
	valueFlag                    string
	gasPriceFlag                 string
	gasLimitFlag                 uint64
	dataFlag                     string
	walletFlag                   string
	stakeInThetaFlag             string
	purposeFlag                  uint8
	sourceFlag                   string
	holderFlag                   string
	asyncFlag                    bool
	onChainFlag                  bool
	sourceSignatureFlag          string
	dryRunFlag                   bool
	beneficiaryFlag              string
	splitBasisPointFlag          uint64
	passwordFlag                 string
	debuggingFlag                bool
)

// TxCmd represents the Tx command
var TxCmd = &cobra.Command{
	Use:   "tx",
	Short: "Manage transactions",
	Long:  `Manage transactions.`,
}

func init() {
	fmt.Println("thetaoffchaingo_tx main.go init called.")
	//	TxCmd.AddCommand(sendCmd)
	// TxCmd.AddCommand(reserveFundCmd)
	// TxCmd.AddCommand(servicePaymentCmd)
	//	TxCmd.AddCommand(releaseFundCmd) // No need for releaseFundCmd since auto-release is already implemented
	//	TxCmd.AddCommand(splitRuleCmd)
	//	TxCmd.AddCommand(smartContractCmd)
	//	TxCmd.AddCommand(depositStakeCmd)
	//	TxCmd.AddCommand(withdrawStakeCmd)
	//	TxCmd.AddCommand(stakeRewardDistributionCmd)
}
