package virtualbox

import (
	"fmt"
	"github.com/mitchellh/multistep"
	"github.com/mitchellh/packer/packer"
	"log"
	"os"
)

type guestAdditionsPathTemplate struct {
	Version string
}

// This step uploads the guest additions ISO to the VM.
type stepUploadGuestAdditions struct{}

func (s *stepUploadGuestAdditions) Run(state multistep.StateBag) multistep.StepAction {
	comm := state.Get("communicator").(packer.Communicator)
	config := state.Get("config").(*config)
	driver := state.Get("driver").(Driver)
	guestAdditionsPath := state.Get("guest_additions_path").(string)
	ui := state.Get("ui").(packer.Ui)

	// If we're attaching then don't do this, since we attached.
	if config.GuestAdditionsAttach {
		log.Println("Not uploading guest additions since we're attaching.")
		return multistep.ActionContinue
	}

	version, err := driver.Version()
	if err != nil {
		state.Put("error", fmt.Errorf("Error reading version for guest additions upload: %s", err))
		return multistep.ActionHalt
	}

	f, err := os.Open(guestAdditionsPath)
	if err != nil {
		state.Put("error", fmt.Errorf("Error opening guest additions ISO: %s", err))
		return multistep.ActionHalt
	}

	tplData := &guestAdditionsPathTemplate{
		Version: version,
	}

	config.GuestAdditionsPath, err = config.tpl.Process(config.GuestAdditionsPath, tplData)
	if err != nil {
		err := fmt.Errorf("Error preparing guest additions path: %s", err)
		state.Put("error", err)
		ui.Error(err.Error())
		return multistep.ActionHalt
	}

	ui.Say("Uploading VirtualBox guest additions ISO...")
	if err := comm.Upload(config.GuestAdditionsPath, f); err != nil {
		state.Put("error", fmt.Errorf("Error uploading guest additions: %s", err))
		return multistep.ActionHalt
	}

	return multistep.ActionContinue
}

func (s *stepUploadGuestAdditions) Cleanup(state multistep.StateBag) {}
