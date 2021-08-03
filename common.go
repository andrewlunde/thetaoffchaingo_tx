package tx

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/thetatoken/theta/cmd/thetacli/cmd/utils"
	//	"github.com/ybbus/jsonrpc"
)

func SetChainID(cmd *cobra.Command, chainstr string) {
	err := cmd.Flags().Set("chain", chainstr)
	if err != nil {
		utils.Error("Failed to set Flag: %v\n", err)
	}
}

func SetConfigPath(cmd *cobra.Command, configpathstr string) {
	err := cmd.Flags().Set("config", configpathstr)
	if err != nil {
		utils.Error("Failed to set Flag: %v\n", err)
	}
}

func SetFee(cmd *cobra.Command, fee string) {
	err := cmd.Flags().Set("fee", fee)
	if err != nil {
		utils.Error("Failed to set Flag: %v\n", err)
	}
}

func SetFrom(cmd *cobra.Command, addr string) {
	err := cmd.Flags().Set("from", addr)
	if err != nil {
		utils.Error("Failed to set Flag: %v\n", err)
	}
}

func SetTo(cmd *cobra.Command, addr string) {
	err := cmd.Flags().Set("to", addr)
	if err != nil {
		utils.Error("Failed to set Flag: %v\n", err)
	}
}

func SetSeq(cmd *cobra.Command, seqnum string) {
	err := cmd.Flags().Set("seq", seqnum)
	if err != nil {
		utils.Error("Failed to set Flag: %v\n", err)
	}
}

func SetPaymentSeq(cmd *cobra.Command, seqnum string) {
	err := cmd.Flags().Set("payment_seq", seqnum)
	if err != nil {
		utils.Error("Failed to set Flag: %v\n", err)
	}
}

func SetReserveSeq(cmd *cobra.Command, seqnum string) {
	err := cmd.Flags().Set("reserve_seq", seqnum)
	if err != nil {
		utils.Error("Failed to set Flag: %v\n", err)
	}
}

// func SetTFuelAmt(cmd *cobra.Command, amt uint64) {
// 	err := cmd.Flags().Set("tfuel", fmt.Sprintf("%d", amt))
// 	if err != nil {
// 		utils.Error("Failed to set Flag: %v\n", err)
// 	}
// }

func SetTFuelAmt(cmd *cobra.Command, amt string) {
	err := cmd.Flags().Set("tfuel", amt)
	if err != nil {
		utils.Error("Failed to set Flag: %v\n", err)
	}
}

func SetResourceID(cmd *cobra.Command, resourceid string) {
	err := cmd.Flags().Set("resource_id", resourceid)
	if err != nil {
		utils.Error("Failed to set Flag: %v\n", err)
	}
}

// func SetResourceIDs(cmd *cobra.Command, resourceIDs string) {
// 	err := cmd.Flags().Set("resource_ids", resourceIDs)
// 	if err != nil {
// 		utils.Error("Failed to set Flag: %v\n", err)
// 	}
// }

func SetSourceSig(cmd *cobra.Command, srcsig string) {
	err := cmd.Flags().Set("src_sig", srcsig)
	if err != nil {
		utils.Error("Failed to set Flag: %v\n", err)
	}
}

func SetOnChain(cmd *cobra.Command, onchain bool) {
	err := cmd.Flags().Set("on_chain", fmt.Sprintf("%v", onchain))
	if err != nil {
		utils.Error("Failed to set Flag: %v\n", err)
	}
}

func SetPassword(cmd *cobra.Command, passstr string) {
	err := cmd.Flags().Set("password", passstr)
	if err != nil {
		utils.Error("Failed to set Flag: %v\n", err)
	}
}

func SetDebugging(cmd *cobra.Command, debuggin bool) {
	err := cmd.Flags().Set("debugging", fmt.Sprintf("%v", debuggin))
	if err != nil {
		utils.Error("Failed to set Flag: %v\n", err)
	}
}

func SetDuration(cmd *cobra.Command, durationnum string) {
	err := cmd.Flags().Set("duration", durationnum)
	if err != nil {
		utils.Error("Failed to set Flag: %v\n", err)
	}
}

var rootCmd = &cobra.Command{Use: "app"}

var Alice = "0x2E833968E5bB786Ae419c4d13189fB081Cc43bab"
var Bob = "0x70f587259738cB626A1720Af7038B8DcDb6a42a0"
var Carol = "0xcd56123D0c5D6C1Ba4D39367b88cba61D93F5405"

func init() {
	fmt.Println("thetaoffchaingo_tx common.go init called.")

	rootCmd.Execute()
}
