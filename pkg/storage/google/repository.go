package google

import (
	"context"
	"io"
	"os"

	"cloud.google.com/go/storage"
	"github.com/rs/zerolog/log"
)

type Storage struct {
	client *storage.Client
	writer *storage.Writer
	ctx    context.Context
}

func NewStorage(ctx context.Context) Storage {
	client, err := storage.NewClient(ctx)
	if err != nil {
		log.Error().Msgf("Cannot get GCloud Storage Bucket: %s", err)
		os.Exit(1)
	}

	return Storage{ctx: ctx, client: client}
}

func (gs *Storage) CreateFile(location, destinationFile string) (io.Writer, error) {

	log.Debug().Msgf("creating %s/%s", location, destinationFile)

	bh := gs.client.Bucket(location)
	// Next check if the bucket exists
	if _, err := bh.Attrs(gs.ctx); err != nil {
		return nil, err
	}

	obj := bh.Object(destinationFile)

	gs.writer = obj.NewWriter(gs.ctx)

	log.Debug().Msgf("file %s/%s created", location, destinationFile)

	return gs.writer, nil
}

func (gs *Storage) CloseFile() {
	if gs.writer != nil {
		err := gs.writer.Close()
		if err != nil {
			log.Error().Msgf("close bucket writer failed: %s", err)
			return
		}
		log.Debug().Msg("closed bucket writer")
	}
}

func (gs Storage) GetReader(file, directory string) (io.ReadCloser, error) {
	readBucket := gs.client.Bucket(directory)
	readObj := readBucket.Object(file)

	storageReader, err := readObj.NewReader(gs.ctx)
	if err != nil {
		log.Error().Msgf("cannot create a reader: %s", err)
		return nil, err
	}

	return storageReader, nil
}

func (gs Storage) GetWriter(file, directory string) io.WriteCloser {
	writeBucket := gs.client.Bucket(directory)
	writeObj := writeBucket.Object(file)

	storageWriter := writeObj.NewWriter(gs.ctx)
	return storageWriter
}
