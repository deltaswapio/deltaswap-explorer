package metrics

import (
	"bufio"
	"context"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/deltaswapio/deltaswap-explorer/analytics/cmd/token"
	"github.com/deltaswapio/deltaswap-explorer/analytics/prices"
	"github.com/deltaswapio/deltaswap-explorer/common/client/parser"
	"github.com/deltaswapio/deltaswap-explorer/common/logger"
	"go.uber.org/zap"
)

type LineParser struct {
	converter *VaaConverter
}

// read a csv file with VAAs and convert into a decoded csv file
// ready to upload to the database
func RunVaaVolumeFromFile(inputFile, outputFile, pricesFile, vaaPayloadParserURL string) {

	ctx := context.Background()
	// build logger
	logger := logger.New("deltaswap-explorer-analytics")

	logger.Info("starting deltaswap-explorer-analytics ...")

	// create a parserVAAAPIClient
	parserVAAAPIClient, err := parser.NewParserVAAAPIClient(10, vaaPayloadParserURL, logger)
	if err != nil {
		logger.Fatal("failed to create parse vaa api client")
	}

	// create a token resolver
	tokenResolver := token.NewTokenResolver(parserVAAAPIClient, logger)

	// open input file
	f, err := os.Open(inputFile)
	if err != nil {
		logger.Fatal("opening input file", zap.Error(err))
	}
	defer f.Close()

	// create missing tokens file
	missingTokensFile := "missing_tokens.csv"
	fmissingTokens, err := os.Create(missingTokensFile)
	if err != nil {
		logger.Fatal("creating missing tokens file", zap.Error(err))
	}
	defer fmissingTokens.Close()

	//open output file for writing
	fout, err := os.Create(outputFile)
	if err != nil {
		logger.Fatal("creating output file", zap.Error(err))
	}
	defer fout.Close()

	// init price cache!
	logger.Info("loading historical prices...")
	priceCache := prices.NewCoinPricesCache(pricesFile)
	priceCache.InitCache()
	converter := NewVaaConverter(priceCache, tokenResolver.GetTransferredTokenByVaa)
	lp := NewLineParser(converter)
	logger.Info("loaded historical prices")

	r := bufio.NewReader(f)

	c := 0
	i := 0
	// read file line by line and send to workpool
	for {
		line, _, err := r.ReadLine() //loading chunk into buffer
		if err != nil {
			if err == io.EOF {
				break
			}
			logger.Fatal("a real error happened here", zap.Error(err))
		}
		nl, err := lp.ParseLine(ctx, line)
		if err != nil {
			//fmt.Printf(",")
		} else {
			c++
			fout.Write([]byte(nl))

			if c == 10000 {
				fmt.Printf(".")
				c = 0
				i := i + 1
				if i == 10 {
					fmt.Printf("\n")
					i = 0
				}
			}
		}

	}

	for k := range converter.MissingTokensCounter {
		fmissingTokens.WriteString(fmt.Sprintf("%s,%s,%d\n", k.String(), converter.MissingTokens[k], converter.MissingTokensCounter[k]))
	}

	logger.Info("missing tokens", zap.Int("count", len(converter.MissingTokens)))

	logger.Info("finished deltaswap-explorer-analytics")

}

func NewLineParser(converter *VaaConverter) *LineParser {
	return &LineParser{
		converter: converter,
	}
}

// ParseLine takes a CSV line as input, and generates a line protocol entry as output.
//
// The format for InfluxDB line protocol is: vaa,tags fields timestamp
func (lp *LineParser) ParseLine(ctx context.Context, line []byte) (string, error) {

	// Parse the VAA and payload
	tt := strings.Split(string(line), ",")
	if len(tt) != 2 {
		return "", fmt.Errorf("expected line to have two tokens, but has %d: %s", len(tt), line)
	}
	vaaBytes, err := hex.DecodeString(tt[1])
	if err != nil {
		return "", fmt.Errorf("error decoding: %v", err)
	}

	return lp.converter.Convert(ctx, vaaBytes)
}
