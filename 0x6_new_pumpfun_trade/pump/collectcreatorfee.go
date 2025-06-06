// Code generated by https://github.com/daog1/solana-anchor-go. DO NOT EDIT.

package pump

import (
	"errors"
	ag_binary "github.com/gagliardetto/binary"
	ag_solanago "github.com/gagliardetto/solana-go"
	ag_format "github.com/gagliardetto/solana-go/text/format"
	ag_treeout "github.com/gagliardetto/treeout"
)

// Collects creator_fee from creator_vault to the coin creator account
type CollectCreatorFee struct {

	// [0] = [WRITE, SIGNER] creator
	//
	// [1] = [WRITE] creator_vault
	//
	// [2] = [] system_program
	//
	// [3] = [] event_authority
	//
	// [4] = [] program
	ag_solanago.AccountMetaSlice `bin:"-"`
}

// NewCollectCreatorFeeInstructionBuilder creates a new `CollectCreatorFee` instruction builder.
func NewCollectCreatorFeeInstructionBuilder() *CollectCreatorFee {
	nd := &CollectCreatorFee{
		AccountMetaSlice: make(ag_solanago.AccountMetaSlice, 5),
	}
	nd.AccountMetaSlice[2] = ag_solanago.Meta(Addresses["11111111111111111111111111111111"])
	return nd
}

// NewCollectCreatorFeeInstructionBuilderExt creates a new `CollectCreatorFee` instruction builder.
func NewCollectCreatorFeeInstructionBuilderExt(remainingAccounts int) *CollectCreatorFee {
	nd := &CollectCreatorFee{
		AccountMetaSlice: make(ag_solanago.AccountMetaSlice, 5+remainingAccounts),
	}
	nd.AccountMetaSlice[2] = ag_solanago.Meta(Addresses["11111111111111111111111111111111"])
	return nd
}

// SetCreatorAccount sets the "creator" account.
func (inst *CollectCreatorFee) SetCreatorAccount(creator ag_solanago.PublicKey) *CollectCreatorFee {
	inst.AccountMetaSlice[0] = ag_solanago.Meta(creator).WRITE().SIGNER()
	return inst
}

// GetCreatorAccount gets the "creator" account.
func (inst *CollectCreatorFee) GetCreatorAccount() *ag_solanago.AccountMeta {
	return inst.AccountMetaSlice.Get(0)
}

// SetCreatorVaultAccount sets the "creator_vault" account.
func (inst *CollectCreatorFee) SetCreatorVaultAccount(creatorVault ag_solanago.PublicKey) *CollectCreatorFee {
	inst.AccountMetaSlice[1] = ag_solanago.Meta(creatorVault).WRITE()
	return inst
}

func (inst *CollectCreatorFee) findFindCreatorVaultAddress(creator ag_solanago.PublicKey, knownBumpSeed uint8) (pda ag_solanago.PublicKey, bumpSeed uint8, err error) {
	var seeds [][]byte
	// const: 0x63726561746f722d7661756c74
	seeds = append(seeds, []byte{byte(0x63), byte(0x72), byte(0x65), byte(0x61), byte(0x74), byte(0x6f), byte(0x72), byte(0x2d), byte(0x76), byte(0x61), byte(0x75), byte(0x6c), byte(0x74)})
	// path: creator
	seeds = append(seeds, creator.Bytes())

	if knownBumpSeed != 0 {
		seeds = append(seeds, []byte{byte(bumpSeed)})
		pda, err = ag_solanago.CreateProgramAddress(seeds, ProgramID)
	} else {
		pda, bumpSeed, err = ag_solanago.FindProgramAddress(seeds, ProgramID)
	}
	return
}

// FindCreatorVaultAddressWithBumpSeed calculates CreatorVault account address with given seeds and a known bump seed.
func (inst *CollectCreatorFee) FindCreatorVaultAddressWithBumpSeed(creator ag_solanago.PublicKey, bumpSeed uint8) (pda ag_solanago.PublicKey, err error) {
	pda, _, err = inst.findFindCreatorVaultAddress(creator, bumpSeed)
	return
}

func (inst *CollectCreatorFee) MustFindCreatorVaultAddressWithBumpSeed(creator ag_solanago.PublicKey, bumpSeed uint8) (pda ag_solanago.PublicKey) {
	pda, _, err := inst.findFindCreatorVaultAddress(creator, bumpSeed)
	if err != nil {
		panic(err)
	}
	return
}

// FindCreatorVaultAddress finds CreatorVault account address with given seeds.
func (inst *CollectCreatorFee) FindCreatorVaultAddress(creator ag_solanago.PublicKey) (pda ag_solanago.PublicKey, bumpSeed uint8, err error) {
	pda, bumpSeed, err = inst.findFindCreatorVaultAddress(creator, 0)
	return
}

func (inst *CollectCreatorFee) MustFindCreatorVaultAddress(creator ag_solanago.PublicKey) (pda ag_solanago.PublicKey) {
	pda, _, err := inst.findFindCreatorVaultAddress(creator, 0)
	if err != nil {
		panic(err)
	}
	return
}

// GetCreatorVaultAccount gets the "creator_vault" account.
func (inst *CollectCreatorFee) GetCreatorVaultAccount() *ag_solanago.AccountMeta {
	return inst.AccountMetaSlice.Get(1)
}

// SetSystemProgramAccount sets the "system_program" account.
func (inst *CollectCreatorFee) SetSystemProgramAccount(systemProgram ag_solanago.PublicKey) *CollectCreatorFee {
	inst.AccountMetaSlice[2] = ag_solanago.Meta(systemProgram)
	return inst
}

// GetSystemProgramAccount gets the "system_program" account.
func (inst *CollectCreatorFee) GetSystemProgramAccount() *ag_solanago.AccountMeta {
	return inst.AccountMetaSlice.Get(2)
}

// SetEventAuthorityAccount sets the "event_authority" account.
func (inst *CollectCreatorFee) SetEventAuthorityAccount(eventAuthority ag_solanago.PublicKey) *CollectCreatorFee {
	inst.AccountMetaSlice[3] = ag_solanago.Meta(eventAuthority)
	return inst
}

func (inst *CollectCreatorFee) findFindEventAuthorityAddress(knownBumpSeed uint8) (pda ag_solanago.PublicKey, bumpSeed uint8, err error) {
	var seeds [][]byte
	// const: 0x5f5f6576656e745f617574686f72697479
	seeds = append(seeds, []byte{byte(0x5f), byte(0x5f), byte(0x65), byte(0x76), byte(0x65), byte(0x6e), byte(0x74), byte(0x5f), byte(0x61), byte(0x75), byte(0x74), byte(0x68), byte(0x6f), byte(0x72), byte(0x69), byte(0x74), byte(0x79)})

	if knownBumpSeed != 0 {
		seeds = append(seeds, []byte{byte(bumpSeed)})
		pda, err = ag_solanago.CreateProgramAddress(seeds, ProgramID)
	} else {
		pda, bumpSeed, err = ag_solanago.FindProgramAddress(seeds, ProgramID)
	}
	return
}

// FindEventAuthorityAddressWithBumpSeed calculates EventAuthority account address with given seeds and a known bump seed.
func (inst *CollectCreatorFee) FindEventAuthorityAddressWithBumpSeed(bumpSeed uint8) (pda ag_solanago.PublicKey, err error) {
	pda, _, err = inst.findFindEventAuthorityAddress(bumpSeed)
	return
}

func (inst *CollectCreatorFee) MustFindEventAuthorityAddressWithBumpSeed(bumpSeed uint8) (pda ag_solanago.PublicKey) {
	pda, _, err := inst.findFindEventAuthorityAddress(bumpSeed)
	if err != nil {
		panic(err)
	}
	return
}

// FindEventAuthorityAddress finds EventAuthority account address with given seeds.
func (inst *CollectCreatorFee) FindEventAuthorityAddress() (pda ag_solanago.PublicKey, bumpSeed uint8, err error) {
	pda, bumpSeed, err = inst.findFindEventAuthorityAddress(0)
	return
}

func (inst *CollectCreatorFee) MustFindEventAuthorityAddress() (pda ag_solanago.PublicKey) {
	pda, _, err := inst.findFindEventAuthorityAddress(0)
	if err != nil {
		panic(err)
	}
	return
}

// GetEventAuthorityAccount gets the "event_authority" account.
func (inst *CollectCreatorFee) GetEventAuthorityAccount() *ag_solanago.AccountMeta {
	return inst.AccountMetaSlice.Get(3)
}

// SetProgramAccount sets the "program" account.
func (inst *CollectCreatorFee) SetProgramAccount(program ag_solanago.PublicKey) *CollectCreatorFee {
	inst.AccountMetaSlice[4] = ag_solanago.Meta(program)
	return inst
}

// GetProgramAccount gets the "program" account.
func (inst *CollectCreatorFee) GetProgramAccount() *ag_solanago.AccountMeta {
	return inst.AccountMetaSlice.Get(4)
}

func (inst *CollectCreatorFee) AddRemainingAccounts(remainingAccounts []ag_solanago.PublicKey) *CollectCreatorFee {
	accounts := 5
	for i, _ := range remainingAccounts {
		index := accounts + i
		inst.AccountMetaSlice[index] = ag_solanago.Meta(remainingAccounts[i]).WRITE()
	}
	return inst
}

func (inst CollectCreatorFee) Build() *Instruction {
	return &Instruction{BaseVariant: ag_binary.BaseVariant{
		Impl:   inst,
		TypeID: Instruction_CollectCreatorFee,
	}}
}

// ValidateAndBuild validates the instruction parameters and accounts;
// if there is a validation error, it returns the error.
// Otherwise, it builds and returns the instruction.
func (inst CollectCreatorFee) ValidateAndBuild() (*Instruction, error) {
	if err := inst.Validate(); err != nil {
		return nil, err
	}
	return inst.Build(), nil
}

func (inst *CollectCreatorFee) Validate() error {
	// Check whether all (required) accounts are set:
	{
		if inst.AccountMetaSlice[0] == nil {
			return errors.New("accounts.Creator is not set")
		}
		if inst.AccountMetaSlice[1] == nil {
			return errors.New("accounts.CreatorVault is not set")
		}
		if inst.AccountMetaSlice[2] == nil {
			return errors.New("accounts.SystemProgram is not set")
		}
		if inst.AccountMetaSlice[3] == nil {
			return errors.New("accounts.EventAuthority is not set")
		}
		if inst.AccountMetaSlice[4] == nil {
			return errors.New("accounts.Program is not set")
		}
	}
	return nil
}

func (inst *CollectCreatorFee) EncodeToTree(parent ag_treeout.Branches) {
	parent.Child(ag_format.Program(ProgramName, ProgramID)).
		//
		ParentFunc(func(programBranch ag_treeout.Branches) {
			programBranch.Child(ag_format.Instruction("CollectCreatorFee")).
				//
				ParentFunc(func(instructionBranch ag_treeout.Branches) {

					// Parameters of the instruction:
					instructionBranch.Child("Params[len=0]").ParentFunc(func(paramsBranch ag_treeout.Branches) {})

					// Accounts of the instruction:
					instructionBranch.Child("Accounts[len=5]").ParentFunc(func(accountsBranch ag_treeout.Branches) {
						accountsBranch.Child(ag_format.Meta("        creator", inst.AccountMetaSlice.Get(0)))
						accountsBranch.Child(ag_format.Meta("  creator_vault", inst.AccountMetaSlice.Get(1)))
						accountsBranch.Child(ag_format.Meta(" system_program", inst.AccountMetaSlice.Get(2)))
						accountsBranch.Child(ag_format.Meta("event_authority", inst.AccountMetaSlice.Get(3)))
						accountsBranch.Child(ag_format.Meta("        program", inst.AccountMetaSlice.Get(4)))
					})
				})
		})
}

func (obj CollectCreatorFee) MarshalWithEncoder(encoder *ag_binary.Encoder) (err error) {
	return nil
}
func (obj *CollectCreatorFee) UnmarshalWithDecoder(decoder *ag_binary.Decoder) (err error) {
	return nil
}

// NewCollectCreatorFeeInstruction declares a new CollectCreatorFee instruction with the provided parameters and accounts.
func NewCollectCreatorFeeInstruction(
	// Accounts:
	creator ag_solanago.PublicKey,
	creatorVault ag_solanago.PublicKey,
	systemProgram ag_solanago.PublicKey,
	eventAuthority ag_solanago.PublicKey,
	program ag_solanago.PublicKey) *CollectCreatorFee {
	return NewCollectCreatorFeeInstructionBuilder().
		SetCreatorAccount(creator).
		SetCreatorVaultAccount(creatorVault).
		SetSystemProgramAccount(systemProgram).
		SetEventAuthorityAccount(eventAuthority).
		SetProgramAccount(program)
}
