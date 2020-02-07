package kubeconfig

import (
	"fmt"

	"github.com/rs/zerolog/log"
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
)

// Save writes out the updated main config to the filesystem
func (k *KConf) Save() error {
	err := clientcmd.WriteToFile(k.Config, MainConfigPath)
	if err != nil {
		log.Debug().Msgf("Error while writing main config: %v", err)
		return err
	}
	return nil
}

// Merge takes a config and combines it into a config file
func (k *KConf) Merge(config *clientcmdapi.Config) error {
	renamedClusters := make(map[string]string)
	renamedUsers := make(map[string]string)

	for clsName, cls := range config.Clusters {
		added, err := k.AddCluster(clsName, cls)
		if err != nil {
			return err
		}
		if added != "" { // this cluster was newly added
			if added != clsName {
				fmt.Printf("Renamed cluster '%s' to '%s'\n", clsName, added)
				renamedClusters[clsName] = added
			}
			fmt.Printf("Added cluster '%s'\n", added)
		}
	}
	for uName, user := range config.AuthInfos {
		added, err := k.AddUser(uName, user)
		if err != nil {
			return err
		}
		if added != "" { // this user was newly added
			if added != uName {
				fmt.Printf("Renamed user '%s' to '%s'\n", uName, added)
				renamedUsers[uName] = added
			}
			fmt.Printf("Added user '%s'\n", added)
		}
	}
	for ctxName, ctx := range config.Contexts {
		// if anything got renamed, we need to update this context
		if renamed, ok := renamedClusters[ctx.Cluster]; ok {
			ctx.Cluster = renamed
		}
		if renamed, ok := renamedUsers[ctx.AuthInfo]; ok {
			ctx.AuthInfo = renamed
		}

		added, err := k.AddContext(ctxName, ctx)
		if err != nil {
			return err
		}
		if added != "" { // this context was newly added
			if added != ctxName {
				fmt.Printf("Renamed context '%s' to '%s'\n", ctxName, added)
			}
			fmt.Printf("Added context '%s'\n", added)
		}
	}
	return nil
}

// Export extracts a single context from a kubeconfig file
func (k *KConf) Export(name string) (*clientcmdapi.Context, error) {
	for ctxName, ctx := range k.Contexts {
		if ctxName == name {
			return ctx, nil
		}
	}
	return nil, fmt.Errorf("Could not find context '%s'", name)
}
