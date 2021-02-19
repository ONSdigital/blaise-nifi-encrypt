package google

import (
	"cloud.google.com/go/storage"
	"context"
	"github.com/rs/zerolog/log"
	"io"
	"os"
	"time"
)

type Storage struct {
	client *storage.Client
	writer *storage.Writer
	ctx    context.Context
}

func NewStorage(ctx context.Context) Storage {
	client, err := storage.NewClient(ctx)
	if err != nil {
		log.Err(err).Msg("Cannot get GCloud Storage Bucket")
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
			log.Err(err).Msg("close bucket writer failed")
			return
		}
		log.Debug().Msg("closed bucket writer")
	}
}

func (gs *Storage) DeleteFile(file, directory string) error {

	ctx, cancel := context.WithTimeout(gs.ctx, time.Second*10)
	defer cancel()

	o := gs.client.Bucket(directory).Object(file)
	if err := o.Delete(ctx); err != nil {
		log.Warn().Msgf("delete of file %s fromm directory: %s failed", file, directory)
		return err
	}

	log.Debug().Msgf("file: %s/%s deleted", directory, file)

	return nil
}

func (gs Storage) GetReader(file, directory string) (io.ReadCloser, error) {
	readBucket := gs.client.Bucket(directory)
	readObj := readBucket.Object(file)

	storageReader, err := readObj.NewReader(gs.ctx)
	if err != nil {
		log.Err(err).Msgf("cannot create a reader")
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
