package main

import (
	"fmt"

	eigentypes "github.com/Layr-Labs/eigenlayer-cli/pkg/types"
	eigenutils "github.com/Layr-Labs/eigenlayer-cli/pkg/utils"
	eigensdktypes "github.com/Layr-Labs/eigensdk-go/types"
)

func main() {
	test := eigentypes.PrivateKeySigner
	fmt.Println(test)
	test1 := eigenutils.EmojiInfo
	fmt.Println(test1)
	test2 := eigensdktypes.EigenPromNamespace
	fmt.Println(test2)
}
