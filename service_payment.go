package tx

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/thetatoken/theta/cmd/thetacli/cmd/utils"
	"github.com/thetatoken/theta/common"
	"github.com/thetatoken/theta/crypto"
	"github.com/thetatoken/theta/ledger/types"
	"github.com/thetatoken/theta/rpc"
	wtypes "github.com/thetatoken/theta/wallet/types"

	"github.com/ybbus/jsonrpc"
	rpcc "github.com/ybbus/jsonrpc"
)

// ServicePaymentCmd represents the send command
// Example:
//		thetacli tx service_payment --chain="privatenet" --from=2E833968E5bB786Ae419c4d13189fB081Cc43bab --to=9F1233798E905E173560071255140b4A8aBd3Ec6 --payment_seq=1 --reserve_seq=1 --resource_id=rid1000001
var ServicePaymentCmd = &cobra.Command{
	Use:     "service_payment",
	Short:   "Make Service Payment from Reserve fund",
	Example: `thetacli tx service_payment --chain="privatenet" --from=2E833968E5bB786Ae419c4d13189fB081Cc43bab --to=9F1233798E905E173560071255140b4A8aBd3Ec6 --payment_seq=1 --reserve_seq=1 --resource_id=rid1000001`,
	//Run:     doServicePaymentCmd,
}

func doServicePaymentCmd(cmd *cobra.Command, args []string) []byte {
	walletType := getWalletType(cmd)
	if walletType == wtypes.WalletTypeSoft && len(fromFlag) == 0 {
		utils.Error("The from address cannot be empty") // we don't need to specify the "from address" for hardware wallets
		return ([]byte("error"))
	}

	if len(toFlag) == 0 {
		utils.Error("The to address cannot be empty")
		return ([]byte("error"))
	}

	if fromFlag == toFlag {
		utils.Error("The from and to address cannot be identical")
		return ([]byte("error"))
	}

	var swallet wtypes.Wallet
	//common.HexToAddress(addressStr)
	var fromAddress = common.HexToAddress(fromFlag)
	var twallet wtypes.Wallet
	var toAddress = common.HexToAddress(toFlag)
	var err error

	if onChainFlag {
		//	if 1 == 1 {
		twallet, toAddress, err = walletUnlockWithPath(cmd, toFlag, pathFlag, passwordFlag)
		if err != nil || twallet == nil {
			return ([]byte("error"))
		}
		defer twallet.Lock(toAddress)
	} else {
		swallet, fromAddress, err = walletUnlockWithPath(cmd, fromFlag, pathFlag, passwordFlag)
		if err != nil || swallet == nil {
			return ([]byte("error"))
		}
		defer swallet.Lock(fromAddress)
	}

	tfuel, ok := types.ParseCoinAmount(tfuelAmountFlag)
	if !ok {
		utils.Error("Failed to parse tfuel amount")
	}

	fee, ok := types.ParseCoinAmount(feeFlag)
	if !ok {
		utils.Error("Failed to parse fee")
	}

	sinput := types.TxInput{
		Address: fromAddress,
		Coins: types.Coins{
			ThetaWei: new(big.Int).SetUint64(0),
			//TFuelWei: new(big.Int).Add(tfuel, fee),
			TFuelWei: tfuel,
			//TFuelWei: new(big.Int).SetUint64(1),
		},
		Sequence: uint64(paymentSeqFlag),
		//Signature:
	}

	tinput := types.TxInput{
		Address: toAddress,
		Coins: types.Coins{
			ThetaWei: new(big.Int).SetUint64(0),
			TFuelWei: tfuel,
			//TFuelWei: new(big.Int).SetUint64(0),
		},
		//Sequence: uint64(paymentSeqFlag),
		//Signature:
	}

	// See theta-protocol-ledger > ledger > types > tx.go : Line 522
	servicePaymentTx := &types.ServicePaymentTx{
		Fee: types.Coins{
			ThetaWei: new(big.Int).SetUint64(0),
			TFuelWei: fee,
		},
		Source:          sinput,
		Target:          tinput,
		PaymentSequence: paymentSeqFlag,
		ReserveSequence: reserveSeqFlag,
		ResourceID:      resourceIDFlag,
	}

	// Set the Source Signature
	if onChainFlag {
		//ssig, err := crypto.UnmarshalJSON([]byte(sourceSignatureFlag))
		ssig, err := crypto.SignatureFromBytes(common.FromHex(sourceSignatureFlag))
		if err != nil {
			utils.Error("Failed to convert passed signature: %v\n", err)
		}
		servicePaymentTx.SetSourceSignature(ssig)
	} else {
		ssig, err := swallet.Sign(fromAddress, servicePaymentTx.SourceSignBytes(chainIDFlag))
		if err != nil {
			utils.Error("Failed to sign source transaction: %v\n", err)
		}
		servicePaymentTx.SetSourceSignature(ssig)
	}

	// Set the Target Signature
	if onChainFlag {
		tsig, err := twallet.Sign(toAddress, servicePaymentTx.TargetSignBytes(chainIDFlag))
		if err != nil {
			utils.Error("Failed to sign target transaction: %v\n", err)
		}
		servicePaymentTx.SetTargetSignature(tsig)
	} else {
		tsig, err := crypto.SignatureFromBytes([]byte("unsigned"))
		if err != nil {
			utils.Error("Failed to convert passed signature: %v\n", err)
		}
		servicePaymentTx.SetTargetSignature(tsig)
	}

	raw, err := types.TxToBytes(servicePaymentTx)
	if err != nil {
		utils.Error("Failed to encode transaction: %v\n", err)
	}
	signedTx := hex.EncodeToString(raw)

	var formatted []byte

	if onChainFlag {
		if dryRunFlag {
			formatted, err = json.MarshalIndent(servicePaymentTx, "", "    ")
			if err != nil {
				utils.Error("Failed to parse off-chain transaction: %v\n", err)
			}
			//fmt.Printf("On-Chain transaction(dry-run):\n%s\n", formatted)
			fmt.Printf("%s\n", formatted)

		} else {

			client := rpcc.NewRPCClient(viper.GetString(utils.CfgRemoteRPCEndpoint))

			var res *jsonrpc.RPCResponse
			if asyncFlag {
				res, err = client.Call("theta.BroadcastRawTransactionAsync", rpc.BroadcastRawTransactionArgs{TxBytes: signedTx})
			} else {
				res, err = client.Call("theta.BroadcastRawTransaction", rpc.BroadcastRawTransactionArgs{TxBytes: signedTx})
			}

			if err != nil {
				utils.Error("Failed to broadcast transaction: %v\n", err)
			}
			if res.Error != nil {
				utils.Error("Server returned error: %v\n", res.Error)
				return ([]byte("Likely Expired."))
			}
			result := &rpc.BroadcastRawTransactionResult{}
			err = res.GetObject(result)
			if err != nil {
				utils.Error("Failed to parse server response: %v\n", err)
			}
			formatted, err = json.MarshalIndent(result, "", "    ")
			if err != nil {
				utils.Error("Failed to parse server response: %v\n", err)
			}
			//fmt.Printf("Successfully broadcasted transaction:\n%s\n", formatted)
			// Verbose output makes parsing json difficult
			// fmt.Printf("%s\n", formatted)
		}
	} else {
		formatted, err = json.MarshalIndent(servicePaymentTx, "", "    ")
		if err != nil {
			utils.Error("Failed to parse off-chain transaction: %v\n", err)
		}
		//fmt.Printf("Off-Chain transaction:\n%s\n", formatted)
		//fmt.Printf("%s\n", formatted)
	}

	return (formatted)

}

func DoServicePayment() []byte {
	return (doServicePaymentCmd(ServicePaymentCmd, make([]string, 0)))
}

func init() {
	fmt.Println("thetaoffchaingo_tx service_payment.go init called.")
	// https://github.com/spf13/cobra/blob/master/user_guide.md#working-with-flags
	ServicePaymentCmd.Flags().StringVar(&chainIDFlag, "chain", "", "Chain ID")
	ServicePaymentCmd.Flags().StringVar(&pathFlag, "config", "./thetacli", "Path to Config")
	ServicePaymentCmd.Flags().StringVar(&fromFlag, "from", "", "Address to send from")
	ServicePaymentCmd.Flags().StringVar(&toFlag, "to", "", "Address to send to")
	//ServicePaymentCmd.Flags().StringVar(&pathFlag, "path", "", "Wallet derivation path")
	ServicePaymentCmd.Flags().Uint64Var(&paymentSeqFlag, "payment_seq", 0, "Payment sequence number of the transaction")
	ServicePaymentCmd.Flags().Uint64Var(&reserveSeqFlag, "reserve_seq", 0, "Reserve sequence number of the transaction")
	ServicePaymentCmd.Flags().StringVar(&tfuelAmountFlag, "tfuel", "0", "TFuel amount")
	ServicePaymentCmd.Flags().StringVar(&resourceIDFlag, "resource_id", "", "Corresponding resourceID")
	ServicePaymentCmd.Flags().StringVar(&feeFlag, "fee", fmt.Sprintf("%dwei", types.MinimumTransactionFeeTFuelWei), "Fee")
	ServicePaymentCmd.Flags().StringVar(&walletFlag, "wallet", "soft", "Wallet type (soft|nano|trezor)")
	ServicePaymentCmd.Flags().StringVar(&sourceSignatureFlag, "src_sig", "unsigned", "Source Signature from prior Off-Chain transaction")
	ServicePaymentCmd.Flags().BoolVar(&onChainFlag, "on_chain", false, "Process transaction On-Chain else return json of what would have been sent")
	ServicePaymentCmd.Flags().BoolVar(&asyncFlag, "async", false, "Block until tx has been included in the blockchain")
	ServicePaymentCmd.Flags().StringVar(&passwordFlag, "password", "", "password to unlock the wallet")
	ServicePaymentCmd.Flags().BoolVar(&dryRunFlag, "dry_run", false, "Dry Run(don't execute) the On-Chain transaction")
	ServicePaymentCmd.Flags().BoolVar(&debuggingFlag, "debugging", false, "Print verbose debugging output")

	ServicePaymentCmd.MarkFlagRequired("chain")
	ServicePaymentCmd.MarkFlagRequired("from")
	ServicePaymentCmd.MarkFlagRequired("to")
	ServicePaymentCmd.MarkFlagRequired("payment_seq")
	ServicePaymentCmd.MarkFlagRequired("reserve_seq")
	ServicePaymentCmd.MarkFlagRequired("resource_id")

	rootCmd.Execute()

}
