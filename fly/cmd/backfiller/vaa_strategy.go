package main

import (
	"context"
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/deltaswapio/deltaswap-explorer/fly/storage"
	"github.com/deltaswapio/deltaswap/sdk/vaa"
)

func workerVaa(ctx context.Context, repo *storage.Repository, line string) error {
	tokens := strings.Split(line, ",")
	//fmt.Printf("bcid %s, emmiter %s, seq %s\n", header[0], header[1], header[2])

	if len(tokens) != 2 {
		//fmt.Printf("invalid line: %s", line)
		return fmt.Errorf("invalid line: %s", line)
	}

	data, err := hex.DecodeString(tokens[1])
	if err != nil {
		return fmt.Errorf("error decoding: %v", err)
	}

	v, err := vaa.Unmarshal(data)
	if err != nil {
		return fmt.Errorf("error unmarshaling vaa: %v", err)
	}

	err = repo.UpsertVaa(ctx, v, data)
	if err != nil {
		return fmt.Errorf("error upserting vaa: %v\n", err)
	}

	return nil
}
