package sunrise

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"testing"

	cli "github.com/ethereum-optimism/optimism/op-alt-da"
	"github.com/ethereum-optimism/optimism/op-service/testlog"
	"github.com/ethereum/go-ethereum/log"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/require"
)

var (
	URL              string
	DataShardCount   int
	ParityShardCount int
)

func Check() error {
	if URL == "" {
		return errors.New("no url provided")
	}
	if DataShardCount == 0 {
		return errors.New("invalid data shard count")
	}
	if ParityShardCount == 0 {
		return errors.New("invalid parity shard count")
	}
	return nil
}

func TestSunriseDAClientService(t *testing.T) {
	logger := testlog.Logger(t, log.LevelDebug)
	if err := godotenv.Load(); err != nil {
		logger.Crit("Error loading .env file: ", err)
	}
	URL = os.Getenv("SUNRISE_SERVER")
	dataShardCount, err := strconv.ParseInt(os.Getenv("SUNRISE_DATA_SHARD_COUNT"), 10, 64)
	if err != nil {
		log.Crit("Error parsing data shard count: ", err)
	}
	DataShardCount = int(dataShardCount)
	parityShardCount, err := strconv.ParseInt(os.Getenv("SUNRISE_PARITY_SHARD_COUNT"), 10, 64)
	if err != nil {
		log.Crit("Error parsing parity shard count: ", err)
	}
	ParityShardCount = int(parityShardCount)

	err = Check()
	if err != nil {
		panic(err)
	}

	store := NewSunriseStore(SunriseConfig{
		URL:              URL,
		DataShardCount:   DataShardCount,
		ParityShardCount: ParityShardCount,
	}, logger)

	ctx := context.Background()

	server := NewSunriseServer("127.0.0.1", 0, store, logger)

	require.NoError(t, server.Start())

	cfg := cli.CLIConfig{
		Enabled:      true,
		DAServerURL:  fmt.Sprintf("http://%s", server.Endpoint()),
		VerifyOnRead: false,
		GenericDA:    true,
	}
	require.NoError(t, cfg.Check())

	client := cfg.NewDAClient()

	rng := rand.New(rand.NewSource(1234))

	input := RandomData(rng, 2000)

	comm, err := client.SetInput(ctx, input)
	println("comm", comm)
	require.NoError(t, err)

	stored, err := client.GetInput(ctx, comm)
	require.NoError(t, err)
	require.Equal(t, stored, input)
}

func RandomData(rng *rand.Rand, size int) []byte {
	out := make([]byte, size)
	rng.Read(out)
	return out
}
