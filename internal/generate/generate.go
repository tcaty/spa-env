package generate

import (
	"fmt"

	"github.com/tcaty/spa-env/internal/common/dotenv"
)

func Generate(workdir, dotenvDev, dotenvProd, keyPrefix, placeholderPrefix string, enableComments bool) error {
	dotenvDevEntries, err := dotenv.Read(workdir, dotenvDev, keyPrefix, placeholderPrefix)
	if err != nil {
		return fmt.Errorf("unable to read %s file: %v", dotenvDev, err)
	}

	if err := dotenv.Write(workdir, dotenvProd, dotenvDevEntries, enableComments); err != nil {
		return fmt.Errorf("unable to write %s file: %v", dotenvProd, err)
	}

	return nil
}
