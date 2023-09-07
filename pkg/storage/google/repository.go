package google

import (
	"context"
	"io"

	"cloud.google.com/go/storage"
	"github.com/rs/zerolog/log"
)

type Storage struct {
	client *storage.Client
	writer *storage.Writer
	ctx    context.Context
}

func NewStorage(ctx context.Context) (Storage, error) {
	client, err := storage.NewClient(ctx)
	if err != nil {
		log.Err(err).Msgf("Cannot get GCloud Storage Bucket")
		return Storage{}, err
	}

	return Storage{ctx: ctx, client: client}, nil
}

func (gs *Storage) CreateFile(location, destinationFile string) (io.Writer, error) {

	log.Debug().Msgf("Creating %s/%s", location, destinationFile)

	bh := gs.client.Bucket(location)
	// Next check if the bucket exists
	if _, err := bh.Attrs(gs.ctx); err != nil {
		return nil, err
	}

	obj := bh.Object(destinationFile)

	gs.writer = obj.NewWriter(gs.ctx)

	log.Debug().Msgf("File %s/%s created", location, destinationFile)

	return gs.writer, nil
}

func (gs *Storage) CloseFile() {
	if gs.writer != nil {
		err := gs.writer.Close()
		if err != nil {
			log.Err(err).Msgf("Close bucket writer failed")
			return
		}
		log.Debug().Msg("Closed bucket writer")
	}
}

func (gs Storage) GetReader(file, directory string) (io.ReadCloser, error) {
	readBucket := gs.client.Bucket(directory)
	readObj := readBucket.Object(file)

	storageReader, err := readObj.NewReader(gs.ctx)
	if err != nil {
		log.Err(err).Msgf("Cannot create a reader")
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
