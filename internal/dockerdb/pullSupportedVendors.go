package dockerdb

import (
	"bufio"
	"context"
	"github.com/docker/docker/api/types"
	"log"
	"quicktables/internal/globals"
)

// Pull downloads all docker images from supported vendors
func Pull() error {
	for _, ven := range globals.DownloadableVendors {
		ctx := context.TODO()
		ddb, err := New(nil)
		if err != nil {
			return err
		}

		pull, err := ddb.cli.ImagePull(ctx, ven, types.ImagePullOptions{})
		if err != nil {
			return err
		}

		scanner := bufio.NewScanner(pull)
		for scanner.Scan() {
			log.Println(scanner.Text())
		}

		err = pull.Close()
		if err != nil {
			return err
		}
	}

	return nil
}
