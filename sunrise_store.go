package sunrise

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	plasma "github.com/ethereum-optimism/optimism/op-plasma"
	"github.com/ethereum/go-ethereum/log"
	api "github.com/sunriselayer/sunrise-data/api"
)

const VersionByte = 0x0c

type SunriseConfig struct {
	URL       string
	Namespace []byte
}

// SunriseStore implements DAStorage with sunrise backend
type SunriseStore struct {
	Log        log.Logger
	Config     SunriseConfig
	GetTimeout time.Duration
	Namespace  []byte
}

// NewSunriseStore returns a sunrise store.
func NewSunriseStore(cfg SunriseConfig) *SunriseStore {
	Log := log.New()

	return &SunriseStore{
		Log:        Log,
		Config:     cfg,
		GetTimeout: time.Minute,
		Namespace:  cfg.Namespace,
	}
}

func (d *SunriseStore) Get(ctx context.Context, key []byte) ([]byte, error) {
	log.Info("sunrise: blob request", "id", hex.EncodeToString(key))
	ctx, cancel := context.WithTimeout(context.Background(), d.GetTimeout)

	resp, err := http.Get(fmt.Sprintf("%s/api/get-blob?metadata_uri=%s", d.Config.URL, key))
	if err != nil {
		return nil, fmt.Errorf("sunrise: failed to get blob: %w", err)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("sunrise: failed to read response body: %w", err)
	}

	blobResp := api.GetBlobResponse{}
	err = json.Unmarshal(body, &blobResp)
	if err != nil {
		return nil, fmt.Errorf("sunrise: failed to unmarshal response body: %w", err)
	}
	blobs := blobResp.Blob

	cancel()
	if err != nil || len(blobs) == 0 {
		return nil, fmt.Errorf("sunrise: failed to resolve frame: %w", err)
	}

	return []byte(blobs), nil
}

func (d *SunriseStore) Put(ctx context.Context, data []byte) ([]byte, error) {
	publishReq := api.PublishRequest{
		Blob:             base64.StdEncoding.EncodeToString(data),
		DataShardCount:   10,
		ParityShardCount: 10,
		Protocol:         "ipfs",
	}
	jsonData, err := json.Marshal(publishReq)
	if err != nil {
		return nil, fmt.Errorf("sunrise: failed to marshal publish request: %w", err)
	}
	resp, err := http.Post(fmt.Sprintf("%s/api/publish", d.Config.URL), "application/json", bytes.NewBuffer(jsonData))

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("sunrise: failed to read response body: %w", err)
	}

	publishResp := api.PublishResponse{}
	err = json.Unmarshal(body, &publishResp)
	if err != nil {
		return nil, fmt.Errorf("sunrise: failed to unmarshal response body: %w", err)
	}

	if err == nil {
		d.Log.Info("sunrise: blob successfully submitted", "uri", publishResp.MetadataUri)
		commitment := plasma.NewGenericCommitment(append([]byte{VersionByte}, []byte(publishResp.MetadataUri)...))
		return commitment.Encode(), nil
	}
	return nil, err
}
