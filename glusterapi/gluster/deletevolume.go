package gluster

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/SchweizerischeBundesbahnen/ssp-backend/glusterapi/models"
)

func deleteVolume(pvName string) error {
	if len(pvName) == 0 {
		return errors.New("Not all input values provided")
	}

	if err := deleteLvOnAllServers(pvName); err != nil {
		return err
	}

	return nil
}

func deleteLvOnAllServers(pvName string) error {
	// Delete the lv on all other gluster servers
	if err := deleteLvOnOtherServers(pvName); err != nil {
		return err
	}

	// Delete the lv locally
	if err := deleteLvLocally(pvName); err != nil {
		return err
	}

	return nil
}

func deleteLvOnOtherServers(pvName string) error {
	remotes, err := getGlusterPeerServers()
	if err != nil {
		return err
	}

	// Execute the commands remote via API
	client := &http.Client{}
	for _, r := range remotes {
		p := models.DeleteVolumeCommand{
			PvName: pvName,
		}
		b := new(bytes.Buffer)

		if err = json.NewEncoder(b).Encode(p); err != nil {
			log.Println("Error encoding json", err.Error())
			return errors.New(commandExecutionError)
		}

		log.Println("Going to delete lv on remote:", r)

		req, _ := http.NewRequest("POST", fmt.Sprintf("http://%v:%v/sec/lv/delete", r, Port), b)
		req.SetBasicAuth("GLUSTER_API", Secret)

		resp, err := client.Do(req)
		if err != nil || resp.StatusCode != http.StatusOK {
			if resp != nil {
				log.Println("Remote did not respond with OK", resp.StatusCode)
			} else {
				log.Println("Connection to remote not possible", r, err.Error())
			}
			return errors.New(commandExecutionError)
		}
		resp.Body.Close()
	}

	return nil
}

func deleteLvLocally(pvName string) error {
	//lvName := fmt.Sprintf("lv_%v", pvName)

	commands := []string{
	// TODO: correct commands
	// Grow lv
	//fmt.Sprintf("lvextend -L %v /dev/%v/%v", newSize, VgName, lvName),

	// Grow file system
	//fmt.Sprintf("xfs_growfs /dev/%v/%v", VgName, lvName),
	}

	if err := executeCommandsLocally(commands); err != nil {
		return err
	}

	return nil
}
