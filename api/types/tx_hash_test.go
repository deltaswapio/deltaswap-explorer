package types

import "testing"

// TestParseTxHash tests the ParseTxHash function.
func TestParseTxHash(t *testing.T) {

	// a table containing several test cases
	tcs := []struct {
		input             string
		output            string
		isSolanaTxHash    bool
		isDeltaswapTxHash bool
	}{
		{
			// Invalid Solana hash - 86 characters (too short)
			input: "VKrJx5ak3amnpY5EXiqfu6pnrzxHTLU95m9vfbYnGSSLQrkRzb4tm4NztCGeLcJxieXQYnqddUwoaEsDRTRh57",
		},
		{
			// Valid Solana hash - 88 characters
			input:          "2maR6uDZzroV7JFF76rp5QR4CFP1PFUe76VRE8gF8QtWRifpGAKJQo4SQDBNs3TAM9RrchJhnJ644jUL2yfagZco",
			output:         "2maR6uDZzroV7JFF76rp5QR4CFP1PFUe76VRE8gF8QtWRifpGAKJQo4SQDBNs3TAM9RrchJhnJ644jUL2yfagZco",
			isSolanaTxHash: true,
		},
		{
			// Valid Solana hash - 87 characters
			input:          "VKrJx5ak3amnpY5EXiqfu6pnrzxHTLU95m9vfbYnGSSLQrkRzb4tm4NztCGeLcJxieXQYnqddUwoaEsDRTRh57R",
			output:         "VKrJx5ak3amnpY5EXiqfu6pnrzxHTLU95m9vfbYnGSSLQrkRzb4tm4NztCGeLcJxieXQYnqddUwoaEsDRTRh57R",
			isSolanaTxHash: true,
		},
		{
			// Invalid Solana hash - 89 characters (too long)
			input: "2maR6uDZzroV7JFF76rp5QR4CFP1PFUe76VRE8gF8QtWRifpGAKJQo4SQDBNs3TAM9RrchJhnJ644jUL2yfagZco2",
		},
		{
			// Invalid Sui hash - 42 characters (too short)
			input: "cVWa8xZtWbTxXQGLQaYquwmChE2sQYxFNGnHmp6oXX",
		},
		{
			// Valid Sui hash - 43 characters
			input:             "cVWa8xZtWbTxXQGLQaYquwmChE2sQYxFNGnHmp6oXXL",
			output:            "cVWa8xZtWbTxXQGLQaYquwmChE2sQYxFNGnHmp6oXXL",
			isDeltaswapTxHash: true,
		},
		{
			// Valid Sui hash - 44 characters
			input:             "9yQWLTNmFkwZ6CdK3QXhk8utKr42n3Eh1CFFBWcdCeJC",
			output:            "9yQWLTNmFkwZ6CdK3QXhk8utKr42n3Eh1CFFBWcdCeJC",
			isDeltaswapTxHash: true,
		},
		{
			// Invalid Sui hash - 45 characters (too long)
			input: "9yQWLTNmFkwZ6CdK3QXhk8utKr42n3Eh1CFFBWcdCeJC9",
		},
		{
			// Invalid Deltaswap hash - 63 characters (too short)
			input: "f77f8b44f35ff047a74ee8235ce007afbab357d4e30010d51b6f6990f921637",
		},
		{
			// Deltaswap hash with 0x prefix
			input:             "0x3f77f8b44f35ff047a74ee8235ce007afbab357d4e30010d51b6f6990f921637",
			output:            "3f77f8b44f35ff047a74ee8235ce007afbab357d4e30010d51b6f6990f921637",
			isDeltaswapTxHash: true,
		},
		{
			// Deltaswap hash with 0X prefix
			input:             "0X3F77F8B44F35FF047A74EE8235CE007AFBAB357D4E30010D51B6F6990F921637",
			output:            "3f77f8b44f35ff047a74ee8235ce007afbab357d4e30010d51b6f6990f921637",
			isDeltaswapTxHash: true,
		},
		{
			// Deltaswap hash with no prefix
			input:             "3f77f8b44f35ff047a74ee8235ce007afbab357d4e30010d51b6f6990f921637",
			output:            "3f77f8b44f35ff047a74ee8235ce007afbab357d4e30010d51b6f6990f921637",
			isDeltaswapTxHash: true,
		},
		{
			// Invalid Deltaswap hash - 65 characters (too long)
			input: "33f77f8b44f35ff047a74ee8235ce007afbab357d4e30010d51b6f6990f921637",
		},
		{
			// A bunch of random characters
			input: "434234i32042oiu08d8sauf0suif",
		},
	}

	// run each test case in the table
	for i := range tcs {
		tc := tcs[i]

		// try to parse the hash
		txHash, err := ParseTxHash(tc.input)
		if tc.output == "" && err == nil {
			t.Fatalf("expected parseTxHash(%s) to fail", tc.input)
		} else if tc.output != "" && err != nil {
			t.Fatalf("parseTxHash(%s) failed with error %v", tc.input, err)
		}

		if tc.output == "" {
			continue
		}

		// make assertions about the output struct
		if tc.output != txHash.String() {
			t.Fatalf("expected TxHash.String()=%s, got %s", tc.output, txHash.String())
		}
		if tc.isSolanaTxHash != txHash.IsSolanaTxHash() {
			t.Fatalf("expected TxHash.IsSolanaHash()=%t, but got %t", tc.isSolanaTxHash, txHash.IsSolanaTxHash())
		}
		if tc.isDeltaswapTxHash != txHash.IsDeltaswapTxHash() {
			t.Fatalf("expected TxHash.IsDeltaswapHash()=%t, but got %t", tc.isDeltaswapTxHash, txHash.IsDeltaswapTxHash())
		}

	}

}
