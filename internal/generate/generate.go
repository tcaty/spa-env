package generate

import (
	"fmt"
	"strings"

	"github.com/tcaty/spa-env/internal/common/dotenv"
	"github.com/tcaty/spa-env/internal/common/log"
)

func Generate(workdir, dotenvDev, dotenvProd, keyPrefix, placeholderPrefix string, enableComments bool) error {
	dotenvDevMap, err := dotenv.Read(workdir, dotenvDev)
	if err != nil {
		return fmt.Errorf("unable to read %s file: %v", dotenvDev, err)
	}

	dotenvProdMap := generateDotenvProdMap(dotenvDevMap, keyPrefix, placeholderPrefix)

	if err := dotenv.Write(dotenvProdMap, workdir, dotenvProd, placeholderPrefix, enableComments); err != nil {
		return fmt.Errorf("unable to write %s file: %v", dotenvProd, err)
	}

	return nil
}

func generateDotenvProdMap(dotenvDevMap map[string]string, keyPrefix, placeholderPrefix string) map[string]string {
	dotenvProdMap := make(map[string]string)

	for key := range dotenvDevMap {
		if strings.HasPrefix(key, keyPrefix) {
			placeholder := fmt.Sprintf("%s%s", placeholderPrefix, strings.TrimPrefix(key, keyPrefix))
			dotenvProdMap[key] = placeholder
		} else {
			log.Debug(
				"skip variable cause it has no prefix",
				"key", key,
				"prefix", keyPrefix,
			)
		}
	}

	return dotenvProdMap
}
