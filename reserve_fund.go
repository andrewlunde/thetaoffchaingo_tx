package tx

import (
	"encoding/hex"
	"encoding/json"
	"fmt"

	// "fmt"
	"math/big"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/thetatoken/theta/cmd/thetacli/cmd/utils"
	"github.com/thetatoken/theta/ledger/types"
	"github.com/thetatoken/theta/rpc"

	//	"github.com/ybbus/jsonrpc"
	rpcc "github.com/ybbus/jsonrpc"
)

// ReserveFundCmd represents the reserve fund command
// Example:
//		thetacli tx reserve --chain="privatenet" --from=2E833968E5bB786Ae419c4d13189fB081Cc43bab --fund=900 --collateral=1203 --seq=6 --duration=1002 --resource_ids=die_another_day,hello
var ReserveFundCmd = &cobra.Command{
	Use:     "reserve",
	Short:   "Reserve fund for an off-chain micropayment",
	Example: `thetacli tx reserve --chain="privatenet" --from=2E833968E5bB786Ae419c4d13189fB081Cc43bab --fund=900 --collateral=1203 --seq=6 --duration=1002 --resource_ids=die_another_day,hello`,
	//Run:     doReserveFundCmd,
}

func doReserveFundCmd(cmd *cobra.Command, args []string) []byte {
	// thetacli tx reserve --chain="privatenet" --async --from=2E833968E5bB786Ae419c4d13189fB081Cc43bab --fund=100 --collateral=101 --duration=100 --resource_ids=rid1000001 --seq=2 --password=qwertyuiop
	//func ReserveFund(cmd *cobra.Command, args []string) []byte {
	wallet, fromAddress, err := walletUnlock(cmd, fromFlag, passwordFlag)
	if err != nil {
		return ([]byte("error"))
	}
	defer wallet.Lock(fromAddress)

	fee, ok := types.ParseCoinAmount(feeFlag)
	if !ok {
		utils.Error("Failed to parse fee")
	}
	fund, ok := types.ParseCoinAmount(reserveFundInTFuelFlag)
	if !ok {
		utils.Error("Failed to parse fund")
	}
	col, ok := types.ParseCoinAmount(reserveCollateralInTFuelFlag)
	if !ok {
		utils.Error("Failed to parse collateral")
	}
	input := types.TxInput{
		Address: fromAddress,
		Coins: types.Coins{
			ThetaWei: new(big.Int).SetUint64(0),
			TFuelWei: fund,
		},
		Sequence: uint64(seqFlag),
	}
	resourceIDs := []string{}
	for _, id := range resourceIDsFlag {
		resourceIDs = append(resourceIDs, id)
	}
	collateral := types.Coins{
		ThetaWei: new(big.Int).SetUint64(0),
		TFuelWei: col,
	}
	if !collateral.IsPositive() {
		utils.Error("Invalid input: collateral must be positive\n")
	}

	reserveFundTx := &types.ReserveFundTx{
		Fee: types.Coins{
			ThetaWei: new(big.Int).SetUint64(0),
			TFuelWei: fee,
		},
		Source:      input,
		ResourceIDs: resourceIDs,
		Collateral:  collateral,
		Duration:    durationFlag,
	}

	sig, err := wallet.Sign(fromAddress, reserveFundTx.SignBytes(chainIDFlag))
	if err != nil {
		utils.Error("Failed to sign transaction: %v\n", err)
	}
	reserveFundTx.SetSignature(fromAddress, sig)

	// reqformatted, err := json.MarshalIndent(reserveFundTx, "", "    ")
	// fmt.Printf("ReserveFundRequest:\n%s\n", reqformatted)

	raw, err := types.TxToBytes(reserveFundTx)
	if err != nil {
		utils.Error("Failed to encode transaction: %v\n", err)
	}
	signedTx := hex.EncodeToString(raw)

	client := rpcc.NewRPCClient(viper.GetString(utils.CfgRemoteRPCEndpoint))

	var res *rpcc.RPCResponse
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
	}
	result := &rpc.BroadcastRawTransactionResult{}
	err = res.GetObject(result)
	if err != nil {
		utils.Error("Failed to parse server response: %v\n", err)
	}
	formatted, err := json.MarshalIndent(result, "", "    ")
	if err != nil {
		utils.Error("Failed to parse server response: %v\n", err)
	}
	//fmt.Printf("Successfully broadcasted transaction:\n%s\n", formatted)
	// Verbose output makes parsing json difficult

	// fmt.Printf("%s\n", formatted)

	return (formatted)

	//fmt.Printf("Successfully broadcasted transaction.\n")
}

func DoReserveFund() []byte {
	return (doReserveFundCmd(ReserveFundCmd, make([]string, 0)))
}

func init() {
	fmt.Println("thetaoffchaingo_tx reserve_fund.go init called.")
	// https://github.com/spf13/cobra/blob/master/user_guide.md#working-with-flags
	ReserveFundCmd.Flags().StringVar(&chainIDFlag, "chain", "", "Chain ID")
	ReserveFundCmd.Flags().StringVar(&pathFlag, "config", "./thetacli", "Path to Config")
	//rootCmd.PersistentFlags().StringVar(&chainIDFlag, "chain", "unknown", "Chain ID")
	ReserveFundCmd.Flags().StringVar(&fromFlag, "from", "", "Address to send from")
	ReserveFundCmd.Flags().Uint64Var(&seqFlag, "seq", 0, "Sequence number of the transaction")
	ReserveFundCmd.Flags().StringVar(&reserveFundInTFuelFlag, "fund", "100", "TFuel amount to reserve")
	ReserveFundCmd.Flags().StringVar(&reserveCollateralInTFuelFlag, "collateral", "101", "TFuel amount as collateral")
	ReserveFundCmd.Flags().StringVar(&feeFlag, "fee", fmt.Sprintf("%dwei", types.MinimumTransactionFeeTFuelWeiJune2021), "Fee")
	ReserveFundCmd.Flags().Uint64Var(&durationFlag, "duration", 1000, "Reserve duration")
	//	ReserveFundCmd.Flags().StringSliceVar(&resourceIDsFlag, "resource_ids", []string{}, "Resource IDs")
	ReserveFundCmd.Flags().StringSliceVar(&resourceIDsFlag, "resource_ids", []string{"rid1000001"}, "Resource IDs")
	ReserveFundCmd.Flags().StringVar(&walletFlag, "wallet", "soft", "Wallet type (soft|nano)")
	ReserveFundCmd.Flags().BoolVar(&asyncFlag, "async", false, "block until tx has been included in the blockchain")
	ReserveFundCmd.Flags().StringVar(&passwordFlag, "password", "", "password to unlock the wallet")

	ReserveFundCmd.MarkFlagRequired("chain")
	ReserveFundCmd.MarkFlagRequired("from")
	ReserveFundCmd.MarkFlagRequired("seq")
	ReserveFundCmd.MarkFlagRequired("duration")
	ReserveFundCmd.MarkFlagRequired("resource_id")

	rootCmd.Execute()
}
